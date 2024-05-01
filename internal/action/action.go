// Package action add package action
package action

import "errors"

// Runner - action runner.
type Runner struct {
	action hasRun
}

type hasRun interface {
	Run() error
}

// Set sets the runner.
func (a *Runner) Set(r hasRun) {
	if a.action == nil {
		a.action = r
	}
}

// SetForce sets the runner.
func (a *Runner) SetForce(r hasRun) {
	a.action = r
}

// Run runs the action.
func (a *Runner) Run() error {
	if a.action == nil {
		return errors.New("action not set")
	}

	return a.action.Run()
}
