package cmd

import (
	"errors"
	"fmt"

	"github.com/klimby/version/internal/changelog"
	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/internal/di"
	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/internal/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nextCmd represents the next command.
var nextCmd = &cobra.Command{
	Use:           "next",
	Short:         "Generate next version",
	Long:          `Generate next version, bump files, generate changelog.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		major, err := cmd.Flags().GetBool("major")
		if err != nil {
			return err
		}

		if major {
			return next(func(options *nextArgs) {
				options.nt = git.NextMajor
			})
		}

		minor, err := cmd.Flags().GetBool("minor")
		if err != nil {
			return err
		}

		if minor {
			return next(func(options *nextArgs) {
				options.nt = git.NextMinor
			})
		}

		patch, err := cmd.Flags().GetBool("patch")
		if err != nil {
			return err
		}

		if patch {
			return next(func(options *nextArgs) {
				options.nt = git.NextPatch
			})
		}

		ver, err := cmd.Flags().GetString("ver")
		if err != nil {
			return err
		}

		if ver != "" {
			nextV := version.V(ver)
			if nextV.Invalid() {
				return fmt.Errorf("invalid version: %s", nextV)
			}

			return next(func(options *nextArgs) {
				options.custom = nextV
				options.nt = git.NextCustom
			})
		}

		if !major && !minor && !patch && ver == "" {
			if err := cmd.Help(); err != nil {
				return err
			}
		}

		return nil
	},
}

// init - init next command.
func init() {
	nextCmd.Flags().Bool("major", false, "next major version")
	nextCmd.Flags().Bool("minor", false, "next minor version")
	nextCmd.Flags().Bool("patch", false, "next patch version")
	nextCmd.Flags().String("ver", "", "next build version in format 1.2.3")
	nextCmd.MarkFlagsMutuallyExclusive("major", "minor", "patch", "ver")

	rootCmd.PersistentFlags().BoolP("backup", "b", false, "backup changed files")

	if err := viper.BindPFlag(config.Backup, rootCmd.PersistentFlags().Lookup("backup")); err != nil {
		viper.Set(config.Backup, false)
	}

	viper.SetDefault(config.Backup, false)

	rootCmd.PersistentFlags().BoolP("force", "f", false, "force mode")

	if err := viper.BindPFlag(config.Force, rootCmd.PersistentFlags().Lookup("force")); err != nil {
		viper.Set(config.Force, false)
	}

	viper.SetDefault(config.Force, false)

	config.SetForce()

	rootCmd.AddCommand(nextCmd)
}

// nextArgs - arguments for next.
type nextArgs struct {
	nt     git.NextType
	custom version.V
	repo   nextArgsRepo
	chGen  nextArgsChGen
	cfg    nextArgsConfig
	f      file.ReadWriter
	bump   nextArgsBump
	cmd    nextArgsCmd
}

// nextArgsRepo - repo interface for nextArgs.
type nextArgsRepo interface {
	IsClean() (bool, error)
	NextVersion(nt git.NextType, custom version.V) (version.V, bool, error)
	CheckDowngrade(v version.V) error
	CommitTag(v version.V) (*git.Commit, error)
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

// next - generate next version.
func next(opts ...func(options *nextArgs)) error {
	a := &nextArgs{
		nt:     git.NextNone,
		custom: version.V(""),
		repo:   di.C.Repo(),
		chGen:  di.C.Changelog(),
		cfg:    di.C.Config(),
		f:      di.C.FS(),
		bump:   di.C.Bump(),
		cmd:    di.C.Cmd(),
	}

	for _, o := range opts {
		o(a)
	}

	if err := checkClean(a.repo); err != nil {
		return err
	}

	nextV, err := nextVersion(a.repo, a.nt, a.custom)
	if err != nil {
		return err
	}

	console.Notice(fmt.Sprintf("Bump version to %s...", nextV.FormatString()))

	if err := checkDowngrade(a.repo, nextV); err != nil {
		return err
	}

	a.bump.Apply(a.cfg.BumpFiles(), nextV)

	if err := writeChangelog(a.chGen, nextV); err != nil {
		if !errors.Is(err, changelog.ErrWarning) {
			return err
		}

		console.Warn(err.Error())
	}

	if err := runCommands(a.cmd, a.cfg.CommandsBefore(), nextV); err != nil {
		return err
	}

	if err := a.repo.AddModified(); err != nil {
		console.Warn(err.Error())
	}

	if _, err := a.repo.CommitTag(nextV); err != nil {
		return err
	}

	if err := runCommands(a.cmd, a.cfg.CommandsAfter(), nextV); err != nil {
		return err
	}

	console.Success(fmt.Sprintf("Version bumped to %s.", nextV.FormatString()))

	return nil
}

// runCommands runs commands.
func runCommands(cmd nextArgsCmd, cs []config.Command, v version.V) error {
	dryMode := viper.GetBool(config.DryRun)

	for _, c := range cs {
		if dryMode && !c.RunInDry {
			console.Verbose(fmt.Sprintf("Skip command %s in dry mode", c.String()))
			continue
		}

		args := c.Args(v)
		if err := cmd.Run(c.Name(), args...); err != nil {
			if c.BreakOnError {
				return err
			}

			console.Warn(err.Error())
		}
	}

	return nil
}

// checkClean checks if the repository is clean.
func checkClean(repo nextArgsRepo) error {
	clean, err := repo.IsClean()
	if err != nil {
		return err
	}

	if !clean {
		if !viper.GetBool(config.AllowCommitDirty) {
			return errors.New("repository is not clean")
		}

		console.Warn("Repository is not clean.")
	}

	return nil
}

// nextVersion returns the next version.
func nextVersion(repo nextArgsRepo, nt git.NextType, custom version.V) (version.V, error) {
	nextV, exists, err := repo.NextVersion(nt, custom)
	if err != nil {
		return custom, err
	}

	if exists {
		if !viper.GetBool(config.AutoGenerateNextPatch) {
			return custom, fmt.Errorf("version %s already exists", nextV.FormatString())
		}

		console.Warn(fmt.Sprintf("Version already exists. Will be generated next patch version: %s", nextV.FormatString()))
	}

	return nextV, nil
}

// checkDowngrade checks if the version is not downgraded.
func checkDowngrade(repo nextArgsRepo, v version.V) error {
	if err := repo.CheckDowngrade(v); err != nil {
		if !viper.GetBool(config.AllowDowngrades) {
			return err
		}

		console.Warn(err.Error())
	}

	return nil
}

// writeChangelog writes the changelog.
func writeChangelog(g nextArgsChGen, v version.V) error {
	if !viper.GetBool(config.GenerateChangelog) {
		return nil
	}

	return g.Add(v)
}
