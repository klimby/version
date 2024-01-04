package di

import (
	"errors"

	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

var C = &Container{}

// Container - DI container.
// Properties are singleton objects.
type Container struct {
	// isInit is true if the container is initialized.
	isInit bool

	// repo object singleton.
	repo *git.Repository

	// ch changelog object singleton.
	ch *changelog.Generator

	// cfg config object singleton.
	cfg *config.C

	f *file.FS
}

// Init initializes the container.
func (c *Container) Init(needUpdateVersion version.V) error {
	if c.isInit {
		panic("container is already initialized")
	}
	c.isInit = true

	c.f = file.NewFS()

	repo, err := git.NewRepository(func(options *git.RepoOptions) {
		options.Path = viper.GetString(config.WorkDir)
	})
	if err != nil {
		return err
	}

	c.repo = repo

	remote, err := repo.RemoteURL()
	if err != nil {
		console.Warn(err.Error())
	}

	config.SetUrlFromGit(remote)

	cfg, err := config.Load(c.f)
	if err != nil {
		return err
	}

	c.cfg = &cfg

	c.ch = changelog.NewGenerator(c.f, repo)

	if err := config.Check(cfg, needUpdateVersion); err != nil {
		if errors.Is(err, config.ErrConfigWarn) {
			console.Warn(err.Error())
		} else {
			return err
		}
	}

	return nil
}

// Repo returns the repo object.
func (c *Container) Repo() *git.Repository {
	if !c.isInit {
		panic("container is not initialized")
	}

	return c.repo
}

// Changelog returns the changelog object.
func (c *Container) Changelog() *changelog.Generator {
	if !c.isInit {
		panic("container is not initialized")
	}

	return c.ch
}

// Config returns the config object.
func (c *Container) Config() *config.C {
	if !c.isInit {
		panic("container is not initialized")
	}

	return c.cfg
}

// FS returns the file system object.
func (c *Container) FS() *file.FS {
	if !c.isInit {
		panic("container is not initialized")
	}

	return c.f
}
