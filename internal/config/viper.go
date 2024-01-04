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
	Force                 bool
	ConfigFile            string
	ChangelogCommitNames  []commitName
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

	setCommitNames(co.ChangelogCommitNames)
}

type FlagOptions struct {
	Silent     bool
	DryRun     bool
	Force      bool
	Backup     bool
	ConfigFile string
}

func SetForce() {
	if viper.GetBool(Force) {
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

// Check checks the configuration.
// Run after config.Load() in main.go.
func Check(c C, needUpdateVersion version.V) error {
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

	if !needUpdateVersion.Empty() && c.Version.LessThen(needUpdateVersion) {
		return fmt.Errorf(`%w: you use older version of config file. For update run "version generate --config"`, ErrConfigWarn)
	}

	return nil
}

func setCommitNames(names []commitName) {
	mp, order := toViperCommitNames(names)
	viper.Set(changelogCommitTypes, mp)
	viper.Set(changelogCommitOrder, order)
}

func CommitNames() []commitName {
	return fromViperCommitNames(viper.GetStringMapString(changelogCommitTypes), viper.GetStringSlice(changelogCommitOrder))
}
