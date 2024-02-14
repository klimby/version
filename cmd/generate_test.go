package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_generateCmd(t *testing.T) {
	helperMock := __newHelpMock()

	generateCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
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
		name     string
		arg      string
		wantCall wantCall
	}{
		{
			name: "call config-file",
			arg:  "--config-file",
			wantCall: wantCall{
				action: true,
			},
		},
		{
			name: "call changelog",
			arg:  "--changelog",
			wantCall: wantCall{
				action: true,
			},
		},
		{
			name: "call without args",
			arg:  "",
			wantCall: wantCall{
				help: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helperMock = __newHelpMock()
				generateCmd.ResetFlags()
				initGenerateCmd()
			})

			runnerMock := __newRunnerMock(nil)
			command.SetForce(runnerMock)

			rootCmd.SetArgs([]string{generateCmd.Use, tt.arg})

			assert.NoError(t, rootCmd.Execute(), "generateCmd()")

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
