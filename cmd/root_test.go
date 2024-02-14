package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_rootFlags(t *testing.T) {
	helperMock := __newHelpMock()

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helperMock.Help()
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	rootCmd.SetArgs([]string{
		"--silent",
		"--dry",
		"--verbose",
		"--config=c.yml",
		"--dir=foo",
	})

	assert.NoError(t, rootCmd.Execute(), "rootCmd()")

	assert.Equal(t, true, viper.GetBool(key.Silent), "silent flag not set")
	assert.Equal(t, true, viper.GetBool(key.DryRun), "dry flag not set")
	assert.Equal(t, true, viper.GetBool(key.Verbose), "verbose flag not set")
	assert.Equal(t, "c.yml", viper.GetString(key.CfgFile), "config-file flag not set")
	assert.Equal(t, "foo", viper.GetString(key.WorkDir), "dir flag not set")

}
