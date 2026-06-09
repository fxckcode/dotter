// Package mcpserv provides an MCP server for domain checking.
// Run with: dotter mcp
package mcpserv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fxckcode/dotter/checker"
	"github.com/fxckcode/dotter/output"
	"github.com/fxckcode/dotter/tlds"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Start runs the MCP stdio server.
func Start(version string) error {
	s := server.NewMCPServer(
		"dotter",
		version,
		server.WithInstructions("Domain availability checker. Use check_domain to check if a domain is available and get alternative TLD suggestions."),
		server.WithLogging(),
	)

	// Tool: check_domain
	checkDomainTool := mcp.NewTool("check_domain",
		mcp.WithDescription("Check if a domain is available (not registered) and suggest alternative TLDs that are free"),
		mcp.WithString("domain",
			mcp.Required(),
			mcp.Description("Domain to check (e.g., \"example.com\" or just \"example\" for .com)"),
		),
		mcp.WithString("tlds",
			mcp.Description("Comma-separated TLDs to scan for alternatives (default: curated list)"),
		),
		mcp.WithBoolean("all",
			mcp.Description("Scan all extended TLDs (~70) instead of the default curated list"),
		),
		mcp.WithBoolean("dns_only",
			mcp.Description("Skip WHOIS, use DNS only (faster but less accurate)"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout per query in seconds (default 5)"),
		),
		mcp.WithNumber("concurrency",
			mcp.Description("Max parallel queries (default 10)"),
		),
	)

	s.AddTool(checkDomainTool, handleCheckDomain)

	return server.ServeStdio(s)
}

func handleCheckDomain(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rawArgs, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments format"), nil
	}

	domain, ok := rawArgs["domain"].(string)
	if !ok || domain == "" {
		return mcp.NewToolResultError("missing required argument: domain"), nil
	}

	tldsList := ""
	if v, ok := rawArgs["tlds"].(string); ok {
		tldsList = v
	}

	allTLDs := false
	if v, ok := rawArgs["all"].(bool); ok {
		allTLDs = v
	}

	dnsOnly := false
	if v, ok := rawArgs["dns_only"].(bool); ok {
		dnsOnly = v
	}

	timeoutSec := 5
	if v, ok := rawArgs["timeout"].(float64); ok {
		timeoutSec = int(v)
	}

	concurrency := 10
	if v, ok := rawArgs["concurrency"].(float64); ok {
		concurrency = int(v)
	}

	// Normalize domain
	name, primaryTLD := tlds.NormalizeDomain(domain)
	if primaryTLD == "" {
		primaryTLD = ".com"
	}
	fullDomain := name + primaryTLD

	// Build TLD list for alternatives
	var scanTLDs []string
	switch {
	case tldsList != "":
		parts := strings.Split(tldsList, ",")
		for _, p := range parts {
			scanTLDs = append(scanTLDs, tlds.PrettyTLD(strings.TrimSpace(p)))
		}
	case allTLDs:
		scanTLDs = tlds.ExtendedTLDs()
	default:
		scanTLDs = tlds.DefaultTLDs()
	}

	// Check primary domain
	checkOpts := checker.Options{
		Timeout: secs(timeoutSec),
		DNSOnly: dnsOnly,
	}
	primary := checker.Check(ctx, fullDomain, checkOpts)

	// Scan alternative TLDs
	var altTLDs []string
	for _, t := range scanTLDs {
		if t != primaryTLD {
			altTLDs = append(altTLDs, t)
		}
	}

	scanOpts := tlds.ScanOptions{
		TLDs:        altTLDs,
		TimeoutSec:  timeoutSec,
		Concurrency: concurrency,
		DNSOnly:     dnsOnly,
	}
	alternatives := tlds.Scan(ctx, name, scanOpts)

	// Build output
	availableCount := 0
	var freeTLDs []string
	for _, r := range alternatives {
		if r.Available {
			availableCount++
			freeTLDs = append(freeTLDs, r.TLD)
		}
	}

	result := map[string]interface{}{
		"domain":          domain,
		"normalized":      fullDomain,
		"available":       primary.Available,
		"method":          primary.Method,
		"alt_tlds_available": availableCount,
		"alternative_tlds": formatAlternatives(alternatives),
		"summary":         output.FormatSummary(name, primary, alternatives),
	}
	if primary.Error != "" {
		result["error"] = primary.Error
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("JSON marshal error: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

func formatAlternatives(results []tlds.ScanResult) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(results))
	for _, r := range results {
		entry := map[string]interface{}{
			"tld":       r.TLD,
			"domain":    r.Domain,
			"available": r.Available,
		}
		if r.Error != "" {
			entry["error"] = r.Error
		}
		out = append(out, entry)
	}
	return out
}

func secs(n int) time.Duration {
	return time.Duration(n) * time.Second
}
