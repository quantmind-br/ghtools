package cmd

import (
	"fmt"
	"strings"

	gitpkg "github.com/diogo/ghtools/internal/git"
	"github.com/diogo/ghtools/internal/runner"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local repositories with remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		all, _ := cmd.Flags().GetBool("all")
		maxDepth, _ := cmd.Flags().GetInt("max-depth")
		return runSync(path, dryRun, all, maxDepth)
	},
}

func init() {
	syncCmd.Flags().String("path", ".", "Base path to scan")
	syncCmd.Flags().Bool("dry-run", false, "Show what would happen without making changes")
	syncCmd.Flags().Bool("all", false, "Sync all repos without selection")
	syncCmd.Flags().Int("max-depth", 3, "Maximum directory depth to scan")
	rootCmd.AddCommand(syncCmd)
}

func runSync(basePath string, dryRun, syncAll bool, maxDepth int) error {
	tui.PrintInfo(fmt.Sprintf("Scanning for git repositories in %s...", basePath))

	dirs, err := gitpkg.FindRepos(basePath, maxDepth)
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		tui.ShowEmptyState("No git repositories found")
		return nil
	}

	tui.PrintInfo(fmt.Sprintf("Found %d repositories", len(dirs)))

	var selectedDirs []string
	if syncAll || yesMode {
		selectedDirs = dirs
	} else {
		items := make([]tui.MultiSelectItem, len(dirs))
		for i, d := range dirs {
			items[i] = tui.MultiSelectItem{Label: d, Value: d}
		}
		selected, err := tui.RunMultiSelect("Select repos to SYNC (--all to skip)", items)
		if err != nil || len(selected) == 0 {
			return nil
		}
		for _, s := range selected {
			selectedDirs = append(selectedDirs, s.Value)
		}
	}

	if dryRun {
		tui.PrintWarning("DRY-RUN mode: no changes will be made")
	}

	tui.PrintInfo(fmt.Sprintf("Syncing %d repositories in parallel...", len(selectedDirs)))

	var tasks []runner.Task
	for _, dir := range selectedDirs {
		d := dir
		tasks = append(tasks, runner.Task{
			Name: d,
			Fn: func() (string, error) {
				if gitpkg.IsDirty(d) {
					return "Skipped (dirty state)", nil
				}

				if dryRun {
					_ = gitpkg.Fetch(d)
					_, behind := gitpkg.AheadBehind(d)
					if behind > 0 {
						return fmt.Sprintf("[DRY-RUN] Would pull %d commits", behind), nil
					}
					return "[DRY-RUN] Already up to date", nil
				}

				output, err := gitpkg.Pull(d)
				if err != nil {
					return "", fmt.Errorf("conflict or diverged")
				}
				if strings.Contains(output, "Already up to date") {
					return "No change", nil
				}
				return "Synced", nil
			},
		})
	}

	r := runner.New(cfg.MaxJobs)
	results := r.Run(tasks, func(done, total int) {
		fmt.Printf("\r  Progress: %d/%d", done, total)
	})
	fmt.Println()

	for _, result := range results {
		if result.Success {
			tui.PrintSuccess(result.Name + ": " + result.Message)
		} else {
			tui.PrintError(result.Name + ": " + result.Message)
		}
	}

	tui.PrintSuccess("Sync completed.")
	return nil
}
