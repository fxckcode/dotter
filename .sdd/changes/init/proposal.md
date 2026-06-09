# Proposal: dotter — Domain Availability CLI

## Intent

CLI tool en Go que verifica si un dominio está disponible (no registrado) y sugiere alternativas con otras extensiones (.dev, .io, .tech, etc.) que sí estén libres. Todo por línea de comandos, rápido, en paralelo.

## Scope

### In
- [x] Verificar disponibilidad de un dominio específico (`dotter diego.com`)
- [x] Si no se especifica TLD, asumir `.com` (`dotter diego`)
- [x] Sugerir alternativas libres con otras extensiones populares
- [x] Lista curada de ~20 TLDs populares por defecto
- [x] Modo `--all` para escanear TLDs completos
- [x] Filter por TLDs específicos con `--tlds .io,.dev,.tech`
- [x] Output claro: colorizado, tabla de disponibles vs ocupados
- [x] Soporte JSON con `--json`
- [x] Timeout por consulta (no dejar colgado)
- [x] Testing completo (unit + integration con mocks)

### Out
- [ ] No es WHOIS detail (no mostrar registrant, dates, etc.)
- [ ] No es DNS brute-forcer (subdominios, etc.)
- [ ] No tiene modo interactivo TUI (solo CLI flag-based)
- [ ] No tiene daemon/server mode
- [ ] No soporta IDN (internationalized domain names) inicialmente

## Approach

### Arquitectura

```
dotter/
├── main.go              # Entry point, flag parsing
├── cmd/
│   └── root.go          # CLI command structure
├── checker/
│   ├── checker.go       # Interface + main check logic (DNS rápido + WHOIS fallback)
│   └── checker_test.go
├── tlds/
│   ├── tlds.go          # Lista curada de TLDs + load from file
│   └── tlds_test.go
├── output/
│   ├── output.go        # Formatting: table, colors, JSON
│   └── output_test.go
├── go.mod
├── go.sum
└── README.md
```

### Algoritmo de detección

1. **DNS check rápido** (100-500ms): resolver NS records del dominio
   - Si tiene NS → registrado (ocupado)
   - Si no tiene NS → posiblemente libre, pasar a WHOIS
2. **WHOIS fallback** (1-3s): consultar WHOIS server del TLD
   - `likexian/whois` library en Go
   - Parsear "No match" / "NOT FOUND" / "Domain not found" en respuesta
3. **Paralelo por TLD**: goroutines + errgroup para verificar TLDs alternativos en simultáneo
   - Timeout individual: 5s por consulta
   - Max paralelo: 10 goroutines (semáforo)

### Tech Stack

| Capa | Tecnología | Justificación |
|------|-----------|---------------|
| Lenguaje | Go 1.24+ | Binario estático, rápido, excelente para CLI y concurrencia |
| CLI flags | `flag` stdlib | Sin dependencias externas innecesarias |
| WHOIS | `github.com/likexian/whois` + `whoisparser` | Única lib WHOIS madura en Go |
| Concurrencia | `golang.org/x/sync/errgroup` | Patrón estándar para fan-out goroutines |
| DNS | `net` stdlib | LookupNS es nativo de Go |
| Tests | `testing` stdlib | Go testing convention + table-driven tests |

### Modules Affected

Todos nuevos — estructura planificada arriba.

### Risks

| Riesgo | Mitigación |
|--------|-----------|
| WHOIS rate limiting | Cache de resultados en memoria, delay entre consultas |
| WHOIS format varies por TLD | Usar `whoisparser` + string matching "not found" patterns |
| DNS falsos positivos (domain registrado sin DNS) | WHOIS fallback confirma |
| Timeout total alto | HTTP-style progress, `--timeout` flag global |
| TLDs nuevos sin WHOIS público | Skip silencioso + warning |

## Skill Resolution

- sdd-workflow (este)
- github-repo-management (crear repo)
- test-driven-development (tests)
- git-setup-skill (post-archive)
