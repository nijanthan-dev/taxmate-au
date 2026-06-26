package taxcalc

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

const (
	SGRate2025_26       = 12.0
	FBTRate2025_26      = 47.0
	FBTType1GrossUp     = 2.0802
	FBTType2GrossUp     = 1.8868
	MedicareLevyDefault = 2.0
)

type Result struct {
	Tool        string            `json:"tool"`
	IncomeYear  string            `json:"income_year"`
	Inputs      map[string]any    `json:"inputs"`
	Outputs     map[string]any    `json:"outputs"`
	Assumptions []string          `json:"assumptions"`
	ReviewFlags []string          `json:"review_flags"`
	Sources     []string          `json:"sources"`
	Official    map[string]string `json:"official_state_sources,omitempty"`
}

func WriteJSON(w io.Writer, result Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func BAS(salesGST, purchaseGST, paygWithheld, fuelTaxCredit, adjustments float64) Result {
	netGST := round2(salesGST - purchaseGST + adjustments)
	netPayable := round2(netGST + paygWithheld - fuelTaxCredit)
	return Result{
		Tool:       "bas",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"gst_collected":    salesGST,
			"gst_credits":      purchaseGST,
			"payg_withheld":    paygWithheld,
			"fuel_tax_credit":  fuelTaxCredit,
			"gst_adjustments":  adjustments,
			"amounts_are_gst":  true,
			"cash_or_accruals": "user supplied",
		},
		Outputs: map[string]any{
			"net_gst_payable":     netGST,
			"estimated_bas_total": netPayable,
			"nil_bas":             netPayable == 0,
		},
		Assumptions: []string{"Inputs are already separated into GST collected, GST credits, PAYG withheld, fuel tax credits, and adjustments."},
		ReviewFlags: []string{"Confirm BAS reporting cycle, accounting basis, labels, and whether GST credits have valid tax invoices."},
		Sources: []string{
			"https://www.ato.gov.au/businesses-and-organisations/preparing-lodging-and-paying/business-activity-statements-bas",
			"https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/claiming-gst-credits",
		},
	}
}

func SuperGuarantee(ote, rate float64) Result {
	if rate == 0 {
		rate = SGRate2025_26
	}
	return Result{
		Tool:       "super",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"ordinary_time_earnings": ote,
			"sg_rate_percent":        rate,
		},
		Outputs: map[string]any{
			"minimum_sg": round2(ote * rate / 100),
		},
		Assumptions: []string{"Uses ordinary time earnings supplied by the user; rate defaults to 12% for payments made from 1 July 2025."},
		ReviewFlags: []string{"Check award/agreement higher rates, OTE classification, quarterly due date, and late-payment SGC exposure."},
		Sources: []string{
			"https://www.ato.gov.au/businesses-and-organisations/super-for-employers/paying-super-contributions/how-much-super-to-pay",
		},
	}
}

func FBT(taxableValue float64, benefitType string) Result {
	grossUp := FBTType2GrossUp
	if strings.EqualFold(benefitType, "type1") || strings.EqualFold(benefitType, "type-1") || benefitType == "1" {
		grossUp = FBTType1GrossUp
		benefitType = "type1"
	} else {
		benefitType = "type2"
	}
	grossedUp := round2(taxableValue * grossUp)
	return Result{
		Tool:       "fbt",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"taxable_value":  taxableValue,
			"benefit_type":   benefitType,
			"gross_up_rate":  grossUp,
			"fbt_rate":       FBTRate2025_26,
			"fbt_year_basis": "year ending 31 March 2026",
		},
		Outputs: map[string]any{
			"grossed_up_taxable_value": grossedUp,
			"estimated_fbt":            round2(grossedUp * FBTRate2025_26 / 100),
		},
		Assumptions: []string{"Taxable value has already been worked out under the relevant FBT benefit rules."},
		ReviewFlags: []string{"Does not determine car statutory formula, operating cost method, exemptions, employee contributions, or reportable fringe benefit treatment."},
		Sources: []string{
			"https://www.ato.gov.au/tax-rates-and-codes/fringe-benefits-tax-rates-and-thresholds",
			"https://www.ato.gov.au/businesses-and-organisations/hiring-and-paying-your-workers/fringe-benefits-tax",
		},
	}
}

func CGT(proceeds, costBase, capitalLosses float64, acquired, disposed string, discount bool) Result {
	rawGain := round2(proceeds - costBase)
	netBeforeDiscount := round2(rawGain - capitalLosses)
	heldMonths := monthsHeld(acquired, disposed)
	discountAllowed := discount && heldMonths >= 12 && netBeforeDiscount > 0
	net := netBeforeDiscount
	discountAmount := 0.0
	if discountAllowed {
		discountAmount = round2(netBeforeDiscount * 0.5)
		net = round2(netBeforeDiscount - discountAmount)
	}
	return Result{
		Tool:       "cgt",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"capital_proceeds": proceeds,
			"cost_base":        costBase,
			"capital_losses":   capitalLosses,
			"acquired":         acquired,
			"disposed":         disposed,
			"discount_claimed": discount,
		},
		Outputs: map[string]any{
			"gross_capital_gain":   rawGain,
			"net_before_discount":  netBeforeDiscount,
			"held_months":          heldMonths,
			"discount_allowed":     discountAllowed,
			"discount_amount":      discountAmount,
			"net_capital_gain_est": net,
		},
		Assumptions: []string{"Cost base, proceeds, and capital losses are user-supplied and already include relevant incidental amounts."},
		ReviewFlags: []string{"Check asset type, main residence exemption, small business concessions, rollovers, foreign-resident rules, AMIT/ETF cost-base adjustments, and carried-forward losses."},
		Sources: []string{
			"https://www.ato.gov.au/individuals-and-families/investments-and-assets/capital-gains-tax/calculating-your-cgt",
			"https://www.ato.gov.au/individuals-and-families/investments-and-assets/capital-gains-tax/cgt-discount",
		},
	}
}

func PAYGEstimate(grossPay float64, periodsPerYear int, taxFreeThreshold bool, medicare bool) Result {
	if periodsPerYear <= 0 {
		periodsPerYear = 52
	}
	annual := grossPay * float64(periodsPerYear)
	tax := residentTax2025_26(annual)
	if medicare {
		tax += annual * MedicareLevyDefault / 100
	}
	if !taxFreeThreshold && annual <= 18200 {
		tax = annual * 0.16
	}
	withhold := round2(tax / float64(periodsPerYear))
	return Result{
		Tool:       "payg-estimate",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"gross_pay":           grossPay,
			"periods_per_year":    periodsPerYear,
			"annualised_pay":      round2(annual),
			"tax_free_threshold":  taxFreeThreshold,
			"medicare_levy_added": medicare,
		},
		Outputs: map[string]any{
			"estimated_withholding_per_period": withhold,
			"estimated_annual_tax":             round2(tax),
		},
		Assumptions: []string{"Annualises regular pay and applies 2025-26 resident tax rates; this is not a substitute for ATO PAYG withholding tax tables."},
		ReviewFlags: []string{"Use official ATO withholding tables for payroll, HELP/STSL, Medicare variations, no-TFN cases, bonuses, commissions, termination payments, allowances, foreign residents, and rounding."},
		Sources: []string{
			"https://www.ato.gov.au/tax-rates-and-codes/tax-rates-australian-residents",
			"https://www.ato.gov.au/tax-rates-and-codes/tax-tables-overview",
			"https://www.ato.gov.au/businesses-and-organisations/hiring-and-paying-your-workers/payg-withholding",
		},
	}
}

func StampDutyRouter(state string, value float64) Result {
	state = strings.ToUpper(strings.TrimSpace(state))
	official := StateRevenueSources()
	source := official[state]
	if source == "" {
		source = "unknown state; use the relevant state or territory revenue office"
	}
	return Result{
		Tool:       "stamp-duty-source-router",
		IncomeYear: "2025-26",
		Inputs: map[string]any{
			"state":          state,
			"dutiable_value": value,
		},
		Outputs: map[string]any{
			"calculation": "not_calculated",
			"reason":      "Stamp duty is state or territory based and must be checked live against the relevant revenue-office calculator/rates.",
			"source":      source,
		},
		Assumptions: []string{"TaxMate does not embed state stamp-duty rate tables because concessions and surcharges change frequently."},
		ReviewFlags: []string{"Check property type, first-home concessions, principal-place-of-residence rules, foreign purchaser surcharge, off-the-plan rules, vacant residential land tax, and transfer date."},
		Sources:     []string{source},
		Official:    official,
	}
}

func StateRevenueSources() map[string]string {
	return map[string]string{
		"ACT": "https://www.revenue.act.gov.au/",
		"NSW": "https://www.revenue.nsw.gov.au/taxes-duties-levies-royalties/transfer-duty",
		"NT":  "https://nt.gov.au/property/buying-and-selling-property/stamp-duty",
		"QLD": "https://qro.qld.gov.au/duties/transfer-duty/",
		"SA":  "https://www.revenuesa.sa.gov.au/stampduty",
		"TAS": "https://www.sro.tas.gov.au/property-transfer-duties",
		"VIC": "https://www.sro.vic.gov.au/land-transfer-duty",
		"WA":  "https://www.wa.gov.au/organisation/department-of-finance/transfer-duty",
	}
}

func residentTax2025_26(income float64) float64 {
	switch {
	case income <= 18200:
		return 0
	case income <= 45000:
		return (income - 18200) * 0.16
	case income <= 135000:
		return 4288 + (income-45000)*0.30
	case income <= 190000:
		return 31288 + (income-135000)*0.37
	default:
		return 51638 + (income-190000)*0.45
	}
}

func monthsHeld(acquired, disposed string) int {
	a, errA := time.Parse("2006-01-02", acquired)
	d, errD := time.Parse("2006-01-02", disposed)
	if errA != nil || errD != nil || d.Before(a) {
		return 0
	}
	return int(d.Sub(a).Hours() / 24 / 30.4375)
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func ValidateTool(name string) error {
	switch name {
	case "bas", "super", "fbt", "cgt", "payg", "stamp-duty":
		return nil
	default:
		return fmt.Errorf("unknown calculator %q", name)
	}
}
