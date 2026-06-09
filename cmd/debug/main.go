package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fxckcode/dotter/checker"
)

func main() {
	fmt.Fprintln(os.Stderr, "DEBUG: starting")

	timeoutSec := flag.Int("timeout", 5, "")
	dnsOnly := flag.Bool("dns-only", false, "")
	flag.Parse()

	fmt.Fprintf(os.Stderr, "DEBUG: timeout=%d dns-only=%v args=%v\n", *timeoutSec, *dnsOnly, flag.Args())

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "DEBUG: no args")
		os.Exit(2)
	}

	domain := args[0]
	fmt.Fprintf(os.Stderr, "DEBUG: checking %s\n", domain)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutSec)*time.Second)
	defer cancel()

	opts := checker.Options{
		Timeout: time.Duration(*timeoutSec) * time.Second,
		DNSOnly: *dnsOnly,
	}

	fmt.Fprintln(os.Stderr, "DEBUG: calling Check...")
	r := checker.Check(ctx, domain, opts)
	fmt.Fprintf(os.Stderr, "DEBUG: result: %+v\n", r)
}
