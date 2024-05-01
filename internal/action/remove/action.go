// Package remove provides remove action.
package remove

import (
	"fmt"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/backup"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/types"
	"github.com/spf13/viper"
)

// ActionType is a file type for remove action.
type ActionType string

// String returns string representation of ActionType.
func (ft ActionType) String() string {
	return string(ft)
}

const (
	// ActionUnknown is an unknown file type.
	ActionUnknown ActionType = "unknown"
	// ActionBackup is a backup file type. RemoveAll backup files.
	ActionBackup ActionType = "backup"
)

// Action - remove action.
type Action struct {
	remover    remover
	cfg        cfgSrv
	actionType ActionType
}

// remover - remove interface.
type remover interface {
	Remove(path ...string)
}

// cfgSrv - config service interface.
type cfgSrv interface {
	BumpFiles() []config.BumpFile
}

// Args - arguments for Action.
type Args struct {
	Remover    remover
	Cfg        cfgSrv
	ActionType ActionType
}

// New creates new Action.
func New(args ...func(arg *Args)) *Action {
	a := &Args{
		Remover:    backup.New(),
		ActionType: ActionUnknown,
	}

	for _, arg := range args {
		arg(a)
	}

	return &Action{
		remover:    a.Remover,
		cfg:        a.Cfg,
		actionType: a.ActionType,
	}
}

// Run action.
func (a Action) Run() error {
	if err := a.validate(); err != nil {
		return err
	}

	switch a.actionType {
	case ActionBackup:
		a.makeBackup()
	default:
		return fmt.Errorf("%w: unknown file type", types.ErrInvalidArguments)
	}

	return nil
}

// makeBackup - remove backup files.
func (a Action) makeBackup() {
	console.Notice("RemoveAll backup files...")

	c := fsys.File(viper.GetString(key.CfgFile)).Path()
	a.remover.Remove(c)

	for _, bmp := range a.cfg.BumpFiles() {
		a.remover.Remove(bmp.File.Path())
	}

	console.Success("Backup files removed.")
}

// validate - validate action.
func (a Action) validate() error {
	if a.cfg == nil {
		return fmt.Errorf("%w: config is nil in remove", types.ErrInvalidArguments)
	}

	return nil
}
