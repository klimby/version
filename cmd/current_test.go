package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
)

func Test_currentCmd(t *testing.T) {
	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	t.Run("current", func(t *testing.T) {
		//t.Parallel()

		runner := &__runner{}
		command.SetForce(runner)

		rootCmd.SetArgs([]string{currentCmd.Use})

		err := rootCmd.Execute()
		if err != nil {
			t.Errorf("currentCmd() error = %v, wantErr %v", err, false)
		}

		if !runner.called {
			t.Errorf("currentCmd() runner.called = %v, want %v", runner.called, true)
		}

	})
}
