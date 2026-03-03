package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/git"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of local repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		maxDepth, _ := cmd.Flags().GetInt("max-depth")
		return runStatus(path, maxDepth)
	},
}

func init() {
	statusCmd.Flags().String("path", ".", "Base path to scan")
	statusCmd.Flags().Int("max-depth", 3, "Maximum directory depth to scan")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(basePath string, maxDepth int) error {
	tui.PrintInfo(fmt.Sprintf("Scanning repositories in %s...", basePath))

	dirs, err := git.FindRepos(basePath, maxDepth)
	if err != nil {
		return err
	}

	if len(dirs) == 0 {
		tui.ShowEmptyState("No git repositories found")
		return nil
	}

	tui.PrintInfo(fmt.Sprintf("Found %d repositories", len(dirs)))
	fmt.Println()

	headers := []string{"REPOSITORY", "BRANCH", "STATUS", "AHEAD", "BEHIND"}
	widths := []int{35, 12, 12, 8, 8}

	var rows [][]string
	for _, dir := range dirs {
		status := git.GetRepoStatus(dir)

		statusText := "clean"
		if status.Dirty && status.Untracked {
			statusText = "dirty+untrk"
		} else if status.Dirty {
			statusText = "dirty"
		} else if status.Untracked {
			statusText = "untracked"
		}

		ahead := "?"
		if status.Ahead >= 0 {
			ahead = fmt.Sprintf("%d", status.Ahead)
		}
		behind := "?"
		if status.Behind >= 0 {
			behind = fmt.Sprintf("%d", status.Behind)
		}

		rows = append(rows, []string{status.Name, status.Branch, statusText, ahead, behind})
	}

	tui.PrintTable(headers, widths, rows)
	fmt.Println()
	fmt.Println(tui.StyleInfo.Render("Legend:") + " clean | dirty/untracked | ahead (unpushed) | behind (needs pull)")
	return nil
}
