package cmd

import (
	"fmt"
	"os"

	"github.com/diogo/ghtools/internal/config"
	"github.com/diogo/ghtools/internal/gh"
	"github.com/diogo/ghtools/internal/git"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

const version = "4.0.0"

var (
	cfg     config.Config
	verbose bool
	quiet   bool
	yesMode bool
)

var rootCmd = &cobra.Command{
	Use:   "ghtools",
	Short: "GitHub repository management tool",
	Long:  "A cross-platform CLI for managing GitHub repositories with an interactive TUI.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		tui.Quiet = quiet

		if msg := config.CheckMigration(); msg != "" {
			tui.PrintWarning(msg)
		}

		cfg = config.Load()

		if err := gh.CheckInstalled(); err != nil {
			return err
		}
		if err := git.CheckInstalled(); err != nil {
			return err
		}
		if err := gh.CheckAuth(); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMenu()
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-error output")
	rootCmd.PersistentFlags().BoolVarP(&yesMode, "yes", "y", false, "Non-interactive mode (auto-confirm with defaults)")
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("ghtools {{.Version}}\n")
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, tui.StyleError.Render("ERROR")+" "+err.Error())
		return err
	}
	return nil
}
