package finance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeCSVConservativeBuckets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "records.csv")
	data := `date,description,amount,gst,owner,purpose,evidence,abn,type
2026-04-02,Developer tools subscription,-33.00,3.00,Taxpayer A,ABN digital product development,invoice,yes,expense
2026-04-03,Private health premium,-120.00,0,Joint,private health insurance,tax statement,,expense
2026-04-04,ETF distribution,42.00,0,Taxpayer A,Broker ETF annual tax statement,statement,,income
`
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		t.Fatal(err)
	}
	report, err := AnalyzeCSV(Input{Path: path, Mode: ModeStrict})
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Findings) != 3 {
		t.Fatalf("findings=%d", len(report.Findings))
	}
	if report.Findings[0].Bucket != "abn_business_software" {
		t.Fatalf("first bucket=%s", report.Findings[0].Bucket)
	}
	if !report.Findings[0].GSTCreditCandidate {
		t.Fatalf("expected gst credit candidate")
	}
	if report.Findings[1].TaxTreatment != "tax_return_info_only" {
		t.Fatalf("private health treatment=%s", report.Findings[1].TaxTreatment)
	}
	if report.Findings[2].Bucket != "investment_income" {
		t.Fatalf("investment bucket=%s", report.Findings[2].Bucket)
	}
	if report.BASSummary.NilBASLikely {
		t.Fatalf("gst candidate means BAS should not be called nil")
	}
}
