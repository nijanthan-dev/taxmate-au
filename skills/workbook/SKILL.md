---
name: workbook
description: Create accountant-facing Excel workbook outputs from reviewed TaxMate AU data. Use for taxpayer/spouse-separated and combined tax expense workbooks, BAS/GST summaries, ETF/super/private-health tabs, evidence checklists, and accountant-review queues.
metadata:
  priority: 4
  promptSignals:
    phrases:
      - "Excel workbook"
      - "spreadsheet"
      - "tax expense workbook"
      - "accountant spreadsheet"
      - "spouse and me"
      - "combined summary"
---

# TaxMate AU Workbook

Use this skill for output rendering only. It creates draft accountant-facing artifacts, not lodgment-ready advice. It must consume reviewed data from `$taxmate-au:finance-review` and tax treatment from `$taxmate-au:research`; it must not create new tax logic.

## Workbook Shape

Default tabs:

- `Read Me`
- `Primary Taxpayer - Employee`
- `Primary Taxpayer - ABN`
- `Spouse or Partner - Employee`
- `Spouse or Partner - ABN`
- `Joint / Combined`
- `GST BAS`
- `ETF / Investments`
- `Super`
- `Private Health`
- `Evidence Checklist`
- `Accountant Review`
- `Source URLs`

## Rules

- Separate employee and ABN/business expenses.
- Separate primary taxpayer, spouse or partner, joint, and entity records.
- Preserve gross, GST, GST-exclusive, claim %, claim amount, evidence, source URL, and review status.
- Put ambiguous rows in `Accountant Review`.
- Do not silently drop rows.
- Do not mark BAS nil when GST credits or GST collected exist.
- Use formulas only for transparent totals; keep source rows visible.

## Invocation

Use `$taxmate-au:workbook`.
