package next

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/changelog"
	"github.com/klimby/version/internal/service/git"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAction_Run(t *testing.T) {
	const nextVersion = version.V("1.2.4")

	type fields struct {
		actionType   ActionType
		repo         *__actionRepoMock
		changelogGen *__actionChGenMock
		cfg          *__actionCfgMock
		bump         *__actionBumpMock
		cmd          *__actionCmdMock
		prepare      bool
	}

	type repoMockArgs struct {
		isCleanErr, nextVersionErr, checkDowngradeErr, addModifiedErr, commitTagErr error
	}

	repoMock := func(a repoMockArgs) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("IsClean").Return(true, a.isCleanErr)
		r.On("NextVersion", git.NextPatch, mock.Anything).Return(nextVersion, false, a.nextVersionErr)
		r.On("CheckDowngrade", nextVersion).Return(a.checkDowngradeErr)
		r.On("AddModified").Return(a.addModifiedErr)
		r.On("CommitTag", nextVersion).Return(a.commitTagErr)
		return r
	}

	bumpMock := func() *__actionBumpMock {
		b := &__actionBumpMock{}
		b.On("Apply", mock.Anything, nextVersion)
		return b
	}

	changelogMock := func(e error) *__actionChGenMock {
		c := &__actionChGenMock{}
		c.On("Add", nextVersion).Return(e)
		return c
	}

	cfgMock := func() *__actionCfgMock {
		c := &__actionCfgMock{}
		c.On("BumpFiles").Return([]config.BumpFile{})
		c.On("CommandsBefore").Return([]config.Command{
			{
				Cmd:          []string{"echo", "test"},
				VersionFlag:  "version",
				BreakOnError: true,
			},
		})
		c.On("CommandsAfter").Return([]config.Command{
			{
				Cmd:          []string{"echo", "test-after"},
				VersionFlag:  "version",
				BreakOnError: true,
			},
		})
		return c
	}

	cmdMock := func(beforeErr, afterErr error) *__actionCmdMock {
		c := &__actionCmdMock{}
		c.On("Run", "echo", "test", "version=1.2.4").Return(beforeErr)
		c.On("Run", "echo", "test-after", "version=1.2.4").Return(afterErr)
		return c
	}

	type wantCalls struct {
		prepare bool
		apply   bool
	}

	tests := []struct {
		name      string
		fields    fields
		wantCalls wantCalls
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal run",
			fields: fields{
				actionType:   ActionPatch,
				repo:         repoMock(repoMockArgs{}),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				bump:         bumpMock(),
				cmd:          cmdMock(nil, nil),
			},
			wantCalls: wantCalls{
				prepare: true,
				apply:   true,
			},
			assertion: assert.NoError,
		},
		{
			name: "prepare run",
			fields: fields{
				actionType:   ActionPatch,
				repo:         repoMock(repoMockArgs{}),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				bump:         bumpMock(),
				cmd:          cmdMock(nil, nil),
				prepare:      true,
			},
			wantCalls: wantCalls{
				prepare: true,
				apply:   false,
			},
			assertion: assert.NoError,
		},
		{
			name: "prepare error",
			fields: fields{
				actionType:   ActionPatch,
				repo:         repoMock(repoMockArgs{}),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				bump:         bumpMock(),
				cmd:          cmdMock(assert.AnError, nil),
			},
			wantCalls: wantCalls{
				prepare: true,
				apply:   false,
			},
			assertion: assert.Error,
		},
		{
			name: "apply error",
			fields: fields{
				actionType:   ActionPatch,
				repo:         repoMock(repoMockArgs{}),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				bump:         bumpMock(),
				cmd:          cmdMock(nil, assert.AnError),
			},
			wantCalls: wantCalls{
				prepare: true,
				apply:   true,
			},
			assertion: assert.Error,
		},
		{
			name: "validate error",
			fields: fields{
				actionType:   ActionUnknown,
				repo:         repoMock(repoMockArgs{}),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				bump:         bumpMock(),
				cmd:          cmdMock(nil, nil),
			},
			wantCalls: wantCalls{
				prepare: false,
				apply:   false,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.GenerateChangelog, true)
			viper.Set(key.Prepare, tt.fields.prepare)

			a := New(func(args *Args) {
				args.ActionType = tt.fields.actionType
				args.Repo = tt.fields.repo
				args.ChangelogGen = tt.fields.changelogGen
				args.Cfg = tt.fields.cfg
				args.Bump = tt.fields.bump
				args.Cmd = tt.fields.cmd
			})

			tt.assertion(t, a.Run(), "Run() error = %v, wantErr %v", tt.assertion, false)

			if tt.wantCalls.prepare {
				// Last command in prepare.
				tt.fields.cmd.AssertCalled(t, "Run", "echo", "test", "version=1.2.4")
			} else {
				tt.fields.cmd.AssertNotCalled(t, "Run", "echo", "test", "version=1.2.4")
			}

			if tt.wantCalls.apply {
				tt.fields.cmd.AssertCalled(t, "Run", "echo", "test-after", "version=1.2.4")
			} else {
				tt.fields.cmd.AssertNotCalled(t, "Run", "echo", "test-after", "version=1.2.4")
			}
		})

	}
}

func TestAction_prepare(t *testing.T) {
	const nextVersion = version.V("1.2.4")

	type fields struct {
		repo *__actionRepoMock
		cfg  *__actionCfgMock
		bump *__actionBumpMock
		cmd  *__actionCmdMock
	}

	type repoMockArgs struct {
		isCleanErr, nextVersionErr, checkDowngradeErr error
	}

	repoMock := func(a repoMockArgs) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("IsClean").Return(true, a.isCleanErr)
		r.On("NextVersion", git.NextPatch, mock.Anything).Return(nextVersion, false, a.nextVersionErr)
		r.On("CheckDowngrade", nextVersion).Return(a.checkDowngradeErr)
		return r
	}

	bumpMock := func() *__actionBumpMock {
		b := &__actionBumpMock{}
		b.On("Apply", mock.Anything, nextVersion)
		return b
	}

	cfgMock := func() *__actionCfgMock {
		c := &__actionCfgMock{}
		c.On("BumpFiles").Return([]config.BumpFile{})
		c.On("CommandsBefore").Return([]config.Command{
			{
				Cmd:          []string{"echo", "test"},
				VersionFlag:  "version",
				BreakOnError: true,
			},
		})
		return c
	}

	cmdMock := func(e error) *__actionCmdMock {
		c := &__actionCmdMock{}
		c.On("Run", "echo", "test", "version=1.2.4").Return(e)
		return c
	}

	type wantCalls struct {
		checkClean     bool
		nextVersion    bool
		checkDowngrade bool
		bump           bool
		runCommands    bool
	}

	tests := []struct {
		name      string
		fields    fields
		wantCalls wantCalls
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal run",
			fields: fields{
				repo: repoMock(repoMockArgs{}),
				cfg:  cfgMock(),
				bump: bumpMock(),
				cmd:  cmdMock(nil),
			},
			wantCalls: wantCalls{
				checkClean:     true,
				nextVersion:    true,
				checkDowngrade: true,
				bump:           true,
				runCommands:    true,
			},
			assertion: assert.NoError,
		},
		{
			name: "checkClean error",
			fields: fields{
				repo: repoMock(repoMockArgs{isCleanErr: assert.AnError}),
				cfg:  cfgMock(),
				bump: bumpMock(),
				cmd:  cmdMock(nil),
			},
			wantCalls: wantCalls{
				checkClean:     true,
				nextVersion:    false,
				checkDowngrade: false,
				bump:           false,
				runCommands:    false,
			},
			assertion: assert.Error,
		},
		{
			name: "nextVersion error",
			fields: fields{
				repo: repoMock(repoMockArgs{nextVersionErr: assert.AnError}),
				cfg:  cfgMock(),
				bump: bumpMock(),
				cmd:  cmdMock(nil),
			},
			wantCalls: wantCalls{
				checkClean:     true,
				nextVersion:    true,
				checkDowngrade: false,
				bump:           false,
				runCommands:    false,
			},
			assertion: assert.Error,
		},
		{
			name: "checkDowngrade error",
			fields: fields{
				repo: repoMock(repoMockArgs{checkDowngradeErr: assert.AnError}),
				cfg:  cfgMock(),
				bump: bumpMock(),
				cmd:  cmdMock(nil),
			},
			wantCalls: wantCalls{
				checkClean:     true,
				nextVersion:    true,
				checkDowngrade: true,
				bump:           false,
				runCommands:    false,
			},
			assertion: assert.Error,
		},
		{
			name: "runCommands error",
			fields: fields{
				repo: repoMock(repoMockArgs{}),
				cfg:  cfgMock(),
				bump: bumpMock(),
				cmd:  cmdMock(assert.AnError),
			},
			wantCalls: wantCalls{
				checkClean:     true,
				nextVersion:    true,
				checkDowngrade: true,
				bump:           true,
				runCommands:    true,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.DryRun, false)
			viper.Set(key.GenerateChangelog, true)
			viper.Set(key.AllowDowngrades, false)
			viper.Set(key.AllowCommitDirty, false)

			a := New(func(args *Args) {
				args.ActionType = ActionPatch
				args.Repo = tt.fields.repo
				args.Cfg = tt.fields.cfg
				args.Bump = tt.fields.bump
				args.Cmd = tt.fields.cmd
			})

			nv, err := a.prepare()
			tt.assertion(t, err, "prepare() error = %v, wantErr %v", err, tt.assertion)

			if err == nil {
				assert.Equal(t, nextVersion, nv, "prepare() next version")
			}

			if tt.wantCalls.checkClean {
				tt.fields.repo.AssertCalled(t, "IsClean")
			} else {
				tt.fields.repo.AssertNotCalled(t, "IsClean")
			}

			if tt.wantCalls.nextVersion {
				tt.fields.repo.AssertCalled(t, "NextVersion", git.NextPatch, mock.Anything)
			} else {
				tt.fields.repo.AssertNotCalled(t, "NextVersion")
			}

			if tt.wantCalls.checkDowngrade {
				tt.fields.repo.AssertCalled(t, "CheckDowngrade", nextVersion)
			} else {
				tt.fields.repo.AssertNotCalled(t, "CheckDowngrade")
			}

			if tt.wantCalls.bump {
				tt.fields.bump.AssertCalled(t, "Apply", mock.Anything, nextVersion)
			} else {
				tt.fields.bump.AssertNotCalled(t, "Apply")
			}

			if tt.wantCalls.runCommands {
				tt.fields.cmd.AssertCalled(t, "Run", "echo", "test", "version=1.2.4")
			} else {
				tt.fields.cmd.AssertNotCalled(t, "Run")
			}
		})

	}
}

func TestAction_apply(t *testing.T) {
	const nextV = version.V("1.2.4")

	type fields struct {
		repo         *__actionRepoMock
		changelogGen *__actionChGenMock
		cfg          *__actionCfgMock
		cmd          *__actionCmdMock
	}

	cfgMock := func() *__actionCfgMock {
		c := &__actionCfgMock{}
		c.On("CommandsAfter").Return([]config.Command{
			{
				Cmd:          []string{"echo", "test"},
				VersionFlag:  "version",
				BreakOnError: true,
			},
		})

		return c
	}

	repoMock := func(addModifiedErr, commitTagErr error) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("AddModified").Return(addModifiedErr)
		r.On("CommitTag", nextV).Return(commitTagErr)
		return r
	}

	changelogMock := func(e error) *__actionChGenMock {
		c := &__actionChGenMock{}
		c.On("Add", nextV).Return(e)
		return c
	}

	cmdMock := func(e error) *__actionCmdMock {
		c := &__actionCmdMock{}
		c.On("Run", "echo", "test", "version=1.2.4").Return(e)
		return c
	}

	type wantCalls struct {
		writeChangelog  bool
		repoAddModified bool
		repoCommitTag   bool
		runCommands     bool
	}

	tests := []struct {
		name      string
		fields    fields
		wantCalls wantCalls
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "normal run",
			fields: fields{
				repo:         repoMock(nil, nil),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				cmd:          cmdMock(nil),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: true,
				repoCommitTag:   true,
				runCommands:     true,
			},
			assertion: assert.NoError,
		},
		{
			name: "writeChangelog error",
			fields: fields{
				repo:         repoMock(nil, nil),
				changelogGen: changelogMock(assert.AnError),
				cfg:          cfgMock(),
				cmd:          cmdMock(nil),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: false,
				repoCommitTag:   false,
				runCommands:     false,
			},
			assertion: assert.Error,
		},
		{
			name: "writeChangelog warning",
			fields: fields{
				repo:         repoMock(nil, nil),
				changelogGen: changelogMock(changelog.ErrWarning),
				cfg:          cfgMock(),
				cmd:          cmdMock(nil),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: true,
				repoCommitTag:   true,
				runCommands:     true,
			},
			assertion: assert.NoError,
		},
		{
			name: "repo.AddModified error",
			fields: fields{
				repo:         repoMock(assert.AnError, nil),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				cmd:          cmdMock(nil),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: true,
				repoCommitTag:   true,
				runCommands:     true,
			},
			assertion: assert.NoError,
		},
		{
			name: "repo.CommitTag error",
			fields: fields{
				repo:         repoMock(nil, assert.AnError),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				cmd:          cmdMock(nil),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: true,
				repoCommitTag:   true,
				runCommands:     false,
			},
			assertion: assert.Error,
		},
		{
			name: "runCommands error",
			fields: fields{
				repo:         repoMock(nil, nil),
				changelogGen: changelogMock(nil),
				cfg:          cfgMock(),
				cmd:          cmdMock(assert.AnError),
			},
			wantCalls: wantCalls{
				writeChangelog:  true,
				repoAddModified: true,
				repoCommitTag:   true,
				runCommands:     true,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.GenerateChangelog, true)

			a := New(func(args *Args) {
				args.ActionType = ActionPatch
				args.Repo = tt.fields.repo
				args.ChangelogGen = tt.fields.changelogGen
				args.Cfg = tt.fields.cfg
				args.Cmd = tt.fields.cmd
			})

			tt.assertion(t, a.apply(nextV), "apply() error = %v, wantErr %v", tt.assertion, false)

			if tt.wantCalls.writeChangelog {
				tt.fields.changelogGen.AssertCalled(t, "Add", nextV)
			} else {
				tt.fields.changelogGen.AssertNotCalled(t, "Add")
			}

			if tt.wantCalls.repoAddModified {
				tt.fields.repo.AssertCalled(t, "AddModified")
			} else {
				tt.fields.repo.AssertNotCalled(t, "AddModified")
			}

			if tt.wantCalls.repoCommitTag {
				tt.fields.repo.AssertCalled(t, "CommitTag", nextV)
			} else {
				tt.fields.repo.AssertNotCalled(t, "CommitTag")
			}

			if tt.wantCalls.runCommands {
				tt.fields.cmd.AssertCalled(t, "Run", "echo", "test", "version=1.2.4")
			} else {
				tt.fields.cmd.AssertNotCalled(t, "Run")
			}
		})

	}
}

func TestAction_validate(t *testing.T) {
	type fields struct {
		actionType    ActionType
		customVersion version.V
		repo          *__actionRepoMock
		changelogGen  *__actionChGenMock
		cfg           *__actionCfgMock
		bump          *__actionBumpMock
		cmd           *__actionCmdMock
	}

	type wantErr struct {
		want     bool
		contains string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr wantErr
	}{
		{
			name: "unknown type",
			fields: fields{
				actionType: ActionUnknown,
			},
			wantErr: wantErr{
				want:     true,
				contains: "action type is unknown",
			},
		},
		{
			name: "custom action invalid version",
			fields: fields{
				actionType:    ActionCustom,
				customVersion: version.V(""),
			},
			wantErr: wantErr{
				want:     true,
				contains: "custom version",
			},
		},
		{
			name: "nil repo",
			fields: fields{
				actionType:   ActionPatch,
				repo:         nil,
				changelogGen: &__actionChGenMock{},
				cfg:          &__actionCfgMock{},
				bump:         &__actionBumpMock{},
				cmd:          &__actionCmdMock{},
			},
			wantErr: wantErr{
				want:     true,
				contains: "repo is nil",
			},
		},
		{
			name: "nil changelog generator",
			fields: fields{
				actionType:   ActionPatch,
				repo:         &__actionRepoMock{},
				changelogGen: nil,
				cfg:          &__actionCfgMock{},
				bump:         &__actionBumpMock{},
				cmd:          &__actionCmdMock{},
			},
			wantErr: wantErr{
				want:     true,
				contains: "changelog generator is nil",
			},
		},
		{
			name: "nil config",
			fields: fields{
				actionType:   ActionPatch,
				repo:         &__actionRepoMock{},
				changelogGen: &__actionChGenMock{},
				cfg:          nil,
				bump:         &__actionBumpMock{},
				cmd:          &__actionCmdMock{},
			},
			wantErr: wantErr{
				want:     true,
				contains: "config is nil",
			},
		},
		{
			name: "nil bump",
			fields: fields{
				actionType:   ActionPatch,
				repo:         &__actionRepoMock{},
				changelogGen: &__actionChGenMock{},
				cfg:          &__actionCfgMock{},
				bump:         nil,
				cmd:          &__actionCmdMock{},
			},
			wantErr: wantErr{
				want:     true,
				contains: "bump is nil",
			},
		},
		{
			name: "cmd bump",
			fields: fields{
				actionType:   ActionPatch,
				repo:         &__actionRepoMock{},
				changelogGen: &__actionChGenMock{},
				cfg:          &__actionCfgMock{},
				bump:         &__actionBumpMock{},
				cmd:          nil,
			},
			wantErr: wantErr{
				want:     true,
				contains: "cmd is nil",
			},
		},
		{
			name: "valid",
			fields: fields{
				actionType:   ActionPatch,
				repo:         &__actionRepoMock{},
				changelogGen: &__actionChGenMock{},
				cfg:          &__actionCfgMock{},
				bump:         &__actionBumpMock{},
				cmd:          &__actionCmdMock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(func(args *Args) {
				args.ActionType = tt.fields.actionType
				args.Version = tt.fields.customVersion

				if tt.fields.repo != nil {
					args.Repo = tt.fields.repo
				}

				if tt.fields.changelogGen != nil {
					args.ChangelogGen = tt.fields.changelogGen
				}

				if tt.fields.cfg != nil {
					args.Cfg = tt.fields.cfg
				}

				if tt.fields.bump != nil {
					args.Bump = tt.fields.bump
				}

				if tt.fields.cmd != nil {
					args.Cmd = tt.fields.cmd
				} else {
					args.Cmd = nil // because of default value
				}
			})

			err := a.validate()
			if (err != nil) != tt.wantErr.want {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr.want)
			}

			if err != nil {
				assert.Contains(t, err.Error(), tt.wantErr.contains)
			}
		})

	}
}

func TestAction_checkClean(t *testing.T) {
	type fields struct {
		repo             *__actionRepoMock
		allowCommitDirty bool
	}

	repoFn := func(clean bool, e error) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("IsClean").Return(clean, e)
		return r
	}

	type wantErr struct {
		want        bool
		contains    string
		containsNot string
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr wantErr
	}{
		{
			name: "clean",
			fields: fields{
				repo: repoFn(true, nil),
			},
		},
		{
			name: "clean error",
			fields: fields{
				repo: repoFn(false, assert.AnError),
			},
			wantErr: wantErr{
				want:        true,
				containsNot: "repository is not clean",
			},
		},
		{
			name: "clean false allow dirty",
			fields: fields{
				repo:             repoFn(false, nil),
				allowCommitDirty: true,
			},
		},
		{
			name: "clean false not allow dirty",
			fields: fields{
				repo: repoFn(false, nil),
			},
			wantErr: wantErr{
				want:     true,
				contains: "repository is not clean",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.AllowCommitDirty, tt.fields.allowCommitDirty)

			a := &Action{
				repo: tt.fields.repo,
			}

			err := a.checkClean()
			if (err != nil) != tt.wantErr.want {
				t.Errorf("checkClean() error = %v, wantErr %v", err, tt.wantErr.want)
			}

			if err != nil {
				if tt.wantErr.contains != "" {
					assert.Contains(t, err.Error(), tt.wantErr.contains)
				}

				if tt.wantErr.containsNot != "" {
					assert.NotContains(t, err.Error(), tt.wantErr.containsNot)
				}
			}

			tt.fields.repo.AssertCalled(t, "IsClean")
		})

	}
}

func TestAction_nextVersion(t *testing.T) {
	const (
		currentVersion = version.V("1.2.3")
		nextVersion    = version.V("1.2.4")
	)

	type fields struct {
		repo                  *__actionRepoMock
		autoGenerateNextPatch bool
	}

	repoFn := func(exists bool, e error) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("NextVersion", git.NextPatch, currentVersion).Return(nextVersion, exists, e)
		return r
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "next version",
			fields: fields{
				repo: repoFn(false, nil),
			},
			assertion: assert.NoError,
		},
		{
			name: "next version error",
			fields: fields{
				repo: repoFn(false, assert.AnError),
			},
			assertion: assert.Error,
		},
		{
			name: "next version exists allow next patch",
			fields: fields{
				repo:                  repoFn(true, nil),
				autoGenerateNextPatch: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "next version exists not allow next patch",
			fields: fields{
				repo: repoFn(true, nil),
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.AutoGenerateNextPatch, tt.fields.autoGenerateNextPatch)

			a := &Action{
				actionType:    ActionPatch,
				repo:          tt.fields.repo,
				customVersion: currentVersion,
			}

			_, err := a.nextVersion()
			tt.assertion(t, err, "nextVersion() error")

			tt.fields.repo.AssertCalled(t, "NextVersion", git.NextPatch, currentVersion)
		})

	}
}

func TestAction_checkDowngrade(t *testing.T) {
	const versionToCheck = version.V("1.2.3")

	type fields struct {
		repo            *__actionRepoMock
		allowDowngrades bool
	}

	repoFn := func(e error) *__actionRepoMock {
		r := &__actionRepoMock{}
		r.On("CheckDowngrade", versionToCheck).Return(e)
		return r
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "check downgrade",
			fields: fields{
				repo: repoFn(nil),
			},
			assertion: assert.NoError,
		},
		{
			name: "check downgrade error allow downgrade",
			fields: fields{
				repo:            repoFn(assert.AnError),
				allowDowngrades: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "check downgrade error not allow downgrade",
			fields: fields{
				repo:            repoFn(assert.AnError),
				allowDowngrades: false,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.AllowDowngrades, tt.fields.allowDowngrades)

			a := &Action{
				repo: tt.fields.repo,
			}

			tt.assertion(t, a.checkDowngrade(versionToCheck), "checkDowngrade() error")

			tt.fields.repo.AssertCalled(t, "CheckDowngrade", versionToCheck)
		})

	}
}

func TestAction_writeChangelog(t *testing.T) {
	const versionToCheck = version.V("1.2.3")

	type fields struct {
		changelogGen      *__actionChGenMock
		generateChangelog bool
	}

	changelogFn := func(e error) *__actionChGenMock {
		c := &__actionChGenMock{}
		c.On("Add", versionToCheck).Return(e)
		return c
	}

	tests := []struct {
		name     string
		fields   fields
		wantCall bool
	}{
		{
			name: "no changelog",
			fields: fields{
				changelogGen:      changelogFn(nil),
				generateChangelog: false,
			},
			wantCall: false,
		},
		{
			name: "changelog",
			fields: fields{
				changelogGen:      changelogFn(nil),
				generateChangelog: true,
			},
			wantCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.GenerateChangelog, tt.fields.generateChangelog)

			a := &Action{
				changelogGen: tt.fields.changelogGen,
			}

			assert.NoError(t, a.writeChangelog(versionToCheck), "writeChangelog() error")

			if tt.wantCall {
				tt.fields.changelogGen.AssertCalled(t, "Add", versionToCheck)
			} else {
				tt.fields.changelogGen.AssertNotCalled(t, "Add")
			}
		})
	}
}

func TestAction_runCommands(t *testing.T) {
	const versionToCheck = version.V("1.2.3")

	type fields struct {
		cmd     *__actionCmdMock
		dryRun  bool
		verbose bool
	}

	type args struct {
		cs []config.Command
	}

	commandsFn := func(breakOnError bool) []config.Command {
		return []config.Command{
			{
				Cmd:          []string{"echo", "test"},
				VersionFlag:  "version",
				BreakOnError: breakOnError,
			},
		}
	}

	cmdFn := func(e error) *__actionCmdMock {
		c := &__actionCmdMock{}
		c.On("Run", "echo", "test", "version=1.2.3").Return(e)
		return c
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCall  bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "dry run",
			fields: fields{
				cmd:    cmdFn(nil),
				dryRun: true,
			},
			args: args{
				cs: commandsFn(false),
			},
			assertion: assert.NoError,
		},
		{
			name: "dry run verbose",
			fields: fields{
				cmd:     cmdFn(nil),
				dryRun:  true,
				verbose: true,
			},
			args: args{
				cs: commandsFn(false),
			},
			assertion: assert.NoError,
		},
		{
			name: "normal command",
			fields: fields{
				cmd: cmdFn(nil),
			},
			args: args{
				cs: commandsFn(false),
			},
			wantCall:  true,
			assertion: assert.NoError,
		},
		{
			name: "command error",
			fields: fields{
				cmd: cmdFn(assert.AnError),
			},
			args: args{
				cs: commandsFn(false),
			},
			wantCall:  true,
			assertion: assert.NoError,
		},
		{
			name: "command break error",
			fields: fields{
				cmd: cmdFn(assert.AnError),
			},
			args: args{
				cs: commandsFn(true),
			},
			wantCall:  true,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.Verbose, tt.fields.verbose)
			viper.Set(key.DryRun, tt.fields.dryRun)

			a := &Action{
				cmd: tt.fields.cmd,
			}

			tt.assertion(t, a.runCommands(tt.args.cs, versionToCheck), "runCommands() error")

			if tt.wantCall {
				tt.fields.cmd.AssertCalled(t, "Run", "echo", "test", "version=1.2.3")
			} else {
				tt.fields.cmd.AssertNotCalled(t, "Run")
			}
		})
	}
}

type __actionRepoMock struct {
	mock.Mock
}

func (m *__actionRepoMock) IsClean() (bool, error) {
	ret := m.Called()

	r0 := ret.Get(0).(bool)
	r1 := ret.Error(1)

	return r0, r1
}

func (m *__actionRepoMock) NextVersion(nt git.NextType, custom version.V) (version.V, bool, error) {
	ret := m.Called(nt, custom)

	r0 := ret.Get(0).(version.V)
	r1 := ret.Get(1).(bool)
	r2 := ret.Error(2)

	return r0, r1, r2
}

func (m *__actionRepoMock) CheckDowngrade(v version.V) error {
	ret := m.Called(v)

	return ret.Error(0)
}

func (m *__actionRepoMock) CommitTag(v version.V) error {
	ret := m.Called(v)

	return ret.Error(0)
}

func (m *__actionRepoMock) AddModified() error {
	ret := m.Called()

	return ret.Error(0)
}

type __actionChGenMock struct {
	mock.Mock
}

func (m *__actionChGenMock) Add(v version.V) error {
	ret := m.Called(v)

	return ret.Error(0)
}

type __actionCfgMock struct {
	mock.Mock
}

func (m *__actionCfgMock) BumpFiles() []config.BumpFile {
	ret := m.Called()

	var r0 []config.BumpFile
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]config.BumpFile)
	}

	return r0
}

func (m *__actionCfgMock) CommandsBefore() []config.Command {
	ret := m.Called()

	var r0 []config.Command
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]config.Command)
	}

	return r0
}

func (m *__actionCfgMock) CommandsAfter() []config.Command {
	ret := m.Called()

	var r0 []config.Command
	if ret.Get(0) != nil {
		r0 = ret.Get(0).([]config.Command)
	}

	return r0
}

type __actionBumpMock struct {
	mock.Mock
}

func (m *__actionBumpMock) Apply(bumps []config.BumpFile, v version.V) {
	m.Called(bumps, v)
}

type __actionCmdMock struct {
	mock.Mock
}

func (m *__actionCmdMock) Run(name string, arg ...string) error {
	args := []any{name}
	for _, a := range arg {
		args = append(args, a)
	}
	ret := m.Called(args...)

	return ret.Error(0)
}

func TestActionType_String(t *testing.T) {
	tests := []struct {
		name string
		a    ActionType
		want string
	}{
		{
			name: "unknown",
			a:    ActionUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.String(), "String()")
		})
	}
}

func TestActionType_gitNextType(t *testing.T) {
	tests := []struct {
		name string
		a    ActionType
		want git.NextType
	}{
		{
			name: "patch",
			a:    ActionPatch,
			want: git.NextPatch,
		},
		{
			name: "minor",
			a:    ActionMinor,
			want: git.NextMinor,
		},
		{
			name: "major",
			a:    ActionMajor,
			want: git.NextMajor,
		},
		{
			name: "custom",
			a:    ActionCustom,
			want: git.NextCustom,
		},
		{
			name: "unknown",
			a:    ActionUnknown,
			want: git.NextNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.gitNextType(), "gitNextType()")
		})
	}
}
