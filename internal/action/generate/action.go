// Package generate provides generate changelog action.
package generate

import (
	"fmt"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/backup"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/types"
	"github.com/spf13/viper"
)

// ActionType is a file type for generate action.
type ActionType string

// String returns string representation of ActionType.
func (ft ActionType) String() string {
	return string(ft)
}

const (
	// FileUnknown is an unknown file type.
	FileUnknown ActionType = "unknown"
	// FileConfig is a config file type.
	FileConfig ActionType = "config-file"
	// FileChangelog is a changelog file type.
	FileChangelog ActionType = "changelog"
)

// Action - generate action.
type Action struct {
	cfgGen       generator
	changelogGen generator
	backup       backupService
	actionType   ActionType
}

// generator - config generator interface.
type generator interface {
	Generate() error
}

// backupService - backup service interface.
type backupService interface {
	Create(path string) error
}

// Args is an Action arguments.
type Args struct {
	CfgGenerator generator
	ChangelogGen generator
	Backup       backupService
	ActionType   ActionType
}

// New creates new Action.
func New(args ...func(arg *Args)) *Action {
	a := &Args{
		Backup:     backup.New(),
		ActionType: FileUnknown,
	}

	for _, arg := range args {
		arg(a)
	}

	return &Action{
		cfgGen:       a.CfgGenerator,
		changelogGen: a.ChangelogGen,
		backup:       a.Backup,
		actionType:   a.ActionType,
	}
}

// Run action.
func (a Action) Run() error {
	if err := a.validate(); err != nil {
		return err
	}

	switch a.actionType {
	case FileConfig:
		return a.config()
	case FileChangelog:
		return a.changelog()
	default:
		return fmt.Errorf("%w: unknown file type", types.ErrInvalidArguments)
	}
}

// validate action.
func (a Action) validate() error {
	if a.actionType == FileConfig && a.cfgGen == nil {
		return fmt.Errorf("%w: config generator is nil in generate", types.ErrInvalidArguments)
	}

	if a.actionType == FileChangelog && a.changelogGen == nil {
		return fmt.Errorf("%w: changelog generator is nil in generate", types.ErrInvalidArguments)
	}

	return nil
}

// config generates config file.
func (a Action) config() error {
	console.Notice("Generate config file...")

	p := fsys.File(viper.GetString(key.CfgFile))

	if err := a.backup.Create(p.Path()); err != nil {
		return err
	}

	if err := a.cfgGen.Generate(); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Config %s created.", p.String()))

	return nil
}

// changelog generates changelog file.
func (a Action) changelog() error {
	if !viper.GetBool(key.GenerateChangelog) {
		console.Info("Changelog generation disabled.")

		return nil
	}

	console.Notice("Generate changelog...")

	return a.changelogGen.Generate()
}
