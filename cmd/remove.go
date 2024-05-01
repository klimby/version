package cmd

import (
	"github.com/klimby/version/internal/action/remove"
	"github.com/klimby/version/internal/di"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command.
var removeCmd = &cobra.Command{
	Use:           "remove",
	Short:         "RemoveAll files",
	Long:          `RemoveAll backup files, if exists`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example:       `version remove --backup`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		actionType := remove.ActionUnknown

		flags := []remove.ActionType{remove.ActionBackup}

		for _, f := range flags {
			b, err := cmd.Flags().GetBool(f.String())
			if err != nil {
				return err
			}

			if b {
				actionType = f
				break
			}
		}

		if actionType == remove.ActionUnknown {
			return cmd.Help()
		}

		action := remove.New(func(args *remove.Args) {
			args.ActionType = actionType
			args.Cfg = di.C.Config
		})

		command.Set(action)

		return command.Run()
	},
}

// init - init remove command.
func init() {
	initRemoveCmd()
	rootCmd.AddCommand(removeCmd)
}

// initRemoveCmd - init remove command.
func initRemoveCmd() {
	removeCmd.Flags().Bool(remove.ActionBackup.String(), false, "remove backup files")
}
