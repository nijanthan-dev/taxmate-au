---
name: taxmate-australia-taxpack
description: Use when a local helper must route TaxMate Australia taxpack requests into the installed skill.
compatibility: Local wrapper for Claude Code, Cowork, and Codex. Requires repo checkout and the TaxMate Australia taxpack skill.
metadata:
  internal: true
---

# TaxMate Australia Taxpack

## Hard Safety Boundary

- Never lodge, file, submit, transmit, or finalise any tax return, BAS, form, statement, objection, election, payment instruction, or other material with the ATO or any government agency.
- Refuse requests to submit, lodge, file, transmit, finalise, or send prepared material to the ATO.
- Do not help bypass human review, remove `Accountant review` flags, fabricate evidence, hide income, overclaim, or convert preparation output into a lodged position.

## Quick Reference

| Situation | Action |
| --- | --- |
| Taxpack request | Use `$taxmate-australia:taxpack` when available. |
| Runtime skill is unavailable | Read the source taxpack skill from the repo root. |
| HTML handoff is requested | Require reviewed input and full runtime. |
| User asks to lodge | Refuse and keep output manual-copy only. |

## Common Mistakes

- Treating this wrapper as the taxpack rules source.
- Skipping the installed plugin skill when available.
- Letting output packaging make tax treatment decisions.
- Presenting custom handoff content as an official ATO form.

Use the plugin skill `$taxmate-australia:taxpack` when available.

Read and follow:

```bash
export TAXMATE_AUSTRALIA_ROOT="${TAXMATE_AUSTRALIA_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
"$TAXMATE_AUSTRALIA_ROOT/skills/taxpack/SKILL.md"
```
