package action

import (
	"errors"
	"fmt"

	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/internal/types"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

// Next - next action.
type Next struct {
	repo  nextArgsRepo
	chGen nextArgsChGen
	cfg   nextArgsConfig
	f     file.ReadWriter
	bump  nextArgsBump
	cmd   nextArgsCmd
}

// nextArgsRepo - repo interface for nextArgs.
type nextArgsRepo interface {
	IsClean() (bool, error)
	NextVersion(nt git.NextType, custom version.V) (version.V, bool, error)
	CheckDowngrade(v version.V) error
	CommitTag(v version.V) error
	AddModified() error
}

// nextArgsChGen - changelog interface for nextArgs.
type nextArgsChGen interface {
	Add(v version.V) error
}

// nextArgsConfig - config interface for nextArgs.
type nextArgsConfig interface {
	BumpFiles() []config.BumpFile
	CommandsBefore() []config.Command
	CommandsAfter() []config.Command
}

// nextArgsBump - bump interface for nextArgs.
type nextArgsBump interface {
	Apply(bumps []config.BumpFile, v version.V)
}

// nextArgsCmd - cmd interface for nextArgs.
type nextArgsCmd interface {
	Run(name string, arg ...string) error
}

// NextArgs - arguments for Next.
type NextArgs struct {
	Repo  nextArgsRepo
	ChGen nextArgsChGen
	Cfg   nextArgsConfig
	F     file.ReadWriter
	Bump  nextArgsBump
	Cmd   nextArgsCmd
}

// NewNext creates new Next.
func NewNext(args ...func(arg *NextArgs)) (*Next, error) {
	a := &NextArgs{
		F: file.NewFS(),
	}

	for _, arg := range args {
		arg(a)
	}

	return &Next{
		repo:  a.Repo,
		chGen: a.ChGen,
		cfg:   a.Cfg,
		f:     a.F,
		bump:  a.Bump,
		cmd:   a.Cmd,
	}, nil
}

// PrepareNextArgs - arguments for Prepare.
type PrepareNextArgs struct {
	NextType git.NextType
	Custom   version.V
}

// Prepare prepares the next version. Run only bump files and commands before.
func (n *Next) Prepare(args ...func(arg *PrepareNextArgs)) (version.V, error) {
	a := PrepareNextArgs{
		NextType: git.NextNone,
		Custom:   version.V(""),
	}

	for _, arg := range args {
		arg(&a)
	}

	if err := n.validate(a); err != nil {
		return version.V(""), err
	}

	if err := n.checkClean(); err != nil {
		return version.V(""), err
	}

	nextV, err := n.nextVersion(a)
	if err != nil {
		return version.V(""), err
	}

	console.Notice(fmt.Sprintf("Bump version to %s...", nextV.FormatString()))

	if err := n.checkDowngrade(nextV); err != nil {
		return version.V(""), err
	}

	n.bump.Apply(n.cfg.BumpFiles(), nextV)

	if err := n.runCommands(n.cfg.CommandsBefore(), nextV); err != nil {
		return version.V(""), err
	}

	return nextV, nil
}

// Apply applies the next version. Run only changelog, add modified files, commit tag and commands after.
func (n *Next) Apply(nextV version.V) error {
	if err := n.writeChangelog(nextV); err != nil {
		if !errors.Is(err, changelog.ErrWarning) {
			return err
		}

		console.Warn(err.Error())
	}

	if err := n.repo.AddModified(); err != nil {
		console.Warn(err.Error())
	}

	if err := n.repo.CommitTag(nextV); err != nil {
		return err
	}

	if err := n.runCommands(n.cfg.CommandsAfter(), nextV); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Version set to %s.", nextV.FormatString()))

	return nil
}

// validate validates the arguments.
func (n *Next) validate(a PrepareNextArgs) error {
	if a.NextType == git.NextNone {
		return fmt.Errorf("next type is not set: %w", types.ErrInvalidArguments)
	}

	if a.NextType == git.NextCustom && a.Custom == "" {
		return fmt.Errorf("custom version is not set: %w", types.ErrInvalidArguments)
	}

	if n.repo == nil {
		return fmt.Errorf("repo is not set: %w", types.ErrNotInitialized)
	}

	if n.chGen == nil {
		return fmt.Errorf("chGen is not set: %w", types.ErrNotInitialized)
	}

	if n.cfg == nil {
		return fmt.Errorf("cfg is not set: %w", types.ErrNotInitialized)
	}

	if n.f == nil {
		return fmt.Errorf("f is not set: %w", types.ErrNotInitialized)
	}

	if n.bump == nil {
		return fmt.Errorf("bump is not set: %w", types.ErrNotInitialized)
	}

	if n.cmd == nil {
		return fmt.Errorf("cmd is not set: %w", types.ErrNotInitialized)
	}

	return nil
}

// checkClean checks if the repository is clean.
func (n *Next) checkClean() error {

	isClean, err := n.repo.IsClean()
	if err != nil {
		return err
	}

	if !isClean {
		if !viper.GetBool(config.AllowCommitDirty) {
			return errors.New("repository is not clean")
		}

		console.Warn("Repository is not clean")
	}

	return nil
}

// nextVersion returns the next version.
func (n *Next) nextVersion(arg PrepareNextArgs) (version.V, error) {
	nextV, exists, err := n.repo.NextVersion(arg.NextType, arg.Custom)
	if err != nil {
		return arg.Custom, err
	}

	if exists {
		if !viper.GetBool(config.AutoGenerateNextPatch) {
			return arg.Custom, fmt.Errorf("version %s already exists", nextV.FormatString())
		}

		console.Warn(fmt.Sprintf("Version already exists. Will be generated next patch version: %s", nextV.FormatString()))
	}

	return nextV, nil
}

// checkDowngrade checks if the version is not downgraded.
func (n *Next) checkDowngrade(v version.V) error {
	if err := n.repo.CheckDowngrade(v); err != nil {
		if !viper.GetBool(config.AllowDowngrades) {
			return err
		}

		console.Warn(err.Error())
	}

	return nil
}

// writeChangelog writes the changelog.
func (n *Next) writeChangelog(v version.V) error {
	if !viper.GetBool(config.GenerateChangelog) {
		return nil
	}

	return n.chGen.Add(v)
}

// runCommands runs commands.
func (n *Next) runCommands(cs []config.Command, v version.V) error {
	dryMode := viper.GetBool(config.DryRun)

	for _, c := range cs {
		if dryMode && !c.RunInDry {
			if viper.GetBool(config.Verbose) {
				console.Info(fmt.Sprintf("Skip command %s in dry mode", c.String()))
			}

			continue
		}

		args := c.Args(v)
		if err := n.cmd.Run(c.Name(), args...); err != nil {
			if c.BreakOnError {
				return err
			}

			console.Warn(err.Error())
		}
	}

	return nil
}
