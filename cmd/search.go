package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Interactive fuzzy search with actions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearch()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch() error {
	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive search")
		return nil
	}

	items := reposToItems(repos)
	selected, err := tui.RunMultiSelect("Search repositories (type to filter)", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	fmt.Println()
	tui.PrintInfo(fmt.Sprintf("Selected %d repositories", len(selected)))
	for _, item := range selected {
		fmt.Println("  " + tui.StyleSecondary.Render(item.Value))
	}
	fmt.Println()

	action, err := tui.RunChoose("Select action:", []string{
		"Clone - Clone to local machine",
		"Browse - Open in browser",
		"Delete - Delete repositories",
		"Cancel",
	})
	if err != nil {
		return nil
	}

	switch {
	case len(action) >= 5 && action[:5] == "Clone":
		for _, item := range selected {
			if err := gh.CloneRepo(item.Value, ""); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Cloned: " + item.Value)
			}
		}
	case len(action) >= 6 && action[:6] == "Browse":
		for _, item := range selected {
			_ = gh.BrowseRepo(item.Value)
		}
		tui.PrintSuccess("Opened in browser")
	case len(action) >= 6 && action[:6] == "Delete":
		tui.PrintWarning("Delete selected repos? This is NOT reversible!")
		confirm, err := tui.RunInput("Type DELETE to confirm", "DELETE", "")
		if err != nil || confirm != "DELETE" {
			tui.PrintInfo("Cancelled")
			return nil
		}
		for _, item := range selected {
			if err := gh.DeleteRepo(item.Value); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Deleted: " + item.Value)
			}
		}
		_ = cache.Clear()
		tui.PrintInfo("Cache cleared (repo deleted)")
	default:
		tui.PrintInfo("Cancelled")
	}

	return nil
}
