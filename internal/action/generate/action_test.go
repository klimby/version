package generate

import (
	"testing"

	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAction_Run(t *testing.T) {
	type fields struct {
		actionType   ActionType
		backup       *__backupServiceMock
		changelogGen *__generatorMock
		cfg          *__generatorMock
	}

	generatorMock := func(e error) *__generatorMock {
		gen := &__generatorMock{}
		gen.On("Generate").Return(e)
		return gen
	}

	backupServiceMock := func(e error) *__backupServiceMock {
		backup := &__backupServiceMock{}
		backup.On("Create", mock.Anything).Return(e)
		return backup
	}

	type wantCall struct {
		config    bool
		changelog bool
	}

	tests := []struct {
		name      string
		fields    fields
		wantCall  wantCall
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "run config",
			fields: fields{
				actionType: FileConfig,
				backup:     backupServiceMock(nil),
				cfg:        generatorMock(nil),
			},
			wantCall: wantCall{
				config: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "run changelog",
			fields: fields{
				actionType:   FileChangelog,
				backup:       backupServiceMock(nil),
				changelogGen: generatorMock(nil),
			},
			wantCall: wantCall{
				changelog: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "unknown error",
			fields: fields{
				actionType:   FileUnknown,
				backup:       backupServiceMock(nil),
				changelogGen: generatorMock(nil),
				cfg:          generatorMock(nil),
			},
			assertion: assert.Error,
		},
		{
			name: "validate error",
			fields: fields{
				actionType: FileConfig,
				backup:     backupServiceMock(nil),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.CfgFile, "config.yaml")
			viper.Set(key.GenerateChangelog, true)

			a := New(func(args *Args) {
				args.ActionType = tt.fields.actionType
				args.Backup = tt.fields.backup

				if tt.fields.cfg != nil {
					args.CfgGenerator = tt.fields.cfg
				}

				if tt.fields.changelogGen != nil {
					args.ChangelogGen = tt.fields.changelogGen
				}
			})

			tt.assertion(t, a.Run(), "Run()")

			if tt.fields.cfg != nil {
				if tt.wantCall.config {
					tt.fields.cfg.AssertExpectations(t)
				} else {
					tt.fields.cfg.AssertNotCalled(t, "Generate")
				}
			}

			if tt.fields.changelogGen != nil {
				if tt.wantCall.changelog {
					tt.fields.changelogGen.AssertExpectations(t)
				} else {
					tt.fields.changelogGen.AssertNotCalled(t, "Generate")
				}
			}
		})
	}
}

func TestAction_validate(t *testing.T) {
	type fields struct {
		actionType ActionType
	}

	tests := []struct {
		name      string
		fields    fields
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "config error",
			fields: fields{
				actionType: FileConfig,
			},
			assertion: assert.Error,
		},
		{
			name: "changelog error",
			fields: fields{
				actionType: FileChangelog,
			},
			assertion: assert.Error,
		},
		{
			name: "ok",
			fields: fields{
				actionType: FileUnknown,
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(func(args *Args) {
				args.ActionType = tt.fields.actionType
			})

			tt.assertion(t, a.validate(), "validate()")
		})
	}
}

func TestAction_config(t *testing.T) {
	type fields struct {
		backup *__backupServiceMock
		cfg    *__generatorMock
	}

	generatorMock := func(e error) *__generatorMock {
		gen := &__generatorMock{}
		gen.On("Generate").Return(e)
		return gen
	}

	backupServiceMock := func(e error) *__backupServiceMock {
		backup := &__backupServiceMock{}
		backup.On("Create", mock.Anything).Return(e)
		return backup
	}

	type wantCall struct {
		config bool
		backup bool
	}

	tests := []struct {
		name      string
		fields    fields
		wantCall  wantCall
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "ok",
			fields: fields{
				backup: backupServiceMock(nil),
				cfg:    generatorMock(nil),
			},
			wantCall: wantCall{
				config: true,
				backup: true,
			},
			assertion: assert.NoError,
		},
		{
			name: "backup err",
			fields: fields{
				backup: backupServiceMock(assert.AnError),
				cfg:    generatorMock(nil),
			},
			wantCall: wantCall{
				config: false,
				backup: true,
			},
			assertion: assert.Error,
		},
		{
			name: "config err",
			fields: fields{
				backup: backupServiceMock(nil),
				cfg:    generatorMock(assert.AnError),
			},
			wantCall: wantCall{
				config: true,
				backup: true,
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.CfgFile, "config.yaml")

			a := New(func(args *Args) {
				args.Backup = tt.fields.backup
				args.CfgGenerator = tt.fields.cfg
			})

			tt.assertion(t, a.config(), "config()")

			if tt.fields.cfg != nil {
				if tt.wantCall.config {
					tt.fields.cfg.AssertExpectations(t)
				} else {
					tt.fields.cfg.AssertNotCalled(t, "Generate")
				}
			}

			if tt.fields.backup != nil {
				if tt.wantCall.backup {
					tt.fields.backup.AssertExpectations(t)
				} else {
					tt.fields.backup.AssertNotCalled(t, "Create")
				}
			}
		})

	}
}

func TestAction_changelog(t *testing.T) {
	type fields struct {
		changelogGen *__generatorMock
	}

	generatorMock := func(e error) *__generatorMock {
		gen := &__generatorMock{}
		gen.On("Generate").Return(e)
		return gen
	}

	tests := []struct {
		name              string
		fields            fields
		generateChangelog bool
		wantCall          bool
	}{
		{
			name: "ok",
			fields: fields{
				changelogGen: generatorMock(nil),
			},
			generateChangelog: true,
			wantCall:          true,
		},
		{
			name: "no generate",
			fields: fields{
				changelogGen: generatorMock(nil),
			},
			wantCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.GenerateChangelog, tt.generateChangelog)

			a := New(func(args *Args) {
				args.ChangelogGen = tt.fields.changelogGen
			})

			assert.NoError(t, a.changelog(), "changelog()")

			if tt.wantCall {
				tt.fields.changelogGen.AssertExpectations(t)
			} else {
				tt.fields.changelogGen.AssertNotCalled(t, "Generate")
			}
		})

	}
}

type __generatorMock struct {
	mock.Mock
}

func (m *__generatorMock) Generate() error {
	args := m.Called()
	return args.Error(0)
}

type __backupServiceMock struct {
	mock.Mock
}

func (m *__backupServiceMock) Create(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func TestActionType_String(t *testing.T) {
	tests := []struct {
		name string
		ft   ActionType
		want string
	}{
		{
			name: "unknown",
			ft:   FileUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.ft.String(), "String()")
		})
	}
}
