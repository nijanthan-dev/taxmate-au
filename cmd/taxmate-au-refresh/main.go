package main

import (
	"flag"
	"os"

	"taxmate-au-skill/internal/atodata"
)

func main() {
	query := flag.String("query", "", "Refresh indexed pages matching a topic query.")
	all := flag.Bool("all", false, "Refresh all indexed pages.")
	recrawl := flag.Bool("recrawl", false, "Rebuild the scoped ATO source pack from seed URLs.")
	limit := flag.Int("limit", 12, "Max query matches to refresh.")
	maxPages := flag.Int("max-pages", 250, "Max pages for --recrawl.")
	var urls multiFlag
	flag.Var(&urls, "url", "Refresh explicit indexed ATO URL. Repeatable.")
	flag.Parse()

	root, err := atodata.SkillRoot()
	if err != nil {
		atodata.Errorf("%v", err)
		os.Exit(1)
	}

	if *recrawl {
		idx, err := atodata.Recrawl(root, *maxPages)
		if err != nil {
			atodata.Errorf("%v", err)
			os.Exit(1)
		}
		_ = atodata.WriteJSON(map[string]any{
			"records":  len(idx.Records),
			"failures": len(idx.Failures),
			"index":    atodata.IndexPath(root),
		})
		return
	}

	idx, err := atodata.LoadIndex(root)
	if err != nil {
		atodata.Errorf("%v", err)
		os.Exit(1)
	}

	var selected []*atodata.Record
	var missing []string
	switch {
	case *all:
		selected = idx.Records
	case len(urls) > 0:
		selected, missing = atodata.SelectByURL(idx.Records, urls)
	case *query != "":
		selected = atodata.SelectByQuery(root, idx.Records, *query, *limit)
	default:
		atodata.Errorf("use --query, --url, --all, or --recrawl")
		os.Exit(2)
	}

	results := make([]atodata.RefreshResult, 0, len(selected)+len(missing))
	for _, rawURL := range missing {
		results = append(results, atodata.RefreshResult{URL: rawURL, Error: "not in index"})
	}
	changed := 0
	for _, rec := range selected {
		result := atodata.RefreshRecord(root, rec)
		if result.Changed {
			changed++
		}
		results = append(results, result)
	}
	if err := atodata.SaveIndex(root, idx); err != nil {
		atodata.Errorf("%v", err)
		os.Exit(1)
	}
	_ = atodata.WriteJSON(map[string]any{
		"matched": len(selected),
		"changed": changed,
		"results": results,
	})
}

type multiFlag []string

func (m *multiFlag) String() string {
	return ""
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}
