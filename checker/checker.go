// Package checker provides domain availability checking via DNS and WHOIS.
package checker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// Result represents the availability status for a single domain.
type Result struct {
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Error     string `json:"error,omitempty"`
	Method    string `json:"method"` // "dns", "whois", or "error"
}

// Options controls checker behaviour.
type Options struct {
	Timeout     time.Duration
	DNSOnly     bool // skip WHOIS fallback, rely only on DNS
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Timeout: 5 * time.Second,
	}
}

// Check determines if a domain is available (not registered).
// It first tries a fast DNS NS lookup. If NS records exist, the domain is taken.
// If no NS records, it falls back to WHOIS for confirmation.
func Check(ctx context.Context, domain string, opts Options) Result {
	if opts.Timeout == 0 {
		opts = DefaultOptions()
	}

	ctx, cancel := context.WithTimeout(ctx, opts.Timeout)
	defer cancel()

	// Phase 1: DNS check (fast)
	status := dnsCheck(ctx, domain)
	if status == statusTaken {
		return Result{Domain: domain, Available: false, Method: "dns"}
	}

	if opts.DNSOnly {
		// DNS-only: if DNS didn't confirm taken, we're uncertain
		return Result{Domain: domain, Available: false, Method: "dns", Error: "dns check inconclusive"}
	}

	// Phase 2: WHOIS fallback
	return whoisCheck(ctx, domain)
}

// dnsCheck returns status based on NS record lookup.
// - statusTaken: NS records exist → domain is registered
// - statusUncertain: no NS records or error → need WHOIS
func dnsCheck(ctx context.Context, domain string) checkStatus {
	resolver := net.Resolver{}
	ns, err := resolver.LookupNS(ctx, domain)
	if err != nil {
		// NXDOMAIN or other DNS errors → uncertain, try WHOIS
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			if dnsErr.IsNotFound {
				return statusUncertain
			}
			if dnsErr.IsTimeout {
				return statusUncertain
			}
		}
		return statusUncertain
	}

	// Filter out empty or temporary NS records that some DNS servers return
	for _, n := range ns {
		if strings.TrimSpace(n.Host) != "" && !strings.HasSuffix(strings.ToLower(n.Host), ".placeholder.invalid.") {
			return statusTaken
		}
	}
	return statusUncertain
}

// whoisCheck performs a WHOIS lookup and parses the result.
// Returns true (available) if WHOIS returns "not found" / domain not registered.
func whoisCheck(ctx context.Context, domain string) Result {
	done := make(chan Result, 1)

	go func() {
		raw, err := whois.Whois(domain)
		if err != nil {
			done <- Result{
				Domain: domain, Available: false, Method: "whois",
				Error: fmt.Sprintf("whois query failed: %v", err),
			}
			return
		}

		_, err = whoisparser.Parse(raw)
		if err != nil {
			if errors.Is(err, whoisparser.ErrNotFoundDomain) {
				// Domain is NOT found in WHOIS → it's available
				done <- Result{Domain: domain, Available: true, Method: "whois"}
				return
			}
			if errors.Is(err, whoisparser.ErrReservedDomain) ||
				errors.Is(err, whoisparser.ErrPremiumDomain) {
				// Reserved or premium → can be registered but may cost more
				done <- Result{Domain: domain, Available: true, Method: "whois", Error: err.Error()}
				return
			}
			if errors.Is(err, whoisparser.ErrDomainLimitExceed) {
				done <- Result{Domain: domain, Available: false, Method: "whois",
					Error: "whois rate limited"}
				return
			}
			// Some other parse error — try text-based fallback
			available := textFallback(raw)
			done <- Result{Domain: domain, Available: available, Method: "whois"}
			return
		}

		// Parse succeeded → domain IS registered
		done <- Result{Domain: domain, Available: false, Method: "whois"}
	}()

	select {
	case r := <-done:
		return r
	case <-ctx.Done():
		return Result{Domain: domain, Available: false, Method: "whois", Error: "timeout"}
	}
}

// textFallback does a simple string match for "not found" patterns
// in raw WHOIS text when the parser fails.
func textFallback(raw string) bool {
	lower := strings.ToLower(raw)
	notFoundKeys := []string{
		"no data found",
		"no match",
		"not found",
		"not registered",
		"domain name not known",
		"no entries found",
		"no matching record",
		"object does not exist",
		"is free",
		"status: free",
		"available",
	}
	for _, key := range notFoundKeys {
		if strings.Contains(lower, key) {
			return true
		}
	}
	return false
}

// checkStatus is an internal status indicator.
type checkStatus int

const (
	statusUncertain checkStatus = iota
	statusTaken
)

// IsAvailable determines if a full domain (e.g., "example.com") is available.
// This is a convenience wrapper around Check.
func IsAvailable(ctx context.Context, domain string, opts Options) (bool, error) {
	r := Check(ctx, domain, opts)
	if r.Error != "" {
		return r.Available, errors.New(r.Error)
	}
	return r.Available, nil
}
