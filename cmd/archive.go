package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive or unarchive repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		unarchive, _ := cmd.Flags().GetBool("unarchive")
		return runArchive(unarchive)
	},
}

func init() {
	archiveCmd.Flags().Bool("unarchive", false, "Unarchive instead of archive")
	rootCmd.AddCommand(archiveCmd)
}

func runArchive(unarchive bool) error {
	mode := "archive"
	if unarchive {
		mode = "unarchive"
	}

	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	// Filter repos based on mode
	var items []tui.MultiSelectItem
	for _, r := range repos {
		if unarchive && !r.IsArchived {
			continue
		}
		if !unarchive && r.IsArchived {
			continue
		}
		desc := r.Description
		if desc == "" {
			desc = "No description"
		}
		label := fmt.Sprintf("%-30s  %s", tui.Truncate(r.NameWithOwner, 30), tui.Truncate(desc, 40))
		items = append(items, tui.MultiSelectItem{Label: label, Value: r.NameWithOwner})
	}

	if len(items) == 0 {
		tui.PrintWarning(fmt.Sprintf("No repositories available to %s", mode))
		return nil
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive selection")
		return nil
	}

	header := fmt.Sprintf("Select repos to %s", mode)
	selected, err := tui.RunMultiSelect(header, items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	tui.PrintWarning(fmt.Sprintf("You are about to %s %d repositories:", mode, len(selected)))
	for _, item := range selected {
		fmt.Println("  - " + item.Value)
	}
	fmt.Println()

	cont, err := tui.RunConfirm("Continue?", false)
	if err != nil || !cont {
		tui.PrintInfo("Cancelled.")
		return nil
	}

	for _, item := range selected {
		tui.PrintInfo(fmt.Sprintf("%sing %s...", mode, item.Value))
		if unarchive {
			if err := gh.UnarchiveRepo(item.Value); err != nil {
				tui.PrintError(fmt.Sprintf("Failed to %s: %s", mode, item.Value))
			} else {
				tui.PrintSuccess(fmt.Sprintf("Unarchived: %s", item.Value))
			}
		} else {
			if err := gh.ArchiveRepo(item.Value); err != nil {
				tui.PrintError(fmt.Sprintf("Failed to %s: %s", mode, item.Value))
			} else {
				tui.PrintSuccess(fmt.Sprintf("Archived: %s", item.Value))
			}
		}
	}

	_ = cache.Clear()
	tui.PrintInfo("Cache cleared (repo states changed)")
	return nil
}
