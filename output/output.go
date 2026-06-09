// Package output formats domain check results for display.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fxckcode/dotter/checker"
	"github.com/fxckcode/dotter/tlds"
)

// FormatResult returns a string describing the primary domain check result.
func FormatResult(r checker.Result) string {
	if r.Error != "" {
		return fmt.Sprintf("⚠  %s — %s", r.Domain, r.Error)
	}
	if r.Available {
		return fmt.Sprintf("✓  %s is FREE!", r.Domain)
	}
	return fmt.Sprintf("✗  %s is TAKEN", r.Domain)
}

// FormatResults renders the primary result + alternative TLD scan as a table.
func FormatResults(domain string, primary checker.Result, alternatives []tlds.ScanResult) string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(FormatResult(primary))
	b.WriteString("\n\n")

	availableCount := 0
	for _, r := range alternatives {
		if r.Available {
			availableCount++
		}
	}

	if availableCount > 0 {
		fmt.Fprintf(&b, "  Alternative TLDs for \"%s\":\n", domain)
		b.WriteString("  TLD         Status\n")
		b.WriteString("  ─────────── ───────\n")

		for _, r := range alternatives {
			if r.Available {
				fmt.Fprintf(&b, "  %-12s ✓ FREE\n", r.TLD)
			} else if r.Error != "" {
				fmt.Fprintf(&b, "  %-12s ⚠ %s\n", r.TLD, r.Error)
			} else {
				fmt.Fprintf(&b, "  %-12s ✗ TAKEN\n", r.TLD)
			}
		}
		b.WriteString("\n")

		suggestion := firstAvailable(alternatives)
		if suggestion != "" {
			fmt.Fprintf(&b, "  → %d TLDs available! Try: %s%s\n",
				availableCount, domain, suggestion)
		} else {
			fmt.Fprintf(&b, "  → %d TLDs available!\n", availableCount)
		}
	} else {
		fmt.Fprintf(&b, "  Alternative TLDs for \"%s\":\n", domain)
		b.WriteString("  TLD         Status\n")
		b.WriteString("  ─────────── ───────\n")

		for _, r := range alternatives {
			if r.Error != "" {
				fmt.Fprintf(&b, "  %-12s ⚠ %s\n", r.TLD, r.Error)
			} else {
				fmt.Fprintf(&b, "  %-12s ✗ TAKEN\n", r.TLD)
			}
		}
		b.WriteString("\n")
		fmt.Fprintf(&b, "  → No TLDs available for \"%s\"\n", domain)
	}

	return b.String()
}

// JSONOutput formats results as JSON.
type JSONOutput struct {
	Domain       string           `json:"domain"`
	Available    bool             `json:"available"`
	Method       string           `json:"method,omitempty"`
	Error        string           `json:"error,omitempty"`
	Alternatives []tlds.ScanResult `json:"alternatives,omitempty"`
}

// FormatJSON serializes results as indented JSON.
func FormatJSON(domain string, primary checker.Result, alternatives []tlds.ScanResult) (string, error) {
	out := JSONOutput{
		Domain:       domain,
		Available:    primary.Available,
		Method:       primary.Method,
		Alternatives: alternatives,
	}
	if primary.Error != "" {
		out.Error = primary.Error
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(data), nil
}

// FormatSummary provides a one-liner for quick output.
func FormatSummary(name string, primary checker.Result, alternatives []tlds.ScanResult) string {
	availCount := countAvailable(alternatives)
	if primary.Available {
		return fmt.Sprintf("%s: ✓ FREE! Also %d alternative TLDs available", name, availCount)
	}

	if availCount > 0 {
		suggestion := firstAvailable(alternatives)
		return fmt.Sprintf("%s: ✗ TAKEN — but %d alt TLDs free! Try: %s%s",
			name, availCount, name, suggestion)
	}
	return fmt.Sprintf("%s: ✗ TAKEN — no alt TLDs free", name)
}

func firstAvailable(results []tlds.ScanResult) string {
	for _, r := range results {
		if r.Available {
			return r.TLD
		}
	}
	return ""
}

func countAvailable(results []tlds.ScanResult) int {
	count := 0
	for _, r := range results {
		if r.Available {
			count++
		}
	}
	return count
}

// PrintError writes an error message to stderr.
func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
}

// PrintUsage prints the CLI usage help.
func PrintUsage() {
	fmt.Println(`Dotter — Domain Availability Checker

Usage:
  dotter <domain> [flags]

Arguments:
  domain    Domain to check (e.g., "example.com" or just "example" for .com)

Flags:
  --tlds string    Comma-separated TLDs to scan (default: curated ~20 TLDs)
  --all            Scan all extended TLDs (~70)
  --json           Output as JSON
  --dns-only       Skip WHOIS, use DNS only (faster but less accurate)
  --timeout int    Timeout per query in seconds (default 5)
  --concurrency int  Max parallel queries (default 10)
  --version        Print version and exit
  -h, --help       Print this help

Examples:
  dotter myproject
  dotter myproject.io
  dotter myproject --tlds .io,.dev,.tech
  dotter myproject --all --json
  dotter myproject --dns-only

Exit codes:
  0   At least one alternative TLD is available
  1   No alternative TLDs available
  2   Argument error
  3   Check error`)
}
