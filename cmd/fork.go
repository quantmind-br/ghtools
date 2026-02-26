package cmd

import (
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var forkCmd = &cobra.Command{
	Use:   "fork [query]",
	Short: "Fork external repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		cloneAfter, _ := cmd.Flags().GetBool("clone")
		query := ""
		if len(args) > 0 {
			query = args[0]
		}
		return runFork(query, cloneAfter)
	},
}

func init() {
	forkCmd.Flags().Bool("clone", false, "Clone after forking")
	rootCmd.AddCommand(forkCmd)
}

func runFork(query string, cloneAfter bool) error {
	if query == "" {
		var err error
		query, err = tui.RunInput("Search for repository to fork", "owner/repo or search terms", "")
		if err != nil || query == "" {
			tui.PrintError("Search query required")
			return nil
		}
	}

	tui.PrintInfo("Searching GitHub for '" + query + "'...")

	results, err := gh.SearchRepos(query, "stars", "", 50)
	if err != nil || len(results) == 0 {
		tui.PrintWarning("No repositories found for '" + query + "'")
		return nil
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive selection")
		return nil
	}

	items := searchResultsToItems(results)
	selected, err := tui.RunMultiSelect("Select repository to FORK", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	for _, item := range selected {
		tui.PrintInfo("Forking " + item.Value + "...")
		if err := gh.ForkRepo(item.Value, cloneAfter); err != nil {
			tui.PrintError("Failed to fork: " + item.Value)
		} else {
			tui.PrintSuccess("Forked: " + item.Value)
			if cloneAfter {
				tui.PrintSuccess("Cloned locally")
			}
		}
	}
	return nil
}
