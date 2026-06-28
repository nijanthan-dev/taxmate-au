# Generated Topic Inputs

Workbook and taxpack are output layers only. They must consume reviewed classifications from topic skills and must not invent tax treatment.

- Preserve `Accountant review` flags.
- If input fields conflict, explicit or review-like `Accountant review` wins over stale evidence, used, ATO-label, skipped, status-kind, tab-kind, or styling fields.
- If explanation fields are blank, review queues must fall back to row number/status instead of rendering blank review items.
- Preserve valid falsey output values such as numeric `0` and boolean `false`.
- Review-feedback fixes must cover parsed input and direct renderer/workbook-row paths before another review is requested.
- Preserve source URLs and checked-at dates.
- Do not turn raw transactions into lodging-ready claims from source extracts alone.
