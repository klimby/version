package di

import (
	"errors"
	"os"

	"github.com/klimby/version/internal/action"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/bump"
	"github.com/klimby/version/internal/service/changelog"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/service/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

// C is the DI container.
var C = &container{}

// container - DI container.
// Properties are singleton objects.
type container struct {
	// IsInit is true if the container is initialized.
	IsInit bool

	// repo object singleton.
	repo *git.Repository

	// ch changelog object singleton.
	ch *changelog.Generator

	// cfg config object singleton.
	cfg *config.C

	f *fsys.FS

	// bump object singleton.
	bump *bump.B

	// cmd object singleton.
	c *console.Cmd

	// ActionRemove object singleton.
	ActionRemove actionRemove

	// ActionGenerate object singleton.
	ActionGenerate actionGenerate

	// ActionNext object singleton.
	ActionNext actionNext

	// ActionCurrent object singleton.
	ActionCurrent actionCurrent
}

type actionRemove interface {
	Backup()
}

type actionGenerate interface {
	Config() error
	Changelog() error
}

type actionNext interface {
	Prepare(args ...func(arg *action.PrepareNextArgs)) (version.V, error)
	Apply(nextV version.V) error
}

type actionCurrent interface {
	Current() (version.V, error)
}

// Init initializes the container.
func (c *container) Init() error {
	if viper.GetBool(key.TestingSkipDIInit) {
		c.IsInit = true

		return nil
	}

	if !viper.GetBool(key.Silent) {
		console.Init(func(args *console.OutArgs) {
			args.Stderr = os.Stderr
			args.Stdout = os.Stdout
			args.Colorize = true
		})
	}

	if c.IsInit {
		return errors.New("container is already initialized")
	}
	c.IsInit = true

	c.f = fsys.NewFS()

	repo, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = viper.GetString(key.WorkDir)
	})
	if err != nil {
		return err
	}

	c.repo = repo

	remote, err := repo.RemoteURL()
	if err != nil {
		console.Warn(err.Error())
	}

	config.SetURLFromGit(remote)

	cfg, err := config.Load(c.f)
	if err != nil {
		return err
	}

	c.cfg = &cfg

	if err := cfg.Validate(); err != nil {
		if !errors.Is(err, config.ErrConfigWarn) {
			return err
		}

		console.Warn(err.Error())
	}

	c.ch = changelog.New(func(options *changelog.Args) {
		options.RW = c.f
		options.Repo = c.repo
		options.ConfigFile = fsys.File(viper.GetString(key.ChangelogFileName))
		options.CommitNames = cfg.CommitTypes()
	})

	c.bump = bump.New(func(arg *bump.Args) {
		arg.RW = c.f
		arg.Repo = c.repo
	})

	c.c = console.NewCmd()

	a, err := action.NewRemove(func(options *action.ArgsRemove) {
		options.Cfg = cfg
	})

	if err != nil {
		return err
	}

	c.ActionRemove = a

	aG, err := action.NewGenerate(func(options *action.GenerateArgs) {
		options.Rw = c.f
		options.CfgGen = cfg
		options.ClogGen = c.ch
	})
	if err != nil {
		return err
	}

	c.ActionGenerate = aG

	aN, err := action.NewNext(func(args *action.NextArgs) {
		args.Repo = c.repo
		args.ChGen = c.ch
		args.Cfg = c.cfg
		args.F = c.f
		args.Bump = c.bump
		args.Cmd = c.c
	})
	if err != nil {
		return err
	}

	c.ActionNext = aN

	c.ActionCurrent = c.repo

	return nil
}
