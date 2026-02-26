package cmd

import (
	"github.com/diogo/ghtools/internal/config"
	"github.com/diogo/ghtools/internal/tui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Initialize/show config file location",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig() error {
	path, err := config.Init()
	if err != nil {
		return err
	}
	tui.PrintInfo("Config file: " + path)
	return nil
}
