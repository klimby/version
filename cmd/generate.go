package cmd

import (
	"github.com/klimby/version/internal/action/generate"
	"github.com/klimby/version/internal/di"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate files",
	Long:          `Generate and rewrite config and changelog files.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Example: `./version generate --config-file
./version generate --changelog`,
	RunE: func(cmd *cobra.Command, args []string) error {
		actionType := generate.FileUnknown

		flags := []generate.ActionType{generate.FileConfig, generate.FileChangelog}

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

		if actionType == generate.FileUnknown {
			return cmd.Help()
		}

		action := generate.New(func(args *generate.Args) {
			args.ActionType = actionType
			args.CfgGenerator = di.C.Config
			args.ChangelogGen = di.C.ChangelogGenerator
		})

		command.Set(action)

		return command.Run()
	},
}

// init - init generate command.
func init() {
	initGenerateCmd()
	rootCmd.AddCommand(generateCmd)
}

// initGenerateCmd - init generate command.
func initGenerateCmd() {
	generateCmd.Flags().Bool(generate.FileConfig.String(), false, "generate config file")
	generateCmd.Flags().Bool(generate.FileChangelog.String(), false, "generate changelog file")
}
