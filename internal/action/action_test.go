package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunner_Set(t *testing.T) {
	runner := &Runner{}

	err := runner.Run()
	if err == nil {
		t.Errorf("Runner.Run() error = %v, wantErr %v", err, true)
	}

	runner.Set(&__hasPun{instance: 1})
	runner.Set(&__hasPun{instance: 2})

	assert.Equal(t, 1, runner.action.(*__hasPun).instance)

	runner.SetForce(&__hasPun{instance: 3})

	assert.Equal(t, 3, runner.action.(*__hasPun).instance)

	err = runner.Run()
	if err != nil {
		t.Errorf("Runner.Run() error = %v, wantErr %v", err, false)
	}
}

type __hasPun struct {
	instance int
}

func (h *__hasPun) Run() error {
	return nil
}
