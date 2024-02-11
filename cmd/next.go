package cmd

import (
	"github.com/klimby/version/internal/action/next"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/di"
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
	Example: `./version next --major
./version next --minor
./version next --patch
./version next --ver=1.2.3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		actionType := next.ActionUnknown
		v := version.V("")

		flags := []next.ActionType{next.ActionMajor, next.ActionMinor, next.ActionPatch}

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

		if actionType == next.ActionUnknown {
			ver, err := cmd.Flags().GetString(next.ActionCustom.String())
			if err != nil {
				return err
			}

			if ver != "" {
				v = version.V(ver)
				actionType = next.ActionCustom
			}
		}

		if actionType == next.ActionUnknown {
			return cmd.Help()
		}

		action := next.New(func(args *next.Args) {
			args.Repo = di.C.Repo
			args.ChangelogGen = di.C.ChangelogGenerator
			args.Cfg = di.C.Config
			args.Bump = di.C.Bump
			args.ActionType = actionType
			args.Version = v
		})

		command.Set(action)

		return command.Run()
	},
}

// init - init next command.
func init() {
	initNextCmd()
	rootCmd.AddCommand(nextCmd)
}

// initNextCmd - init next command.
func initNextCmd() {
	nextCmd.Flags().Bool(next.ActionMajor.String(), false, "next major version")
	nextCmd.Flags().Bool(next.ActionMinor.String(), false, "next minor version")
	nextCmd.Flags().Bool(next.ActionPatch.String(), false, "next patch version")

	nextCmd.Flags().String(next.ActionCustom.String(), "", "next build version in format 1.2.3")

	nextCmd.Flags().Bool("prepare", false, "run only bump files and commands before")

	if err := viper.BindPFlag(key.Prepare, nextCmd.Flags().Lookup("prepare")); err != nil {
		viper.Set(key.Prepare, false)
	}

	viper.SetDefault(key.Prepare, false)

	nextCmd.Flags().BoolP("backup", "b", false, "backup changed files")

	if err := viper.BindPFlag(key.Backup, nextCmd.Flags().Lookup("backup")); err != nil {
		viper.Set(key.Backup, false)
	}

	viper.SetDefault(key.Backup, false)

	nextCmd.Flags().BoolP("force", "f", false, "force mode")

	if err := viper.BindPFlag(key.Force, nextCmd.Flags().Lookup("force")); err != nil {
		viper.Set(key.Force, false)
	}

	viper.SetDefault(key.Force, false)

	config.SetForce()
}
