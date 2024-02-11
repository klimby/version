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
	helper := &__helpCaller{}

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helper.called = true
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

	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("rootCmd() error = %v, wantErr %v", err, false)
	}

	assert.Equal(t, true, viper.GetBool(key.Silent), "silent flag not set")
	assert.Equal(t, true, viper.GetBool(key.DryRun), "dry flag not set")
	assert.Equal(t, true, viper.GetBool(key.Verbose), "verbose flag not set")
	assert.Equal(t, "c.yml", viper.GetString(key.CfgFile), "config-file flag not set")
	assert.Equal(t, "foo", viper.GetString(key.WorkDir), "dir flag not set")

}
