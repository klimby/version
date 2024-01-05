package di

import (
	"errors"

	"github.com/klimby/version/internal/bump"
	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
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

	f *file.FS

	// bump object singleton.
	bump *bump.B
}

// Init initializes the container.
func (c *container) Init(needUpdateVersion version.V) error {
	if c.IsInit {
		return errors.New("container is already initialized")
	}
	c.IsInit = true

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

	config.SetURLFromGit(remote)

	cfg, err := config.Load(c.f)
	if err != nil {
		return err
	}

	c.cfg = &cfg

	c.ch = changelog.NewGenerator(c.f, repo)

	if err := config.Check(cfg, needUpdateVersion); err != nil {
		if !errors.Is(err, config.ErrConfigWarn) {
			return err
		}

		console.Warn(err.Error())
	}

	c.bump = bump.New(func(arg *bump.Args) {
		arg.RW = c.f
		arg.Repo = c.repo
	})

	return nil
}

// Repo returns the repo object.
func (c *container) Repo() *git.Repository {
	return c.repo
}

// Changelog returns the changelog object.
func (c *container) Changelog() *changelog.Generator {
	return c.ch
}

// Config returns the config object.
func (c *container) Config() *config.C {
	return c.cfg
}

// FS returns the file system object.
func (c *container) FS() *file.FS {
	return c.f
}

// Bump returns the bump object.
func (c *container) Bump() *bump.B {
	return c.bump
}
