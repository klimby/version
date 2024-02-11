package cmd

import (
	"github.com/klimby/version/internal/action/current"
	"github.com/klimby/version/internal/di"
	"github.com/spf13/cobra"
)

// currentCmd represents the current command.
var currentCmd = &cobra.Command{
	Use:           "current",
	Short:         "Current version",
	SilenceErrors: true,
	SilenceUsage:  true,
	Example:       `./version current`,
	RunE: func(cmd *cobra.Command, args []string) error {
		action := current.New(func(args *current.Args) {
			args.Repo = di.C.Repo
		})

		command.Set(action)

		return command.Run()
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
