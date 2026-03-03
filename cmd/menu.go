package cmd

import (
	"fmt"

	"github.com/diogo/ghtools/internal/tui"
)

func runMenu() error {
	tui.ShowHeader("GHTOOLS", fmt.Sprintf("GitHub Repository Manager v%s", version))

	options := []string{
		"[L] List Repositories",
		"[S] Search My Repos",
		"[E] Explore GitHub",
		"[T] Trending Repos",
		"[D] Statistics Dashboard",
		"[C] Clone Repositories",
		"[Y] Sync Local Repos",
		"[O] Local Repo Status",
		"[F] Fork Repository",
		"[R] Create Repository",
		"[X] Delete Repositories",
		"[A] Archive/Unarchive",
		"[V] Change Visibility",
		"[B] Browse in Browser",
		"[P] Pull Requests",
		"[G] Config",
		"[M] Refresh Cache",
		"[Q] Exit",
	}

	for {
		choice, err := tui.RunChoose("Select an action:", options)
		if err != nil {
			return nil
		}

		switch choice {
		case "[L] List Repositories":
			runList(false, "", "")
		case "[S] Search My Repos":
			runSearch()
		case "[E] Explore GitHub":
			runExplore("", "stars", "", 100)
		case "[T] Trending Repos":
			runTrending("", "daily")
		case "[D] Statistics Dashboard":
			runStats()
		case "[C] Clone Repositories":
			runClone("")
		case "[Y] Sync Local Repos":
			runSync(".", false, false, 3)
		case "[O] Local Repo Status":
			runStatus(".", 3)
		case "[F] Fork Repository":
			runFork("", false)
		case "[R] Create Repository":
			runCreate()
		case "[X] Delete Repositories":
			runDelete()
		case "[A] Archive/Unarchive":
			runArchive(false)
		case "[V] Change Visibility":
			runVisibility("")
		case "[B] Browse in Browser":
			runBrowse()
		case "[P] Pull Requests":
			runPRList()
		case "[G] Config":
			runConfig()
		case "[M] Refresh Cache":
			runRefresh()
		case "[Q] Exit":
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
