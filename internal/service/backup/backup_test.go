package backup

import (
	"io/fs"
	"testing"

	"github.com/klimby/version/internal/service/console"
	"github.com/stretchr/testify/assert"
)

func TestService_Remove(t *testing.T) {
	type fields struct {
		rw *__rwMock
	}

	type wantConsole struct {
		error   assert.BoolAssertionFunc
		success assert.BoolAssertionFunc
	}

	tests := []struct {
		name        string
		fields      fields
		wantConsole wantConsole
	}{
		{
			name: "normal",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
			},
			wantConsole: wantConsole{
				success: assert.True,
				error:   assert.False,
			},
		},
		{
			name: "remove error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{removeErr: assert.AnError}, __rwMockArgs{}),
			},
			wantConsole: wantConsole{
				success: assert.False,
				error:   assert.True,
			},
		},
		{
			name: "remove ErrNotExist error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{removeErr: fs.ErrNotExist}, __rwMockArgs{}),
			},
			wantConsole: wantConsole{
				success: assert.False,
				error:   assert.False,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(func(args *Args) {
				args.RW = tt.fields.rw
			})

			stdOut := &__testWriter{}
			stdErr := &__testWriter{}

			console.Init(func(options *console.OutArgs) {
				options.Stdout = stdOut
				options.Stderr = stdErr
				options.Colorize = false
			})

			s.Remove("test")

			tt.wantConsole.error(t, len(stdErr.String()) > 0, "stdErr error")
			tt.wantConsole.success(t, len(stdOut.String()) > 0, "stdOut error")

		})
	}
}
