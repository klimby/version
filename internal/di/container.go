// Package di provides dependency injection container.
package di

import (
	"errors"
	"os"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/bump"
	"github.com/klimby/version/internal/service/changelog"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/internal/service/git"
	"github.com/spf13/viper"
)

// C is the DI container.
var C = &container{}

// container - DI container.
// Properties are singleton objects.
type container struct {
	// IsInit is true if the container is initialized.
	IsInit bool

	// Repo object singleton.
	Repo *git.Repository

	// ChangelogGenerator changelog object singleton.
	ChangelogGenerator *changelog.Generator

	// Config object singleton.
	Config *config.C

	// Bump object singleton.
	Bump *bump.B
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

	repo, err := git.NewRepository()
	if err != nil {
		return err
	}

	c.Repo = repo

	remote, err := repo.RemoteURL()
	if err != nil {
		console.Warn(err.Error())
	}

	config.SetURLFromGit(remote)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	c.Config = &cfg

	if err := cfg.Validate(); err != nil {
		if !errors.Is(err, config.ErrConfigWarn) {
			return err
		}

		console.Warn(err.Error())
	}

	c.ChangelogGenerator = changelog.New(func(options *changelog.Args) {
		options.Repo = c.Repo
		options.ConfigFile = fsys.File(viper.GetString(key.ChangelogFileName))
		options.CommitNames = cfg.CommitTypes()
	})

	c.Bump = bump.New(func(arg *bump.Args) {
		arg.Repo = c.Repo
	})

	return nil
}
