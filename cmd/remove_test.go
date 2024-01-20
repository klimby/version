package cmd

import (
	"testing"

	"github.com/klimby/version/internal/config"
)

func Test_removeCmd(t *testing.T) {
	config.Init(func(options *config.Options) {
		options.TestingSkipDIInit = true
	})

	h := &__hasRemoveBackup{}

	callRemoveBackup = func() (hasRemoveBackup, error) {
		return h, nil
	}

	rootCmd.SetArgs([]string{"remove", "--backup"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf("removeCmd() error = %v", err)
	}

	if !h.called {
		t.Errorf("removeCmd() error = %v", "not called")
	}
}

type __hasRemoveBackup struct {
	called bool
}

func (h *__hasRemoveBackup) Backup() {
	h.called = true
}
