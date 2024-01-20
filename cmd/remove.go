package cmd

import (
	"github.com/klimby/version/internal/action"
	"github.com/klimby/version/internal/di"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command.
var removeCmd = &cobra.Command{
	Use:           "remove",
	Short:         "Remove files",
	Long:          `Remove backup files, if exists`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		backup, err := cmd.Flags().GetBool("backup")
		if err != nil {
			return err
		}

		if backup {
			a, err := callRemoveBackup()
			if err != nil {
				return err
			}

			a.Backup()

			return nil
		}

		if !backup {
			if err := cmd.Help(); err != nil {
				return err
			}
		}

		return nil
	},
}

// init - init remove command.
func init() {
	removeCmd.Flags().Bool("backup", false, "remove backup files")
	rootCmd.AddCommand(removeCmd)
}

type hasRemoveBackup interface {
	Backup()
}

var callRemoveBackup = func() (hasRemoveBackup, error) {
	return action.NewRemove(func(options *action.ArgsRemove) {
		options.Cfg = di.C.Config()
	})
}
