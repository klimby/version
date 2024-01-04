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

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:           "generate",
	Short:         "Generate files",
	Long:          `Generate and rewrite config and changelog files.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := cmd.Flags().GetBool("config")
		if err != nil {
			return err
		}

		if config {
			return generateConfig()
		}

		changelog, err := cmd.Flags().GetBool("changelog")
		if err != nil {
			return err
		}

		if changelog {
			return generateChangelog()
		}

		if !config && !changelog {
			if err := cmd.Help(); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().Bool("config", false, "generate config file")
	generateCmd.Flags().Bool("changelog", false, "generate changelog file")

	rootCmd.AddCommand(generateCmd)
}

type generateConfigArgs struct {
	f file.ReadWriter
	c generateConfigArgsConfig
}

type generateConfigArgsConfig interface {
	Generate(file.Writer) error
}

func generateConfig(opts ...func(options *generateConfigArgs)) error {
	a := &generateConfigArgs{
		f: di.C.FS(),
		c: di.C.Config(),
	}

	for _, opt := range opts {
		opt(a)
	}

	console.Notice("Generate config file\n")

	p := config.File(viper.GetString(config.CfgFile))

	if err := backup.Create(a.f, p.Path()); err != nil {
		return err
	}

	if err := a.c.Generate(a.f); err != nil {
		return err
	}

	console.Success("Config file generated.")

	return nil
}

type generateChangelogArgs struct {
	chGen generateChangelogArgsChGen
}

type generateChangelogArgsChGen interface {
	Generate() error
}

func generateChangelog(opts ...func(options *generateChangelogArgs)) error {
	a := &generateChangelogArgs{
		chGen: di.C.Changelog(),
	}

	for _, opt := range opts {
		opt(a)
	}

	console.Notice("Generate changelog\n")

	if err := a.chGen.Generate(); err != nil {
		return err
	}

	console.Success("Changelog generated.")

	return nil
}
