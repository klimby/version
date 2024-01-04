package action

import (
	"errors"
	"fmt"

	"github.com/klimby/version/internal/backup"
	"github.com/klimby/version/internal/bump"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

var (
	ErrAction = errors.New("action error")
)

type A struct {
	repo  actRepo
	cfg   config.C
	f     file.ReadWriter
	chGen changelogGenerator
}

type actRepo interface {
	IsClean() (bool, error)
	NextVersion(nt git.NextType, custom version.V) (version.V, bool, error)
	CheckDowngrade(v version.V) error
}

type changelogGenerator interface {
	Generate() error
	Add(from version.V) error
}

type Args struct {
	ReadWriter file.ReadWriter
	Repo       actRepo
	Config     config.C
	ChGen      changelogGenerator
}

func New(opt ...func(o *Args)) *A {
	args := &Args{
		ReadWriter: file.NewFS(),
	}

	for _, o := range opt {
		o(args)
	}

	return &A{
		repo:  args.Repo,
		cfg:   args.Config,
		f:     args.ReadWriter,
		chGen: args.ChGen,
	}
}

// GenerateConfig generates a config file.
func (a *A) GenerateConfig() error {
	p := config.File(viper.GetString(config.ConfigFile))

	if err := backup.Create(a.f, p.Path()); err != nil {
		return err
	}

	return config.Generate(a.f, a.cfg)
}

// GenerateChangelog generates a changelog file.
func (a *A) GenerateChangelog() error {
	return a.chGen.Generate()
}

// RemoveBackup removes backup files.
func (a *A) RemoveBackup(f file.Remover) {
	p := config.File(viper.GetString(config.ConfigFile))
	backup.Remove(f, p.Path())

	for _, bmp := range a.cfg.Bump {
		backup.Remove(f, bmp.File.Path())
	}
}

// NextVersion returns the next version.
func (a *A) NextVersion(nt git.NextType, custom version.V) error {
	if nt == git.NextCustom && custom.Invalid() {
		return fmt.Errorf("%w: invalid version %s", ErrAction, custom.FormatString())
	}

	if err := a.checkClean(); err != nil {
		return err
	}

	nextV, err := a.nextVersion(nt, custom)
	if err != nil {
		return err
	}

	if err := a.checkDowngrade(nextV); err != nil {
		return err
	}

	bump.Apply(a.f, a.cfg.Bump, nextV)

	return nil
}

// checkClean checks if the repository is clean.
func (a *A) checkClean() error {
	clean, err := a.repo.IsClean()
	if err != nil {
		return err
	}

	if !clean {
		if !viper.GetBool(config.AllowCommitDirty) {
			return fmt.Errorf("%w: repository is not clean", ErrAction)
		}

		console.Warn("Repository is not clean.")
	}

	return nil
}

// nextVersion returns the next version.
func (a *A) nextVersion(nt git.NextType, custom version.V) (version.V, error) {
	next, exists, err := a.repo.NextVersion(nt, custom)
	if err != nil {
		return custom, err
	}

	if exists {
		if !viper.GetBool(config.AutoGenerateNextPatch) {
			return custom, fmt.Errorf("%w: version %s already exists", ErrAction, next.FormatString())
		}

		console.Warn(fmt.Sprintf("Version already exists. Will be generated next patch version: %s", next.FormatString()))
	}

	return next, nil
}

// checkDowngrade checks if the version is not downgraded.
func (a *A) checkDowngrade(v version.V) error {
	if err := a.repo.CheckDowngrade(v); err != nil {
		if !viper.GetBool(config.AllowDowngrades) {
			return err
		}

		console.Warn(err.Error())
	}

	return nil
}

// writeChangelog writes the changelog.
func (a *A) writeChangelog(v version.V) error {
	if !viper.GetBool(config.GenerateChangelog) {
		return nil
	}

	if err := a.chGen.Add(v); err != nil {
		return err
	}

	return nil
}
