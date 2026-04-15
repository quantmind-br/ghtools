package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/diogo/ghtools/internal/types"
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
	var repos []types.Repo
	err := tui.RunWithSpinner("Fetching repositories...", func() error {
		var err error
		repos, err = gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
		return err
	})
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
		cloned := make(map[string]bool)
		for _, item := range selected {
			if err := gh.CloneRepo(item.Value, ""); err != nil {
				tui.PrintError("Failed: " + item.Value)
				cloned[item.Value] = false
			} else {
				tui.PrintSuccess("Cloned: " + item.Value)
				cloned[item.Value] = true
			}
		}
		// Retry prompt
		failed := 0
		for _, v := range cloned {
			if !v {
				failed++
			}
		}
		if failed > 0 {
			retry, _ := tui.RunConfirm(fmt.Sprintf("%d operations failed. Retry?", failed), false)
			if retry {
				for _, item := range selected {
					if !cloned[item.Value] {
						if err := gh.CloneRepo(item.Value, ""); err != nil {
							tui.PrintError("Failed again: " + item.Value)
						} else {
							tui.PrintSuccess("Cloned on retry: " + item.Value)
						}
					}
				}
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
			tui.PrintWarning("Input did not match 'DELETE'. Cancelled.")
			return nil
		}
		deleted := make(map[string]bool)
		for _, item := range selected {
			if err := gh.DeleteRepo(item.Value); err != nil {
				tui.PrintError("Failed: " + item.Value)
				deleted[item.Value] = false
			} else {
				tui.PrintSuccess("Deleted: " + item.Value)
				deleted[item.Value] = true
			}
		}
		// Retry prompt
		failed := 0
		for _, v := range deleted {
			if !v {
				failed++
			}
		}
		if failed > 0 {
			retry, _ := tui.RunConfirm(fmt.Sprintf("%d operations failed. Retry?", failed), false)
			if retry {
				for _, item := range selected {
					if !deleted[item.Value] {
						if err := gh.DeleteRepo(item.Value); err != nil {
							tui.PrintError("Failed again: " + item.Value)
						} else {
							tui.PrintSuccess("Deleted on retry: " + item.Value)
						}
					}
				}
			}
		}
		_ = cache.Clear()
		tui.PrintInfo("Cache cleared (repo deleted)")
	default:
		tui.PrintInfo("Cancelled")
	}

	return nil
}
