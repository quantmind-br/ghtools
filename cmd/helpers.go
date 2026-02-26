package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/tui"
	"github.com/diogo/ghtools/internal/types"
)

func reposToItems(repos []types.Repo) []tui.MultiSelectItem {
	var items []tui.MultiSelectItem
	for _, r := range repos {
		desc := r.Description
		if desc == "" {
			desc = "No description"
		}
		label := fmt.Sprintf("%-30s  %s", tui.Truncate(r.NameWithOwner, 30), tui.Truncate(desc, 40))
		items = append(items, tui.MultiSelectItem{
			Label: label,
			Value: r.NameWithOwner,
		})
	}
	return items
}

func searchResultsToItems(results []types.SearchResult) []tui.MultiSelectItem {
	var items []tui.MultiSelectItem
	for _, r := range results {
		lang := r.Language
		if lang == "" {
			lang = "-"
		}
		label := fmt.Sprintf("%-30s  %d*  %s  %s",
			tui.Truncate(r.FullName, 30),
			r.StargazersCount,
			lang,
			tui.Truncate(r.Description, 30))
		items = append(items, tui.MultiSelectItem{
			Label: label,
			Value: r.FullName,
		})
	}
	return items
}
