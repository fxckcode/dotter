# Dotter

> Domain availability checker — CLI tool to check if a domain is free and find alternative TLDs.

```
$ dotter myproject
✗ myproject.com is TAKEN

  Alternative TLDs for "myproject":
  TLD         Status
  ─────────── ───────
  .dev        ✓ FREE
  .io         ✓ FREE
  .tech       ✓ FREE
  .app        ✓ FREE
  .ai         ✓ FREE
  .net        ✓ FREE

  → 6 TLDs available! Try: myproject.dev
```

## Install

```bash
go install github.com/fxckcode/dotter@latest
```

Or download a pre-built binary from the [releases page](https://github.com/fxckcode/dotter/releases).

## Usage

```bash
# Check a specific domain
dotter example.com

# Just give a name (assumes .com)
dotter myproject

# Scan alternative TLDs
dotter myproject --tlds .io,.dev,.tech
dotter myproject --all

# JSON output (pipeable)
dotter myproject --json

# Fast mode (DNS only, no WHOIS)
dotter myproject --dns-only
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--tlds` | (curated ~20) | Comma-separated TLDs to scan |
| `--all` | false | Scan all extended TLDs (~70) |
| `--json` | false | Output as JSON |
| `--dns-only` | false | Skip WHOIS (faster but less accurate) |
| `--timeout` | 5s | Timeout per query |
| `--concurrency` | 10 | Max parallel queries |
| `--version` | — | Show version |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | At least one TLD is available |
| 1 | No TLDs available |
| 2 | Argument error |
| 3 | Check/network error |

## How It Works

1. **DNS check** (fast) — resolves NS records. If NS exist → domain is taken.
2. **WHOIS fallback** (accurate) — if DNS is inconclusive, queries WHOIS servers.
3. **Parallel TLD scan** — goroutines check multiple TLDs simultaneously.

## License

MIT
