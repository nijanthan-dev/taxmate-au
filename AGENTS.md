# Agent Instructions

- Keep changes conservative and source-backed.
- Do not invent tax treatment.
- Mark ambiguous, mixed-use, pre-revenue, home-business, FBT, CGT, GST/BAS, non-commercial-loss, and business-versus-hobby items as `Accountant review` unless sources clearly resolve them.
- Keep tax logic in `skills/research`, `skills/finance-review`, `skills/calculators`, and shared Python runtime.
- Keep `skills/workbook` and `skills/taxpack` as output layers only.
- Do not commit private user tax records.
- Do not commit built binaries from `bin/`.
- Before PR/merge, run `python3 -m py_compile scripts/*.py`, `./scripts/taxmate validate`, `./scripts/taxmate skills generate --check`, `./scripts/taxmate skills audit --check`, `scripts/check-publication-ready.sh`, and run a secret scan.
