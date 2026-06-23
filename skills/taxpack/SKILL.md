---
name: taxpack
description: Prepare accountant handoff packs and future PDF/form outputs from reviewed TaxMate AU data. Use for summary PDFs, checklists, source bundles, and later tax-form drafts from accountant-reviewed structured data.
metadata:
  priority: 3
  promptSignals:
    phrases:
      - "tax pack"
      - "PDF tax form"
      - "accountant pack"
      - "handoff pack"
      - "tax summary PDF"
---

# TaxMate AU Taxpack

Use this skill for final handoff packaging only. It creates draft preparation artifacts, not official lodgment forms or professional advice. It consumes reviewed data from research, finance review, calculators, and workbook outputs.

## Rules

- Do not fill final tax forms from raw records.
- Do not make independent tax treatment decisions.
- Keep source URLs, evidence status, and accountant-review flags visible.
- Treat PDF/form filling as draft preparation only unless the user explicitly asks for a final accountant-ready copy.
- For any official lodgment form, require reviewed structured data and exact income-year labels.

## Invocation

Use `$taxmate-au:taxpack`.
