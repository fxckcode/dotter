# Tasks: dotter — Initial Implementation

## Dependency Order
T1 ← T2, T3 ← T4 (T1 es scaffold, T2+T3 son independientes, T4 depende de T2+T3)

## Tasks

### T1: Project Scaffold + GitHub Repo (AFK)
- **Files:**
  - Create: `go.mod`
  - Create: `main.go` (esqueleto con flag parsing)
  - Create: `.gitignore`
  - Create: `README.md` (básico)
  - GitHub: crear repo `dotter` público
  - Git: init, commit, push initial
- **Acceptance:** `go build ./...` compila sin errores, repo existe en GitHub
- **Dependencies:** ninguna
- **Estimated size:** small

### T2: Core Checker — DNS + WHOIS (AFK)
- **Files:**
  - Create: `checker/checker.go`
  - Create: `checker/checker_test.go`
  - Modify: `go.mod` (add `likexian/whois`, `golang.org/x/sync`)
- **Acceptance:**
  - `Check(domain string) (bool, error)` — retorna `available=true/false`
  - DNS check rápido: `net.LookupNS`
  - WHOIS fallback con `likexian/whois`
  - Tests: mock DNS + mock WHOIS, table-driven
  - `go test ./checker/...` pasa
- **Dependencies:** T1
- **Estimated size:** medium

### T3: TLD List + Alternative Scanning (AFK)
- **Files:**
  - Create: `tlds/tlds.go`
  - Create: `tlds/tlds_test.go`
  - Create: `tlds/data.go` (lista curada de ~50 TLDs)
- **Acceptance:**
  - `GetTLDs(mode string) []string` — retorna lista default o todas
  - `ParseTLDsFlag(s string) []string` — parsea `--tlds .io,.dev`
  - `ScanAlternatives(name string, tlds []string, opts Options) []Result`
  - Goroutines paralelas con errgroup + semáforo de concurrencia
  - Timeout por consulta
  - Tests: mock checker, verificar paralelismo
  - `go test ./tlds/...` pasa
- **Dependencies:** T1
- **Estimated size:** medium

### T4: CLI Output + Integration (AFK)
- **Files:**
  - Create: `output/output.go`
  - Create: `output/output_test.go`
  - Modify: `main.go` (wiring completo)
- **Acceptance:**
  - `FormatTable(results []Result, domain string)` — tabla colorizada
  - `FormatJSON(results []Result, domain string)` — JSON output
  - `FormatText(res Result)` — línea simple para dominio principal
  - Exit codes correctos (0=hay libres, 1=no hay, 2=arg error)
  - `go build ./...` compila
  - `go test ./...` pasa
- **Dependencies:** T2, T3
- **Estimated size:** medium
