package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol structures
type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type MCPContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type MCPCallToolResult struct {
	Content []MCPContentItem `json:"content"`
	IsError bool             `json:"isError,omitempty"`
}

// EAM backend DTOs
type EAMToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

type EAMListToolsResponse struct {
	Tools []EAMToolDefinition `json:"tools"`
}

type EAMCallToolResponse struct {
	Content interface{} `json:"content"`
	IsError bool        `json:"is_error"`
}

type Config struct {
	BaseURL string
	Token   string
}

func main() {
	baseURL := os.Getenv("EAM_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	token := os.Getenv("EAM_TOKEN")

	cfg := Config{
		BaseURL: baseURL,
		Token:   token,
	}

	scanner := bufio.NewScanner(os.Stdin)
	// Support larger payloads (up to 4MB)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 4*1024*1024)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			sendError(nil, -32700, "Parse error: "+err.Error())
			continue
		}

		handleRequest(&req, cfg, client)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdio: %v\n", err)
	}
}

func handleRequest(req *JSONRPCRequest, cfg Config, client *http.Client) {
	// Notifications don't need a response
	if req.ID == nil && strings.HasPrefix(req.Method, "notifications/") {
		return
	}

	switch req.Method {
	case "initialize":
		result := map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "eam-agent-mcp",
				"version": "1.0.0",
			},
		}
		sendResult(req.ID, result)

	case "ping":
		sendResult(req.ID, map[string]interface{}{})

	case "tools/list":
		tools, err := fetchEAMTools(cfg, client)
		if err != nil {
			sendError(req.ID, -32603, fmt.Sprintf("Failed to fetch tools from EAM: %v", err))
			return
		}
		sendResult(req.ID, map[string]interface{}{
			"tools": tools,
		})

	case "tools/call":
		var callParams struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &callParams); err != nil {
			sendError(req.ID, -32602, "Invalid params: "+err.Error())
			return
		}

		res, err := callEAMTool(cfg, client, callParams.Name, callParams.Arguments)
		if err != nil {
			sendResult(req.ID, MCPCallToolResult{
				Content: []MCPContentItem{
					{Type: "text", Text: "Error executing tool: " + err.Error()},
				},
				IsError: true,
			})
			return
		}
		sendResult(req.ID, res)

	default:
		if req.ID != nil {
			sendError(req.ID, -32601, fmt.Sprintf("Method '%s' not found", req.Method))
		}
	}
}

func fetchEAMTools(cfg Config, client *http.Client) ([]MCPTool, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("EAM_TOKEN is required. Please set EAM_TOKEN environment variable")
	}

	reqURL := cfg.BaseURL + "/api/v1/agent/tools"
	httpReq, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.Token)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("EAM server returned status %d: %s", resp.StatusCode, string(body))
	}

	var listResp EAMListToolsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode tools: %w", err)
	}

	var mcpTools []MCPTool
	for _, t := range listResp.Tools {
		mcpTools = append(mcpTools, MCPTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.InputSchema,
		})
	}
	return mcpTools, nil
}

func callEAMTool(cfg Config, client *http.Client, name string, args map[string]interface{}) (*MCPCallToolResult, error) {
	if cfg.Token == "" {
		return nil, fmt.Errorf("EAM_TOKEN is required. Please set EAM_TOKEN environment variable")
	}

	reqBody := map[string]interface{}{
		"name":      name,
		"arguments": args,
	}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	reqURL := cfg.BaseURL + "/api/v1/agent/tools/call"
	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+cfg.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &MCPCallToolResult{
			Content: []MCPContentItem{
				{Type: "text", Text: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))},
			},
			IsError: true,
		}, nil
	}

	var callResp EAMCallToolResponse
	if err := json.Unmarshal(bodyBytes, &callResp); err != nil {
		return &MCPCallToolResult{
			Content: []MCPContentItem{
				{Type: "text", Text: string(bodyBytes)},
			},
			IsError: false,
		}, nil
	}

	var outputText string
	switch v := callResp.Content.(type) {
	case string:
		outputText = v
	default:
		prettyJSON, err := json.MarshalIndent(v, "", "  ")
		if err == nil {
			outputText = string(prettyJSON)
		} else {
			outputText = fmt.Sprintf("%v", v)
		}
	}

	return &MCPCallToolResult{
		Content: []MCPContentItem{
			{Type: "text", Text: outputText},
		},
		IsError: callResp.IsError,
	}, nil
}

func sendResult(id interface{}, result interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	send(resp)
}

func sendError(id interface{}, code int, message string) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	send(resp)
}

func send(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON-RPC: %v\n", err)
		return
	}
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}
