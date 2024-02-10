package action

import (
	"fmt"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/backup"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/types"
	"github.com/spf13/viper"
)

// Generate - generate action.
type Generate struct {
	rw      fsys.ReadWriter
	cfgGen  configGenerator
	clogGen changelogGenerator
}

// configGenerator - config generator interface.
type configGenerator interface {
	Generate(fsys.Writer) error
}

// changelogGenerator - changelog generator interface.
type changelogGenerator interface {
	Generate() error
}

// GenerateArgs - arguments for Generate.
type GenerateArgs struct {
	Rw      fsys.ReadWriter
	CfgGen  configGenerator
	ClogGen changelogGenerator
}

// NewGenerate creates new Generate.
func NewGenerate(args ...func(arg *GenerateArgs)) (*Generate, error) {
	a := &GenerateArgs{
		Rw: fsys.NewFS(),
	}

	for _, arg := range args {
		arg(a)
	}

	if a.CfgGen == nil {
		return nil, fmt.Errorf("%w: config generator is nil in generate", types.ErrInvalidArguments)
	}

	if a.ClogGen == nil {
		return nil, fmt.Errorf("%w: changelog generator is nil in generate", types.ErrInvalidArguments)
	}

	return &Generate{
		rw:      a.Rw,
		cfgGen:  a.CfgGen,
		clogGen: a.ClogGen,
	}, nil
}

// Config generates config file.
func (g *Generate) Config() error {
	console.Notice("Generate config file...")

	p := fsys.File(viper.GetString(key.CfgFile))

	if err := backup.Create(g.rw, p.Path()); err != nil {
		return err
	}

	if err := g.cfgGen.Generate(g.rw); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Config %s created.", p.String()))

	return nil
}

// Changelog generates changelog file.
func (g *Generate) Changelog() error {
	if !viper.GetBool(key.GenerateChangelog) {
		console.Info("Changelog generation disabled.")

		return nil
	}

	console.Notice("Generate changelog...")

	return g.clogGen.Generate()
}
