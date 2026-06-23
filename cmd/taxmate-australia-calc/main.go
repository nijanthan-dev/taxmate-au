package main

import (
	"flag"
	"fmt"
	"os"

	"taxmate-au-skill/internal/taxcalc"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	tool := os.Args[1]
	if err := taxcalc.ValidateTool(tool); err != nil {
		fmt.Fprintln(os.Stderr, err)
		usage()
		os.Exit(2)
	}
	var result taxcalc.Result
	switch tool {
	case "bas":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		gstCollected := fs.Float64("gst-collected", 0, "GST on sales.")
		gstCredits := fs.Float64("gst-credits", 0, "GST credits on purchases.")
		payg := fs.Float64("payg-withheld", 0, "PAYG withheld to report.")
		fuel := fs.Float64("fuel-tax-credit", 0, "Fuel tax credits.")
		adjustments := fs.Float64("adjustments", 0, "GST adjustments.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.BAS(*gstCollected, *gstCredits, *payg, *fuel, *adjustments)
	case "super":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		ote := fs.Float64("ote", 0, "Ordinary time earnings.")
		rate := fs.Float64("rate", 0, "SG rate percent. Defaults to 12 for 2025-26.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.SuperGuarantee(*ote, *rate)
	case "fbt":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		taxableValue := fs.Float64("taxable-value", 0, "FBT taxable value already worked out.")
		benefitType := fs.String("type", "type2", "type1 or type2 gross-up.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.FBT(*taxableValue, *benefitType)
	case "cgt":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		proceeds := fs.Float64("proceeds", 0, "Capital proceeds.")
		costBase := fs.Float64("cost-base", 0, "Cost base.")
		losses := fs.Float64("capital-losses", 0, "Capital losses applied before discount.")
		acquired := fs.String("acquired", "", "Acquisition date YYYY-MM-DD.")
		disposed := fs.String("disposed", "", "Disposal date YYYY-MM-DD.")
		discount := fs.Bool("discount", false, "Apply 50% discount only if eligible.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.CGT(*proceeds, *costBase, *losses, *acquired, *disposed, *discount)
	case "payg":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		gross := fs.Float64("gross-pay", 0, "Regular gross pay per period.")
		periods := fs.Int("periods", 52, "Pay periods per year.")
		tft := fs.Bool("tax-free-threshold", true, "Tax-free threshold claimed.")
		medicare := fs.Bool("medicare", false, "Add a simple 2% Medicare levy estimate.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.PAYGEstimate(*gross, *periods, *tft, *medicare)
	case "stamp-duty":
		fs := flag.NewFlagSet(tool, flag.ExitOnError)
		state := fs.String("state", "", "State or territory abbreviation, for example VIC.")
		value := fs.Float64("value", 0, "Dutiable value.")
		_ = fs.Parse(os.Args[2:])
		result = taxcalc.StampDutyRouter(*state, *value)
	}
	if err := taxcalc.WriteJSON(os.Stdout, result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: taxmate-australia-calc <bas|super|fbt|cgt|payg|stamp-duty> [flags]")
}
