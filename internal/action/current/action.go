package current

import (
	"fmt"

	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/types"
	"github.com/klimby/version/pkg/version"
)

// Action - current action.
type Action struct {
	repo actionRepo
}

// actionRepo - repo interface.
type actionRepo interface {
	Current() (version.V, error)
}

// Args - action arguments.
type Args struct {
	Repo actionRepo
}

// New creates new action.
func New(args ...func(arg *Args)) *Action {
	a := &Args{}

	for _, arg := range args {
		arg(a)
	}

	return &Action{
		repo: a.Repo,
	}
}

// Run action.
func (a *Action) Run() error {
	if err := a.validate(); err != nil {
		return err
	}

	v, err := a.repo.Current()
	if err != nil {
		return err
	}

	console.Notice(fmt.Sprintf("Current version: %s", v.FormatString()))

	return nil
}

// validate action.
func (a *Action) validate() error {
	if a.repo == nil {
		return fmt.Errorf("%w: repo is nil in current", types.ErrInvalidArguments)
	}

	return nil
}
