package cmd

import (
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Open repositories in browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runBrowse()
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}

func runBrowse() error {
	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	items := reposToItems(repos)
	selected, err := tui.RunMultiSelect("Select repos to open in browser", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	for _, item := range selected {
		tui.PrintInfo("Opening " + item.Value + " in browser...")
		_ = gh.BrowseRepo(item.Value)
	}
	tui.PrintSuccess("Opened selected repositories in browser")
	return nil
}
