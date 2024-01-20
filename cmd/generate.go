package cmd

import (
	"github.com/klimby/version/internal/action"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		gen, err := callGenerator()
		if err != nil {
			return err
		}

		c, err := cmd.Flags().GetBool("config-file")
		if err != nil {
			return err
		}

		if c {
			return gen.Config()
		}

		changelog, err := cmd.Flags().GetBool("changelog")
		if err != nil {
			return err
		}

		if changelog {
			return gen.Changelog()
		}

		if !c && !changelog {
			if err := cmd.Help(); err != nil {
				return err
			}
		}

		return nil
	},
}

// init - init generate command.
func init() {
	generateCmd.Flags().Bool("config-file", false, "generate config file")
	generateCmd.Flags().Bool("changelog", false, "generate changelog file")

	rootCmd.AddCommand(generateCmd)
}

type canGenerate interface {
	Config() error
	Changelog() error
}

var callGenerator = func() (canGenerate, error) {
	return action.NewGenerate(func(arg *action.GenerateArgs) {
		arg.Rw = di.C.FS()
		arg.CfgGen = di.C.Config()
		arg.ClogGen = di.C.Changelog()
	})
}
