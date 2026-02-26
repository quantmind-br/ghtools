package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete repositories interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDelete()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDelete() error {
	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive deletion")
		return nil
	}

	items := reposToItems(repos)
	selected, err := tui.RunMultiSelect("SELECT REPOS TO DELETE (NO UNDO)", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	tui.PrintWarning(fmt.Sprintf("You are about to delete %d repositories:", len(selected)))
	for _, item := range selected {
		fmt.Println("  - " + tui.StyleError.Render(item.Value))
	}
	fmt.Println()

	confirm, err := tui.RunInput("Type 'DELETE' to confirm, or anything else for Dry Run", "DELETE", "")
	if err != nil {
		return nil
	}

	dryRun := confirm != "DELETE"

	for _, item := range selected {
		if dryRun {
			tui.PrintInfo("[DRY-RUN] Would delete: " + item.Value)
		} else {
			fmt.Println(tui.StyleError.Render("[DELETING]") + " " + item.Value)
			if err := gh.DeleteRepo(item.Value); err != nil {
				tui.PrintError("Failed to delete: " + item.Value)
			} else {
				tui.PrintSuccess("Deleted: " + item.Value)
			}
		}
	}

	if !dryRun {
		_ = cache.Clear()
		tui.PrintInfo("Cache cleared (repo deleted)")
	}

	return nil
}
