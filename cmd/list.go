package cmd

import (
	"fmt"
	"strings"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/diogo/ghtools/internal/types"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories with optional filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		refresh, _ := cmd.Flags().GetBool("refresh")
		lang, _ := cmd.Flags().GetString("lang")
		org, _ := cmd.Flags().GetString("org")
		return runList(refresh, lang, org)
	},
}

func init() {
	listCmd.Flags().Bool("refresh", false, "Force refresh cache")
	listCmd.Flags().String("lang", "", "Filter by language")
	listCmd.Flags().String("org", "", "Filter by organization")
	rootCmd.AddCommand(listCmd)
}

func getTableWidths() []int {
	width, _ := tui.GetTerminalSize()
	if width < 100 {
		return []int{25, 30, 8, 8, 10}
	}
	return []int{30, 40, 10, 10, 12}
}

func runList(refresh bool, lang, org string) error {
	if org == "" {
		org = cfg.DefaultOrg
	}

	var repos []types.Repo
	err := tui.RunWithSpinner("Fetching repositories...", func() error {
		var err error
		repos, err = gh.FetchRepos(refresh, cfg.CacheTTL, org)
		return err
	})
	if err != nil {
		return err
	}

	if lang != "" {
		repos = filterByLang(repos, lang)
	}

	if len(repos) == 0 {
		tui.ShowEmptyState("No repositories found matching your criteria")
		return nil
	}

	headers := []string{"NAME", "DESCRIPTION", "VIS", "LANG", "UPDATED"}
	widths := getTableWidths()

	var rows [][]string
	for _, r := range repos {
		updated := r.UpdatedAt.Format("2006-01-02")
		vis := r.Visibility
		l := r.Lang()
		if l == "" {
			l = "-"
		}
		desc := r.Description
		if desc == "" {
			desc = "No description"
		}
		rows = append(rows, []string{r.NameWithOwner, desc, vis, l, updated})
	}

	fmt.Println()
	tui.PrintTable(headers, widths, rows)
	fmt.Println()
	tui.PrintInfo(fmt.Sprintf("Total: %d repositories", len(repos)))
	return nil
}

func filterByLang(repos []types.Repo, lang string) []types.Repo {
	lang = strings.ToLower(lang)
	var filtered []types.Repo
	for _, r := range repos {
		if strings.EqualFold(r.Lang(), lang) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
