// Package tlds provides TLD listing and domain scanning.
package tlds

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/fxckcode/dotter/checker"
	"golang.org/x/sync/errgroup"
)

// ScanResult represents the availability of a domain with a specific TLD.
type ScanResult struct {
	TLD       string `json:"tld"`
	Domain    string `json:"domain"`
	Available bool   `json:"available"`
	Error     string `json:"error,omitempty"`
	Method    string `json:"method,omitempty"`
}

// ScanOptions configures the TLD scanning behaviour.
type ScanOptions struct {
	TLDs        []string // TLDs to check (including the dot)
	TimeoutSec  int      // timeout per query in seconds
	Concurrency int      // max parallel goroutines
	DNSOnly     bool     // skip WHOIS (faster but less accurate)
}

// DefaultScanOptions returns sensible defaults.
func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		TLDs:        DefaultTLDs(),
		TimeoutSec:  5,
		Concurrency: 10,
	}
}

// Scan performs a parallel scan of a domain name across multiple TLDs.
// The `name` parameter should be the bare name without TLD (e.g., "diego").
// It checks name + each TLD in the list.
func Scan(ctx context.Context, name string, opts ScanOptions) []ScanResult {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 10
	}
	if opts.TimeoutSec <= 0 {
		opts.TimeoutSec = 5
	}

	results := make([]ScanResult, 0, len(opts.TLDs))
	var mu sync.Mutex

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(opts.Concurrency)

	for _, tld := range opts.TLDs {
		tld := tld
		domain := name + tld
		g.Go(func() error {
			checkOpts := checker.Options{
				Timeout: timeDur(opts.TimeoutSec),
				DNSOnly: opts.DNSOnly,
			}
			r := checker.Check(ctx, domain, checkOpts)
			mu.Lock()
			results = append(results, ScanResult{
				TLD:       tld,
				Domain:    domain,
				Available: r.Available,
				Error:     r.Error,
				Method:    r.Method,
			})
			mu.Unlock()
			return nil
		})
	}

	_ = g.Wait() // ignore errors — each result carries its own error
	return results
}

func timeDur(sec int) time.Duration {
	return time.Duration(sec) * time.Second
}

// PrettyTLD returns the TLD with a leading dot, ensuring consistent formatting.
func PrettyTLD(tld string) string {
	tld = strings.TrimSpace(tld)
	if !strings.HasPrefix(tld, ".") {
		tld = "." + tld
	}
	return strings.ToLower(tld)
}

// NormalizeDomain takes a domain input and returns the bare name + primary TLD.
// E.g., "diego.com" → ("diego", ".com")
// E.g., "diego" → ("diego", "")
func NormalizeDomain(input string) (name string, primaryTLD string) {
	input = strings.TrimSpace(strings.ToLower(input))
	if idx := strings.LastIndex(input, "."); idx > 0 {
		return input[:idx], input[idx:]
	}
	return input, ""
}

// Available returns the list of scan results that are available (free).
func Available(results []ScanResult) []ScanResult {
	var out []ScanResult
	for _, r := range results {
		if r.Available {
			out = append(out, r)
		}
	}
	return out
}

// Taken returns the list of scan results that are taken (registered).
func Taken(results []ScanResult) []ScanResult {
	var out []ScanResult
	for _, r := range results {
		if !r.Available {
			out = append(out, r)
		}
	}
	return out
}
