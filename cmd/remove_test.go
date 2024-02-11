package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/spf13/cobra"
)

func Test_removeCmd(t *testing.T) {
	helper := &__helpCaller{}

	removeCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		t.Helper()
		helper.called = true
	})

	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	tests := []struct {
		name           string
		arg            string
		wantCallBackup bool
		wantCallHelp   bool
		wantErr        bool
	}{
		{
			name:           "call backup",
			arg:            "--backup",
			wantCallBackup: true,
			wantErr:        false,
		},
		{
			name:         "call without args",
			arg:          "",
			wantCallHelp: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()

			t.Cleanup(func() {
				helper.called = false
				removeCmd.ResetFlags()
				initRemoveCmd()
			})

			runner := &__runner{}
			command.SetForce(runner)

			rootCmd.SetArgs([]string{"remove", tt.arg})

			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("removeCmd() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantCallBackup != runner.called {
				t.Errorf("removeCmd() error = %v, wantErr %v", "not called", tt.wantCallBackup)
			}

			if tt.wantCallHelp != helper.called {
				t.Errorf("removeCmd() error = %v, wantErr %v", "not called", tt.wantCallHelp)
			}
		})
	}

}

type __runner struct {
	called bool
}

func (r *__runner) Run() error {
	r.called = true

	return nil
}

type __helpCaller struct {
	called bool
}

func (h *__helpCaller) Help() {
	h.called = true
}
