# ATO FY2025-26 Scope Summary

Fetched official ATO pages for the tax tracker on 2026-06-23.

## Coverage

- Individual tax return setup, records, deductions, income to declare, tax offsets.
- Employee deductions, including working from home, fixed-rate method, actual-cost method, occupancy expenses, tools, software, phone, internet, depreciating assets, travel, clothing, self-education, memberships, and tax affairs costs.
- ABN/sole-trader topics, including assessable business income, business deductions, business losses, non-commercial-loss indicators, business-versus-hobby boundaries, PSI, home-based business expenses, home-business CGT implications, and small-business depreciation.
- GST/BAS, including registration, GST credits, tax invoices, when to charge GST, effect of GST credits on income-tax deductions, BAS, due dates, PAYG instalments, PAYG withholding, adjustments, TPAR, and GST labels.
- Employer/reporting topics, including Single Touch Payroll, income statements, PAYG withholding annual reporting, and contractor-payment reporting through TPAR.
- Investments, including investment income, shares/units, ETF-like unit holdings, crypto records, rental-property records, CGT events, CGT discount, records, dividend reinvestment plans, trust non-assessable payments, share investing versus trading, and capital losses.
- Calculators/scaffolds, including PAYG estimate-only, BAS arithmetic, FBT gross-up arithmetic, CGT gain/discount arithmetic, SG minimum contributions, and stamp-duty source routing.
- Super, including personal super contributions, deductions/records around super, and employer super guarantee.
- Private health, including private health insurance statements, rebate, Medicare levy, Medicare levy surcharge, thresholds, family/dependants, and tax return treatment.

## Key Tracker Defaults From Sources

- FY2025-26 employee WFH fixed rate is 70 cents per work hour.
- WFH fixed-rate method still needs actual records of hours worked from home and one record for each relevant running-expense type.
- Fixed-rate WFH covers energy, phone, internet, stationery, and computer consumables; separate claims need to stay outside those covered costs.
- Business and employee claims must be separated; private portions must be excluded.
- If GST credits are claimed, income-tax deductions generally use GST-exclusive amounts.
- GST/BAS cannot be assumed nil if GST credits are being claimed; accountant review required.
- Home-based business occupancy claims can affect main-residence CGT, so default tracker should not claim occupancy without accountant review.
- ETF/share records should track income statement data, cost base, disposals, DRP, trust non-assessable payments, CGT events, and capital losses.
- Crypto records should track acquisition, disposal, swaps, wallet/exchange records, cost base, proceeds, and CGT events.
- Rental-property records should separate repairs, capital works/improvements, interest, agent fees, insurance, occupancy/private-use, and disposal-related CGT records.
- PAYG estimates are not payroll withholding tables; payroll, HELP/STSL, bonus, termination, and rounding cases require official ATO tables and accountant/payroll review.
- TPAR and STP/income-statement obligations depend on role, payment type, contractor status, and reporting channel.
- Non-commercial losses and business-versus-hobby questions default to accountant review unless the facts clearly satisfy ATO criteria.

## Known Limits

- This pack is ATO-first. Stamp duty is source-routed to state revenue offices only; it does not include embedded state duty rate tables.
- This pack does not include ASIC, ABR, private insurers, brokers, super funds, app marketplaces, software vendors, or accountant-specific guidance.
- Ten seed URLs returned 404, but the crawl found current moved alternatives for the main affected topics: tools/equipment, investment income, GST credit effects, personal super contributions, and business deductions.
- ATO search returned a 502 for `business or hobby`; use business-losses, assessable-income, PSI, and home-business pages for now.
- This is a source pack, not advice. Any ambiguous or material classification should be marked `Accountant review`.

## Files

- `source_index.json`: machine-readable source index.
- `README.md`: source list with links and ATO last-updated dates when available.
- `raw/`: fetched HTML.
- `text/`: cleaned searchable text.
