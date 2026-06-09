# Archive: dotter — Domain Availability CLI

## Summary
- **Proposal:** `.sdd/changes/init/proposal.md`
- **Spec:** `.sdd/changes/init/spec.md`
- **Tasks:** 4 tasks, all AFK
- **Verification:** PASS

## Files Changed
| File | Purpose |
|------|---------|
| `main.go` | CLI entry point, flag parsing, orchestration |
| `checker/checker.go` | DNS + WHOIS availability check |
| `checker/checker_test.go` | Tests for checker (integration via real DNS) |
| `tlds/tlds.go` | Parallel TLD scanning with errgroup |
| `tlds/data.go` | Curated list of ~50 TLDs |
| `tlds/tlds_test.go` | Tests for TLD normalization, filtering |
| `output/output.go` | Table, JSON, summary formatters + CLI usage |
| `output/output_test.go` | Tests for all output formats |
| `go.mod` / `go.sum` | Go module dependencies |
| `README.md` | Project documentation |
| `.gitignore` | Go binary ignores |
| `.gitattributes` | LF normalization |
| `.github/workflows/ci.yml` | Go CI build + test |
| `.sdd/changes/init/*.md` | SDD documentation |

## What Was Learned
- WHOIS library `likexian/whois` works well but some TLDs don't have public WHOIS servers → good fallback to DNS-only mode
- `errgroup.SetLimit` is effective for bounding concurrent WHOIS/DNS queries
- Go's `net.LookupNS` returns NXDOMAIN for unregistered domains quickly (<100ms)
- WHOIS queries can be slow (3-5s) for some TLDs — DNS-first strategy is critical

## Repo Metadata
- **Repo:** github.com/fxckcode/dotter
- **Topics:** cli, go, domain-checker, whois, dns, tld, devtools
- **License:** MIT

## Next Steps
- Add `go install` instructions once GitHub releases are tagged
- Optionally add a `--interactive` mode for browsing results
- Consider adding WHOIS result caching (in-memory TTL-based)
