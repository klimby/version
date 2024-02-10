package cmd

import (
	"fmt"

	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/service/console"
	"github.com/spf13/cobra"
)

// currentCmd represents the current command.
var currentCmd = &cobra.Command{
	Use:           "current",
	Short:         "Current version",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		act := di.C.ActionCurrent

		current, err := act.Current()
		if err != nil {
			return err
		}

		console.Notice(fmt.Sprintf("Current version: %s", current.FormatString()))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
