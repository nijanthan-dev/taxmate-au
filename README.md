# TaxMate AU

TaxMate AU is an Australian tax-prep plugin for Codex. It combines official ATO source refresh, conservative tax treatment rules, transaction review, calculation scaffolds, and accountant-facing output workflows.

TaxMate AU is a preparation aid, not professional tax, legal, accounting, financial, BAS-agent, or registered-tax-agent advice. It is not affiliated with, sponsored by, endorsed by, or approved by the Australian Taxation Office or any government agency. Read [DISCLAIMER.md](DISCLAIMER.md) before using it.

Ambiguous, material, mixed-use, pre-revenue, home-business, FBT, CGT, GST/BAS, non-commercial-loss, and business-versus-hobby items should stay marked `Accountant review` unless the facts and current official guidance clearly resolve them.

## What It Does

- Refreshes and searches an official ATO source pack.
- Reviews structured transaction CSVs for claim candidates, GST candidates, evidence gaps, and accountant-review flags.
- Runs bounded calculators for PAYG, BAS, CGT, FBT, super, and stamp-duty source routing.
- Defines workbook and taxpack output workflows without duplicating tax logic.

## Plugin Layout

- `.codex-plugin/plugin.json`: Codex plugin metadata.
- `skills/research`: official ATO research and conservative tax treatment.
- `skills/finance-review`: transaction and evidence review.
- `skills/calculators`: bounded calculation scaffolds.
- `skills/workbook`: accountant-facing workbook workflow.
- `skills/taxpack`: handoff pack and future PDF/form workflow.
- `bin/`: shared Go binaries.
- `cmd/`, `internal/`: shared Go backend.
- `data/ato_knowledge_base/`: official ATO source pack.
- `wrappers/`: compatibility skills for agents that do not yet load local plugin skills directly.

## Agent Support

Codex is the first supported runtime. The skill files are plain Markdown with frontmatter and the backend is a portable Go CLI, so Claude or other agent runtimes can add their own thin wrappers without changing tax logic.

## Install

Set `TAXMATE_AU_ROOT` to the plugin root:

```bash
export TAXMATE_AU_ROOT="/path/to/taxmate-au"
```

Build:

```bash
cd "$TAXMATE_AU_ROOT"
go test ./...
go build -o bin/taxmate-au-refresh ./cmd/taxmate-au-refresh
go build -o bin/taxmate-au-validate ./cmd/taxmate-au-validate
go build -o bin/taxmate-au-finance ./cmd/taxmate-au-finance
go build -o bin/taxmate-au-calc ./cmd/taxmate-au-calc
```

Validate:

```bash
"$TAXMATE_AU_ROOT/bin/taxmate-au-validate"
```

## Boundaries

Tax treatment belongs in `research`, `finance-review`, and `calculators`. Output skills such as `workbook` and `taxpack` consume reviewed data only; they must not invent independent tax rules.

## Sources

The bundled source pack is ATO-first. Stamp duty is source-routed to official state or territory revenue offices. Non-ATO commercial sources are out of scope unless a user explicitly asks for them.

ATO and Commonwealth material remains subject to the notices and terms published by the relevant official source. TaxMate AU must not imply official endorsement.
