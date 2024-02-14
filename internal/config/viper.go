package config

import (
	"errors"
	"net/url"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/spf13/viper"
)

var (
	errConfig = errors.New("config error")
	// ErrConfigWarn is a config warning.
	ErrConfigWarn = errors.New("config warning")
)

// Options is a configuration options.
type Options struct {
	WorkDir               string
	Version               string
	AllowCommitDirty      bool
	AutoGenerateNextPatch bool
	AllowDowngrades       bool
	GenerateChangelog     bool
	ChangelogFileName     string
	ChangelogTitle        string
	ChangelogIssueHref    string
	ChangelogShowAuthor   bool
	ChangelogShowBody     bool
	Silent                bool
	DryRun                bool
	Backup                bool
	Force                 bool
	Verbose               bool
	ConfigFile            string

	TestingSkipDIInit bool
}

// Init initializes the configuration.
func Init(opts ...func(options *Options)) {
	co := &Options{
		Version:               _Version,
		AllowCommitDirty:      _AllowCommitDirty,
		AutoGenerateNextPatch: _AutoGenerateNextPatch,
		AllowDowngrades:       _AllowDowngrades,
		GenerateChangelog:     _GenerateChangelog,
		ChangelogFileName:     _ChangelogFileName,
		ChangelogTitle:        _ChangelogTitle,
		ConfigFile:            DefaultConfigFile,
		ChangelogShowAuthor:   _ChangelogShowAuthor,
		ChangelogShowBody:     _ChangelogShowBody,
	}

	for _, opt := range opts {
		opt(co)
	}

	viper.Set(key.AppName, _AppName)
	viper.Set(key.Version, co.Version)

	if co.WorkDir != "" {
		viper.Set(key.WorkDir, co.WorkDir)
	}

	viper.Set(key.AllowCommitDirty, co.AllowCommitDirty)
	viper.Set(key.AutoGenerateNextPatch, co.AutoGenerateNextPatch)
	viper.Set(key.AllowDowngrades, co.AllowDowngrades)

	viper.Set(key.GenerateChangelog, co.GenerateChangelog)
	viper.Set(key.ChangelogFileName, co.ChangelogFileName)
	viper.Set(key.ChangelogTitle, co.ChangelogTitle)
	viper.Set(key.ChangelogIssueURL, co.ChangelogIssueHref)
	viper.Set(key.ChangelogShowAuthor, co.ChangelogShowAuthor)
	viper.Set(key.ChangelogShowBody, co.ChangelogShowBody)

	if co.ConfigFile != DefaultConfigFile {
		viper.Set(key.CfgFile, co.ConfigFile)
	}

	if co.Force {
		viper.Set(key.Force, co.Force)
		SetForce()
	}

	if co.Silent {
		viper.Set(key.Silent, co.Silent)
	}

	if co.DryRun {
		viper.Set(key.DryRun, co.DryRun)
	}

	if co.Backup {
		viper.Set(key.Backup, co.Backup)
	}

	if co.Verbose {
		viper.Set(key.Verbose, co.Verbose)
	}

	if co.TestingSkipDIInit {
		viper.Set(key.TestingSkipDIInit, co.TestingSkipDIInit)
	}
}

// SetForce sets force mode.
func SetForce() {
	if viper.GetBool(key.Force) {
		viper.Set(key.AllowCommitDirty, true)
		viper.Set(key.AutoGenerateNextPatch, true)
		viper.Set(key.AllowDowngrades, true)
	}
}

// SetURLFromGit sets the remote repository URL from git.
func SetURLFromGit(u string) {
	if u != "" {
		viper.Set(key.RemoteURL, u)

		iU, err := url.JoinPath(u, "/issues/")
		if err != nil {
			return
		}

		viper.Set(key.ChangelogIssueURL, iU)
	}
}

// LoadArgs is a load configuration arguments.
type LoadArgs struct {
	RW configRW
}

// Load loads the configuration from config.yaml.
// Run after config.Init() in main.go.
func Load(opts ...func(options *LoadArgs)) (C, error) {
	lo := &LoadArgs{
		RW: fsys.New(),
	}

	for _, opt := range opts {
		opt(lo)
	}

	c, err := newConfig(lo.RW)
	if err != nil {
		return c, err
	}

	if c.IsFileConfig {
		if c.GitOptions.AllowCommitDirty {
			viper.Set(key.AllowCommitDirty, c.GitOptions.AllowCommitDirty)
		}

		if c.GitOptions.AutoGenerateNextPatch {
			viper.Set(key.AutoGenerateNextPatch, c.GitOptions.AutoGenerateNextPatch)
		}

		if c.GitOptions.AllowDowngrades {
			viper.Set(key.AllowDowngrades, c.GitOptions.AllowDowngrades)
		}

		viper.Set(key.GenerateChangelog, c.ChangelogOptions.Generate)
		viper.Set(key.ChangelogFileName, c.ChangelogOptions.FileName)
		viper.Set(key.ChangelogTitle, c.ChangelogOptions.Title)
		viper.Set(key.ChangelogShowAuthor, c.ChangelogOptions.ShowAuthor)
		viper.Set(key.ChangelogShowBody, c.ChangelogOptions.ShowBody)

		if c.Backup {
			viper.Set(key.Backup, c.Backup)
		}

		if c.ChangelogOptions.IssueURL != "" {
			viper.Set(key.ChangelogIssueURL, c.ChangelogOptions.IssueURL)
		}

		if c.GitOptions.RemoteURL != "" {
			viper.Set(key.RemoteURL, c.GitOptions.RemoteURL)
		}
	}

	return c, nil
}
