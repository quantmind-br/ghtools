package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/tui"
	"github.com/diogo/ghtools/internal/types"
)

func reposToItems(repos []types.Repo) []tui.MultiSelectItem {
	termWidth, _ := tui.GetTerminalSize()
	nameWidth := termWidth / 3
	if nameWidth > 30 {
		nameWidth = 30
	}
	if nameWidth < 10 {
		nameWidth = 10
	}
	descWidth := termWidth - nameWidth - 8
	if descWidth < 10 {
		descWidth = 10
	}

	var items []tui.MultiSelectItem
	for _, r := range repos {
		desc := r.Description
		if desc == "" {
			desc = "No description"
		}
		label := fmt.Sprintf("%-*s  %s", nameWidth, tui.Truncate(r.NameWithOwner, nameWidth), tui.Truncate(desc, descWidth))
		items = append(items, tui.MultiSelectItem{
			Label: label,
			Value: r.NameWithOwner,
		})
	}
	return items
}

func searchResultsToItems(results []types.SearchResult) []tui.MultiSelectItem {
	termWidth, _ := tui.GetTerminalSize()
	nameWidth := termWidth / 3
	if nameWidth > 30 {
		nameWidth = 30
	}
	if nameWidth < 10 {
		nameWidth = 10
	}
	descWidth := termWidth - nameWidth - 12
	if descWidth < 10 {
		descWidth = 10
	}

	var items []tui.MultiSelectItem
	for _, r := range results {
		lang := r.Language
		if lang == "" {
			lang = "-"
		}
		label := fmt.Sprintf("%-*s  %d*  %s  %s",
			nameWidth,
			tui.Truncate(r.FullName, nameWidth),
			r.StargazersCount,
			lang,
			tui.Truncate(r.Description, descWidth))
		items = append(items, tui.MultiSelectItem{
			Label: label,
			Value: r.FullName,
		})
	}
	return items
}
