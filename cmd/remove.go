package cmd

import (
	"fmt"

	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/types"
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
		remove := di.C.ActionRemove
		if remove == nil {
			return fmt.Errorf("action remove is nil: %w", types.ErrNotInitialized)
		}

		backup, err := cmd.Flags().GetBool("backup")
		if err != nil {
			return err
		}

		if backup {
			remove.Backup()

			return nil
		}

		return cmd.Help()
	},
}

// init - init remove command.
func init() {
	initRemoveCmd()
	rootCmd.AddCommand(removeCmd)
}

// initRemoveCmd - init remove command.
func initRemoveCmd() {
	removeCmd.Flags().Bool("backup", false, "remove backup files")
}
