package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "version",
	Short: "CLI tool for versioning",
	Long:  `CLI tool for versioning, generate changelog, bumping version.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !di.C.IsInit {
			return errors.New("container is not initialized")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	rootCmd.Version = viper.GetString(config.Version)

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("silent", "s", false, "silent run")

	if err := viper.BindPFlag(config.Silent, rootCmd.PersistentFlags().Lookup("silent")); err != nil {
		viper.Set(config.Silent, false)
	}

	viper.SetDefault(config.Silent, false)

	rootCmd.PersistentFlags().BoolP("dry", "d", false, "dry run")

	if err := viper.BindPFlag(config.DryRun, rootCmd.PersistentFlags().Lookup("dry")); err != nil {
		viper.Set(config.DryRun, false)
	}

	viper.SetDefault(config.DryRun, false)

	rootCmd.PersistentFlags().StringP("config", "c", config.DefaultConfigFile, "config file path")

	if err := viper.BindPFlag(config.CfgFile, rootCmd.PersistentFlags().Lookup("config-file")); err != nil {
		viper.Set(config.CfgFile, config.DefaultConfigFile)
	}

	viper.SetDefault(config.CfgFile, config.DefaultConfigFile)

	rootCmd.PersistentFlags().String("dir", "", "working directory, default - current")

	if err := viper.BindPFlag(config.WorkDir, rootCmd.PersistentFlags().Lookup("dir")); err != nil {
		viper.Set(config.WorkDir, "")
	}

	viper.SetDefault(config.WorkDir, "")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	if err := viper.BindPFlag(config.Verbose, rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		viper.Set(config.Verbose, false)
	}

	viper.SetDefault(config.Verbose, false)
}

// initConfig reads in config file and initializes di.
// It run after init(), main() and before Execute().
func initConfig() {
	if err := di.C.Init(); err != nil {
		console.Error(err.Error())

		//nolint:revive
		os.Exit(1)
	}

	console.Notice(fmt.Sprintf("CLI tool for versioning Version v%s.", viper.GetString(config.Version)))
	console.Notice("See https://github.com/klimby/version for more information.\n")
}
