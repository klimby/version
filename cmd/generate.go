package cmd

import (
	"fmt"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/file"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command.
var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate files",
	Long:          `Generate and rewrite config and changelog files.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cmd.Flags().GetBool("config-file")
		if err != nil {
			return err
		}

		if c {
			return generateConfig()
		}

		changelog, err := cmd.Flags().GetBool("changelog")
		if err != nil {
			return err
		}

		if changelog {
			return generateChangelog()
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

// generateConfigArgs - arguments for generateConfig.
type generateConfigArgs struct {
	f file.ReadWriter
	c generateConfigArgsConfig
}

// generateConfigArgsConfig - config interface for generateConfig.
type generateConfigArgsConfig interface {
	Generate(file.Writer) error
}

// generateConfig generates config file.
func generateConfig(opts ...func(options *generateConfigArgs)) error {
	a := &generateConfigArgs{
		f: di.C.FS(),
		c: di.C.Config(),
	}

	for _, opt := range opts {
		opt(a)
	}

	console.Notice("Generate config file...")

	p := config.File(viper.GetString(config.CfgFile))

	if err := backup.Create(a.f, p.Path()); err != nil {
		return err
	}

	if err := a.c.Generate(a.f); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Config %s created.", p.String()))

	return nil
}

// generateChangelogArgs - arguments for generateChangelog.
type generateChangelogArgs struct {
	chGen generateChangelogArgsChGen
}

// generateChangelogArgsChGen - changelog interface for generateChangelog.
type generateChangelogArgsChGen interface {
	Generate() error
}

// generateChangelog generates changelog file.
func generateChangelog(opts ...func(options *generateChangelogArgs)) error {
	a := &generateChangelogArgs{
		chGen: di.C.Changelog(),
	}

	for _, opt := range opts {
		opt(a)
	}

	if !viper.GetBool(config.GenerateChangelog) {
		console.Info("Changelog generation disabled.")

		return nil
	}

	console.Notice("Generate changelog...")

	if err := a.chGen.Generate(); err != nil {
		return err
	}

	return nil
}
