package console

import (
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/spf13/viper"
)

func TestError(t *testing.T) {
	type args struct {
		s        string
		colorize bool
		silent   bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "not silent not colorize",
			args: args{
				s: "test",
			},
			want: "test\n",
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "not silent colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(err) + "test" + string(reset) + "\n",
		},
		{
			name: "silent not colorize",
			args: args{
				s:      "test",
				silent: true,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			viper.Set(config.Silent, tt.args.silent)

			defer viper.Set(config.Silent, false)

			Error(tt.args.s)

			if got := stderr.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stdout.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(warning) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			Warn(tt.args.s)

			if got := stdout.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(success) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			Success(tt.args.s)

			if got := stdout.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestNotice(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(notice) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			Notice(tt.args.s)

			if got := stdout.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		s        string
		colorize bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: string(info) + "test" + string(reset) + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			Info(tt.args.s)

			if got := stdout.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

func TestVerbose(t *testing.T) {
	type args struct {
		s        string
		colorize bool
		verbose  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "colorize",
			args: args{
				s:        "test",
				colorize: true,
				verbose:  true,
			},
			want: string(info) + "test" + string(reset) + "\n",
		},
		{
			name: "no verbose",
			args: args{
				s:        "test",
				colorize: true,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &__fakeWriter{}
			stderr := &__fakeWriter{}

			SetOutput(func(arg *OutArgs) {
				arg.Stdout = stdout
				arg.Stderr = stderr
				arg.Colorize = tt.args.colorize
			})

			viper.Set(config.Verbose, tt.args.verbose)

			defer viper.Set(config.Verbose, false)

			Verbose(tt.args.s)

			if got := stdout.p; string(got) != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}

			if got := stderr.p; string(got) != "" {
				t.Errorf("Error() = %v, want %v", got, "")
			}
		})
	}
}

type __fakeWriter struct {
	p []byte
}

func (f *__fakeWriter) Write(p []byte) (int, error) {
	f.p = p
	return len(p), nil
}
