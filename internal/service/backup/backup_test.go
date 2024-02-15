package backup

import (
	"io/fs"
	"testing"

	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/console"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestService_Create(t *testing.T) {
	type fields struct {
		rw *__rwMock
	}

	type viperArgs struct {
		backup bool
		dryRun bool
	}

	type wantCall struct {
		readSource  bool
		writeTarget bool
	}

	data := []byte("test")

	const patch = __patch

	tests := []struct {
		name      string
		fields    fields
		viperArgs viperArgs
		wantCall  wantCall
		wantCopy  bool
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: true,
			},
			wantCopy: true,
			wantErr:  assert.NoError,
		},
		{
			name: "no backup",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{},
			wantCall: wantCall{
				readSource:  false,
				writeTarget: false,
			},
			wantCopy: false,
			wantErr:  assert.NoError,
		},
		{
			name: "dry run",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
				dryRun: true,
			},
			wantCall: wantCall{
				readSource:  false,
				writeTarget: false,
			},
			wantCopy: false,
			wantErr:  assert.NoError,
		},
		{
			name: "read error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data, readErr: assert.AnError}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: false,
			},
			wantCopy: false,
			wantErr:  assert.Error,
		},
		{
			name: "ErrNotExist error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data, readErr: fs.ErrNotExist}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: false,
			},
			wantCopy: false,
			wantErr:  assert.NoError,
		},
		{
			name: "close reader error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data, closerErr: assert.AnError}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: true,
			},
			wantCopy: true,
			wantErr:  assert.Error,
		},
		{
			name: "writer error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data}, __rwMockArgs{writeErr: assert.AnError}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: true,
			},
			wantCopy: false,
			wantErr:  assert.Error,
		},
		{
			name: "close writer error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data}, __rwMockArgs{closerErr: assert.AnError}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: true,
			},
			wantCopy: true,
			wantErr:  assert.Error,
		},
		{
			name: "io copy",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{data: data, readerErr: assert.AnError}, __rwMockArgs{}),
			},
			viperArgs: viperArgs{
				backup: true,
			},
			wantCall: wantCall{
				readSource:  true,
				writeTarget: true,
			},
			wantCopy: false,
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.Backup, tt.viperArgs.backup)
			viper.Set(key.DryRun, tt.viperArgs.dryRun)

			s := New(func(args *Args) {
				args.RW = tt.fields.rw
			})

			err := s.Create(patch)

			tt.wantErr(t, err, "Create() error")

			if tt.wantCall.readSource {
				tt.fields.rw.AssertCalled(t, "Read", patch)
			} else {
				tt.fields.rw.AssertNotCalled(t, "Read", patch)
			}

			if tt.wantCall.writeTarget {
				tt.fields.rw.AssertCalled(t, "Write", patch+_suffix, mock.Anything)
			} else {
				tt.fields.rw.AssertNotCalled(t, "Write", patch+_suffix, mock.Anything)
			}

			writeData := tt.fields.rw.rwcTarget.buf.Bytes()

			if tt.wantCopy {
				assert.Equal(t, data, writeData, "write data")
			} else {
				assert.Empty(t, writeData, "write data")
			}
		})
	}
}
