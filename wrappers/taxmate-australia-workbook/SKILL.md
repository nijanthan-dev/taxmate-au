---
name: taxmate-australia-workbook
description: Use when a local helper must route TaxMate Australia workbook requests into the installed skill.
compatibility: Local wrapper for Claude Code, Cowork, and Codex. Requires repo checkout and the TaxMate Australia workbook skill.
metadata:
  internal: true
---

# TaxMate Australia Workbook

## Hard Safety Boundary

- Never lodge, file, submit, transmit, or finalise any tax return, BAS, form, statement, objection, election, payment instruction, or other material with the ATO or any government agency.
- Refuse requests to submit, lodge, file, transmit, finalise, or send prepared material to the ATO.
- Do not help bypass human review, remove `Accountant review` flags, fabricate evidence, hide income, overclaim, or convert preparation output into a lodged position.

## Quick Reference

| Situation | Action |
| --- | --- |
| Workbook request | Use `$taxmate-australia:workbook` when available. |
| Runtime skill is unavailable | Read the source workbook skill from the repo root. |
| Root path is unknown | Resolve `TAXMATE_AUSTRALIA_ROOT` first. |
| User asks to submit | Refuse and keep output prep-only. |

## Common Mistakes

- Treating this wrapper as the workbook rules source.
- Skipping the installed plugin skill when available.
- Reading fallback paths before resolving the repo root.
- Removing review flags from output guidance.

Use the plugin skill `$taxmate-australia:workbook` when available.

Read and follow:

```bash
export TAXMATE_AUSTRALIA_ROOT="${TAXMATE_AUSTRALIA_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
"$TAXMATE_AUSTRALIA_ROOT/skills/workbook/SKILL.md"
```
