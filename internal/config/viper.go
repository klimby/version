package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/pkg/convert"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
)

var (
	errConfig     = errors.New("config error")
	ErrConfigWarn = errors.New("config warning")
)

// ConfigOptions is a configuration options.
type ConfigOptions struct {
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
	ConfigFile            string
	ChangelogCommitNames  []CommitName
}

// Init initializes the configuration.
func Init(opts ...func(options *ConfigOptions)) {
	co := &ConfigOptions{
		WorkDir:               _WorkDir,
		Version:               _Version,
		AllowCommitDirty:      _AllowCommitDirty,
		AutoGenerateNextPatch: _AutoGenerateNextPatch,
		AllowDowngrades:       _AllowDowngrades,
		GenerateChangelog:     _GenerateChangelog,
		ChangelogFileName:     _ChangelogFileName,
		ChangelogTitle:        _ChangelogTitle,
		ConfigFile:            _ConfigFile,
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

	viper.Set(AppName, _AppName)
	viper.Set(Version, co.Version)
	viper.Set(WorkDir, convert.EnvToString(WorkDir, co.WorkDir))

	viper.Set(AllowCommitDirty, co.AllowCommitDirty)
	viper.Set(AutoGenerateNextPatch, co.AutoGenerateNextPatch)
	viper.Set(AllowDowngrades, co.AllowDowngrades)

	viper.Set(GenerateChangelog, co.GenerateChangelog)
	viper.Set(ChangelogFileName, co.ChangelogFileName)
	viper.Set(ChangelogTitle, co.ChangelogTitle)
	viper.Set(ChangelogIssueURL, co.ChangelogIssueHref)
	viper.Set(ChangelogShowAuthor, co.ChangelogShowAuthor)
	viper.Set(ChangelogShowBody, co.ChangelogShowBody)

	viper.Set(ConfigFile, co.ConfigFile)

	viper.Set(Silent, co.Silent)
	viper.Set(DryRun, co.DryRun)
	viper.Set(Backup, co.Backup)

	setCommitNames(co.ChangelogCommitNames)

}

type FlagOptions struct {
	Silent     bool
	DryRun     bool
	Force      bool
	Backup     bool
	ConfigFile string
}

// SetFlags sets the flags.
// Run in main.go after Parse().
func SetFlags(opts ...func(options *FlagOptions)) {
	fo := &FlagOptions{
		Silent:     false,
		DryRun:     false,
		Force:      false,
		Backup:     false,
		ConfigFile: _ConfigFile,
	}

	for _, opt := range opts {
		opt(fo)
	}

	viper.Set(Silent, fo.Silent)
	viper.Set(DryRun, fo.DryRun)
	viper.Set(Backup, fo.Backup)
	viper.Set(ConfigFile, fo.ConfigFile)

	if fo.Force {
		viper.Set(AllowCommitDirty, true)
		viper.Set(AutoGenerateNextPatch, true)
		viper.Set(AllowDowngrades, true)
	}
}

// SetUrlFromGit sets the remote repository URL from git.
func SetUrlFromGit(u string) {
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
		viper.Set(AllowCommitDirty, c.GitOptions.AllowCommitDirty)
		viper.Set(AutoGenerateNextPatch, c.GitOptions.AutoGenerateNextPatch)
		viper.Set(AllowDowngrades, c.GitOptions.AllowDowngrades)

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

// Check checks the configuration.
// Run after config.Load() in main.go.
func Check(c C) error {
	if !c.IsFileConfig {
		return nil
	}

	for _, f := range c.Bump {
		// check if file exists
		if _, err := os.Stat(f.File.Path()); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf(`%w: file %s does not exist`, errConfig, f.File)
			} else {
				return fmt.Errorf(`%w: file %s error: %v`, errConfig, f.File, err)
			}
		}

		if f.IsPredefinedJSON() {
			continue
		}

		if f.Start > f.End {
			return fmt.Errorf(`%w: file %s start position is greater than end position`, errConfig, f.File)
		}

		if len(f.RegExp) > 0 {
			for _, r := range f.RegExp {
				if _, err := regexp.Compile(r); err != nil {
					return fmt.Errorf(`%w: file %s regexp %s error: %w`, errConfig, f.File, r, err)
				}
			}
		}
	}

	appVersion := version.V(viper.GetString(Version))
	if appVersion.Compare(c.Version) != 0 {
		return fmt.Errorf(`%w: you use older version of config file. For update run "version --generate-config"`, ErrConfigWarn)
	}

	return nil
}

func setCommitNames(names []CommitName) {
	mp, order := toViperCommitNames(names)
	viper.Set(ChangelogCommitTypes, mp)
	viper.Set(ChangelogCommitOrder, order)
}

func CommitNames() []CommitName {
	return fromViperCommitNames(viper.GetStringMapString(ChangelogCommitTypes), viper.GetStringSlice(ChangelogCommitOrder))
}
