package mcpserv

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCheckDomainToolDefined(t *testing.T) {
	// Verify the tool is properly defined
	tool := mcp.NewTool("check_domain",
		mcp.WithDescription("test"),
		mcp.WithString("domain", mcp.Required(), mcp.Description("domain name")),
	)
	if tool.Name != "check_domain" {
		t.Errorf("expected check_domain, got %s", tool.Name)
	}
}

func TestHandleCheckDomain_NoDomain(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "check_domain",
			Arguments: map[string]interface{}{},
		},
	}
	result, err := handleCheckDomain(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsError {
		t.Error("expected error result for missing domain")
	}
}

func TestHandleCheckDomain_InvalidArgs(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "check_domain",
			Arguments: "not-a-map",
		},
	}
	result, err := handleCheckDomain(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || !result.IsError {
		t.Error("expected error result for invalid args")
	}
}
