package cmd

import (
	"fmt"

	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/types"
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
		gen := di.C.ActionGenerate
		if gen == nil {
			return fmt.Errorf("action generate is nil: %w", types.ErrNotInitialized)
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

		return cmd.Help()
	},
}

// init - init generate command.
func init() {
	initGenerateCmd()
	rootCmd.AddCommand(generateCmd)
}

// initGenerateCmd - init generate command.
func initGenerateCmd() {
	generateCmd.Flags().Bool("config-file", false, "generate config file")
	generateCmd.Flags().Bool("changelog", false, "generate changelog file")
}
