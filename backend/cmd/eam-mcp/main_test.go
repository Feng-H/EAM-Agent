package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleInitialize(t *testing.T) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	cfg := Config{BaseURL: "http://localhost:8080", Token: "test_token"}
	client := &http.Client{}

	// Should not panic and send valid response
	handleRequest(&req, cfg, client)
}

func TestFetchEAMToolsWithMock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/agent/tools" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer valid_token" {
			t.Errorf("Unexpected auth header: %s", auth)
		}

		resp := EAMListToolsResponse{
			Tools: []EAMToolDefinition{
				{
					Name:        "search_equipment",
					Description: "Search equipment by keyword",
					InputSchema: map[string]interface{}{"type": "object"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	cfg := Config{
		BaseURL: mockServer.URL,
		Token:   "valid_token",
	}

	tools, err := fetchEAMTools(cfg, http.DefaultClient)
	if err != nil {
		t.Fatalf("fetchEAMTools failed: %v", err)
	}

	if len(tools) != 1 || tools[0].Name != "search_equipment" {
		t.Fatalf("Unexpected tools returned: %+v", tools)
	}
}

func TestCallEAMToolWithMock(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/agent/tools/call" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		var reqBody map[string]interface{}
		json.NewDecoder(r.Body).Decode(&reqBody)

		if reqBody["name"] != "search_equipment" {
			t.Errorf("Unexpected tool name: %v", reqBody["name"])
		}

		resp := EAMCallToolResponse{
			Content: map[string]interface{}{
				"found": true,
				"name":  "Exposure Machine #3",
			},
			IsError: false,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	cfg := Config{
		BaseURL: mockServer.URL,
		Token:   "valid_token",
	}

	res, err := callEAMTool(cfg, http.DefaultClient, "search_equipment", map[string]interface{}{"keyword": "3"})
	if err != nil {
		t.Fatalf("callEAMTool failed: %v", err)
	}

	if res.IsError {
		t.Fatalf("Expected isError = false")
	}

	if len(res.Content) == 0 || res.Content[0].Type != "text" {
		t.Fatalf("Unexpected content format: %+v", res.Content)
	}
}
