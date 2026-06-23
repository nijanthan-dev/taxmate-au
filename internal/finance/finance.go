package finance

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Mode string

const (
	ModeStrict   Mode = "strict"
	ModeAssisted Mode = "assisted"
	ModeReview   Mode = "review"
)

type Input struct {
	Path string
	Mode Mode
}

type Transaction struct {
	Row         int               `json:"row"`
	Date        string            `json:"date,omitempty"`
	Description string            `json:"description,omitempty"`
	Amount      float64           `json:"amount"`
	GST         float64           `json:"gst,omitempty"`
	Direction   string            `json:"direction"`
	Owner       string            `json:"owner,omitempty"`
	Account     string            `json:"account,omitempty"`
	Category    string            `json:"category,omitempty"`
	Purpose     string            `json:"purpose,omitempty"`
	Evidence    string            `json:"evidence,omitempty"`
	ABN         string            `json:"abn,omitempty"`
	Source      string            `json:"source,omitempty"`
	Asset       string            `json:"asset,omitempty"`
	Units       float64           `json:"units,omitempty"`
	Raw         map[string]string `json:"raw,omitempty"`
}

type Finding struct {
	Row                int      `json:"row"`
	Owner              string   `json:"owner,omitempty"`
	Description        string   `json:"description,omitempty"`
	Amount             float64  `json:"amount"`
	Direction          string   `json:"direction"`
	Bucket             string   `json:"bucket"`
	TaxTreatment       string   `json:"tax_treatment"`
	ClaimPercent       float64  `json:"claim_percent"`
	ClaimAmount        float64  `json:"claim_amount"`
	GSTCreditCandidate bool     `json:"gst_credit_candidate"`
	GSTCreditAmount    float64  `json:"gst_credit_amount,omitempty"`
	Confidence         string   `json:"confidence"`
	Reasons            []string `json:"reasons"`
	RecordsNeeded      []string `json:"records_needed,omitempty"`
	AccountantReview   bool     `json:"accountant_review"`
}

type HealthCheck struct {
	Name     string   `json:"name"`
	Passed   bool     `json:"passed"`
	Severity string   `json:"severity"`
	Detail   string   `json:"detail,omitempty"`
	Rows     []int    `json:"rows,omitempty"`
	Advice   []string `json:"advice,omitempty"`
}

type SummaryLine struct {
	Owner        string  `json:"owner"`
	Bucket       string  `json:"bucket"`
	Treatment    string  `json:"treatment"`
	GrossAmount  float64 `json:"gross_amount"`
	ClaimAmount  float64 `json:"claim_amount"`
	GSTCandidate float64 `json:"gst_candidate"`
	Rows         int     `json:"rows"`
}

type Report struct {
	GeneratedAt    string        `json:"generated_at"`
	Mode           Mode          `json:"mode"`
	Source         string        `json:"source"`
	Transactions   []Transaction `json:"transactions"`
	Findings       []Finding     `json:"findings"`
	Summary        []SummaryLine `json:"summary"`
	BASSummary     BASSummary    `json:"bas_summary"`
	ScenarioChecks []Scenario    `json:"scenario_checks"`
	HealthChecks   []HealthCheck `json:"health_checks"`
	Caveats        []string      `json:"caveats"`
	ATOQueries     []string      `json:"ato_refresh_queries"`
}

type BASSummary struct {
	BusinessExpenseGross  float64 `json:"business_expense_gross"`
	GSTCreditCandidate    float64 `json:"gst_credit_candidate"`
	BusinessIncomeGross   float64 `json:"business_income_gross"`
	GSTCollectedCandidate float64 `json:"gst_collected_candidate"`
	NilBASLikely          bool    `json:"nil_bas_likely"`
	ReviewNote            string  `json:"review_note"`
}

type Scenario struct {
	Name       string  `json:"name"`
	BaseAmount float64 `json:"base_amount"`
	WhatIf     string  `json:"what_if"`
	Result     float64 `json:"result"`
	ReviewNote string  `json:"review_note"`
}

var (
	normaliseRE = regexp.MustCompile(`[^a-z0-9]+`)

	privateHealthTerms    = []string{"private health", "health insurance", "medicare levy surcharge"}
	investmentTerms       = []string{"broker", "trading platform", "etf", "managed fund", "index fund", "shares", "dividend", "distribution", "amit", "drp", "cgt"}
	softwareTerms         = []string{"developer program", "developer tools", "software", "hosting", "domain", "api", "subscription", "saas", "source control", "build tool"}
	employmentIncomeTerms = []string{"salary", "wage", "payroll", "employer", "payg"}
	investmentIncomeTerms = []string{"dividend", "distribution", "interest", "etf", "broker", "trading platform", "amit"}
)

func AnalyzeCSV(input Input) (*Report, error) {
	if input.Mode == "" {
		input.Mode = ModeStrict
	}
	if input.Mode != ModeStrict && input.Mode != ModeAssisted && input.Mode != ModeReview {
		return nil, fmt.Errorf("invalid mode %q", input.Mode)
	}
	file, err := os.Open(input.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	transactions, err := readCSV(file)
	if err != nil {
		return nil, err
	}
	report := &Report{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Mode:        input.Mode,
		Source:      input.Path,
		Caveats: []string{
			"This is a preparation aid, not tax advice.",
			"Use official ATO refresh before answering final tax questions.",
			"Ambiguous, mixed-use, pre-revenue, private, capital, and GST items stay flagged for accountant review.",
		},
		ATOQueries: []string{
			"working from home fixed rate method 2025-26",
			"claiming GST credits tax invoices",
			"effect of GST credits on income tax deductions",
			"deductions for digital product expenses",
			"business losses pre revenue expenses",
			"shares funds trusts ETF annual tax statement AMIT CGT",
			"private health insurance rebate Medicare levy surcharge",
			"personal super contributions notice of intent",
		},
		Transactions: transactions,
	}
	for _, tx := range transactions {
		report.Findings = append(report.Findings, classify(tx, input.Mode))
	}
	report.Summary = summarise(report.Findings)
	report.BASSummary = basSummary(report.Findings)
	report.ScenarioChecks = scenarios(report.Findings)
	report.HealthChecks = health(transactions, report.Findings)
	return report, nil
}

func WriteJSON(w io.Writer, report *Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func WriteMarkdown(w io.Writer, report *Report) error {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# TaxMate AU Finance Review\n\n")
	fmt.Fprintf(&b, "- Mode: `%s`\n", report.Mode)
	fmt.Fprintf(&b, "- Source: `%s`\n", report.Source)
	fmt.Fprintf(&b, "- Generated UTC: `%s`\n\n", report.GeneratedAt)

	fmt.Fprintf(&b, "## Summary\n\n")
	fmt.Fprintf(&b, "| Owner | Bucket | Treatment | Rows | Gross | Claim | GST candidate |\n")
	fmt.Fprintf(&b, "|---|---|---:|---:|---:|---:|---:|\n")
	for _, line := range report.Summary {
		fmt.Fprintf(&b, "| %s | %s | %s | %d | %.2f | %.2f | %.2f |\n",
			escapeMD(line.Owner), escapeMD(line.Bucket), escapeMD(line.Treatment), line.Rows,
			line.GrossAmount, line.ClaimAmount, line.GSTCandidate)
	}
	fmt.Fprintf(&b, "\n## BAS/GST\n\n")
	fmt.Fprintf(&b, "- Business expense gross: %.2f\n", report.BASSummary.BusinessExpenseGross)
	fmt.Fprintf(&b, "- GST credit candidate: %.2f\n", report.BASSummary.GSTCreditCandidate)
	fmt.Fprintf(&b, "- Business income gross: %.2f\n", report.BASSummary.BusinessIncomeGross)
	fmt.Fprintf(&b, "- GST collected candidate: %.2f\n", report.BASSummary.GSTCollectedCandidate)
	fmt.Fprintf(&b, "- Nil BAS likely: `%t`\n", report.BASSummary.NilBASLikely)
	fmt.Fprintf(&b, "- Review: %s\n\n", report.BASSummary.ReviewNote)

	fmt.Fprintf(&b, "## Findings\n\n")
	fmt.Fprintf(&b, "| Row | Owner | Description | Bucket | Treatment | Claim %% | Claim | GST | Review | Reasons |\n")
	fmt.Fprintf(&b, "|---:|---|---|---|---|---:|---:|---:|---|---|\n")
	for _, item := range report.Findings {
		fmt.Fprintf(&b, "| %d | %s | %s | %s | %s | %.0f | %.2f | %.2f | %t | %s |\n",
			item.Row, escapeMD(item.Owner), escapeMD(item.Description), escapeMD(item.Bucket),
			escapeMD(item.TaxTreatment), item.ClaimPercent, item.ClaimAmount,
			item.GSTCreditAmount, item.AccountantReview, escapeMD(strings.Join(item.Reasons, "; ")))
	}

	fmt.Fprintf(&b, "\n## Health Checks\n\n")
	for _, check := range report.HealthChecks {
		state := "fail"
		if check.Passed {
			state = "pass"
		}
		fmt.Fprintf(&b, "- `%s` [%s/%s]: %s\n", check.Name, state, check.Severity, check.Detail)
	}
	fmt.Fprintf(&b, "\n## ATO Refresh Queries\n\n")
	for _, q := range report.ATOQueries {
		fmt.Fprintf(&b, "- `%s`\n", q)
	}
	_, err := w.Write(b.Bytes())
	return err
}

func readCSV(r io.Reader) ([]Transaction, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true
	rows, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty csv")
	}
	header := map[string]int{}
	for i, name := range rows[0] {
		header[norm(name)] = i
	}
	var out []Transaction
	for i, row := range rows[1:] {
		if allBlank(row) {
			continue
		}
		tx := Transaction{
			Row:         i + 2,
			Date:        get(row, header, "date", "transactiondate", "posteddate"),
			Description: get(row, header, "description", "merchant", "payee", "memo", "details"),
			Owner:       get(row, header, "owner", "person", "taxpayer"),
			Account:     get(row, header, "account", "accountname"),
			Category:    get(row, header, "category", "class"),
			Purpose:     get(row, header, "purpose", "businesspurpose", "notes", "note"),
			Evidence:    get(row, header, "evidence", "receipt", "invoice", "document"),
			ABN:         get(row, header, "abn", "business", "entity"),
			Source:      get(row, header, "source", "statement"),
			Asset:       get(row, header, "asset", "symbol", "ticker", "security"),
			Raw:         map[string]string{},
		}
		for name, idx := range header {
			if idx < len(row) {
				tx.Raw[name] = strings.TrimSpace(row[idx])
			}
		}
		tx.Amount = parseMoney(firstNonEmpty(
			get(row, header, "amount", "value", "netamount"),
			signedAmount(row, header),
		))
		tx.GST = parseMoney(get(row, header, "gst", "gstamount", "tax", "taxamount"))
		tx.Units = parseMoney(get(row, header, "units", "quantity"))
		tx.Direction = direction(row, header, tx.Amount)
		out = append(out, tx)
	}
	return out, nil
}

func classify(tx Transaction, mode Mode) Finding {
	text := strings.ToLower(strings.Join([]string{tx.Description, tx.Category, tx.Purpose, tx.Account, tx.ABN, tx.Source, tx.Asset}, " "))
	f := Finding{
		Row:              tx.Row,
		Owner:            firstNonEmpty(tx.Owner, "unassigned"),
		Description:      tx.Description,
		Amount:           tx.Amount,
		Direction:        tx.Direction,
		Bucket:           "uncategorised",
		TaxTreatment:     "accountant_review",
		Confidence:       "low",
		Reasons:          []string{"insufficient facts for automatic tax treatment"},
		RecordsNeeded:    []string{"receipt or invoice", "business or work purpose note", "owner"},
		AccountantReview: true,
	}
	if tx.Direction == "income" {
		classifyIncome(&f, tx, text)
		return f
	}
	if containsAny(text, privateHealthTerms...) {
		set(&f, "private_health", "tax_return_info_only", 0, "medium", true, "private health insurance is usually tax-statement information, not a deduction")
		f.RecordsNeeded = []string{"private health insurance tax statement"}
		return f
	}
	if containsAny(text, "super", "superannuation", "personal contribution", "notice of intent") {
		set(&f, "super", "accountant_review", 0, "medium", true, "personal super deduction needs eligibility and notice-of-intent evidence")
		f.RecordsNeeded = []string{"fund acknowledgement", "notice of intent", "contribution statement"}
		return f
	}
	if containsAny(text, investmentTerms...) || tx.Asset != "" || tx.Units != 0 {
		set(&f, "investment", "record_for_income_cgt", 0, "medium", true, "investment records affect distributions, AMIT cost base, DRP, disposals, and CGT")
		f.RecordsNeeded = []string{"annual tax statement", "buy/sell contract notes", "DRP statement", "AMIT cost-base adjustments"}
		return f
	}
	if containsAny(text, softwareTerms...) {
		if isBusiness(tx, text) {
			set(&f, "abn_business_software", "deduction_candidate", 100, "medium", mode == ModeStrict, "software or developer cost appears connected to ABN activity")
			f.RecordsNeeded = []string{"tax invoice", "business purpose", "GST status", "private-use apportionment note"}
			applyGST(&f, tx, true)
			return f
		}
		set(&f, "software_or_subscription", "accountant_review", 0, "low", true, "could be employee, ABN, private, or mixed-use; entity and purpose required")
		return f
	}
	if containsAny(text, "work from home", "wfh", "internet", "phone", "electricity", "stationery", "computer consumables") {
		set(&f, "employee_wfh", "fixed_rate_or_actual_method_review", 0, "medium", true, "WFH claim needs method choice and work-hour records; fixed rate may already cover this cost")
		f.RecordsNeeded = []string{"WFH hours", "method choice", "invoice", "private-use apportionment"}
		return f
	}
	if containsAny(text, "laptop", "monitor", "keyboard", "mouse", "desk", "chair", "equipment", "tool") {
		set(&f, "work_or_business_asset", "depreciation_or_immediate_deduction_review", 0, "medium", true, "asset treatment depends on cost, date, effective life, entity, and private use")
		f.RecordsNeeded = []string{"tax invoice", "purchase date", "private-use percentage", "employee or ABN use"}
		applyGST(&f, tx, isBusiness(tx, text))
		return f
	}
	if containsAny(text, "meal", "coffee", "restaurant", "entertainment", "grocery", "clothes", "fitness", "gym", "medical", "commute", "parking fine") {
		set(&f, "private_or_excluded", "not_claimable", 0, "medium", false, "private or commonly excluded category")
		f.RecordsNeeded = []string{"only keep if accountant asks"}
		return f
	}
	if isBusiness(tx, text) {
		set(&f, "abn_business_expense", "accountant_review", 0, "low", true, "business tag present but expense type is not specific enough")
		f.RecordsNeeded = []string{"tax invoice", "business purpose", "GST status", "private-use apportionment"}
		applyGST(&f, tx, true)
	}
	return f
}

func classifyIncome(f *Finding, tx Transaction, text string) {
	if isBusiness(tx, text) {
		set(f, "abn_business_income", "assessable_income_review", 0, "medium", true, "income appears connected to ABN or side activity")
		if tx.GST > 0 {
			f.GSTCreditAmount = round2(math.Abs(tx.GST))
			f.Reasons = append(f.Reasons, "GST collected candidate present")
		}
		return
	}
	if containsAny(text, employmentIncomeTerms...) {
		set(f, "employment_income", "income_statement_record", 0, "medium", false, "employment income belongs outside ABN expense tracking")
		return
	}
	if containsAny(text, investmentIncomeTerms...) {
		set(f, "investment_income", "tax_statement_record", 0, "medium", true, "investment income needs annual tax statement and franking/AMIT details")
		return
	}
	set(f, "income", "accountant_review", 0, "low", true, "income source needs classification")
}

func set(f *Finding, bucket, treatment string, percent float64, confidence string, review bool, reason string) {
	f.Bucket = bucket
	f.TaxTreatment = treatment
	f.ClaimPercent = percent
	f.ClaimAmount = round2(math.Abs(f.Amount) * percent / 100)
	f.Confidence = confidence
	f.AccountantReview = review
	f.Reasons = []string{reason}
}

func applyGST(f *Finding, tx Transaction, business bool) {
	if business && tx.GST > 0 {
		f.GSTCreditCandidate = true
		f.GSTCreditAmount = round2(math.Abs(tx.GST))
		f.Reasons = append(f.Reasons, "GST credit candidate only if valid tax invoice and creditable purpose")
		if !containsAny(strings.ToLower(tx.Evidence), "tax invoice", "invoice", "receipt") {
			f.RecordsNeeded = appendMissing(f.RecordsNeeded, "valid tax invoice")
		}
	}
}

func summarise(findings []Finding) []SummaryLine {
	type key struct{ owner, bucket, treatment string }
	lines := map[key]*SummaryLine{}
	for _, f := range findings {
		k := key{f.Owner, f.Bucket, f.TaxTreatment}
		if lines[k] == nil {
			lines[k] = &SummaryLine{Owner: f.Owner, Bucket: f.Bucket, Treatment: f.TaxTreatment}
		}
		line := lines[k]
		line.Rows++
		line.GrossAmount = round2(line.GrossAmount + math.Abs(f.Amount))
		line.ClaimAmount = round2(line.ClaimAmount + f.ClaimAmount)
		line.GSTCandidate = round2(line.GSTCandidate + f.GSTCreditAmount)
	}
	out := make([]SummaryLine, 0, len(lines))
	for _, line := range lines {
		out = append(out, *line)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Owner != out[j].Owner {
			return out[i].Owner < out[j].Owner
		}
		if out[i].Bucket != out[j].Bucket {
			return out[i].Bucket < out[j].Bucket
		}
		return out[i].Treatment < out[j].Treatment
	})
	return out
}

func basSummary(findings []Finding) BASSummary {
	var s BASSummary
	for _, f := range findings {
		if strings.HasPrefix(f.Bucket, "abn_business") && f.Direction == "expense" {
			s.BusinessExpenseGross = round2(s.BusinessExpenseGross + math.Abs(f.Amount))
			s.GSTCreditCandidate = round2(s.GSTCreditCandidate + f.GSTCreditAmount)
		}
		if f.Bucket == "abn_business_income" {
			s.BusinessIncomeGross = round2(s.BusinessIncomeGross + math.Abs(f.Amount))
			s.GSTCollectedCandidate = round2(s.GSTCollectedCandidate + f.GSTCreditAmount)
		}
	}
	s.NilBASLikely = s.BusinessIncomeGross == 0 && s.GSTCreditCandidate == 0 && s.GSTCollectedCandidate == 0
	if s.NilBASLikely {
		s.ReviewNote = "No business income or GST-credit candidates detected in supplied rows; confirm complete records before nil BAS."
	} else {
		s.ReviewNote = "Not nil if GST collected or GST credits are being claimed; accountant should review GST labels."
	}
	return s
}

func scenarios(findings []Finding) []Scenario {
	var review, claim float64
	for _, f := range findings {
		if f.AccountantReview {
			review += math.Abs(f.Amount)
		}
		claim += f.ClaimAmount
	}
	return []Scenario{
		{
			Name:       "strict_claims_only",
			BaseAmount: round2(claim),
			WhatIf:     "Only rows with explicit claim percentage are counted.",
			Result:     round2(claim),
			ReviewNote: "Use this for conservative accountant handoff totals.",
		},
		{
			Name:       "review_queue_value",
			BaseAmount: round2(review),
			WhatIf:     "Total value still needing accountant judgement.",
			Result:     round2(review),
			ReviewNote: "Large review value means evidence or entity labels are missing.",
		},
	}
}

func health(txs []Transaction, findings []Finding) []HealthCheck {
	var checks []HealthCheck
	missingOwner, missingEvidence, reviewRows, gstInvoiceRows := []int{}, []int{}, []int{}, []int{}
	seen := map[string]int{}
	dupRows := []int{}
	for _, tx := range txs {
		if strings.TrimSpace(tx.Owner) == "" {
			missingOwner = append(missingOwner, tx.Row)
		}
		if strings.TrimSpace(tx.Evidence) == "" {
			missingEvidence = append(missingEvidence, tx.Row)
		}
		if tx.GST > 0 && !containsAny(strings.ToLower(tx.Evidence), "tax invoice", "invoice", "receipt") {
			gstInvoiceRows = append(gstInvoiceRows, tx.Row)
		}
		key := strings.Join([]string{tx.Date, norm(tx.Description), fmt.Sprintf("%.2f", tx.Amount)}, "|")
		if first, ok := seen[key]; ok {
			dupRows = append(dupRows, first, tx.Row)
		} else {
			seen[key] = tx.Row
		}
	}
	for _, f := range findings {
		if f.AccountantReview {
			reviewRows = append(reviewRows, f.Row)
		}
	}
	checks = append(checks,
		check("owner_present", len(missingOwner) == 0, "medium", "every row should identify taxpayer/spouse/joint/entity", missingOwner),
		check("evidence_present", len(missingEvidence) == 0, "high", "receipt/invoice evidence should be linked before claiming", missingEvidence),
		check("gst_tax_invoice_support", len(gstInvoiceRows) == 0, "high", "GST credit candidates need valid tax invoices", unique(gstInvoiceRows)),
		check("duplicate_scan", len(dupRows) == 0, "medium", "same date/description/amount appears more than once", unique(dupRows)),
		check("accountant_review_queue", len(reviewRows) == 0, "info", "rows flagged for accountant judgement", reviewRows),
	)
	return checks
}

func check(name string, passed bool, severity, detail string, rows []int) HealthCheck {
	return HealthCheck{Name: name, Passed: passed, Severity: severity, Detail: detail, Rows: firstRows(unique(rows), 25)}
}

func direction(row []string, header map[string]int, amount float64) string {
	raw := strings.ToLower(firstNonEmpty(get(row, header, "direction", "type", "transactiontype"), get(row, header, "debitcredit")))
	switch {
	case strings.Contains(raw, "debit"), strings.Contains(raw, "expense"), strings.Contains(raw, "withdrawal"), strings.Contains(raw, "purchase"):
		return "expense"
	case strings.Contains(raw, "credit"), strings.Contains(raw, "income"), strings.Contains(raw, "deposit"):
		return "income"
	case amount < 0:
		return "expense"
	default:
		return "income"
	}
}

func signedAmount(row []string, header map[string]int) string {
	debit := parseMoney(get(row, header, "debit", "withdrawal", "spent"))
	credit := parseMoney(get(row, header, "credit", "deposit", "received"))
	if debit != 0 {
		return fmt.Sprintf("%.2f", -math.Abs(debit))
	}
	if credit != 0 {
		return fmt.Sprintf("%.2f", math.Abs(credit))
	}
	return ""
}

func get(row []string, header map[string]int, names ...string) string {
	for _, name := range names {
		if idx, ok := header[norm(name)]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
	}
	return ""
}

func norm(s string) string {
	return normaliseRE.ReplaceAllString(strings.ToLower(strings.TrimSpace(s)), "")
}

func parseMoney(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	neg := strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")")
	s = strings.Trim(s, "()")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "AUD", "")
	s = strings.TrimSpace(s)
	v, _ := strconv.ParseFloat(s, 64)
	if neg {
		v = -v
	}
	return round2(v)
}

func isBusiness(tx Transaction, text string) bool {
	return strings.TrimSpace(tx.ABN) != "" || containsAny(text, "abn", "sole trader", "business", "side hustle", "app business")
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(s, needle) {
			return true
		}
	}
	return false
}

func allBlank(row []string) bool {
	for _, v := range row {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func appendMissing(values []string, value string) []string {
	for _, item := range values {
		if item == value {
			return values
		}
	}
	return append(values, value)
}

func unique(values []int) []int {
	seen := map[int]bool{}
	var out []int
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Ints(out)
	return out
}

func firstRows(values []int, n int) []int {
	if len(values) <= n {
		return values
	}
	return values[:n]
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func escapeMD(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}
