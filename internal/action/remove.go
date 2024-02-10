package action

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

// Remove - remove action.
type Remove struct {
	fr  fsys.Remover
	cfg removeArgsConfig
}

// ArgsRemove - arguments for Remove.
type ArgsRemove struct {
	Fr  fsys.Remover
	Cfg removeArgsConfig
}

// removeArgsConfig - config interface for removeArgs.
type removeArgsConfig interface {
	BumpFiles() []config.BumpFile
}

// NewRemove creates new Remove.
func NewRemove(args ...func(arg *ArgsRemove)) (*Remove, error) {
	a := &ArgsRemove{
		Fr: fsys.NewFS(),
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

	p := fsys.File(viper.GetString(key.CfgFile))

	backup.Remove(r.fr, p.Path())

	for _, bmp := range r.cfg.BumpFiles() {
		backup.Remove(r.fr, bmp.File.Path())
	}

	console.Success("Backup files removed.")
}
