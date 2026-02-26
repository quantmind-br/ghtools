package cmd

import (
	"fmt"
	"sort"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show repository statistics dashboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runStats()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func runStats() error {
	repos, err := gh.FetchRepos(true, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	tui.ShowHeader("REPOSITORY STATISTICS", "Your GitHub Overview")

	total := len(repos)
	public, private, archived := 0, 0, 0
	totalStars, totalForks, totalSize := 0, 0, 0
	langCount := make(map[string]int)

	for _, r := range repos {
		switch r.Visibility {
		case "PUBLIC":
			public++
		case "PRIVATE":
			private++
		}
		if r.IsArchived {
			archived++
		}
		totalStars += r.StargazerCount
		totalForks += r.ForkCount
		totalSize += r.DiskUsage

		lang := r.Lang()
		if lang == "" {
			lang = "Unknown"
		}
		langCount[lang]++
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorSecondary).
		Padding(1, 2).MarginLeft(2)

	stats := fmt.Sprintf(
		"Total Repositories:  %d\nPublic:              %d\nPrivate:             %d\nArchived:            %d",
		total, public, private, archived)
	fmt.Println(box.Render(stats))
	fmt.Println()

	metricsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorAccent).
		Padding(1, 2).MarginLeft(2)

	sizeMB := totalSize / 1024
	metrics := fmt.Sprintf(
		"Total Stars:   %d\nTotal Forks:   %d\nTotal Size:    %d MB",
		totalStars, totalForks, sizeMB)
	fmt.Println(metricsBox.Render(metrics))
	fmt.Println()

	// Language breakdown
	fmt.Println(tui.StyleMuted.Render("--- Languages Breakdown ---"))
	fmt.Println()

	type langStat struct {
		Name  string
		Count int
	}
	var langs []langStat
	for name, count := range langCount {
		langs = append(langs, langStat{name, count})
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].Count > langs[j].Count })
	for i, l := range langs {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s: %d\n", tui.StyleSecondary.Render(l.Name), l.Count)
	}
	fmt.Println()

	// Top repos by stars
	fmt.Println(tui.StyleMuted.Render("--- Top Repositories (by Stars) ---"))
	fmt.Println()
	sort.Slice(repos, func(i, j int) bool { return repos[i].StargazerCount > repos[j].StargazerCount })
	for i, r := range repos {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s  %s\n",
			tui.StyleWarning.Render(fmt.Sprintf("%-4d", r.StargazerCount)),
			tui.StyleSecondary.Render(r.NameWithOwner))
	}
	fmt.Println()

	// Recently updated
	fmt.Println(tui.StyleMuted.Render("--- Recently Updated (last 5) ---"))
	fmt.Println()
	sort.Slice(repos, func(i, j int) bool { return repos[i].UpdatedAt.After(repos[j].UpdatedAt) })
	for i, r := range repos {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s  %s\n",
			tui.StyleMuted.Render(r.UpdatedAt.Format("2006-01-02")),
			tui.StyleSecondary.Render(r.NameWithOwner))
	}
	fmt.Println()

	return nil
}
