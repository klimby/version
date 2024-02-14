package config

import (
	"testing"

	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type args struct {
		rw *__rwMock
	}

	data := []byte(__fakeConfigYaml)

	tests := []struct {
		name string
		args args

		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data}),
			},
			wantErr: assert.NoError,
		},
		{
			name: "load error",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data, readErr: assert.AnError}),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(func(options *LoadArgs) {
				options.RW = tt.args.rw
			})

			if !tt.wantErr(t, err, "Load() config") {
				return
			}

			if err != nil {
				return
			}

			assert.True(t, viper.GetBool(key.AllowCommitDirty), "AllowCommitDirty")
			assert.True(t, viper.GetBool(key.AutoGenerateNextPatch), "AutoGenerateNextPatch")
			assert.True(t, viper.GetBool(key.AllowDowngrades), "AllowDowngrades")
			assert.True(t, viper.GetBool(key.Backup), "Backup")

			assert.Equal(t, got.ChangelogOptions.Generate, viper.GetBool(key.GenerateChangelog), "GenerateChangelog")
			assert.Equal(t, got.ChangelogOptions.FileName.String(), viper.GetString(key.ChangelogFileName), "ChangelogFileName")
			assert.Equal(t, got.ChangelogOptions.Title, viper.GetString(key.ChangelogTitle), "ChangelogFileName")
			assert.Equal(t, got.ChangelogOptions.ShowAuthor, viper.GetBool(key.ChangelogShowAuthor), "ChangelogShowAuthor")
			assert.Equal(t, got.ChangelogOptions.ShowBody, viper.GetBool(key.ChangelogShowBody), "ChangelogShowBody")

			assert.Equal(t, got.ChangelogOptions.IssueURL, viper.GetString(key.ChangelogIssueURL), "ChangelogIssueURL")
			assert.Equal(t, got.GitOptions.RemoteURL, viper.GetString(key.RemoteURL), "RemoteURL")

		})
	}
}

func TestSetURLFromGit(t *testing.T) {
	type args struct {
		u string
	}

	type expected struct {
		remoteURL         string
		changelogIssueURL string
	}

	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "ok",
			args: args{
				u: "https://github.com/foo/bar",
			},
			expected: expected{
				remoteURL:         "https://github.com/foo/bar",
				changelogIssueURL: "https://github.com/foo/bar/issues/",
			},
		},
		{
			name: "invalid url",
			args: args{
				u: "https://github.com:foo/bar",
			},
			expected: expected{
				remoteURL:         "https://github.com:foo/bar",
				changelogIssueURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.RemoteURL, "")
			viper.Set(key.ChangelogIssueURL, "")

			SetURLFromGit(tt.args.u)

			assert.Equal(t, tt.expected.remoteURL, viper.GetString(key.RemoteURL))
			assert.Equal(t, tt.expected.changelogIssueURL, viper.GetString(key.ChangelogIssueURL))
		})
	}
}

func TestSetForce(t *testing.T) {
	tests := []struct {
		name      string
		force     bool
		assertion assert.BoolAssertionFunc
	}{
		{
			name:      "force",
			force:     true,
			assertion: assert.True,
		},
		{
			name:      "not force",
			force:     false,
			assertion: assert.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.AllowCommitDirty, false)
			viper.Set(key.AutoGenerateNextPatch, false)
			viper.Set(key.AllowDowngrades, false)

			viper.Set(key.Force, tt.force)
			SetForce()

			tt.assertion(t, viper.GetBool(key.AllowCommitDirty))
			tt.assertion(t, viper.GetBool(key.AutoGenerateNextPatch))
			tt.assertion(t, viper.GetBool(key.AllowDowngrades))
		})
	}
}

func TestInit(t *testing.T) {
	o := Options{
		WorkDir:               "foo",
		Version:               "bar",
		AllowCommitDirty:      true,
		AutoGenerateNextPatch: true,
		AllowDowngrades:       true,
		GenerateChangelog:     true,
		ChangelogFileName:     "c.md",
		ChangelogTitle:        "cc",
		ChangelogIssueHref:    "asd",
		ChangelogShowAuthor:   true,
		ChangelogShowBody:     true,
		Silent:                true,
		DryRun:                true,
		Backup:                true,
		Force:                 true,
		Verbose:               true,
		ConfigFile:            "c.yaml",

		TestingSkipDIInit: true,
	}

	Init(func(options *Options) {
		options.WorkDir = o.WorkDir
		options.Version = o.Version
		options.AllowCommitDirty = o.AllowCommitDirty
		options.AutoGenerateNextPatch = o.AutoGenerateNextPatch
		options.AllowDowngrades = o.AllowDowngrades
		options.GenerateChangelog = o.GenerateChangelog
		options.ChangelogFileName = o.ChangelogFileName
		options.ChangelogTitle = o.ChangelogTitle
		options.ChangelogIssueHref = o.ChangelogIssueHref
		options.ChangelogShowAuthor = o.ChangelogShowAuthor
		options.ChangelogShowBody = o.ChangelogShowBody
		options.Silent = o.Silent
		options.DryRun = o.DryRun
		options.Backup = o.Backup
		options.Force = o.Force
		options.Verbose = o.Verbose
		options.ConfigFile = o.ConfigFile

		options.TestingSkipDIInit = o.TestingSkipDIInit
	})

	assert.Equal(t, _AppName, viper.GetString(key.AppName))
	assert.Equal(t, o.Version, viper.GetString(key.Version))
	assert.Equal(t, o.WorkDir, viper.GetString(key.WorkDir))

	assert.Equal(t, o.AllowCommitDirty, viper.GetBool(key.AllowCommitDirty))
	assert.Equal(t, o.AutoGenerateNextPatch, viper.GetBool(key.AutoGenerateNextPatch))
	assert.Equal(t, o.AllowDowngrades, viper.GetBool(key.AllowDowngrades))

	assert.Equal(t, o.GenerateChangelog, viper.GetBool(key.GenerateChangelog))
	assert.Equal(t, o.ChangelogFileName, viper.GetString(key.ChangelogFileName))
	assert.Equal(t, o.ChangelogTitle, viper.GetString(key.ChangelogTitle))
	assert.Equal(t, o.ChangelogIssueHref, viper.GetString(key.ChangelogIssueURL))
	assert.Equal(t, o.ChangelogShowAuthor, viper.GetBool(key.ChangelogShowAuthor))
	assert.Equal(t, o.ChangelogShowBody, viper.GetBool(key.ChangelogShowBody))

	assert.Equal(t, o.Silent, viper.GetBool(key.Silent))
	assert.Equal(t, o.DryRun, viper.GetBool(key.DryRun))
	assert.Equal(t, o.Backup, viper.GetBool(key.Backup))
	assert.Equal(t, o.Force, viper.GetBool(key.Force))
	assert.Equal(t, o.Verbose, viper.GetBool(key.Verbose))
	assert.Equal(t, o.ConfigFile, viper.GetString(key.CfgFile))

	assert.Equal(t, o.TestingSkipDIInit, viper.GetBool(key.TestingSkipDIInit))

}
