package cmd

import (
	"fmt"
	"sort"

	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/diogo/ghtools/internal/types"
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
	var repos []types.Repo
	err := tui.RunWithSpinner("Fetching repositories...", func() error {
		var err error
		repos, err = gh.FetchRepos(true, cfg.CacheTTL, cfg.DefaultOrg)
		return err
	})
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

	termWidth, _ := tui.GetTerminalSize()
	boxWidth := termWidth - 4
	if boxWidth < 20 {
		boxWidth = 20
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorSecondary).
		Padding(1, 2).MarginLeft(2).
		MaxWidth(boxWidth)

	stats := fmt.Sprintf(
		"Total Repositories:  %d\nPublic:              %d\nPrivate:             %d\nArchived:            %d",
		total, public, private, archived)
	fmt.Println(box.Render(stats))
	fmt.Println()

	metricsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(tui.ColorAccent).
		Padding(1, 2).MarginLeft(2).
		MaxWidth(boxWidth)

	sizeMB := totalSize / 1024
	metrics := fmt.Sprintf(
		"Total Stars:   %d\nTotal Forks:   %d\nTotal Size:    %d MB",
		totalStars, totalForks, sizeMB)
	fmt.Println(metricsBox.Render(metrics))
	fmt.Println()

	// Language breakdown
	type langStat struct {
		Name  string
		Count int
	}
	var langs []langStat
	for name, count := range langCount {
		langs = append(langs, langStat{name, count})
	}
	sort.Slice(langs, func(i, j int) bool { return langs[i].Count > langs[j].Count })
	var langContent string
	for i, l := range langs {
		if i >= 10 {
			break
		}
		langContent += fmt.Sprintf("  %s: %d\n", tui.StyleSecondary.Render(l.Name), l.Count)
	}
	tui.ShowSection("Languages", langContent)

	// Top repos by stars
	sort.Slice(repos, func(i, j int) bool { return repos[i].StargazerCount > repos[j].StargazerCount })
	var topReposContent string
	for i, r := range repos {
		if i >= 5 {
			break
		}
		topReposContent += fmt.Sprintf("  %s  %s\n",
			tui.StyleWarning.Render(fmt.Sprintf("%-4d", r.StargazerCount)),
			tui.StyleSecondary.Render(r.NameWithOwner))
	}
	tui.ShowSection("Top Repositories", topReposContent)

	// Recently updated
	sort.Slice(repos, func(i, j int) bool { return repos[i].UpdatedAt.After(repos[j].UpdatedAt) })
	var recentContent string
	for i, r := range repos {
		if i >= 5 {
			break
		}
		recentContent += fmt.Sprintf("  %s  %s\n",
			tui.StyleMuted.Render(r.UpdatedAt.Format("2006-01-02")),
			tui.StyleSecondary.Render(r.NameWithOwner))
	}
	tui.ShowSection("Recently Updated", recentContent)

	return nil
}
