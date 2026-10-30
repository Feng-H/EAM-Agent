package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ems/backend/internal/agent/dto"
	"github.com/ems/backend/internal/middleware"
	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/database"
	"github.com/ems/backend/pkg/trace"
	"github.com/gin-gonic/gin"
)

type mcpSession struct {
	ID       string
	User     model.User
	Scopes   []string
	MsgChan  chan string
	DoneChan chan struct{}
}

type mcpSessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*mcpSession
}

var globalMCPManager = &mcpSessionManager{
	sessions: make(map[string]*mcpSession),
}

func (m *mcpSessionManager) createSession(user model.User, scopes []string) *mcpSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := trace.GenerateTraceID()
	sess := &mcpSession{
		ID:       id,
		User:     user,
		Scopes:   scopes,
		MsgChan:  make(chan string, 32),
		DoneChan: make(chan struct{}),
	}
	m.sessions[id] = sess
	return sess
}

func (m *mcpSessionManager) getSession(id string) (*mcpSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sess, ok := m.sessions[id]
	return sess, ok
}

func (m *mcpSessionManager) removeSession(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess, ok := m.sessions[id]; ok {
		close(sess.DoneChan)
		delete(m.sessions, id)
	}
}

// MCPSSEHandler handles GET /api/v1/mcp/sse for streaming MCP responses
func (ctrl *AgentController) MCPSSEHandler(c *gin.Context) {
	userID, _, ok := requireAuth(c)
	if !ok {
		return
	}

	var user model.User
	if err := database.GetDB().First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	scopes, _ := middleware.GetAPIKeyScopes(c)
	sess := globalMCPManager.createSession(user, scopes)
	defer globalMCPManager.removeSession(sess.ID)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	// Emit initial endpoint event per MCP spec
	endpointURL := fmt.Sprintf("/api/v1/mcp/message?session_id=%s", sess.ID)
	fmt.Fprintf(c.Writer, "event: endpoint\ndata: %s\n\n", endpointURL)
	c.Writer.Flush()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			// Keep alive ping
			fmt.Fprintf(c.Writer, ": ping\n\n")
			c.Writer.Flush()
		case msg, ok := <-sess.MsgChan:
			if !ok {
				return
			}
			fmt.Fprintf(c.Writer, "event: message\ndata: %s\n\n", msg)
			c.Writer.Flush()
		}
	}
}

// MCPMessageHandler handles POST /api/v1/mcp/message?session_id=...
func (ctrl *AgentController) MCPMessageHandler(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id query parameter is required"})
		return
	}

	sess, ok := globalMCPManager.getSession(sessionID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP session not found or expired"})
		return
	}

	var rpcReq struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}

	if err := c.ShouldBindJSON(&rpcReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON-RPC payload: " + err.Error()})
		return
	}

	c.Status(http.StatusAccepted)

	go func() {
		var responseData map[string]interface{}

		switch rpcReq.Method {
		case "initialize":
			responseData = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"result": map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{},
					},
					"serverInfo": map[string]interface{}{
						"name":    "eam-agent-sse",
						"version": "1.0.0",
					},
				},
			}
		case "ping":
			responseData = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      rpcReq.ID,
				"result":  map[string]interface{}{},
			}
		case "tools/list":
			tools, err := ctrl.agentService.ListTools(sess.User)
			if err != nil {
				responseData = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      rpcReq.ID,
					"error": map[string]interface{}{
						"code":    -32603,
						"message": err.Error(),
					},
				}
			} else {
				type mcpToolItem struct {
					Name        string      `json:"name"`
					Description string      `json:"description"`
					InputSchema interface{} `json:"inputSchema"`
				}
				var mcpTools []mcpToolItem
				for _, t := range tools {
					mcpTools = append(mcpTools, mcpToolItem{
						Name:        t.Name,
						Description: t.Description,
						InputSchema: t.InputSchema,
					})
				}
				responseData = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      rpcReq.ID,
					"result": map[string]interface{}{
						"tools": mcpTools,
					},
				}
			}
		case "tools/call":
			var callParams struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal(rpcReq.Params, &callParams); err != nil {
				responseData = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      rpcReq.ID,
					"error": map[string]interface{}{
						"code":    -32602,
						"message": "Invalid params: " + err.Error(),
					},
				}
			} else {
				req := dto.CallToolRequest{
					Name:      callParams.Name,
					Arguments: callParams.Arguments,
				}
				callRes, err := ctrl.agentService.CallTool(sess.User, &req, sess.Scopes)
				if err != nil {
					responseData = map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      rpcReq.ID,
						"result": map[string]interface{}{
							"content": []map[string]interface{}{
								{"type": "text", "text": "Error: " + err.Error()},
							},
							"isError": true,
						},
					}
				} else {
					var textContent string
					switch v := callRes.Content.(type) {
					case string:
						textContent = v
					default:
						if b, err := json.MarshalIndent(v, "", "  "); err == nil {
							textContent = string(b)
						} else {
							textContent = fmt.Sprintf("%v", v)
						}
					}
					responseData = map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      rpcReq.ID,
						"result": map[string]interface{}{
							"content": []map[string]interface{}{
								{"type": "text", "text": textContent},
							},
							"isError": callRes.IsError,
						},
					}
				}
			}
		default:
			if rpcReq.ID != nil {
				responseData = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      rpcReq.ID,
					"error": map[string]interface{}{
						"code":    -32601,
						"message": fmt.Sprintf("Method '%s' not found", rpcReq.Method),
					},
				}
			}
		}

		if responseData != nil {
			if jsonBytes, err := json.Marshal(responseData); err == nil {
				select {
				case sess.MsgChan <- string(jsonBytes):
				case <-sess.DoneChan:
				case <-time.After(5 * time.Second):
				}
			}
		}
	}()
}
