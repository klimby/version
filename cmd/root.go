// Package cmd provides CLI commands.
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/klimby/version/internal/action"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/service/console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var command = &action.Runner{}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "version",
	Short: "CLI tool for versioning",
	Long:  `CLI tool for versioning, generate changelog, bumping version.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if !di.C.IsInit {
			return errors.New("container is not initialized")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	rootCmd.Version = viper.GetString(key.Version)

	console.Notice(fmt.Sprintf("CLI tool for versioning Version v%s.", viper.GetString(key.Version)))
	console.Notice("See https://github.com/klimby/version for more information.\n")

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("silent", "s", false, "silent run")

	if err := viper.BindPFlag(key.Silent, rootCmd.PersistentFlags().Lookup("silent")); err != nil {
		viper.Set(key.Silent, false)
	}

	viper.SetDefault(key.Silent, false)

	rootCmd.PersistentFlags().BoolP("dry", "d", false, "dry run")

	if err := viper.BindPFlag(key.DryRun, rootCmd.PersistentFlags().Lookup("dry")); err != nil {
		viper.Set(key.DryRun, false)
	}

	viper.SetDefault(key.DryRun, false)

	rootCmd.PersistentFlags().StringP("config", "c", config.DefaultConfigFile, "config file path")

	if err := viper.BindPFlag(key.CfgFile, rootCmd.PersistentFlags().Lookup("config")); err != nil {
		viper.Set(key.CfgFile, config.DefaultConfigFile)
	}

	viper.SetDefault(key.CfgFile, config.DefaultConfigFile)

	rootCmd.PersistentFlags().String("dir", "", "working directory, default - current")

	if err := viper.BindPFlag(key.WorkDir, rootCmd.PersistentFlags().Lookup("dir")); err != nil {
		viper.Set(key.WorkDir, "")
	}

	viper.SetDefault(key.WorkDir, "")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	if err := viper.BindPFlag(key.Verbose, rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		viper.Set(key.Verbose, false)
	}

	viper.SetDefault(key.Verbose, false)
}

// initConfig reads in config file and initializes di.
// It run after init(), main() and before Execute().
func initConfig() {
	if err := di.C.Init(); err != nil {
		console.Error(err.Error())

		//nolint:revive
		os.Exit(1)
	}
}
