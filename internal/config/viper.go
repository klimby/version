package config

import (
	"errors"
	"net/url"

	"github.com/klimby/version/internal/file"
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
	ChangelogCommitNames  []CommitName
}

// Init initializes the configuration.
func Init(opts ...func(options *Options)) {
	co := &Options{
		// WorkDir:               _WorkDir,
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
		// Silent:                false,
		// DryRun:                false,
		// Backup:                false,
		ChangelogCommitNames: _defaultCommitNames,
	}

	for _, opt := range opts {
		opt(co)
	}

	viper.Set(appName, _AppName)
	viper.Set(Version, co.Version)

	if co.WorkDir != "" {
		viper.Set(WorkDir, co.WorkDir)
	}

	viper.Set(AllowCommitDirty, co.AllowCommitDirty)
	viper.Set(AutoGenerateNextPatch, co.AutoGenerateNextPatch)
	viper.Set(AllowDowngrades, co.AllowDowngrades)

	viper.Set(GenerateChangelog, co.GenerateChangelog)
	viper.Set(ChangelogFileName, co.ChangelogFileName)
	viper.Set(ChangelogTitle, co.ChangelogTitle)
	viper.Set(ChangelogIssueURL, co.ChangelogIssueHref)
	viper.Set(ChangelogShowAuthor, co.ChangelogShowAuthor)
	viper.Set(ChangelogShowBody, co.ChangelogShowBody)

	if co.ConfigFile != DefaultConfigFile {
		viper.Set(CfgFile, co.ConfigFile)
	}

	if co.Force {
		viper.Set(Force, co.Force)
		SetForce()
	}

	if co.Silent {
		viper.Set(Silent, co.Silent)
	}

	if co.DryRun {
		viper.Set(DryRun, co.DryRun)
	}

	if co.Backup {
		viper.Set(Backup, co.Backup)
	}

	if co.Verbose {
		viper.Set(Verbose, co.Verbose)
	}

	setCommitNames(co.ChangelogCommitNames)
}

func SetForce() {
	if viper.GetBool(Force) {
		viper.Set(AllowCommitDirty, true)
		viper.Set(AutoGenerateNextPatch, true)
		viper.Set(AllowDowngrades, true)
	}
}

// SetURLFromGit sets the remote repository URL from git.
func SetURLFromGit(u string) {
	if u != "" {
		viper.Set(RemoteURL, u)

		iU, err := url.JoinPath(u, "/issues/")
		if err != nil {
			return
		}

		viper.Set(ChangelogIssueURL, iU)
	}
}

// Load loads the configuration from config.yaml.
// Run after config.Init() in main.go.
func Load(f file.Reader) (C, error) {
	c, err := newConfig(f)
	if err != nil {
		return c, err
	}

	if c.IsFileConfig {
		if c.GitOptions.AllowCommitDirty {
			viper.Set(AllowCommitDirty, c.GitOptions.AllowCommitDirty)
		}

		if c.GitOptions.AutoGenerateNextPatch {
			viper.Set(AutoGenerateNextPatch, c.GitOptions.AutoGenerateNextPatch)
		}

		if c.GitOptions.AllowDowngrades {
			viper.Set(AllowDowngrades, c.GitOptions.AllowDowngrades)
		}

		viper.Set(GenerateChangelog, c.ChangelogOptions.Generate)
		viper.Set(ChangelogFileName, c.ChangelogOptions.FileName)
		viper.Set(ChangelogTitle, c.ChangelogOptions.Title)
		viper.Set(ChangelogShowAuthor, c.ChangelogOptions.ShowAuthor)
		viper.Set(ChangelogShowBody, c.ChangelogOptions.ShowBody)

		if c.Backup {
			viper.Set(Backup, c.Backup)
		}

		if c.ChangelogOptions.IssueURL != "" {
			viper.Set(ChangelogIssueURL, c.ChangelogOptions.IssueURL)
		}

		if c.GitOptions.RemoteURL != "" {
			viper.Set(RemoteURL, c.GitOptions.RemoteURL)
		}
	}

	return c, nil
}

func setCommitNames(names []CommitName) {
	mp, order := toViperCommitNames(names)
	viper.Set(changelogCommitTypes, mp)
	viper.Set(changelogCommitOrder, order)
}

// CommitNames returns a list of commit names.
func CommitNames() []CommitName {
	return fromViperCommitNames(viper.GetStringMapString(changelogCommitTypes), viper.GetStringSlice(changelogCommitOrder))
}
