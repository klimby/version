package console

import (
	"io"
	"reflect"
	"testing"
)

func Test_commandString(t *testing.T) {
	type args struct {
		name string
		arg  []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "command string",
			args: args{
				name: "test",
				arg:  []string{"1", "2"},
			},
			want: "test 1 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commandString(tt.args.name, tt.args.arg...); got != tt.want {
				t.Errorf("commandString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeArgs(t *testing.T) {
	type args struct {
		name string
		arg  []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []string
		wantErr bool
	}{
		{
			name: "normalize args",
			args: args{
				name: "test 1 2",
			},
			want:  "test",
			want1: []string{"1", "2"},
		},
		{
			name: "normalize args",
			args: args{
				name: "test",
				arg:  []string{"1", "2"},
			},
			want:  "test",
			want1: []string{"1", "2"},
		},
		{
			name: "normalize args with empty name",
			args: args{
				name: "",
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := normalizeArgs(tt.args.name, tt.args.arg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeArgs() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("normalizeArgs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_stdErrOutput_Write(t *testing.T) {

	type args struct {
		p []byte
	}
	tests := []struct {
		name       string
		args       args
		want       int
		wantStderr string
	}{
		{
			name: "write to stderr",
			args: args{
				p: []byte("test"),
			},
			want:       4,
			wantStderr: "test\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = false
			})

			s := &stdErrOutput{}

			got, _ := s.Write(tt.args.p)
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != tt.wantStderr {
				t.Errorf("Error() = %v, want %v", got, tt.wantStderr)
			}

			if !s.isError() {
				t.Errorf("isError() = %v, want %v", s.isError(), true)
			}
		})
	}
}

func Test_stdStdOutput_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantStd string
	}{
		{
			name: "write to stdout",
			args: args{
				p: []byte("test"),
			},
			want:    4,
			wantStd: "test\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = false
			})

			s := &stdOutput{}

			got, _ := s.Write(tt.args.p)
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}

			if got := stdout.p; string(got) != tt.wantStd {
				t.Errorf("Error() = %v, want %v", got, tt.wantStd)
			}

		})
	}
}

func TestNewCmd(t *testing.T) {
	c := NewCmd()

	if c == nil {
		t.Errorf("NewCmd() = %v, want %v", c, nil)
	}
}

func TestCmd_Run(t *testing.T) {
	type args struct {
		name string
		arg  []string
	}
	tests := []struct {
		name           string
		runnerBehavior int // 0 - ok, 1 - error on run, 2 - error in stdErr
		args           args
		wantErr        bool
	}{
		{
			name: "normalize args error",
			args: args{
				name: "",
			},
			wantErr: true,
		},
		{
			name:           "runner error",
			runnerBehavior: 1,
			args: args{
				name: "test",
			},
			wantErr: true,
		},
		{
			name:           "runner error in stderr",
			runnerBehavior: 2,
			args: args{
				name: "test",
			},
			wantErr: true,
		},
		{
			name:           "runner ok",
			runnerBehavior: 0,
			args: args{
				name: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sO := &__fakeCmdStdout{}
			eO := &__fakeCmdStderr{}
			c := &Cmd{
				commandFactory: func(name string, arg ...string) runner {
					cmd := &__fakeRunner{
						stdOut:   sO,
						stdErr:   eO,
						behavior: tt.runnerBehavior,
					}

					return cmd
				},
				stdout: sO,
				stderr: eO,
			}

			if err := c.Run(tt.args.name, tt.args.arg...); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.runnerBehavior == 0 && tt.wantErr == false {
				if got := sO.p; string(got) != "test" {
					t.Errorf("Run() = %v, want %v", got, "test")
				}
			}

			if tt.runnerBehavior == 2 {
				if got := eO.p; string(got) != "test" {
					t.Errorf("Run() = %v, want %v", got, "test")
				}
			}

		})
	}
}

type __fakeRunner struct {
	stdOut   io.Writer
	stdErr   io.Writer
	behavior int // 0 - ok, 1 - error on run, 2 - error in stdErr
}

func (f *__fakeRunner) Run() error {
	if f.behavior == 1 {
		return io.EOF
	}

	if f.behavior == 2 {
		if _, err := f.stdErr.Write([]byte("test")); err != nil {
			return err
		}

		return nil
	}

	if _, err := f.stdOut.Write([]byte("test")); err != nil {
		return err
	}

	return nil
}

type __fakeCmdStderr struct {
	p []byte
}

func (f *__fakeCmdStderr) Write(p []byte) (int, error) {
	f.p = p
	return len(p), nil
}

func (f *__fakeCmdStderr) isError() bool {
	return len(f.p) > 0
}

type __fakeCmdStdout struct {
	p []byte
}

func (f *__fakeCmdStdout) Write(p []byte) (int, error) {
	f.p = p
	return len(p), nil
}
