package bump

import (
	"fmt"
	"io"
	"regexp"
	"testing"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/internal/service/fsys"
	"github.com/klimby/version/pkg/version"
	"github.com/stretchr/testify/assert"
)

func Test_process_PredefinedJSON(t *testing.T) {
	type args struct {
		r   io.Reader
		bmp config.BumpFile
		v   version.V
	}

	ioReader := func(data string, readError error) io.Reader {
		b := __RWC{
			readError: readError,
		}

		_, _ = b.Write([]byte(data))

		return &b
	}

	tests := []struct {
		name        string
		args        args
		want        []string
		wantChanged assert.BoolAssertionFunc
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				r: ioReader(`{"version": "1.0.0"}`, nil),
				bmp: config.BumpFile{
					File: fsys.File("composer.json"),
				},
				v: version.V("1.0.1"),
			},

			want:        []string{`{"version": "1.0.1"}`},
			wantChanged: assert.True,
			wantErr:     assert.NoError,
		},
		{
			name: "read error",
			args: args{
				r: ioReader(`{"version": "1.0.0"}`, assert.AnError),
				bmp: config.BumpFile{
					File: fsys.File("composer.json"),
				},
			},

			want:        nil,
			wantChanged: assert.False,
			wantErr:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := process{}
			got, gotChanged, err := pr.PredefinedJSON(tt.args.r, tt.args.bmp, tt.args.v)
			if !tt.wantErr(t, err, fmt.Sprintf("PredefinedJSON(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)) {
				return
			}
			assert.Equalf(t, tt.want, got, "PredefinedJSON(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)

			tt.wantChanged(t, gotChanged, "PredefinedJSON(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)
		})
	}
}

func Test_process_CustomFile(t *testing.T) {
	type args struct {
		r   io.Reader
		bmp config.BumpFile
		v   version.V
	}

	ioReader := func(data string, readError error) io.Reader {
		b := __RWC{
			readError: readError,
		}

		_, _ = b.Write([]byte(data))

		return &b
	}

	tests := []struct {
		name        string
		args        args
		want        []string
		wantChanged assert.BoolAssertionFunc
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "normal",
			args: args{
				r: ioReader("1.0.0", nil),
				bmp: config.BumpFile{
					File: fsys.File("file"),
				},
				v: version.V("1.0.1"),
			},

			want:        []string{"1.0.1"},
			wantChanged: assert.True,
			wantErr:     assert.NoError,
		},
		{
			name: "regexp",
			args: args{
				r: ioReader(`![Version v1.0.9]`, nil),
				bmp: config.BumpFile{
					File:   fsys.File("file"),
					RegExp: []string{`^!\[Version.*$`},
				},
				v: version.V("1.0.1"),
			},

			want:        []string{`![Version v1.0.1]`},
			wantChanged: assert.True,
			wantErr:     assert.NoError,
		},
		{
			name: "read error",
			args: args{
				r: ioReader("1.0.0", assert.AnError),
				bmp: config.BumpFile{
					File: fsys.File("file"),
				},
				v: version.V("1.0.1"),
			},

			want:        nil,
			wantChanged: assert.False,
			wantErr:     assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := process{}
			got, gotChanged, err := pr.CustomFile(tt.args.r, tt.args.bmp, tt.args.v)
			if !tt.wantErr(t, err, fmt.Sprintf("CustomFile(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CustomFile(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)

			tt.wantChanged(t, gotChanged, "CustomFile(%v, %v, %v)", tt.args.r, tt.args.bmp, tt.args.v)
		})
	}
}

func Test_handleBumpFile(t *testing.T) {
	type args struct {
		bmp config.BumpFile
	}
	tests := []struct {
		name      string
		args      args
		wantStart int
		wantEnd   int
		wantRegs  []regexp.Regexp
	}{
		{
			name: "normal",
			args: args{
				bmp: config.BumpFile{
					File:  fsys.File("file"),
					Start: 1,
					End:   2,
					RegExp: []string{
						"regexp",
					},
				},
			},
			wantStart: 1,
			wantEnd:   2,
			wantRegs: []regexp.Regexp{
				*regexp.MustCompile("regexp"),
			},
		},
		{
			name: "invalid regexp",
			args: args{
				bmp: config.BumpFile{
					File: fsys.File("file"),

					RegExp: []string{
						"(",
					},
				},
			},
			wantStart: 0,
			wantEnd:   int(^uint(0) >> 1),
			wantRegs:  []regexp.Regexp{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, gotRegs := handleBumpFile(tt.args.bmp)
			assert.Equalf(t, tt.wantStart, gotStart, "handleBumpFile(%v)", tt.args.bmp)
			assert.Equalf(t, tt.wantEnd, gotEnd, "handleBumpFile(%v)", tt.args.bmp)
			assert.Equalf(t, tt.wantRegs, gotRegs, "handleBumpFile(%v)", tt.args.bmp)
		})
	}
}
