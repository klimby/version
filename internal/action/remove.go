package action

import (
	"fmt"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/types"
	"github.com/spf13/viper"
)

// Remove - remove action.
type Remove struct {
	fr  file.Remover
	cfg removeArgsConfig
}

// ArgsRemove - arguments for Remove.
type ArgsRemove struct {
	Fr  file.Remover
	Cfg removeArgsConfig
}

// removeArgsConfig - config interface for removeArgs.
type removeArgsConfig interface {
	BumpFiles() []config.BumpFile
}

// NewRemove creates new Remove.
func NewRemove(args ...func(arg *ArgsRemove)) (*Remove, error) {
	a := &ArgsRemove{
		Fr: file.NewFS(),
	}

	for _, arg := range args {
		arg(a)
	}

	if a.Cfg == nil {
		return nil, fmt.Errorf("%w: config is nil in remove", types.ErrInvalidArguments)
	}

	return &Remove{
		fr:  a.Fr,
		cfg: a.Cfg,
	}, nil
}

func (r *Remove) Backup() {
	console.Notice("Remove backup files...")

	p := config.File(viper.GetString(config.CfgFile))

	backup.Remove(r.fr, p.Path())

	for _, bmp := range r.cfg.BumpFiles() {
		backup.Remove(r.fr, bmp.File.Path())
	}

	console.Success("Backup files removed.")
}
