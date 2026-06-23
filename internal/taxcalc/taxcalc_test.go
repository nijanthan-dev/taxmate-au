package taxcalc

import "testing"

func TestCoreCalculators(t *testing.T) {
	if got := SuperGuarantee(1000, 0).Outputs["minimum_sg"]; got != 120.0 {
		t.Fatalf("sg=%v", got)
	}
	if got := FBT(1000, "type1").Outputs["estimated_fbt"]; got != 977.69 {
		t.Fatalf("fbt=%v", got)
	}
	if got := CGT(25000, 10000, 0, "2024-01-01", "2026-02-01", true).Outputs["net_capital_gain_est"]; got != 7500.0 {
		t.Fatalf("cgt=%v", got)
	}
	if got := PAYGEstimate(1500, 52, true, false).Outputs["estimated_withholding_per_period"]; got != 272.85 {
		t.Fatalf("payg=%v", got)
	}
}

func TestStampDutyRouterDoesNotCalculate(t *testing.T) {
	result := StampDutyRouter("VIC", 800000)
	if result.Outputs["calculation"] != "not_calculated" {
		t.Fatalf("unexpected stamp duty calculation: %v", result.Outputs)
	}
	if result.Official["VIC"] == "" {
		t.Fatalf("missing VIC source")
	}
}
