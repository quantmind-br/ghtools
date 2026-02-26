package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/template"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCreate()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate() error {
	tui.ShowHeader("CREATE REPOSITORY", "Set up a new GitHub repository")

	name, err := tui.RunInput("Repository name", "my-repo", "")
	if err != nil || name == "" {
		tui.PrintError("Name required")
		return nil
	}

	desc, _ := tui.RunInput("Description (optional)", "", "")

	vis, err := tui.RunChoose("Visibility:", []string{"public", "private"})
	if err != nil {
		vis = "private"
	}

	tpl, err := tui.RunChoose("Template:", []string{"none", "python", "node", "go"})
	if err != nil {
		tpl = "none"
	}

	fmt.Println()
	tui.PrintInfo(fmt.Sprintf("Creating repository: %s (%s, template: %s)", name, vis, tpl))

	_, err = gh.CreateRepo(name, desc, vis, true)
	if err != nil {
		tui.PrintError("Failed to create repository.")
		return err
	}

	tui.PrintSuccess("Repository created: " + name)

	if tpl != "none" {
		dir := name
		if _, statErr := os.Stat(dir); statErr == nil {
			tui.PrintInfo("Applying template: " + tpl)
			if err := template.Apply(dir, tpl); err != nil {
				tui.PrintError("Failed to apply template: " + err.Error())
			} else {
				// Commit template
				gitAdd := exec.Command("git", "add", ".")
				gitAdd.Dir = dir
				_ = gitAdd.Run()

				gitCommit := exec.Command("git", "commit", "-m", fmt.Sprintf("Initial commit (Template: %s)", tpl))
				gitCommit.Dir = dir
				_ = gitCommit.Run()

				push, err := tui.RunConfirm("Push initial commit to origin?", true)
				if err == nil && push {
					gitPush := exec.Command("git", "push", "origin", "HEAD")
					gitPush.Dir = filepath.Dir(dir)
					gitPush.Dir = dir
					if err := gitPush.Run(); err == nil {
						tui.PrintSuccess("Template applied and pushed")
					}
				} else {
					tui.PrintSuccess("Template applied (not pushed). Push manually with: cd " + name + " && git push")
				}
			}
		}
	}

	return nil
}
