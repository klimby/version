package config

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	// _VersionWarningUpdate is a version for warning update.
	_VersionWarningUpdate = version.V("")
	// _VersionCriticalUpdate is a version for critical update.
	_VersionCriticalUpdate = version.V("")
)

// C is a configuration file.
type C struct {
	// Version is a version of the application.
	Version version.V `yaml:"version"`
	// IsFileConfig is a flag that indicates that the configuration is from a file.
	IsFileConfig bool `yaml:"-"`
	// BackupChangedFiles is a flag that indicates that the changed files are backed up.
	Backup bool `yaml:"backupChanged"`
	// Before is a list of commands that are executed before the main command.
	Before []Command `yaml:"before"`
	// After is a list of commands that are executed after the main command.
	After []Command `yaml:"after"`
	// GitOptions is a git options.
	GitOptions gitOptions `yaml:"git"`
	// ChangelogOptions is a changelog options.
	ChangelogOptions changelogOptions `yaml:"changelog"`
	// Bump is a list of files for bump.
	Bump []BumpFile `yaml:"bump"`
}

// newConfig returns a new configuration.
func newConfig(f file.Reader) (_ C, err error) {
	c := C{
		Version: version.V(viper.GetString(Version)),
		Backup:  viper.GetBool(Backup),
		Before:  []Command{},
		After:   []Command{},
		GitOptions: gitOptions{
			AllowCommitDirty:      viper.GetBool(AllowCommitDirty),
			AutoGenerateNextPatch: viper.GetBool(AutoGenerateNextPatch),
			AllowDowngrades:       viper.GetBool(AllowDowngrades),
			RemoteURL:             viper.GetString(RemoteURL),
		},
		ChangelogOptions: changelogOptions{
			Generate:    viper.GetBool(GenerateChangelog),
			FileName:    File(viper.GetString(ChangelogFileName)),
			Title:       viper.GetString(ChangelogTitle),
			IssueURL:    viper.GetString(ChangelogIssueURL),
			ShowAuthor:  viper.GetBool(ChangelogShowAuthor),
			ShowBody:    viper.GetBool(ChangelogShowBody),
			CommitTypes: _defaultCommitNames,
		},
	}

	cfg := File(viper.GetString(CfgFile))

	r, err := f.Read(cfg.Path())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return c, nil
		}

		return c, fmt.Errorf("open config file error: %w", err)
	}

	defer func() {
		if e := r.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close config file error: %w", e)
			}
		}
	}()

	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return c, fmt.Errorf("decode config file error: %w", err)
	}

	c.IsFileConfig = true

	return c, nil
}

// BumpFiles returns a list of files for bump.
func (c C) BumpFiles() []BumpFile {
	return c.Bump
}

// CommandsBefore returns a list of commands that are executed before the main command.
func (c C) CommandsBefore() []Command {
	return c.Before
}

// CommandsAfter returns a list of commands that are executed after the main command.
func (c C) CommandsAfter() []Command {
	return c.After
}

// CommitTypes returns a commit types for changelog.
func (c C) CommitTypes() []CommitName {
	return c.ChangelogOptions.CommitTypes
}

// Generate generates the configuration file.
func (c C) Generate(f file.Writer) error {
	p := File(viper.GetString(CfgFile))

	w, err := f.Write(p.Path(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("open config file error: %w", err)
	}

	defer func() {
		if e := w.Close(); e != nil {
			if err == nil {
				err = fmt.Errorf("close config file error: %w", e)
			}
		}
	}()

	c.Version = version.V(viper.GetString(Version))

	tmpl, err := template.New("config").Parse(_configYamlTemplate)
	if err != nil {
		return fmt.Errorf("parse config template error: %w", err)
	}

	if err := tmpl.Execute(w, c); err != nil {
		return fmt.Errorf("execute config template error: %w", err)
	}

	return nil
}

// Validate validates the configuration.
func (c C) Validate() error {
	if !c.IsFileConfig {
		return nil
	}

	for _, f := range c.Before {
		if err := f.validate(); err != nil {
			return err
		}
	}

	for _, f := range c.After {
		if err := f.validate(); err != nil {
			return err
		}
	}

	if err := c.ChangelogOptions.validate(); err != nil {
		return err
	}

	for _, f := range c.Bump {
		if err := f.validate(); err != nil {
			return err
		}
	}

	return validateVersion(c.Version, _VersionWarningUpdate, _VersionCriticalUpdate)
}

// validateVersion validates the version.
func validateVersion(current, warning, critical version.V) error {
	if !critical.Empty() && current.LessThen(critical) {
		return fmt.Errorf(`%w: you use older version of config file. For update run "version generate --config-file"`, errConfig)
	}

	if !warning.Empty() && current.LessThen(warning) {
		return fmt.Errorf(`%w: you use older version of config file. For update run "version generate --config-file"`, ErrConfigWarn)
	}

	return nil
}

// Command is a command for run before or after git commit.
type Command struct {
	// Cmd is a command.
	// Example: ["go", "build", "-o", "bin/app", "."]
	Cmd []string `yaml:"cmd"`
	// Flag to send bumped version in format 1.2.3 to command. Optional.
	VersionFlag string `yaml:"versionFlag"`
	// BreakOnError is a flag that indicates that the command is stopped if an error occurs.
	BreakOnError bool `yaml:"breakOnError"`
	// RunInDry is a flag that indicates that the command is run in dry mode.
	RunInDry bool `yaml:"runInDry"`
}

// String returns a command string.
func (c Command) String() string {
	return strings.Join(c.Cmd, " ")
}

// Name returns a command name.
func (c Command) Name() string {
	return c.Cmd[0]
}

// Args returns a command args.
func (c Command) Args(v version.V) []string {
	var args []string

	if len(c.Cmd) > 1 {
		args = append(args, c.Cmd[1:]...)
	}

	if c.VersionFlag != "" {
		args = append(args, c.VersionFlag+"="+v.FormatString())
	}

	return args
}

// validate the command.
func (c Command) validate() error {
	if len(c.Cmd) == 0 {
		return fmt.Errorf("%w: empty command", errConfig)
	}

	return nil
}

// gitOptions is a git options.
type gitOptions struct {
	// AllowCommitDirty is a flag that indicates that the commit is allowed in a dirty repository.
	AllowCommitDirty bool `yaml:"commitDirty"`
	// AutoGenerateNextPatch is a flag that indicates that the next patch version is automatically generated if the version exists.
	AutoGenerateNextPatch bool `yaml:"autoNextPatch"`
	// AllowDowngrades is a flag that indicates that version downgrades are allowed with the --version flag.
	AllowDowngrades bool `yaml:"allowDowngrades"`
	// RemoteURL is a remote repository URL.
	RemoteURL string `yaml:"remoteUrl"`
}

// changelogOptions is a changelog options.
type changelogOptions struct {
	// Generate is a flag that indicates that the changelog is generated.
	Generate bool `yaml:"generate"`
	// FileName is a changelog file name.
	FileName File `yaml:"file"`
	// Title is a changelog title.
	Title string `yaml:"title"`
	// Issue href template.
	IssueURL string `yaml:"issueUrl"`
	// ShowAuthor is a flag that indicates that the author is shown in the changelog.
	ShowAuthor bool `yaml:"showAuthor"`
	// ShowBody is a flag that indicates that the body is shown in the changelog comment.
	ShowBody bool `yaml:"showBody"`
	// CommitTypes is a commit types for changelog.
	CommitTypes []CommitName `yaml:"commitTypes"`
}

// validate validates the changelog options.
func (c changelogOptions) validate() error {
	if !c.Generate {
		return nil
	}

	if c.FileName.empty() || c.FileName.IsAbs() {
		return fmt.Errorf(`%w: changelog file name is empty or absolute path`, errConfig)
	}

	for _, t := range c.CommitTypes {
		if t.Type == "" || t.Name == "" {
			return fmt.Errorf(`%w: commit type is empty`, errConfig)
		}
	}

	return nil
}

// BumpFile is a file for bump.
type BumpFile struct {
	// File path.
	File File `yaml:"file"`
	// RegExp for string for search version.
	// If not will be found version regexp in strings from start to end.
	RegExp []string `yaml:"regexp"`
	// Start string for search version. Default: 0 (start of file).
	Start int `yaml:"start"`
	// End string for search version. If 0 will be searched to end of file.
	End int `yaml:"end"`
}

// HasPositions returns true if the file has start and end positions.
func (f BumpFile) HasPositions() bool {
	return f.End != 0 && f.End >= f.Start
}

// IsPredefinedJSON returns true if the file is predefined (composer.json or package.json).
// Path contains the file path.
func (f BumpFile) IsPredefinedJSON() bool {
	n := filepath.Base(f.File.String())
	return n == "composer.json" || n == "package.json"
}

// validate BumpFile validates the file for bump.
func (f BumpFile) validate() error {
	// check if file exists
	if _, err := os.Stat(f.File.Path()); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(`%w: file %s does not exist`, errConfig, f.File)
		}

		return fmt.Errorf(`%w: file %s error: %w`, errConfig, f.File, err)
	}

	if f.IsPredefinedJSON() {
		return nil
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

	return nil
}
