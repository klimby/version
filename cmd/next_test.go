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
	helperMock := __newHelpMock()

	nextCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helperMock.Help()
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	type wantCall struct {
		action bool
		help   bool
	}

	tests := []struct {
		name      string
		arg       string
		wantCall  wantCall
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "call major",
			arg:  "--major",
			wantCall: wantCall{
				action: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "call minor",
			arg:  "--minor",
			wantCall: wantCall{
				action: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "call patch",
			arg:  "--patch",
			wantCall: wantCall{
				action: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "call custom",
			arg:  "--ver=1.2.3",
			wantCall: wantCall{
				action: true,
			},
			assertion: assert.NoError,
		},
		{
			name:      "call custom invalid",
			arg:       "--ver",
			assertion: assert.Error,
		},
		{
			name: "call without args",
			arg:  "",
			wantCall: wantCall{
				help: true,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helperMock = __newHelpMock()
				nextCmd.ResetFlags()
				initNextCmd()
			})

			runnerMock := __newRunnerMock(nil)
			command.SetForce(runnerMock)

			rootCmd.SetArgs([]string{nextCmd.Use, tt.arg})

			tt.assertion(t, rootCmd.Execute(), "rootCmd.Execute() error = %v, wantErr %v", tt.assertion, false)

			if tt.wantCall.action {
				runnerMock.AssertCalled(t, "Run")
			} else {
				runnerMock.AssertNotCalled(t, "Run")
			}

			if tt.wantCall.help {
				helperMock.AssertCalled(t, "Help")
			} else {
				helperMock.AssertNotCalled(t, "Help")
			}
		})
	}
}

func Test_nextCmdFlags(t *testing.T) {
	helperMock := __newHelpMock()

	nextCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helperMock.Help()
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	type want struct {
		prepare bool
		backup  bool
		force   bool
	}

	tests := []struct {
		name      string
		arg       string
		want      want
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "call prepare",
			arg:  "--prepare",
			want: want{
				prepare: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "call backup",
			arg:  "--backup",
			want: want{
				backup: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "call force",
			arg:  "--force",
			want: want{
				force: true,
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helperMock = __newHelpMock()
				nextCmd.ResetFlags()
				initNextCmd()
			})

			runnerMock := __newRunnerMock(nil)
			command.SetForce(runnerMock)

			rootCmd.SetArgs([]string{nextCmd.Use, "--patch", tt.arg})

			tt.assertion(t, rootCmd.Execute(), nextCmd.Use+" --patch "+tt.arg+" error = %v, wantErr %v", tt.assertion, false)

			assert.Equal(t, tt.want.prepare, viper.GetBool(key.Prepare), "prepare flag not set")
			assert.Equal(t, tt.want.backup, viper.GetBool(key.Backup), "backup flag not set")
			assert.Equal(t, tt.want.force, viper.GetBool(key.Force), "force flag not set")

		})
	}
}
