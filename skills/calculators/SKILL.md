---
name: calculators
description: Run TaxMate AU bounded calculator scaffolds for PAYG estimates, BAS arithmetic, CGT gains, FBT gross-up arithmetic, super guarantee minimums, and stamp-duty source routing. Use only for estimates and accountant-review workflows.
metadata:
  priority: 4
  promptSignals:
    phrases:
      - "calculate PAYG"
      - "calculate CGT"
      - "BAS estimate"
      - "FBT gross-up"
      - "super guarantee"
      - "stamp duty"
---

# TaxMate AU Calculators

Use this skill for bounded tax-prep calculations. It is not professional advice, payroll advice, lodgment support, or confirmation of entitlement. Keep outputs labelled as estimates or scaffolds where applicable.

Run:

```bash
export TAXMATE_AU_ROOT="${TAXMATE_AU_ROOT:-$(pwd)}"
"$TAXMATE_AU_ROOT/bin/taxmate-au-calc" <bas|super|fbt|cgt|payg|stamp-duty> [flags]
```

## Rules

- PAYG is estimate-only; use official ATO withholding tables for payroll.
- BAS arithmetic is not lodgment advice; GST labels need accountant review.
- CGT output depends on user-supplied cost base, proceeds, dates, losses, and eligibility.
- FBT output assumes taxable value is already known.
- Super guarantee uses ordinary time earnings and the SG rate for the payment date.
- Stamp duty routes to official state/territory sources; no embedded rate tables.
- Keep all review flags in the final answer.

## Invocation

Use `$taxmate-au:calculators`.
