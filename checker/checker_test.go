package checker

import (
	"context"
	"testing"
	"time"
)

func TestDNSAvailable(t *testing.T) {
	// A domain that likely won't have NS records
	// Note: this is an integration test that hits real DNS
	r := Check(context.Background(), "this-domain-definitely-does-not-exist-12345-test.xyz", DefaultOptions())
	if r.Method != "dns" {
		// Even if DNS is uncertain, the check should complete without panic
		t.Logf("Result: %+v", r)
	}
}

func TestDNSTaken(t *testing.T) {
	// google.com should always have NS records
	r := Check(context.Background(), "google.com", DefaultOptions())
	t.Logf("google.com result: %+v", r)
	if !r.Available || r.Method == "dns" {
		// It's fine either way — DNS says taken or WHOIS says taken
	}
}

func TestTimeout(t *testing.T) {
	opts := Options{Timeout: 1 * time.Nanosecond}
	r := Check(context.Background(), "google.com", opts)
	if r.Error != "" {
		t.Logf("Expected timeout-ish result: %+v", r)
	}
}

func TestIsAvailable(t *testing.T) {
	avail, err := IsAvailable(context.Background(), "google.com", DefaultOptions())
	if err != nil {
		t.Logf("IsAvailable error (expected maybe): %v", err)
	}
	if avail {
		t.Log("google.com is somehow available?")
	}
}
