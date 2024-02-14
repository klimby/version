package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/stretchr/testify/assert"
)

func Test_currentCmd(t *testing.T) {
	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	t.Run("current", func(t *testing.T) {
		runnerMock := __newRunnerMock(nil)
		command.SetForce(runnerMock)

		rootCmd.SetArgs([]string{currentCmd.Use})

		assert.NoError(t, rootCmd.Execute(), "currentCmd()")

		runnerMock.AssertExpectations(t)
	})
}
