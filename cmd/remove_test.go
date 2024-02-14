package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_removeCmd(t *testing.T) {
	helperMock := __newHelpMock()

	removeCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helperMock.Help()
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	type wantCall struct {
		backup bool
		help   bool
	}

	tests := []struct {
		name      string
		arg       string
		wantCall  wantCall
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "call backup",
			arg:  "--backup",
			wantCall: wantCall{
				backup: true,
			},
			assertion: assert.NoError,
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helperMock = __newHelpMock()
				removeCmd.ResetFlags()
				initRemoveCmd()
			})

			runnerMock := __newRunnerMock(nil)
			command.SetForce(runnerMock)

			rootCmd.SetArgs([]string{"remove", tt.arg})

			tt.assertion(t, rootCmd.Execute(), "rootCmd.Execute() error = %v, wantErr %v", tt.assertion, false)

			if tt.wantCall.backup {
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
