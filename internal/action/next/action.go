package next

import (
	"errors"
	"fmt"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/changelog"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/git"
	"github.com/klimby/version/internal/types"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

// ActionType is a file type for next action.
type ActionType string

// String returns string representation of ActionType.
func (a ActionType) String() string {
	return string(a)
}

// gitNextType - git next type.
func (a ActionType) gitNextType() git.NextType {
	switch a {
	case ActionMajor:
		return git.NextMajor
	case ActionMinor:
		return git.NextMinor
	case ActionPatch:
		return git.NextPatch
	case ActionCustom:
		return git.NextCustom
	default:
		return git.NextNone
	}
}

// ActionType values.
const (
	ActionUnknown ActionType = "unknown" // ActionUnknown is an unknown action type.
	ActionMajor   ActionType = "major"   // ActionMajor is next major version.
	ActionMinor   ActionType = "minor"   // ActionMinor is next minor version.
	ActionPatch   ActionType = "patch"   // ActionPatch is next patch version.
	ActionCustom  ActionType = "ver"     // ActionCustom is custom version. Requires custom version.
)

// Action - next action.
type Action struct {
	actionType    ActionType
	repo          actionRepo
	changelogGen  actionChGen
	cfg           actionCfg
	bump          actionBump
	cmd           actionCmd
	customVersion version.V
}

// actionRepo - repo interface for nextArgs.
type actionRepo interface {
	IsClean() (bool, error)
	NextVersion(nt git.NextType, custom version.V) (version.V, bool, error)
	CheckDowngrade(v version.V) error
	CommitTag(v version.V) error
	AddModified() error
}

// actionChGen - changelog interface for nextArgs.
type actionChGen interface {
	Add(v version.V) error
}

// actionCfg - config interface for nextArgs.
type actionCfg interface {
	BumpFiles() []config.BumpFile
	CommandsBefore() []config.Command
	CommandsAfter() []config.Command
}

// actionBump - bump interface for nextArgs.
type actionBump interface {
	Apply(bumps []config.BumpFile, v version.V)
}

// actionCmd - cmd interface for nextArgs.
type actionCmd interface {
	Run(name string, arg ...string) error
}

// Args - arguments for Next.
type Args struct {
	Repo         actionRepo
	ChangelogGen actionChGen
	Cfg          actionCfg
	Bump         actionBump
	Cmd          actionCmd
	ActionType   ActionType
	Version      version.V
}

// New creates new Action.
func New(args ...func(arg *Args)) *Action {
	a := &Args{
		ActionType: ActionUnknown,
		Cmd:        console.NewCmd(),
	}

	for _, arg := range args {
		arg(a)
	}

	return &Action{
		repo:          a.Repo,
		changelogGen:  a.ChangelogGen,
		cfg:           a.Cfg,
		bump:          a.Bump,
		cmd:           a.Cmd,
		actionType:    a.ActionType,
		customVersion: a.Version,
	}
}

// Run action.
func (a *Action) Run() error {
	if err := a.validate(); err != nil {
		return err
	}

	nextV, err := a.prepare()
	if err != nil {
		return err
	}

	if viper.GetBool(key.Prepare) {
		console.Success(fmt.Sprintf("Prepare complete, next version is %s", nextV.FormatString()))

		return nil
	}

	if err := a.apply(nextV); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Version set to %s.", nextV.FormatString()))

	return nil
}

// prepare the next version. Run only bump files and commands before.
func (a *Action) prepare() (version.V, error) {

	if err := a.checkClean(); err != nil {
		return "", err
	}

	nextV, err := a.nextVersion()
	if err != nil {
		return "", err
	}

	console.Notice(fmt.Sprintf("Bump version to %s...", nextV.FormatString()))

	if err := a.checkDowngrade(nextV); err != nil {
		return nextV, err
	}

	a.bump.Apply(a.cfg.BumpFiles(), nextV)

	if err := a.runCommands(a.cfg.CommandsBefore(), nextV); err != nil {
		return nextV, err
	}

	return nextV, nil
}

// apply the next version. Run only changelog, add modified files, commit tag and commands after.
func (a *Action) apply(nextV version.V) error {
	if err := a.writeChangelog(nextV); err != nil {
		if !errors.Is(err, changelog.ErrWarning) {
			return err
		}

		console.Warn(err.Error())
	}

	if err := a.repo.AddModified(); err != nil {
		console.Warn(err.Error())
	}

	if err := a.repo.CommitTag(nextV); err != nil {
		return err
	}

	if err := a.runCommands(a.cfg.CommandsAfter(), nextV); err != nil {
		return err
	}

	return nil
}

// validate action.
func (a *Action) validate() error {
	if a.actionType == ActionUnknown {
		return fmt.Errorf("%w: action type is unknown", types.ErrInvalidArguments)
	}

	if a.actionType == ActionCustom && a.customVersion.Invalid() {
		return fmt.Errorf("%w: custom version is empty or invalid", types.ErrInvalidArguments)
	}

	if a.repo == nil {
		return fmt.Errorf("%w: repo is nil", types.ErrInvalidArguments)
	}

	if a.changelogGen == nil {
		return fmt.Errorf("%w: changelog generator is nil", types.ErrInvalidArguments)
	}

	if a.cfg == nil {
		return fmt.Errorf("%w: config is nil", types.ErrInvalidArguments)
	}

	if a.bump == nil {
		return fmt.Errorf("%w: bump is nil", types.ErrInvalidArguments)
	}

	if a.cmd == nil {
		return fmt.Errorf("%w: cmd is nil", types.ErrInvalidArguments)
	}

	return nil
}

// checkClean checks if the repository is clean.
func (a *Action) checkClean() error {
	isClean, err := a.repo.IsClean()
	if err != nil {
		return err
	}

	if !isClean {
		const msg = "repository is not clean"

		if !viper.GetBool(key.AllowCommitDirty) {
			return errors.New(msg)
		}

		console.Warn(msg)
	}

	return nil
}

// nextVersion returns the next version.
func (a *Action) nextVersion() (version.V, error) {
	nextV, exists, err := a.repo.NextVersion(a.actionType.gitNextType(), a.customVersion)
	if err != nil {
		return "", err
	}

	if exists {
		if !viper.GetBool(key.AutoGenerateNextPatch) {
			return "", fmt.Errorf("version %s already exists", nextV.FormatString())
		}

		console.Warn(fmt.Sprintf("Version already exists. Will be generated next patch version: %s", nextV.FormatString()))
	}

	return nextV, nil
}

// checkDowngrade checks if the version is not downgraded.
func (a *Action) checkDowngrade(v version.V) error {
	if err := a.repo.CheckDowngrade(v); err != nil {
		if !viper.GetBool(key.AllowDowngrades) {
			return err
		}

		console.Warn(err.Error())
	}

	return nil
}

// writeChangelog writes the changelog.
func (a *Action) writeChangelog(v version.V) error {
	if !viper.GetBool(key.GenerateChangelog) {
		return nil
	}

	return a.changelogGen.Add(v)
}

// runCommands runs commands.
func (a *Action) runCommands(cs []config.Command, v version.V) error {
	dryMode := viper.GetBool(key.DryRun)

	for _, c := range cs {
		if dryMode && !c.RunInDry {
			if viper.GetBool(key.Verbose) {
				console.Info(fmt.Sprintf("Skip command %s in dry mode", c.String()))
			}

			continue
		}

		args := c.Args(v)
		if err := a.cmd.Run(c.Name(), args...); err != nil {
			if c.BreakOnError {
				return err
			}

			console.Warn(err.Error())
		}
	}

	return nil
}
