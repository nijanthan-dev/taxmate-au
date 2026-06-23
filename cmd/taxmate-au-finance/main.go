package main

import (
	"flag"
	"fmt"
	"os"

	"taxmate-au-skill/internal/finance"
)

func main() {
	input := flag.String("input", "", "CSV file of expenses, income, investments, GST, super, or private-health records.")
	format := flag.String("format", "json", "Output format: json or markdown.")
	mode := flag.String("mode", string(finance.ModeStrict), "Analysis mode: strict, assisted, or review.")
	output := flag.String("output", "", "Optional output file.")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "use --input file.csv")
		os.Exit(2)
	}
	report, err := finance.AnalyzeCSV(finance.Input{Path: *input, Mode: finance.Mode(*mode)})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := os.Stdout
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
		w = file
	}
	switch *format {
	case "json":
		err = finance.WriteJSON(w, report)
	case "markdown", "md":
		err = finance.WriteMarkdown(w, report)
	default:
		fmt.Fprintf(os.Stderr, "invalid format %q\n", *format)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
