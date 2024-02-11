package console

import (
	"errors"
	"reflect"
	"testing"
)

func Test_cmd(t *testing.T) {
	type args struct {
		name string
		arg  []string
	}

	tests := []struct {
		name     string
		args     args
		isError  bool
		wantCall bool
		wantErr  bool
	}{
		{
			name: "run command",
			args: args{
				name: "test",
				arg:  []string{"arg1", "arg2"},
			},
			wantCall: true,
		},
		{
			name:    "incorrect args",
			wantErr: true,
		},
		{
			name: "cmd error",
			args: args{
				name: "test",
			},
			isError:  true,
			wantCall: true,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &__runner{
				isError: tt.isError,
			}

			c := NewCmd(func(args *CmdArgs) {
				args.CF = func(name string, arg ...string) runner {
					return r
				}
			})

			err := c.Run(tt.args.name, tt.args.arg...)

			if (err != nil) != tt.wantErr {
				t.Errorf("cmd() error = %v, wantErr %v", err, tt.wantErr)
			}

			if r.called != tt.wantCall {
				t.Errorf("cmd() runner.called = %v, want %v", r.called, tt.wantCall)
			}
		})
	}

}

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
				arg:  []string{"arg1", "arg2"},
			},
			want: "test arg1 arg2",
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
				name: "test",
				arg:  []string{"arg1", "arg2"},
			},
			want:    "test",
			want1:   []string{"arg1", "arg2"},
			wantErr: false,
		},
		{
			name: "normalize args 2",
			args: args{
				name: "test arg1 arg2",
				arg:  []string{"arg3"},
			},
			want:    "test",
			want1:   []string{"arg1", "arg2", "arg3"},
			wantErr: false,
		},
		{
			name:    "normalize args 3",
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

type __runner struct {
	called  bool
	isError bool
}

// reset runner
func (r *__runner) Reset() {
	r.called = false
	r.isError = false
}

func (r *__runner) Run() error {
	r.called = true
	if r.isError {
		return errors.New("error")
	}
	return nil
}
