package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func Test_nextCmd(t *testing.T) {
	helper := &__helpCaller{}

	nextCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helper.called = true
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	tests := []struct {
		name           string
		arg            string
		wantCallAction bool
		wantCallHelp   bool
		wantErr        bool
	}{
		{
			name:           "call major",
			arg:            "--major",
			wantCallAction: true,
			wantErr:        false,
		},
		{
			name:           "call minor",
			arg:            "--minor",
			wantCallAction: true,
			wantErr:        false,
		},
		{
			name:           "call patch",
			arg:            "--patch",
			wantCallAction: true,
			wantErr:        false,
		},
		{
			name:           "call custom",
			arg:            "--ver=1.2.3",
			wantCallAction: true,
			wantErr:        false,
		},
		{
			name:         "call custom invalid",
			arg:          "--ver",
			wantCallHelp: false,
			wantErr:      true,
		},
		{
			name:         "call without args",
			arg:          "",
			wantCallHelp: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helper.called = false
				nextCmd.ResetFlags()
				initNextCmd()
			})

			runner := &__runner{}
			command.SetForce(runner)

			rootCmd.SetArgs([]string{nextCmd.Use, tt.arg})

			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("nextCmd() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantCallAction != runner.called {
				t.Errorf("nextCmd() error = %v, wantErr %v", "not called", tt.wantCallAction)
			}

			if tt.wantCallHelp != helper.called {
				t.Errorf("nextCmd() error = %v, wantErr %v", "not called", tt.wantCallHelp)
			}

		})
	}
}

func Test_nextCmdFlags(t *testing.T) {
	helper := &__helpCaller{}

	nextCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helper.called = true
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	tests := []struct {
		name        string
		arg         string
		wantPrepare bool
		wantBackup  bool
		wantForce   bool
		wantErr     bool
	}{
		{
			name:        "call prepare",
			arg:         "--prepare",
			wantPrepare: true,
		},
		{
			name:       "call backup",
			arg:        "--backup",
			wantBackup: true,
		},
		{
			name:      "call force",
			arg:       "--force",
			wantForce: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helper.called = false
				nextCmd.ResetFlags()
				initNextCmd()
			})

			runner := &__runner{}
			command.SetForce(runner)

			rootCmd.SetArgs([]string{nextCmd.Use, "--patch", tt.arg})

			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("nextCmd() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.wantPrepare, viper.GetBool(key.Prepare), "prepare flag not set")
			assert.Equal(t, tt.wantBackup, viper.GetBool(key.Backup), "backup flag not set")
			assert.Equal(t, tt.wantForce, viper.GetBool(key.Force), "force flag not set")

		})
	}
}
