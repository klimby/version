package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/spf13/cobra"
)

func Test_generateCmd(t *testing.T) {
	helper := &__helpCaller{}

	generateCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
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
			name:           "call config-file",
			arg:            "--config-file",
			wantCallAction: true,
			wantErr:        false,
		},
		{
			name:           "call changelog",
			arg:            "--changelog",
			wantCallAction: true,
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
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				helper.called = false
				generateCmd.ResetFlags()
				initGenerateCmd()
			})

			runner := &__runner{}
			command.SetForce(runner)

			rootCmd.SetArgs([]string{generateCmd.Use, tt.arg})

			err := rootCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("generateCmd() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantCallAction != runner.called {
				t.Errorf("generateCmd() error = %v, wantErr %v", "not called", tt.wantCallAction)
			}

			if tt.wantCallHelp != helper.called {
				t.Errorf("generateCmd() error = %v, wantErr %v", "not called", tt.wantCallHelp)
			}

		})
	}
}
