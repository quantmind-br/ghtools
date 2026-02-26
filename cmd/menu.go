package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/tui"
)

func runMenu() error {
	tui.ShowHeader("GHTOOLS", fmt.Sprintf("GitHub Repository Manager v%s", version))

	options := []string{
		"List Repositories",
		"Search My Repos",
		"Explore GitHub",
		"Trending Repos",
		"Statistics Dashboard",
		"Clone Repositories",
		"Sync Local Repos",
		"Local Repo Status",
		"Fork Repository",
		"Create Repository",
		"Delete Repositories",
		"Archive/Unarchive",
		"Change Visibility",
		"Browse in Browser",
		"Pull Requests",
		"Config",
		"Refresh Cache",
		"Exit",
	}

	for {
		choice, err := tui.RunChoose("Select an action:", options)
		if err != nil {
			return nil
		}

		switch choice {
		case "List Repositories":
			runList(false, "", "")
		case "Search My Repos":
			runSearch()
		case "Explore GitHub":
			runExplore("", "stars", "", 100)
		case "Trending Repos":
			runTrending("", "daily")
		case "Statistics Dashboard":
			runStats()
		case "Clone Repositories":
			runClone("")
		case "Sync Local Repos":
			runSync(".", false, false, 3)
		case "Local Repo Status":
			runStatus(".", 3)
		case "Fork Repository":
			runFork("", false)
		case "Create Repository":
			runCreate()
		case "Delete Repositories":
			runDelete()
		case "Archive/Unarchive":
			runArchive(false)
		case "Change Visibility":
			runVisibility("")
		case "Browse in Browser":
			runBrowse()
		case "Pull Requests":
			runPRList()
		case "Config":
			runConfig()
		case "Refresh Cache":
			runRefresh()
		case "Exit":
			return nil
		default:
			return nil
		}

		if yesMode {
			return nil
		}

		fmt.Println()
		cont, err := tui.RunConfirm("Continue?", true)
		if err != nil || !cont {
			return nil
		}
	}
}
