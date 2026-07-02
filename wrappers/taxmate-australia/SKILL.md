---
name: taxmate-australia
description: Use when a local helper must route general TaxMate Australia requests into the full runtime.
compatibility: Local wrapper for Claude Code, Cowork, and Codex. Requires repo checkout and the full TaxMate Australia runtime.
metadata:
  internal: true
---

# TaxMate Australia

## Hard Safety Boundary

- Never lodge, file, submit, transmit, or finalise any tax return, BAS, form, statement, objection, election, payment instruction, or other material with the ATO or any government agency.
- Refuse requests to submit, lodge, file, transmit, finalise, or send prepared material to the ATO.
- Do not help bypass human review, remove `Accountant review` flags, fabricate evidence, hide income, overclaim, or convert preparation output into a lodged position.

## Quick Reference

| Situation | Action |
| --- | --- |
| General TaxMate request | Load the full-runtime research skill. |
| Root path is unknown | Resolve `TAXMATE_AUSTRALIA_ROOT` first. |
| Requested skill exists in plugin runtime | Prefer `$taxmate-australia:*` runtime skill. |
| User asks to lodge or finalise | Refuse and keep output prep-only. |

## Common Mistakes

- Treating this wrapper as the source of tax rules.
- Skipping root resolution before reading runtime skills.
- Using wrapper fallback when the plugin runtime skill is available.
- Removing `Accountant review` or prep-only language.

Use the plugin skill `$taxmate-australia:research` when available.

Resolve the local plugin root from `TAXMATE_AUSTRALIA_ROOT`, or from a colocated checkout when this wrapper is copied into a larger plugin bundle:

```bash
export TAXMATE_AUSTRALIA_ROOT="${TAXMATE_AUSTRALIA_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
```

Read:

```bash
"$TAXMATE_AUSTRALIA_ROOT/runtime/skills/research/SKILL.md"
```

Follow that skill exactly. This wrapper exists for Codex installations that load `~/.agents/skills` before local plugin-cache skills.
