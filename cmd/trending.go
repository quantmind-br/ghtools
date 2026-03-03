package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var trendingCmd = &cobra.Command{
	Use:   "trending",
	Short: "Show trending repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		since, _ := cmd.Flags().GetString("since")
		return runTrending(lang, since)
	},
}

func init() {
	trendingCmd.Flags().String("lang", "", "Filter by language")
	trendingCmd.Flags().String("since", "daily", "Time range: daily, weekly, monthly")
	rootCmd.AddCommand(trendingCmd)
}

func runTrending(lang, since string) error {
	tui.ShowHeader("TRENDING REPOSITORIES", "Hot projects this week")

	daysAgo := 7
	switch since {
	case "monthly":
		daysAgo = 30
	case "daily":
		daysAgo = 1
	}

	dateStr := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")
	query := fmt.Sprintf("stars:>100 pushed:>%s", dateStr)
	if lang != "" {
		query += " language:" + lang
	}

	tui.PrintInfo(fmt.Sprintf("Fetching trending repositories%s...", langSuffix(lang)))

	results, err := gh.SearchRepos(query, "stars", "", 30)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		tui.ShowEmptyState("No trending repositories found")
		return nil
	}

	if yesMode {
		for _, r := range results {
			fmt.Printf("  %s  %d*  %s\n", r.FullName, r.StargazersCount, r.Language)
		}
		return nil
	}

	items := searchResultsToItems(results)
	selected, err := tui.RunMultiSelect(fmt.Sprintf("Trending%s (Tab to select)", langSuffix(lang)), items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	action, err := tui.RunChoose("Select action:", []string{
		"Clone - Clone to local machine",
		"Fork - Fork to your account",
		"Browse - Open in browser",
		"Star - Star the repository",
		"Cancel",
	})
	if err != nil {
		return nil
	}

	switch {
	case strings.HasPrefix(action, "Clone"):
		for _, item := range selected {
			if err := gh.CloneRepo(item.Value, ""); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Cloned: " + item.Value)
			}
		}
	case strings.HasPrefix(action, "Fork"):
		for _, item := range selected {
			if err := gh.ForkRepo(item.Value, false); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Forked: " + item.Value)
			}
		}
	case strings.HasPrefix(action, "Browse"):
		for _, item := range selected {
			_ = gh.BrowseRepo(item.Value)
		}
		tui.PrintSuccess("Opened in browser")
	case strings.HasPrefix(action, "Star"):
		for _, item := range selected {
			if err := gh.StarRepo(item.Value); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Starred: " + item.Value)
			}
		}
	default:
		tui.PrintInfo("Cancelled")
	}

	return nil
}

func langSuffix(lang string) string {
	if lang != "" {
		return " for " + lang
	}
	return ""
}
