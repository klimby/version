package remove

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAction_Run(t *testing.T) {
	type fields struct {
		actionType ActionType
		remover    *__removerMock
		cfg        *__cfgSrvMock
	}

	removerMock := func() *__removerMock {
		remover := &__removerMock{}
		remover.On("Remove", "config.yaml")
		remover.On("Remove", "file")
		return remover
	}

	cfgMock := func() *__cfgSrvMock {
		cfg := &__cfgSrvMock{}
		cfg.On("BumpFiles").Return([]config.BumpFile{
			{
				File: fsys.File("file"),
			},
		})
		return cfg
	}

	tests := []struct {
		name      string
		fields    fields
		wantCall  bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "run",
			fields: fields{
				actionType: ActionBackup,
				remover:    removerMock(),
				cfg:        cfgMock(),
			},
			wantCall:  true,
			assertion: assert.NoError,
		},
		{
			name: "unknown type",
			fields: fields{
				actionType: ActionUnknown,
				remover:    removerMock(),
				cfg:        cfgMock(),
			},
			wantCall:  false,
			assertion: assert.Error,
		},
		{
			name: "validate error",
			fields: fields{
				actionType: ActionBackup,
				remover:    removerMock(),
			},
			wantCall:  false,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.CfgFile, "config.yaml")

			a := New(func(arg *Args) {
				arg.ActionType = tt.fields.actionType

				if tt.fields.cfg != nil {
					arg.Cfg = tt.fields.cfg
				}

				if tt.fields.remover != nil {
					arg.Remover = tt.fields.remover
				} else {
					arg.Remover = nil
				}

			})

			tt.assertion(t, a.Run(), tt.name)

			if tt.fields.remover != nil {
				if tt.wantCall {
					tt.fields.remover.AssertCalled(t, "Remove", "file")
				} else {
					tt.fields.remover.AssertNotCalled(t, "Remove", "file")
				}
			}
		})

	}
}

type __removerMock struct {
	mock.Mock
}

func (m *__removerMock) Remove(path ...string) {
	args := make([]any, len(path))
	for i, p := range path {
		args[i] = p
	}
	m.Called(args...)
}

type __cfgSrvMock struct {
	mock.Mock
}

func (m *__cfgSrvMock) BumpFiles() []config.BumpFile {
	ret := m.Called()

	var r0 []config.BumpFile
	if rf, ok := ret.Get(0).(func() []config.BumpFile); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).([]config.BumpFile)
	}

	return r0
}

func TestActionType_String(t *testing.T) {
	tests := []struct {
		name string
		ft   ActionType
		want string
	}{
		{
			name: "unknown",
			ft:   ActionUnknown,
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ft.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
