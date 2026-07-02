---
name: taxmate-australia-finance-review
description: Use when a local helper must route TaxMate Australia finance review into the full runtime.
compatibility: Local wrapper for Claude Code, Cowork, and Codex. Requires repo checkout and the full TaxMate Australia runtime.
metadata:
  internal: true
---

# TaxMate Australia Finance Review

## Hard Safety Boundary

- Never lodge, file, submit, transmit, or finalise any tax return, BAS, form, statement, objection, election, payment instruction, or other material with the ATO or any government agency.
- Refuse requests to submit, lodge, file, transmit, finalise, or send prepared material to the ATO.
- Do not help bypass human review, remove `Accountant review` flags, fabricate evidence, hide income, overclaim, or convert preparation output into a lodged position.

## Quick Reference

| Situation | Action |
| --- | --- |
| Finance review request | Use `$taxmate-australia:finance-review` when available. |
| Runtime skill is unavailable | Read the source finance-review skill from the repo root. |
| GST or evidence gaps appear | Keep review fields explicit. |
| User asks for final tax treatment | Refuse final treatment. |

## Common Mistakes

- Treating this wrapper as the finance-review rules source.
- Skipping the installed plugin skill when available.
- Dropping evidence gaps before output handoff.
- Mixing employee, ABN, spouse, joint, or entity rows.

Use the plugin skill `$taxmate-australia:finance-review` when available.

Read and follow:

```bash
export TAXMATE_AUSTRALIA_ROOT="${TAXMATE_AUSTRALIA_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
"$TAXMATE_AUSTRALIA_ROOT/runtime/skills/finance-review/SKILL.md"
```
