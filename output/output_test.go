package output

import (
	"strings"
	"testing"

	"github.com/fxckcode/dotter/checker"
	"github.com/fxckcode/dotter/tlds"
)

func TestFormatResultTaken(t *testing.T) {
	r := checker.Result{Domain: "google.com", Available: false, Method: "dns"}
	got := FormatResult(r)
	if !strings.Contains(got, "TAKEN") {
		t.Errorf("Expected TAKEN in output, got: %s", got)
	}
}

func TestFormatResultFree(t *testing.T) {
	r := checker.Result{Domain: "example-free.io", Available: true, Method: "whois"}
	got := FormatResult(r)
	if !strings.Contains(got, "FREE") {
		t.Errorf("Expected FREE in output, got: %s", got)
	}
}

func TestFormatResultError(t *testing.T) {
	r := checker.Result{Domain: "test.com", Available: false, Method: "whois", Error: "timeout"}
	got := FormatResult(r)
	if !strings.Contains(got, "timeout") {
		t.Errorf("Expected error message in output, got: %s", got)
	}
}

func TestFormatJSON(t *testing.T) {
	primary := checker.Result{Domain: "test.com", Available: false, Method: "dns"}
	alts := []tlds.ScanResult{
		{TLD: ".dev", Domain: "test.dev", Available: true, Method: "dns"},
		{TLD: ".io", Domain: "test.io", Available: false, Method: "dns"},
	}
	json, err := FormatJSON("test", primary, alts)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}
	if !strings.Contains(json, `"domain": "test"`) {
		t.Errorf("JSON missing domain field: %s", json)
	}
	if !strings.Contains(json, `"available": false`) {
		t.Errorf("JSON missing available field: %s", json)
	}
}

func TestFormatSummary(t *testing.T) {
	tests := []struct {
		name     string
		primary  checker.Result
		alts     []tlds.ScanResult
		contains string
	}{
		{
			name:     "freetest",
			primary:  checker.Result{Domain: "freetest.com", Available: true, Method: "dns"},
			alts:     []tlds.ScanResult{{TLD: ".dev", Available: true}},
			contains: "FREE",
		},
		{
			name:     "taken",
			primary:  checker.Result{Domain: "taken.com", Available: false, Method: "dns"},
			alts:     []tlds.ScanResult{{TLD: ".dev", Available: false}},
			contains: "TAKEN",
		},
	}
	for _, tt := range tests {
		got := FormatSummary(tt.name, tt.primary, tt.alts)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("FormatSummary(%q) = %q, want containing %q", tt.name, got, tt.contains)
		}
	}
}

func TestFormatResults(t *testing.T) {
	primary := checker.Result{Domain: "example.com", Available: false, Method: "dns"}
	alts := []tlds.ScanResult{
		{TLD: ".dev", Domain: "example.dev", Available: true, Method: "dns"},
		{TLD: ".io", Domain: "example.io", Available: false, Method: "dns"},
		{TLD: ".tech", Domain: "example.tech", Available: true, Method: "dns"},
	}
	got := FormatResults("example", primary, alts)
	if !strings.Contains(got, "TAKEN") {
		t.Errorf("Expected TAKEN in results")
	}
	if !strings.Contains(got, "2 TLDs available") {
		t.Errorf("Expected count of available TLDs, got: %s", got)
	}
	if !strings.Contains(got, "FREE") {
		t.Errorf("Expected FREE markers in results")
	}
}

func TestPrintUsage(t *testing.T) {
	// Should not panic, return something reasonable
	PrintUsage()
}

func TestFirstAvailable(t *testing.T) {
	want := firstAvailable([]tlds.ScanResult{
		{TLD: ".io", Available: false},
		{TLD: ".dev", Available: true},
		{TLD: ".tech", Available: true},
	})
	if want != ".dev" {
		t.Errorf("firstAvailable = %q, want .dev", want)
	}
}
