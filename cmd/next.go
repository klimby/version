package cmd

import (
	"fmt"

	"github.com/klimby/version/internal/action"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/internal/types"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nextCmd represents the next command.
var nextCmd = &cobra.Command{
	Use:           "next",
	Short:         "Generate next version",
	Long:          `Generate next version, bump files, generate changelog.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		next := di.C.ActionNext
		if next == nil {
			return fmt.Errorf("action next is nil: %w", types.ErrNotInitialized)
		}

		na := action.PrepareNextArgs{
			NextType: git.NextNone,
		}

		major, err := cmd.Flags().GetBool("major")
		if err != nil {
			return err
		}

		if major {
			na.NextType = git.NextMajor
		}

		if na.NextType == git.NextNone {
			minor, err := cmd.Flags().GetBool("minor")
			if err != nil {
				return err
			}

			if minor {
				na.NextType = git.NextMinor
			}
		}

		if na.NextType == git.NextNone {
			patch, err := cmd.Flags().GetBool("patch")
			if err != nil {
				return err
			}

			if patch {
				na.NextType = git.NextPatch
			}
		}

		if na.NextType == git.NextNone {
			ver, err := cmd.Flags().GetString("ver")
			if err != nil {
				return err
			}

			if ver != "" {
				nextV := version.V(ver)
				if nextV.Invalid() {
					return fmt.Errorf("invalid version: %s", nextV)
				}

				na.NextType = git.NextCustom
				na.Custom = nextV
			}
		}

		if na.NextType == git.NextNone {
			if err := cmd.Help(); err != nil {
				return err
			}
		}

		prepare, err := cmd.Flags().GetBool("prepare")
		if err != nil {
			return err
		}

		nv, err := next.Prepare(func(args *action.PrepareNextArgs) {
			args.NextType = na.NextType
			args.Custom = na.Custom
		})
		if err != nil {
			return err
		}

		if prepare {
			console.Success(fmt.Sprintf("Prepare complete, next version is %s", nv.FormatString()))
			return nil
		}

		return next.Apply(nv)
	},
}

// init - init next command.
func init() {
	initNextCmd()
	rootCmd.AddCommand(nextCmd)
}

// initNextCmd - init next command.
func initNextCmd() {
	nextCmd.Flags().Bool("major", false, "next major version")
	nextCmd.Flags().Bool("minor", false, "next minor version")
	nextCmd.Flags().Bool("patch", false, "next patch version")
	nextCmd.Flags().String("ver", "", "next build version in format 1.2.3")
	nextCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "ver")

	nextCmd.Flags().String("prepare", "", "run only bump files and commands before")

	rootCmd.PersistentFlags().BoolP("backup", "b", false, "backup changed files")

	if err := viper.BindPFlag(config.Backup, rootCmd.PersistentFlags().Lookup("backup")); err != nil {
		viper.Set(config.Backup, false)
	}

	viper.SetDefault(config.Backup, false)

	rootCmd.PersistentFlags().BoolP("force", "f", false, "force mode")

	if err := viper.BindPFlag(config.Force, rootCmd.PersistentFlags().Lookup("force")); err != nil {
		viper.Set(config.Force, false)
	}

	viper.SetDefault(config.Force, false)

	config.SetForce()
}
