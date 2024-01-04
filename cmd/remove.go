package cmd

import (
	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			return removeBackup()
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

// removeArgs - arguments for removeBackup.
type removeArgs struct {
	f   file.Remover
	cfg removeArgsConfig
}

// removeArgsConfig - config interface for removeArgs.
type removeArgsConfig interface {
	BumpFiles() []config.BumpFile
}

// removeBackup removes backup files.
func removeBackup(opts ...func(options *removeArgs)) error {
	a := &removeArgs{
		f:   di.C.FS(),
		cfg: di.C.Config(),
	}

	for _, opt := range opts {
		opt(a)
	}

	console.Notice("Remove backup files\n")

	p := config.File(viper.GetString(config.CfgFile))

	backup.Remove(a.f, p.Path())

	for _, bmp := range a.cfg.BumpFiles() {
		backup.Remove(a.f, bmp.File.Path())
	}

	console.Success("Backup files removed.")

	return nil
}
