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

func runList(refresh bool, lang, org string) error {
	if org == "" {
		org = cfg.DefaultOrg
	}

	tui.PrintInfo("Fetching repositories...")
	repos, err := gh.FetchRepos(refresh, cfg.CacheTTL, org)
	if err != nil {
		return err
	}

	if lang != "" {
		repos = filterByLang(repos, lang)
	}

	if len(repos) == 0 {
		tui.PrintWarning("No repositories found")
		return nil
	}

	headers := []string{"NAME", "DESCRIPTION", "VIS", "LANG", "UPDATED"}
	widths := []int{30, 40, 10, 10, 12}

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
