---
name: calculators
description: Use when full-runtime TaxMate Australia estimate scaffolds are needed for PAYG, BAS, CGT, FBT, super, or stamp-duty source routing.
compatibility: Full-runtime skill for Claude Code, Cowork, and Codex. Requires repo checkout, bash, and Python 3.9+.
metadata:
  internal: true
  priority: 4
---

# TaxMate Australia Calculators

Runtime requirements:

- Bash
- Python 3.9+

Use this full-runtime skill for bounded tax-prep calculations. It is not professional advice, payroll advice, lodgment support, or confirmation of entitlement. Keep outputs labelled as estimates or scaffolds where applicable.

## Hard Safety Boundary

- Never lodge, file, submit, transmit, or finalise any tax return, BAS, form, statement, objection, election, payment instruction, or other material with the ATO or any government agency.
- Refuse requests to submit, lodge, file, transmit, finalise, or send prepared material to the ATO.
- Do not help bypass human review, remove `Accountant review` flags, fabricate evidence, hide income, overclaim, or convert preparation output into a lodged position.

## Quick Reference

| Situation | Action |
| --- | --- |
| Estimate scaffold is requested | Run the relevant `./scripts/taxmate calc` command. |
| Inputs are missing or uncertain | Keep output as estimate-only and flag review. |
| State/territory rate table is needed | Route to official sources instead of embedding rates. |
| User asks for entitlement confirmation | Refuse final treatment and recommend review. |

## Common Mistakes

- Presenting estimates as payroll, lodgment, or entitlement advice.
- Running CGT math without cost base, proceeds, date, and loss facts.
- Treating BAS arithmetic as BAS lodgment support.
- Removing review flags from calculated output.

Run:

```bash
export TAXMATE_AUSTRALIA_ROOT="${TAXMATE_AUSTRALIA_ROOT:-$(pwd)}"
"$TAXMATE_AUSTRALIA_ROOT/scripts/taxmate" calc bas [flags]
```

## Rules

- PAYG is estimate-only; use official ATO withholding tables for payroll.
- BAS arithmetic is not lodgment advice; GST labels need accountant review.
- CGT output depends on user-supplied cost base, proceeds, dates, losses, and eligibility.
- FBT output assumes taxable value is already known.
- Super guarantee uses ordinary time earnings and the SG rate for the payment date.
- Stamp duty routes to official state/territory sources; no embedded rate tables.
- Keep all review flags in the final answer.
