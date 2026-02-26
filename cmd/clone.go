package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/runner"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone repositories interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		return runClone(path)
	},
}

func init() {
	cloneCmd.Flags().String("path", "", "Clone destination path")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(clonePath string) error {
	if clonePath == "" {
		clonePath = cfg.DefaultClonePath
	}
	if clonePath == "" {
		clonePath, _ = os.Getwd()
	}

	if info, err := os.Stat(clonePath); err != nil || !info.IsDir() {
		return fmt.Errorf("clone path does not exist: %s", clonePath)
	}

	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	tui.PrintInfo("Clone destination: " + clonePath)
	items := reposToItems(repos)

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive selection")
		return nil
	}

	selected, err := tui.RunMultiSelect("Select repositories to CLONE (Tab to toggle)", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	tui.PrintInfo(fmt.Sprintf("Cloning %d repositories in parallel (%d threads)...", len(selected), cfg.MaxJobs))

	var tasks []runner.Task
	for _, item := range selected {
		repo := item.Value
		repoName := filepath.Base(repo)
		targetDir := filepath.Join(clonePath, repoName)

		tasks = append(tasks, runner.Task{
			Name: repo,
			Fn: func() (string, error) {
				if _, err := os.Stat(targetDir); err == nil {
					return "Skipped (directory exists)", nil
				}
				if err := gh.CloneRepo(repo, targetDir); err != nil {
					return "", err
				}
				return "Cloned", nil
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

	tui.PrintSuccess("All clone operations completed.")
	return nil
}
