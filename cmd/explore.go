package cmd

import (
	"fmt"
	"strings"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var exploreCmd = &cobra.Command{
	Use:   "explore [query]",
	Short: "Search external GitHub repositories",
	RunE: func(cmd *cobra.Command, args []string) error {
		sort, _ := cmd.Flags().GetString("sort")
		lang, _ := cmd.Flags().GetString("lang")
		limit, _ := cmd.Flags().GetInt("limit")
		query := strings.Join(args, " ")
		return runExplore(query, sort, lang, limit)
	},
}

func init() {
	exploreCmd.Flags().String("sort", "stars", "Sort by: stars, forks, updated")
	exploreCmd.Flags().String("lang", "", "Filter by language")
	exploreCmd.Flags().Int("limit", 100, "Maximum results")
	rootCmd.AddCommand(exploreCmd)
}

func runExplore(query, sort, lang string, limit int) error {
	if query == "" {
		tui.ShowHeader("EXPLORE GITHUB", "Search repositories worldwide")

		var err error
		query, err = tui.RunInput("Search query (e.g., 'machine learning')", "search terms", "")
		if err != nil || query == "" {
			tui.PrintError("Search query required")
			return nil
		}

		sortChoice, err := tui.RunChoose("Sort by:", []string{"stars", "forks", "updated"})
		if err == nil {
			sort = sortChoice
		}

		langInput, _ := tui.RunInput("Filter by language (leave empty for all)", "language", "")
		if langInput != "" {
			lang = langInput
		}
	}

	tui.PrintInfo(fmt.Sprintf("Searching GitHub for '%s' (sorted by %s)...", query, sort))

	results, err := gh.SearchRepos(query, sort, lang, limit)
	if err != nil || len(results) == 0 {
		tui.PrintWarning("No repositories found for '" + query + "'")
		return nil
	}

	tui.PrintInfo(fmt.Sprintf("Found %d repositories", len(results)))

	if yesMode {
		// In yes mode, just print results
		for _, r := range results {
			fmt.Printf("  %s  %d*  %s\n", r.FullName, r.StargazersCount, r.Language)
		}
		return nil
	}

	items := searchResultsToItems(results)
	selected, err := tui.RunMultiSelect("Explore results (Tab to select, Enter to confirm)", items)
	if err != nil || len(selected) == 0 {
		return nil
	}

	fmt.Println()
	tui.PrintInfo(fmt.Sprintf("Selected %d repositories", len(selected)))
	fmt.Println()

	action, err := tui.RunChoose("Select action:", []string{
		"Clone - Clone to local machine",
		"Fork - Fork to your account",
		"Fork + Clone - Fork and clone",
		"Browse - Open in browser",
		"Star - Star the repository",
		"Info - Show detailed info",
		"Cancel",
	})
	if err != nil {
		return nil
	}

	return executeExploreAction(action, selected)
}

func executeExploreAction(action string, selected []tui.MultiSelectItem) error {
	switch {
	case strings.HasPrefix(action, "Clone"):
		for _, item := range selected {
			tui.PrintInfo("Cloning " + item.Value + "...")
			if err := gh.CloneRepo(item.Value, ""); err != nil {
				tui.PrintError("Failed to clone: " + item.Value)
			} else {
				tui.PrintSuccess("Cloned: " + item.Value)
			}
		}
	case action == "Fork - Fork to your account":
		for _, item := range selected {
			tui.PrintInfo("Forking " + item.Value + "...")
			if err := gh.ForkRepo(item.Value, false); err != nil {
				tui.PrintError("Failed to fork: " + item.Value)
			} else {
				tui.PrintSuccess("Forked: " + item.Value)
			}
		}
	case strings.HasPrefix(action, "Fork + Clone"):
		for _, item := range selected {
			tui.PrintInfo("Forking and cloning " + item.Value + "...")
			if err := gh.ForkRepo(item.Value, true); err != nil {
				tui.PrintError("Failed: " + item.Value)
			} else {
				tui.PrintSuccess("Forked and cloned: " + item.Value)
			}
		}
	case strings.HasPrefix(action, "Browse"):
		for _, item := range selected {
			tui.PrintInfo("Opening " + item.Value + " in browser...")
			_ = gh.BrowseRepo(item.Value)
		}
		tui.PrintSuccess("Opened in browser")
	case strings.HasPrefix(action, "Star"):
		for _, item := range selected {
			tui.PrintInfo("Starring " + item.Value + "...")
			if err := gh.StarRepo(item.Value); err != nil {
				tui.PrintError("Failed to star: " + item.Value)
			} else {
				tui.PrintSuccess("Starred: " + item.Value)
			}
		}
	case strings.HasPrefix(action, "Info"):
		for _, item := range selected {
			tui.ShowHeader(item.Value, "Repository Details")
			info, err := gh.ViewRepo(item.Value)
			if err != nil {
				tui.PrintError("Failed to fetch info")
			} else {
				fmt.Println(info)
			}
			fmt.Println()
		}
	default:
		tui.PrintInfo("Cancelled")
	}
	return nil
}
