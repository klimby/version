package config

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klimby/version/internal/file"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// C is a configuration file.
type C struct {
	// Version is a version of the application.
	Version version.V `yaml:"version"`
	// IsFileConfig is a flag that indicates that the configuration is from a file.
	IsFileConfig bool `yaml:"-"`
	// BackupChangedFiles is a flag that indicates that the changed files are backed up.
	Backup bool `yaml:"backupChanged"`
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
			CommitTypes: CommitNames(),
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

	tmpl, err := template.New("config").Parse(_configYamlTemplate)
	if err != nil {
		return fmt.Errorf("parse config template error: %w", err)
	}

	if err := tmpl.Execute(w, c); err != nil {
		return fmt.Errorf("execute config template error: %w", err)
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
