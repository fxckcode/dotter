package tlds

import (
	"testing"
)

func TestDefaultTLDs(t *testing.T) {
	tlds := DefaultTLDs()
	if len(tlds) == 0 {
		t.Fatal("DefaultTLDs returned empty list")
	}
	// Should include common TLDs
	expected := map[string]bool{".com": false, ".io": false, ".dev": false, ".app": false}
	for _, tld := range tlds {
		if _, ok := expected[tld]; ok {
			expected[tld] = true
		}
	}
	for tld, found := range expected {
		if !found {
			t.Errorf("Expected %s in default TLDs", tld)
		}
	}
}

func TestExtendedTLDs(t *testing.T) {
	ext := ExtendedTLDs()
	if len(ext) <= len(DefaultTLDs()) {
		t.Error("ExtendedTLDs should have more than DefaultTLDs")
	}
}

func TestPrettyTLD(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{".com", ".com"},
		{"com", ".com"},
		{".io", ".io"},
		{"io", ".io"},
		{".DEV", ".dev"},
		{" Dev ", ".dev"},
	}
	for _, tt := range tests {
		got := PrettyTLD(tt.input)
		if got != tt.expected {
			t.Errorf("PrettyTLD(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestNormalizeDomain(t *testing.T) {
	tests := []struct {
		input       string
		wantName    string
		wantPrimary string
	}{
		{"diego.com", "diego", ".com"},
		{"example.io", "example", ".io"},
		{"diego", "diego", ""},
		{"test.co.uk", "test.co", ".uk"}, // multi-part TLD — last dot split
		{"  MY-SITE.dev  ", "my-site", ".dev"},
	}
	for _, tt := range tests {
		name, primary := NormalizeDomain(tt.input)
		if name != tt.wantName || primary != tt.wantPrimary {
			t.Errorf("NormalizeDomain(%q) = (%q, %q), want (%q, %q)",
				tt.input, name, primary, tt.wantName, tt.wantPrimary)
		}
	}
}

func TestAvailableTaken(t *testing.T) {
	results := []ScanResult{
		{Domain: "test.com", Available: true},
		{Domain: "test.io", Available: false},
		{Domain: "test.dev", Available: true},
	}
	avail := Available(results)
	if len(avail) != 2 {
		t.Errorf("Available: got %d, want 2", len(avail))
	}
	taken := Taken(results)
	if len(taken) != 1 {
		t.Errorf("Taken: got %d, want 1", len(taken))
	}
}

func TestScanOptionsDefaults(t *testing.T) {
	opts := DefaultScanOptions()
	if opts.Concurrency != 10 {
		t.Errorf("Default concurrency: got %d, want 10", opts.Concurrency)
	}
	if opts.TimeoutSec != 5 {
		t.Errorf("Default timeout: got %d, want 5", opts.TimeoutSec)
	}
	if len(opts.TLDs) == 0 {
		t.Error("Default TLDs should not be empty")
	}
}
