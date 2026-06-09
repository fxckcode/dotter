# Spec: dotter — Domain Availability CLI

## Requirements

### Functional

- [ ] REQ-F1: Aceptar un nombre de dominio como argumento posicional (`dotter diego.com`)
- [ ] REQ-F2: Si el argumento no incluye TLD, asumir `.com` (`dotter diego` → check `diego.com`)
- [ ] REQ-F3: Verificar si el dominio está registrado (ocupado) o disponible (libre)
- [ ] REQ-F4: Mostrar una tabla de TLDs alternativos con estado (libre/ocupado)
- [ ] REQ-F5: Flag `--tlds` para especificar TLDs a verificar (`dotter diego --tlds .io,.dev`)
- [ ] REQ-F6: Flag `--all` para verificar todos los TLDs curados (~50)
- [ ] REQ-F7: Flag `--json` para output en JSON (pipeable)
- [ ] REQ-F8: Flag `--timeout` para tiempo máximo por consulta (default 5s)
- [ ] REQ-F9: Flag `--concurrency` para límite de goroutines paralelas (default 10)
- [ ] REQ-F10: Exit code 0 si al menos un TLD alternativo está libre, 1 si no hay libres
- [ ] REQ-F11: Version flag (`dotter --version`)

### Non-Functional

- [ ] REQ-NF1: Consulta única (dominio principal) debe responder en < 3s típicamente
- [ ] REQ-NF2: Escaneo de TLDs alternativos debe ser paralelo
- [ ] REQ-NF3: Timeout individual por consulta no debe exceder `--timeout`
- [ ] REQ-NF4: Output legible en terminal (colores, sin dependencias externas)
- [ ] REQ-NF5: Binario estático sin runtime dependencies
- [ ] REQ-NF6: Tests: >= 80% code coverage en lógica core

## Scenarios

### Happy Path — Dominio libre

```
$ dotter di-e-go-2342342342
✗ di-e-go-2342342342.com is TAKEN

  Alternative TLDs for "di-e-go-2342342342":
  TLD       Status
  ───────── ───────
  .dev      ✓ FREE
  .io       ✓ FREE
  .tech     ✓ FREE
  .app      ✓ FREE
  .ai       ✓ FREE
  .net      ✓ FREE
  .org      ✓ FREE
  .co       ✓ FREE
  .me       ✓ FREE
  .xyz      ✓ FREE

  → 10 TLDs available! Try: di-e-go-2342342342.dev
```

### Happy Path — Dominio ocupado sin alternativas

```
$ dotter google
✗ google.com is TAKEN

  Alternative TLDs for "google":
  TLD       Status
  ───────── ───────
  .dev      ✗ TAKEN
  .io       ✗ TAKEN
  .tech     ✗ TAKEN
  .app      ✗ TAKEN
  .ai       ✗ TAKEN
  .net      ✗ TAKEN
  .org      ✗ TAKEN
  .co       ✗ TAKEN
  .me       ✗ TAKEN
  .xyz      ✗ TAKEN

  → No TLDs available for "google"
```

### Happy Path — JSON output

```
$ dotter mysite --json
{"domain":"mysite.com","available":false,"alternatives":[
  {"tld":".dev","available":true},
  {"tld":".io","available":false}
]}
```

### Flag — TLDs específicos

```
$ dotter myproject --tlds .io,.dev,.tech,.app
✗ myproject.com is TAKEN

  Alternative TLDs for "myproject":
  TLD       Status
  ───────── ───────
  .dev      ✓ FREE
  .io       ✗ TAKEN
  .tech     ✓ FREE
  .app      ✓ FREE

  → 3 TLDs available!
```

### Edge Cases

- Dominio vacío → error con Usage
- Caracteres inválidos → error descriptivo
- TLD inexistente → warning "unknown TLD, checking anyway"
- Timeout general → mensaje "some queries timed out"
- Sin conectividad → mensaje claro "check your connection"

### Error Cases

- Dominio no especificado: `dotter` → "Usage: dotter <domain> [flags]"
- Red caída: "error: DNS lookup failed: connection refused"
- WHOIS server no responde: fallback a "unable to determine" en vez de crash

## Interface Changes

### CLI Interface

```
Usage: dotter <domain> [flags]

Arguments:
  domain                  Domain name to check (e.g., diego.com or just "diego" for .com)

Flags:
  --tlds string          Comma-separated TLDs to check (default: curated list)
  --all                  Check all known TLDs
  --json                 Output as JSON
  --timeout duration     Timeout per query (default 5s)
  --concurrency int      Max parallel queries (default 10)
  --version              Show version
  -h, --help             Show help
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Al menos un TLD alternativo disponible |
| 1 | No hay TLDs alternativos disponibles |
| 2 | Error de argumentos |
| 3 | Error de red/checks |
