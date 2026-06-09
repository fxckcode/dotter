package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fxckcode/dotter/checker"
	"github.com/fxckcode/dotter/mcpserv"
	"github.com/fxckcode/dotter/output"
	"github.com/fxckcode/dotter/tlds"
)

const version = "0.2.0"

func main() {
	// Check for MCP subcommand before flag parsing
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		if err := mcpserv.Start(version); err != nil {
			fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
			os.Exit(3)
		}
		return
	}
	// Flags
	showVersion := flag.Bool("version", false, "Print version and exit")
	tldsFlag := flag.String("tlds", "", "Comma-separated TLDs to check")
	allFlag := flag.Bool("all", false, "Scan all extended TLDs")
	jsonFlag := flag.Bool("json", false, "Output as JSON")
	dnsOnly := flag.Bool("dns-only", false, "Skip WHOIS, use DNS only (faster)")
	timeoutSec := flag.Int("timeout", 5, "Timeout per query in seconds")
	concurrency := flag.Int("concurrency", 10, "Max parallel queries")

	flag.Usage = output.PrintUsage
	flag.Parse()

	if *showVersion {
		fmt.Printf("dotter v%s\n", version)
		os.Exit(0)
	}

	// Parse domain argument
	args := flag.Args()
	if len(args) < 1 || args[0] == "" {
		output.PrintError("missing domain argument")
		output.PrintUsage()
		os.Exit(2)
	}

	domain := args[0]
	name, primaryTLD := tlds.NormalizeDomain(domain)

	// If no TLD, assume .com
	if primaryTLD == "" {
		primaryTLD = ".com"
	}
	fullDomain := name + primaryTLD

	// Build TLD list for alternatives
	var scanTLDs []string
	if *tldsFlag != "" {
		parts := strings.Split(*tldsFlag, ",")
		for _, p := range parts {
			scanTLDs = append(scanTLDs, tlds.PrettyTLD(strings.TrimSpace(p)))
		}
	} else if *allFlag {
		scanTLDs = tlds.ExtendedTLDs()
	} else {
		scanTLDs = tlds.DefaultTLDs()
	}

	ctx := context.Background()

	// Step 1: Check primary domain
	checkOpts := checker.Options{
		Timeout: time.Duration(*timeoutSec) * time.Second,
		DNSOnly: *dnsOnly,
	}
	primary := checker.Check(ctx, fullDomain, checkOpts)

	// Step 2: Scan alternative TLDs (skip the primary TLD)
	var altTLDs []string
	for _, t := range scanTLDs {
		if t != primaryTLD {
			altTLDs = append(altTLDs, t)
		}
	}

	scanOpts := tlds.ScanOptions{
		TLDs:        altTLDs,
		TimeoutSec:  *timeoutSec,
		Concurrency: *concurrency,
		DNSOnly:     *dnsOnly,
	}
	alternatives := tlds.Scan(ctx, name, scanOpts)

	// Step 3: Output
	if *jsonFlag {
		json, err := output.FormatJSON(name, primary, alternatives)
		if err != nil {
			output.PrintError(fmt.Sprintf("JSON formatting failed: %v", err))
			os.Exit(3)
		}
		fmt.Println(json)
	} else {
		fmt.Print(output.FormatResults(name, primary, alternatives))
	}

	// Exit code
	availableCount := countAltAvailable(alternatives)
	if primary.Available || availableCount > 0 {
		os.Exit(0)
	}
	os.Exit(1)
}

func countAltAvailable(results []tlds.ScanResult) int {
	count := 0
	for _, r := range results {
		if r.Available {
			count++
		}
	}
	return count
}
