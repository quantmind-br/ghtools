package cmd

import (
	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Clear the repository cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRefresh()
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

func runRefresh() error {
	_ = cache.Clear()
	tui.PrintSuccess("Cache cleared.")
	return nil
}
