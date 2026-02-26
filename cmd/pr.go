package cmd

import (
	"fmt"
	"time"

	"github.com/diogo/ghtools/internal/gh"
	gitpkg "github.com/diogo/ghtools/internal/git"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Pull request operations",
}

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List PRs for a repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPRList()
	},
}

var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create PR from current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPRCreate()
	},
}

func init() {
	prCmd.AddCommand(prListCmd)
	prCmd.AddCommand(prCreateCmd)
	rootCmd.AddCommand(prCmd)
}

func runPRList() error {
	repos, err := gh.FetchRepos(false, cfg.CacheTTL, cfg.DefaultOrg)
	if err != nil {
		return err
	}

	if yesMode {
		tui.PrintInfo("Yes mode: skipping interactive selection")
		return nil
	}

	items := make([]string, len(repos))
	for i, r := range repos {
		items[i] = r.NameWithOwner
	}

	selected, err := tui.RunChoose("Select repository to list PRs", items)
	if err != nil {
		return nil
	}

	tui.PrintInfo("Fetching PRs for " + selected + "...")
	fmt.Println()

	prs, err := gh.PRList(selected, 20)
	if err != nil {
		return fmt.Errorf("failed to fetch PRs: %w", err)
	}

	if len(prs) == 0 {
		tui.PrintInfo("No open PRs found")
		return nil
	}

	for _, pr := range prs {
		stateStyle := tui.StyleSuccess
		if pr.State == "CLOSED" {
			stateStyle = tui.StyleError
		} else if pr.State == "MERGED" {
			stateStyle = tui.StyleAccent
		}

		ago := timeAgo(pr.CreatedAt)
		fmt.Printf("  #%-4d %s  %s  (by %s, %s)\n",
			pr.Number,
			stateStyle.Render(fmt.Sprintf("[%-6s]", pr.State)),
			pr.Title,
			pr.Author.Login,
			ago)
	}
	fmt.Println()
	return nil
}

func runPRCreate() error {
	if !gitpkg.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	branch := gitpkg.CurrentBranch(".")
	if branch == "detached" {
		return fmt.Errorf("you are in 'detached HEAD' state. Please checkout a branch first")
	}
	if branch == "main" || branch == "master" {
		tui.PrintWarning("You're on " + branch + " branch. Create a feature branch first.")
		return nil
	}

	tui.PrintInfo("Creating PR from branch: " + branch)

	title, err := tui.RunInput("PR Title", branch, branch)
	if err != nil || title == "" {
		title = branch
	}

	draft := false
	if !yesMode {
		d, err := tui.RunConfirm("Create as draft?", false)
		if err == nil {
			draft = d
		}
	}

	// Push branch if needed
	if !gitpkg.HasRemoteBranch(".", branch) {
		tui.PrintInfo("Pushing branch to origin...")
		if err := gitpkg.Push(".", branch); err != nil {
			return fmt.Errorf("failed to push branch: %w", err)
		}
	}

	if err := gh.PRCreate(title, "", draft); err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	tui.PrintSuccess("PR created successfully")
	return nil
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
