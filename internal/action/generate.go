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

type Generate struct {
	rw      file.ReadWriter
	cfgGen  configGenerator
	clogGen changelogGenerator
}

type configGenerator interface {
	Generate(file.Writer) error
}

type changelogGenerator interface {
	Generate() error
}

type GenerateArgs struct {
	Rw      file.ReadWriter
	CfgGen  configGenerator
	ClogGen changelogGenerator
}

func NewGenerate(args ...func(arg *GenerateArgs)) (*Generate, error) {
	a := &GenerateArgs{
		Rw: file.NewFS(),
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

	p := config.File(viper.GetString(config.CfgFile))

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
	if !viper.GetBool(config.GenerateChangelog) {
		console.Info("Changelog generation disabled.")

		return nil
	}

	console.Notice("Generate changelog...")

	return g.clogGen.Generate()
}
