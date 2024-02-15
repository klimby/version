package bump

import (
	"io"
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/config/key"
	"github.com/klimby/version/internal/service/console"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/pkg/version"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestB_Apply(t *testing.T) {
	const (
		ver          = version.V("1.0.0")
		backupError  = "create backup file composer.json error"
		applyWarning = "bump file composer.json error"
		repoWarning  = "add file composer.json to git error"
	)
	bmps := []config.BumpFile{
		{
			File: fsys.File("composer.json"),
		},
	}

	type fields struct {
		rw      *__rwMock
		process *__processMock
		repo    *__repoMock
		bcp     *__backupSrvMock
	}

	type wantCall struct {
		applyToFile bool
		repo        bool
	}

	type wantConsole struct {
		bumpError    bool
		applyWarning bool
		repoWarning  bool
	}

	tests := []struct {
		name        string
		fields      fields
		wantCall    wantCall
		wantConsole wantConsole
	}{
		{
			name: "normal",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
				repo: __newRepoMock(nil),
				bcp:  __newBackupSrvMock(nil),
			},
			wantCall: wantCall{
				applyToFile: true,
				repo:        true,
			},
		},
		{
			name: "no change",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: false,
				}),
				repo: __newRepoMock(nil),
				bcp:  __newBackupSrvMock(nil),
			},
			wantCall: wantCall{
				applyToFile: true,
				repo:        false,
			},
		},
		{
			name: "backup error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
				repo: __newRepoMock(nil),
				bcp:  __newBackupSrvMock(assert.AnError),
			},
			wantCall: wantCall{
				applyToFile: true,
				repo:        true,
			},
			wantConsole: wantConsole{
				bumpError: true,
			},
		},
		{
			name: "apply to file error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{}, // error generated
					changed: true,
				}),
				repo: __newRepoMock(nil),
				bcp:  __newBackupSrvMock(nil),
			},
			wantCall: wantCall{
				applyToFile: true,
			},
			wantConsole: wantConsole{
				applyWarning: true,
			},
		},
		{
			name: "repo error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
				repo: __newRepoMock(assert.AnError),
				bcp:  __newBackupSrvMock(nil),
			},
			wantCall: wantCall{
				applyToFile: true,
				repo:        true,
			},
			wantConsole: wantConsole{
				repoWarning: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(func(args *Args) {
				args.RW = tt.fields.rw
				args.Repo = tt.fields.repo
				args.Backup = tt.fields.bcp
				args.Proc = tt.fields.process
			})

			stdOut := &__consoleWriter{}
			stdErr := &__consoleWriter{}

			console.Init(func(options *console.OutArgs) {
				options.Stdout = stdOut
				options.Stderr = stdErr
				options.Colorize = false
			})

			b.Apply(bmps, ver)

			if tt.wantCall.applyToFile {
				tt.fields.rw.AssertCalled(t, "Read", mock.Anything)
			} else {
				tt.fields.rw.AssertNotCalled(t, "Read", mock.Anything)
			}

			if tt.wantCall.repo {
				tt.fields.repo.AssertCalled(t, "Add", mock.Anything)
			} else {
				tt.fields.repo.AssertNotCalled(t, "Add", mock.Anything)
			}

			stdErrStr := stdErr.String()
			stdOutStr := stdOut.String()

			if tt.wantConsole.bumpError {
				assert.Contains(t, stdErrStr, backupError)
			} else {
				assert.NotContains(t, stdErrStr, backupError)
			}

			if tt.wantConsole.applyWarning {
				assert.Contains(t, stdOutStr, applyWarning)
			} else {
				assert.NotContains(t, stdOutStr, applyWarning)
			}

			if tt.wantConsole.repoWarning {
				assert.Contains(t, stdOutStr, repoWarning)
			} else {
				assert.NotContains(t, stdOutStr, repoWarning)
			}

		})
	}
}

func TestB_applyToFile(t *testing.T) {
	const ver = version.V("1.0.0")
	bmp := config.BumpFile{
		File: fsys.File("composer.json"),
	}

	type fields struct {
		rw      *__rwMock
		process *__processMock
	}

	type wantCall struct {
		write bool
	}

	type viperArgs struct {
		dryRun bool
	}

	tests := []struct {
		name      string
		fields    fields
		viperArgs viperArgs
		wantCall  wantCall
		want      assert.BoolAssertionFunc
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
			},
			wantCall: wantCall{
				write: true,
			},
			want:    assert.True,
			wantErr: assert.NoError,
		},
		{
			name: "dry run",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
			},
			viperArgs: viperArgs{
				dryRun: true,
			},
			want:    assert.True,
			wantErr: assert.NoError,
		},
		{
			name: "write error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{writeErr: assert.AnError}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
			},
			wantCall: wantCall{
				write: true,
			},
			wantErr: assert.Error,
		},
		{
			name: "read error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{readErr: assert.AnError}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: true,
				}),
			},
			wantErr: assert.Error,
		},
		{
			name: "nil content error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{},
					changed: true,
				}),
			},
			wantErr: assert.Error,
		},
		{
			name: "no change",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{
					data:    []string{"foo", "bar"},
					changed: false,
				}),
			},
			want:    assert.False,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set(key.DryRun, tt.viperArgs.dryRun)

			b := New(func(args *Args) {
				args.RW = tt.fields.rw
				args.Proc = tt.fields.process
			})

			got, err := b.applyToFile(bmp, ver)

			tt.wantErr(t, err, "B.applyToFile() error = %v, wantErr %v", err, tt.wantErr)

			if tt.wantCall.write {
				tt.fields.rw.AssertCalled(t, "Write", mock.Anything, mock.Anything)
			} else {
				tt.fields.rw.AssertNotCalled(t, "Write", mock.Anything, mock.Anything)
			}

			if err == nil {
				tt.want(t, got, "B.applyToFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestB_read(t *testing.T) {
	type fields struct {
		rw      *__rwMock
		process *__processMock
	}

	type args struct {
		bmp config.BumpFile
	}

	type wantCall struct {
		predefinedJSON bool
		customFile     bool
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCall wantCall
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "call predefined JSON",
			fields: fields{
				rw:      __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{}),
			},
			args: args{
				bmp: config.BumpFile{
					File: fsys.File("composer.json"),
				},
			},
			wantCall: wantCall{
				predefinedJSON: true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "call custom file",
			fields: fields{
				rw:      __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{}),
			},
			args: args{
				bmp: config.BumpFile{
					File: fsys.File("custom.txt"),
				},
			},
			wantCall: wantCall{
				customFile: true,
			},
			wantErr: assert.NoError,
		},
		{
			name: "read error",
			fields: fields{
				rw:      __newRWMock(__rwMockArgs{readErr: assert.AnError}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{}),
			},
			args: args{
				bmp: config.BumpFile{
					File: fsys.File("composer.json"),
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "close error",
			fields: fields{
				rw:      __newRWMock(__rwMockArgs{closerErr: assert.AnError}, __rwMockArgs{}),
				process: __newProcessMock(__processMockArgs{}),
			},
			args: args{
				bmp: config.BumpFile{
					File: fsys.File("composer.json"),
				},
			},
			wantCall: wantCall{
				predefinedJSON: true,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(func(args *Args) {
				args.RW = tt.fields.rw
				args.Proc = tt.fields.process
			})

			_, _, err := b.read(tt.args.bmp, version.V("1.0.0"))

			tt.wantErr(t, err, "B.read() error = %v, wantErr %v", err, tt.wantErr)

			if tt.wantCall.predefinedJSON {
				tt.fields.process.AssertCalled(t, "PredefinedJSON", mock.Anything, tt.args.bmp, version.V("1.0.0"))
			} else {
				tt.fields.process.AssertNotCalled(t, "PredefinedJSON", mock.Anything, tt.args.bmp, version.V("1.0.0"))
			}

			if tt.wantCall.customFile {
				tt.fields.process.AssertCalled(t, "CustomFile", mock.Anything, tt.args.bmp, version.V("1.0.0"))
			} else {
				tt.fields.process.AssertNotCalled(t, "CustomFile", mock.Anything, tt.args.bmp, version.V("1.0.0"))
			}

		})
	}
}

func TestB_write(t *testing.T) {
	type fields struct {
		rw *__rwMock
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{}),
			},
			wantErr: assert.NoError,
		},
		{
			name: "write error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{writeErr: assert.AnError}),
			},
			wantErr: assert.Error,
		},
		{
			name: "close error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{closerErr: assert.AnError}),
			},
			wantErr: assert.Error,
		},
		{
			name: "writer error",
			fields: fields{
				rw: __newRWMock(__rwMockArgs{}, __rwMockArgs{writerErr: assert.AnError}),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(func(args *Args) {
				args.RW = tt.fields.rw
			})

			err := b.write("test", []string{"foo", "bar"})

			tt.wantErr(t, err, "B.write() error = %v, wantErr %v", err, tt.wantErr)

			if err == nil {
				res := tt.fields.rw.rwcTarget.buf.String()
				assert.Equal(t, "foo\nbar\n", res, "B.write() = %v, want %v", res, "foo\nbar\n")
			}
		})
	}
}

type __repoMock struct {
	mock.Mock
}

func (m *__repoMock) Add(files ...fsys.File) error {
	ret := m.Called(files)
	return ret.Error(0)
}

func __newRepoMock(err error) *__repoMock {
	m := &__repoMock{}
	m.On("Add", mock.Anything).Return(err)

	return m
}

type __backupSrvMock struct {
	mock.Mock
}

func (m *__backupSrvMock) Create(path string) error {
	ret := m.Called(path)
	return ret.Error(0)
}

func __newBackupSrvMock(err error) *__backupSrvMock {
	m := &__backupSrvMock{}
	m.On("Create", mock.Anything).Return(err)

	return m
}

type __processMock struct {
	mock.Mock
}

// PredefinedJSON
func (m *__processMock) PredefinedJSON(r io.Reader, bmp config.BumpFile, v version.V) ([]string, bool, error) {
	ret := m.Called(r, bmp, v)
	return ret.Get(0).([]string), ret.Bool(1), ret.Error(2)
}

// CustomFile
func (m *__processMock) CustomFile(r io.Reader, bmp config.BumpFile, v version.V) ([]string, bool, error) {
	ret := m.Called(r, bmp, v)
	return ret.Get(0).([]string), ret.Bool(1), ret.Error(2)
}

type __processMockArgs struct {
	data    []string
	changed bool
	err     error
}

func __newProcessMock(a __processMockArgs) *__processMock {
	m := &__processMock{}
	m.On("PredefinedJSON", mock.Anything, mock.Anything, mock.Anything).Return(a.data, a.changed, a.err)
	m.On("CustomFile", mock.Anything, mock.Anything, mock.Anything).Return(a.data, a.changed, a.err)

	return m
}
