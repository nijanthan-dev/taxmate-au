---
name: finance-review
description: Review Australian tax records and transaction CSVs for accountant handoff using TaxMate AU. Use for receipts, invoices, bank exports, ETF records, super statements, private-health statements, GST candidates, claim percentages, evidence health, and accountant-review queues.
metadata:
  priority: 5
  promptSignals:
    phrases:
      - "receipt"
      - "invoice"
      - "CSV"
      - "bank export"
      - "expense tracker"
      - "GST credit"
      - "accountant handoff"
      - "tax spreadsheet"
---

# TaxMate AU Finance Review

Use this skill to review structured financial records before workbook or accountant output. It is a preparation aid, not professional advice or official lodgment support. It does not replace TaxMate AU research; refresh ATO pages before final tax treatment.

Run:

```bash
export TAXMATE_AU_ROOT="${TAXMATE_AU_ROOT:-$(pwd)}"
"$TAXMATE_AU_ROOT/bin/taxmate-au-finance" --input "<records.csv>" --format markdown --output "<review.md>"
```

For machine-readable output:

```bash
"$TAXMATE_AU_ROOT/bin/taxmate-au-finance" --input "<records.csv>" --format json --output "<review.json>"
```

Accepted headers include `date`, `description`, `amount`, `gst`, `owner`, `purpose`, `evidence`, `abn`, `category`, `account`, `asset`, `units`, and `type`.

## Rules

- Keep employee and ABN/business items separate.
- Keep spouse, joint, and entity ownership explicit.
- Do not mark BAS as nil if GST credits or GST collected are present.
- Treat private, mixed-use, pre-revenue, capital, home-business, FBT, CGT, PSI, business-vs-hobby, and non-commercial-loss cases as `Accountant review` unless the facts are clear.
- Refresh ATO pages listed in `ato_refresh_queries` before finalising treatment.
- Pass reviewed JSON/Markdown to output skills; do not let output skills make new tax calls.

## Invocation

Use `$taxmate-au:finance-review`.
