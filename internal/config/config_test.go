package config

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestC_Getters(t *testing.T) {
	c := C{
		Bump: []BumpFile{
			{
				File: "file",
			},
		},
		Before: []Command{
			{
				Cmd: []string{"cmd"},
			},
		},
		After: []Command{
			{
				Cmd: []string{"cmd2"},
			},
		},
		ChangelogOptions: changelogOptions{
			CommitTypes: []CommitName{
				{
					Type: "feat",
					Name: "Features",
				},
			},
		},
	}

	assert.Equal(t, "file", c.BumpFiles()[0].File.String(), "should be equal")
	assert.Equal(t, "cmd", c.CommandsBefore()[0].Cmd[0], "should be equal")
	assert.Equal(t, "cmd2", c.CommandsAfter()[0].Cmd[0], "should be equal")
	assert.Equal(t, "feat", c.CommitTypes()[0].Type, "should be equal")

}

func TestC_Generate(t *testing.T) {
	type fields struct {
		rw *__rwMock
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}),
			},
			assertion: assert.NoError,
		},
		{
			name: "write error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{writeErr: assert.AnError}),
			},
			assertion: assert.Error,
		},
		{
			name: "writer error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{writerErr: assert.AnError}),
			},
			assertion: assert.Error,
		},
		{
			name: "closer error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{closerErr: assert.AnError}),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.CfgFile, "config.yaml")
			viper.Set(key.Version, "1.2.4")
			c := C{
				Version: version.V("1.2.3"),
				Backup:  true,
				Before: []Command{
					{
						Cmd:          []string{"echo", "before commit"},
						VersionFlag:  "--version",
						BreakOnError: true,
						RunInDry:     true,
					},
				},
				After: []Command{
					{
						Cmd:          []string{"echo", "after commit"},
						VersionFlag:  "--version",
						BreakOnError: true,
						RunInDry:     true,
					},
				},
				GitOptions: gitOptions{
					RemoteURL: "https://github.com/klimby/version",
				},
				ChangelogOptions: changelogOptions{
					Generate: true,
					FileName: "CHANGELOG.md",
					Title:    "Changelog",
					IssueURL: "https://github.com/klimby/version/issues/",
					CommitTypes: []CommitName{
						{
							Type: "feat",
							Name: "Features",
						},
					},
				},
				Bump: []BumpFile{
					{
						File: "composer.json",
					},
					{
						File:   "README.md",
						RegExp: []string{`^Version:.+$`},
						Start:  0,
						End:    10,
					},
				},
				rw: tt.fields.rw,
			}

			tt.assertion(t, c.Generate())

			//fmt.Println(tt.fields.rw.rwc.buf.String())
		})
	}
}

func TestC_Validate(t *testing.T) {
	type fields struct {
		Version          version.V
		IsFileConfig     bool
		Backup           bool
		Before           []Command
		After            []Command
		GitOptions       gitOptions
		ChangelogOptions changelogOptions
		Bump             []BumpFile
		rw               *__rwMock
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "ok",
			fields:    fields{},
			assertion: assert.NoError,
		},
		{
			name: "ok full",
			fields: fields{
				IsFileConfig: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "invalid before",
			fields: fields{
				IsFileConfig: true,
				Before:       []Command{{}},
			},
			assertion: assert.Error,
		},
		{
			name: "invalid after",
			fields: fields{
				IsFileConfig: true,
				After:        []Command{{}},
			},
			assertion: assert.Error,
		},
		{
			name: "invalid changelog",
			fields: fields{
				IsFileConfig: true,
				ChangelogOptions: changelogOptions{
					Generate: true,
				},
			},
			assertion: assert.Error,
		},
		{
			name: "invalid bump",
			fields: fields{
				IsFileConfig: true,
				Bump:         []BumpFile{{}},
				rw:           __newRWMock(__rwMockArgs{exists: false}),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := C{
				Version:          tt.fields.Version,
				IsFileConfig:     tt.fields.IsFileConfig,
				Backup:           tt.fields.Backup,
				Before:           tt.fields.Before,
				After:            tt.fields.After,
				GitOptions:       tt.fields.GitOptions,
				ChangelogOptions: tt.fields.ChangelogOptions,
				Bump:             tt.fields.Bump,
				rw:               tt.fields.rw,
			}

			tt.assertion(t, c.Validate())
		})
	}
}

func Test_validateVersion(t *testing.T) {
	type args struct {
		current  version.V
		warning  version.V
		critical version.V
	}
	tests := []struct {
		name      string
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			args: args{
				current:  version.V("1.2.1"),
				warning:  version.V("1.2.2"),
				critical: version.V("1.2.3"),
			},
			assertion: assert.NoError,
		},
		{
			name: "warning",
			args: args{
				current:  version.V("1.2.3"),
				warning:  version.V("1.2.2"),
				critical: version.V("1.2.3"),
			},
			assertion: assert.Error,
		},
		{
			name: "critical",
			args: args{
				current:  version.V("1.2.4"),
				warning:  version.V("1.2.2"),
				critical: version.V("1.2.3"),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t, validateVersion(tt.args.current, tt.args.warning, tt.args.critical))
		})
	}
}

func TestCommand(t *testing.T) {
	c := Command{
		Cmd:         []string{"cmd", "foo"},
		VersionFlag: "flag",
	}

	assert.Equal(t, "cmd foo", c.String(), "should be equal")
	assert.Equal(t, "cmd", c.Name(), "should be equal")

	ver := version.V("1.2.3")
	assert.Equal(t, []string{"foo", "flag=1.2.3"}, c.Args(ver), "should be equal")

}

func TestCommand_validate(t *testing.T) {
	type fields struct {
		Cmd          []string
		VersionFlag  string
		BreakOnError bool
		RunInDry     bool
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "empty cmd",
			fields: fields{
				Cmd: []string{},
			},
			assertion: assert.Error,
		},
		{
			name: "ok",
			fields: fields{
				Cmd: []string{"cmd"},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Cmd:          tt.fields.Cmd,
				VersionFlag:  tt.fields.VersionFlag,
				BreakOnError: tt.fields.BreakOnError,
				RunInDry:     tt.fields.RunInDry,
			}

			tt.assertion(t, c.validate())
		})
	}
}

func Test_changelogOptions_validate(t *testing.T) {
	type fields struct {
		Generate    bool
		FileName    fsys.File
		Title       string
		IssueURL    string
		ShowAuthor  bool
		ShowBody    bool
		CommitTypes []CommitName
	}
	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "not generate",
			fields: fields{
				Generate: false,
			},
			assertion: assert.NoError,
		},
		{
			name: "empty file name",
			fields: fields{
				Generate: true,
				FileName: fsys.File(""),
			},
			assertion: assert.Error,
		},
		{
			name: "absolute file name",
			fields: fields{
				Generate: true,
				FileName: fsys.File("/file"),
			},
			assertion: assert.Error,
		},
		{
			name: "empty commit type",
			fields: fields{
				Generate:    true,
				FileName:    fsys.File("file"),
				CommitTypes: []CommitName{{}},
			},
			assertion: assert.Error,
		},
		{
			name: "ok",
			fields: fields{
				Generate:    true,
				FileName:    fsys.File("file"),
				CommitTypes: []CommitName{{Type: "type", Name: "name"}},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := changelogOptions{
				Generate:    tt.fields.Generate,
				FileName:    tt.fields.FileName,
				Title:       tt.fields.Title,
				IssueURL:    tt.fields.IssueURL,
				ShowAuthor:  tt.fields.ShowAuthor,
				ShowBody:    tt.fields.ShowBody,
				CommitTypes: tt.fields.CommitTypes,
			}

			tt.assertion(t, c.validate())
		})
	}
}

func TestBumpFile_HasPositions(t *testing.T) {
	type fields struct {
		File   fsys.File
		RegExp []string
		Start  int
		End    int
	}
	tests := []struct {
		name      string
		fields    fields
		want      bool
		assertion assert.BoolAssertionFunc
	}{
		{
			name: "has positions",
			fields: fields{
				Start: 1,
				End:   2,
			},
			want:      true,
			assertion: assert.True,
		},
		{
			name: "no positions",
			fields: fields{
				Start: 0,
				End:   0,
			},
			want:      false,
			assertion: assert.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := BumpFile{
				File:   tt.fields.File,
				RegExp: tt.fields.RegExp,
				Start:  tt.fields.Start,
				End:    tt.fields.End,
			}

			tt.assertion(t, f.HasPositions())
		})
	}
}

func TestBumpFile_validate(t *testing.T) {
	type fields struct {
		File   fsys.File
		RegExp []string
		Start  int
		End    int
	}

	type args struct {
		rw *__rwMock
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "file does not exist",
			fields: fields{
				File: fsys.File("file"),
			},
			args: args{
				rw: __newRWMock(__rwMockArgs{exists: false}),
			},
			assertion: assert.Error,
		},
		{
			name: "predefined json",
			fields: fields{
				File: fsys.File("composer.json"),
			},
			args: args{
				rw: __newRWMock(__rwMockArgs{exists: true}),
			},
			assertion: assert.NoError,
		},
		{
			name: "start > end",
			fields: fields{
				File:  fsys.File("file"),
				Start: 1,
				End:   0,
			},
			args: args{
				rw: __newRWMock(__rwMockArgs{exists: true}),
			},
			assertion: assert.Error,
		},
		{
			name: "regexp error",
			fields: fields{
				File:   fsys.File("file"),
				RegExp: []string{"["},
			},
			args: args{
				rw: __newRWMock(__rwMockArgs{exists: true}),
			},
			assertion: assert.Error,
		},
		{
			name: "regexp ok",
			fields: fields{
				File:   fsys.File("file"),
				RegExp: []string{"\\d+"},
			},
			args: args{
				rw: __newRWMock(__rwMockArgs{exists: true}),
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := BumpFile{
				File:   tt.fields.File,
				RegExp: tt.fields.RegExp,
				Start:  tt.fields.Start,
				End:    tt.fields.End,
			}

			tt.assertion(t, f.validate(tt.args.rw))
		})
	}
}

const __fakeConfigYaml = `
version: 1.2.3
backupChanged: true
before:
  - cmd: [ "echo", "before" ]
    versionFlag: "VERSION"
    breakOnError: true
    runInDry: false
after:
  - cmd: [ "echo", "after" ]
    versionFlag: "VERSION"
    breakOnError: true
    runInDry: false
git:
  commitDirty: true
  autoNextPatch: true
  allowDowngrades: true
  remoteUrl: https://github.com/klimby/version
changelog:
  generate: true
  file: CHANGELOG.md
  title: "Changelog"
  issueUrl: https://github.com/klimby/version/issues/
  showAuthor: false
  showBody: true
  commitTypes:
    - type: "feat"
      name: "Features"
bump:
  - file: README.md
    start: 0
    end: 5
    regexp:
      - ^!\[Version.*$	
`

func Test_newConfig(t *testing.T) {
	type args struct {
		rw *__rwMock
	}

	data := []byte(__fakeConfigYaml)

	targetFile := C{
		Version: version.V("1.2.3"),
		Backup:  true,
		Before: []Command{
			{
				Cmd:          []string{"echo", "before"},
				VersionFlag:  "VERSION",
				BreakOnError: true,
				RunInDry:     false,
			},
		},
		After: []Command{
			{
				Cmd:          []string{"echo", "after"},
				VersionFlag:  "VERSION",
				BreakOnError: true,
				RunInDry:     false,
			},
		},
		GitOptions: gitOptions{
			AllowCommitDirty:      true,
			AutoGenerateNextPatch: true,
			AllowDowngrades:       true,
			RemoteURL:             "https://github.com/klimby/version",
		},
		ChangelogOptions: changelogOptions{
			Generate:   true,
			FileName:   fsys.File("CHANGELOG.md"),
			Title:      "Changelog",
			IssueURL:   "https://github.com/klimby/version/issues/",
			ShowAuthor: false,
			ShowBody:   true,
			CommitTypes: []CommitName{
				{
					Type: "feat",
					Name: "Features",
				},
			},
		},
		Bump: []BumpFile{
			{
				File:   "README.md",
				Start:  0,
				End:    5,
				RegExp: []string{`^!\[Version.*$`},
			},
		},
	}

	tests := []struct {
		name         string
		args         args
		target       *C
		isFileConfig assert.BoolAssertionFunc
		assertion    assert.ErrorAssertionFunc
	}{
		{
			name: "load from file",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data}),
			},
			target:       &targetFile,
			isFileConfig: assert.True,
			assertion:    assert.NoError,
		},
		{
			name: "read error",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data, readErr: assert.AnError}),
			},
			assertion: assert.Error,
		},
		{
			name: "close error",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data, closerErr: assert.AnError}),
			},
			target:       &targetFile,
			isFileConfig: assert.True,
			assertion:    assert.Error,
		},
		{
			name: "ErrNotExist error",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: data, readErr: fs.ErrNotExist}),
			},
			isFileConfig: assert.False,
			assertion:    assert.NoError,
		},
		{
			name: "decoder error",
			args: args{
				rw: __newRWMock(__rwMockArgs{data: nil}),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig(tt.args.rw)

			tt.assertion(t, err, fmt.Sprintf("newConfig()"))

			if err != nil {
				return
			}

			if tt.target != nil {
				assert.Equal(t, tt.target.Version.String(), got.Version.String(), "version should be equal")
				assert.Equal(t, tt.target.Backup, got.Backup, "backup should be equal")

				assert.Equal(t, len(tt.target.Before), len(got.Before), "before should be equal")

				if len(tt.target.Before) == len(got.Before) && len(tt.target.Before) > 0 {
					assert.Equal(t, tt.target.Before[0].Cmd, got.Before[0].Cmd, "before Cmd should be equal")
					assert.Equal(t, tt.target.Before[0].VersionFlag, got.Before[0].VersionFlag, "before VersionFlag should be equal")
					assert.Equal(t, tt.target.Before[0].BreakOnError, got.Before[0].BreakOnError, "before BreakOnError should be equal")
					assert.Equal(t, tt.target.Before[0].RunInDry, got.Before[0].RunInDry, "before RunInDry should be equal")
				}

				assert.Equal(t, len(tt.target.After), len(got.After), "after should be equal")

				if len(tt.target.After) == len(got.After) && len(tt.target.After) > 0 {
					assert.Equal(t, tt.target.After[0].Cmd, got.After[0].Cmd, "after Cmd should be equal")
					assert.Equal(t, tt.target.After[0].VersionFlag, got.After[0].VersionFlag, "after VersionFlag should be equal")
					assert.Equal(t, tt.target.After[0].BreakOnError, got.After[0].BreakOnError, "after BreakOnError should be equal")
					assert.Equal(t, tt.target.After[0].RunInDry, got.After[0].RunInDry, "after RunInDry should be equal")
				}

				assert.Equal(t, tt.target.GitOptions.AllowCommitDirty, got.GitOptions.AllowCommitDirty, "git AllowCommitDirty should be equal")
				assert.Equal(t, tt.target.GitOptions.AutoGenerateNextPatch, got.GitOptions.AutoGenerateNextPatch, "git AutoGenerateNextPatch should be equal")
				assert.Equal(t, tt.target.GitOptions.AllowDowngrades, got.GitOptions.AllowDowngrades, "git AllowDowngrades should be equal")
				assert.Equal(t, tt.target.GitOptions.RemoteURL, got.GitOptions.RemoteURL, "git RemoteURL should be equal")

				assert.Equal(t, tt.target.ChangelogOptions.Generate, got.ChangelogOptions.Generate, "changelog Generate should be equal")
				assert.Equal(t, tt.target.ChangelogOptions.FileName.String(), got.ChangelogOptions.FileName.String(), "changelog FileName should be equal")
				assert.Equal(t, tt.target.ChangelogOptions.Title, got.ChangelogOptions.Title, "changelog Title should be equal")
				assert.Equal(t, tt.target.ChangelogOptions.IssueURL, got.ChangelogOptions.IssueURL, "changelog IssueURL should be equal")
				assert.Equal(t, tt.target.ChangelogOptions.ShowAuthor, got.ChangelogOptions.ShowAuthor, "changelog ShowAuthor should be equal")
				assert.Equal(t, tt.target.ChangelogOptions.ShowBody, got.ChangelogOptions.ShowBody, "changelog ShowBody should be equal")

				assert.Equal(t, len(tt.target.ChangelogOptions.CommitTypes), len(got.ChangelogOptions.CommitTypes), "changelog CommitTypes[0].Type should be equal")

				if len(tt.target.ChangelogOptions.CommitTypes) == len(got.ChangelogOptions.CommitTypes) && len(tt.target.ChangelogOptions.CommitTypes) > 0 {
					assert.Equal(t, tt.target.ChangelogOptions.CommitTypes[0].Type, got.ChangelogOptions.CommitTypes[0].Type, "changelog CommitTypes[0].Type should be equal")
					assert.Equal(t, tt.target.ChangelogOptions.CommitTypes[0].Name, got.ChangelogOptions.CommitTypes[0].Name, "changelog CommitTypes[0].Name should be equal")
				}

				assert.Equal(t, len(tt.target.Bump), len(got.Bump), "bump should be equal")

				if len(tt.target.Bump) == len(got.Bump) && len(tt.target.Bump) > 0 {
					assert.Equal(t, tt.target.Bump[0].File.String(), got.Bump[0].File.String(), "bump File should be equal")
					assert.Equal(t, tt.target.Bump[0].Start, got.Bump[0].Start, "bump Start should be equal")
					assert.Equal(t, tt.target.Bump[0].End, got.Bump[0].End, "bump End should be equal")
					assert.Equal(t, tt.target.Bump[0].RegExp[0], got.Bump[0].RegExp[0], "bump RegExp should be equal")
				}
			}

			tt.isFileConfig(t, got.IsFileConfig, "IsFileConfig wrong")
		})
	}
}
