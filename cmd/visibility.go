package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var visibilityCmd = &cobra.Command{
	Use:   "visibility",
	Short: "Change repository visibility",
	RunE: func(cmd *cobra.Command, args []string) error {
		public, _ := cmd.Flags().GetBool("public")
		private, _ := cmd.Flags().GetBool("private")
		target := ""
		if public {
			target = "public"
		} else if private {
			target = "private"
		}
		return runVisibility(target)
	},
}

func init() {
	visibilityCmd.Flags().Bool("public", false, "Make selected repos public")
	visibilityCmd.Flags().Bool("private", false, "Make selected repos private")
	rootCmd.AddCommand(visibilityCmd)
}

func runVisibility(targetVis string) error {
	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	var items []tui.MultiSelectItem
	headerText := "Select repos to toggle visibility"

	for _, r := range repos {
		if targetVis == "public" && r.Visibility != "PRIVATE" {
			continue
		}
		if targetVis == "private" && r.Visibility != "PUBLIC" {
			continue
		}
		label := fmt.Sprintf("%-30s  %s", tui.Truncate(r.NameWithOwner, 30), r.Visibility)
		items = append(items, tui.MultiSelectItem{Label: label, Value: r.NameWithOwner + "\t" + r.Visibility})
	}

	if targetVis == "public" {
		headerText = "Select PRIVATE repos to make PUBLIC"
	} else if targetVis == "private" {
		headerText = "Select PUBLIC repos to make PRIVATE"
	}

	if len(items) == 0 {
		tui.PrintWarning("No repositories match the criteria")
		return nil
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive selection")
		return nil
	}

	selected, err := tui.RunMultiSelect(headerText, items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	tui.PrintWarning("Visibility changes:")
	type change struct {
		repo   string
		newVis string
	}
	var changes []change

	for _, item := range selected {
		parts := splitFirst(item.Value, "\t")
		repo := parts[0]
		curVis := ""
		if len(parts) > 1 {
			curVis = parts[1]
		}

		newVis := targetVis
		if newVis == "" {
			if curVis == "PUBLIC" {
				newVis = "private"
			} else {
				newVis = "public"
			}
		}

		fmt.Printf("  %s: %s -> %s\n", repo, curVis, newVis)
		changes = append(changes, change{repo, newVis})
	}

	fmt.Println()
	cont, err := tui.RunConfirm("Continue?", false)
	if err != nil || !cont {
		tui.PrintInfo("Cancelled.")
		return nil
	}

	for _, c := range changes {
		tui.PrintInfo(fmt.Sprintf("Changing %s to %s...", c.repo, c.newVis))
		if err := gh.SetVisibility(c.repo, c.newVis); err != nil {
			tui.PrintError("Failed: " + c.repo)
		} else {
			tui.PrintSuccess("Updated: " + c.repo)
		}
	}

	_ = cache.Clear()
	return nil
}

func splitFirst(s, sep string) []string {
	i := 0
	for i < len(s) {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
		i++
	}
	return []string{s}
}
