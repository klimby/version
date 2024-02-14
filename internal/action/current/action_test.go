package current

import (
	"testing"

	"github.com/klimby/version/pkg/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAction_Run(t *testing.T) {
	const ver = version.V("1.2.4")

	type fields struct {
		repo *__repoMock
	}

	repoMock := func(e error) *__repoMock {
		repo := &__repoMock{}
		repo.On("Current").Return(ver, e)
		return repo
	}

	tests := []struct {
		name      string
		fields    fields
		wantCall  bool
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "validate error",
			fields:    fields{},
			assertion: assert.Error,
		},
		{
			name: "repo error",
			fields: fields{
				repo: repoMock(assert.AnError),
			},
			wantCall:  true,
			assertion: assert.Error,
		},
		{
			name: "repo ok",
			fields: fields{
				repo: repoMock(nil),
			},
			wantCall:  true,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(func(args *Args) {
				if tt.fields.repo != nil {
					args.Repo = tt.fields.repo
				}
			})

			tt.assertion(t, a.Run(), tt.name)

			if tt.fields.repo != nil {
				if tt.wantCall {
					tt.fields.repo.AssertExpectations(t)
				} else {
					tt.fields.repo.AssertNotCalled(t, "Current")
				}
			}
		})
	}
}

type __repoMock struct {
	mock.Mock
}

func (m *__repoMock) Current() (version.V, error) {
	args := m.Called()
	return args.Get(0).(version.V), args.Error(1)
}
