---
name: taxmate-au
description: Compatibility wrapper for TaxMate AU Research. Use for Australian ATO tax-prep research, source refresh, conservative treatment, and accountant-review flags.
---

# TaxMate AU

Use the plugin skill `$taxmate-au:research` when available.

Resolve the local plugin root from `TAXMATE_AU_ROOT`, or from a colocated checkout when this wrapper is copied into a larger plugin bundle:

```bash
export TAXMATE_AU_ROOT="${TAXMATE_AU_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
```

Read:

```bash
"$TAXMATE_AU_ROOT/skills/research/SKILL.md"
```

Follow that skill exactly. This wrapper exists for Codex installations that load `~/.agents/skills` before local plugin-cache skills.
