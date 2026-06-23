# Contributing

TaxMate AU is conservative by design. Contributions must preserve source-backed behaviour and clear accountant-review boundaries.

Do not present TaxMate AU output as tax, legal, accounting, financial, BAS-agent, or registered-tax-agent advice. TaxMate AU is not affiliated with, sponsored by, endorsed by, or approved by the Australian Taxation Office or any government agency.

## Before Opening A PR

Run:

```bash
go test ./...
mkdir -p bin
go build -o bin/taxmate-au-refresh ./cmd/taxmate-au-refresh
go build -o bin/taxmate-au-validate ./cmd/taxmate-au-validate
go build -o bin/taxmate-au-finance ./cmd/taxmate-au-finance
go build -o bin/taxmate-au-calc ./cmd/taxmate-au-calc
bin/taxmate-au-validate
scripts/check-publication-ready.sh
```

## Rules

- Prefer official ATO sources.
- Do not loosen accountant-review defaults for high-risk topics without source-backed tests.
- Keep output skills separate from tax logic.
- Do not commit user tax records or private documents.
- Do not commit built binaries.
- Keep plugin docs portable; avoid private machine paths in public docs and plugin skills.
