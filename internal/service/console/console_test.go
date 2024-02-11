package console

import (
	"bytes"
	"testing"

	"github.com/klimby/version/pkg/convert"
)

func Test_console(t *testing.T) {
	tests := []struct {
		name       string
		colorize   bool
		s          string
		wantStdOut string
		wantStdErr string
		fn         func(string)
	}{
		{
			name:       "Empty colorize",
			colorize:   true,
			s:          "",
			wantStdOut: "",
			fn:         Info,
		},
		{
			name:       "Info",
			s:          "test",
			wantStdOut: "test\n",
			fn:         Info,
		},
		{
			name:       "Info colorize",
			colorize:   true,
			s:          "test",
			wantStdOut: "\033[37mtest\033[0m\n",
			fn:         Info,
		},
		{
			name:       "Notice",
			s:          "test",
			wantStdOut: "test\n",
			fn:         Notice,
		},
		{
			name:       "Notice colorize",
			colorize:   true,
			s:          "test",
			wantStdOut: "\033[36mtest\033[0m\n",
			fn:         Notice,
		},
		{
			name:       "Success",
			s:          "test",
			wantStdOut: "test\n",
			fn:         Success,
		},
		{
			name:       "Success colorize",
			colorize:   true,
			s:          "test",
			wantStdOut: "\033[32mtest\033[0m\n",
			fn:         Success,
		},
		{
			name:       "Warn",
			s:          "test",
			wantStdOut: "test\n",
			fn:         Warn,
		},
		{
			name:       "Warn colorize",
			colorize:   true,
			s:          "test",
			wantStdOut: "\033[33mtest\033[0m\n",
			fn:         Warn,
		},
		{
			name:       "Error",
			s:          "test",
			wantStdErr: "test\n",
			fn:         Error,
		},
		{
			name:       "Error colorize",
			colorize:   true,
			s:          "test",
			wantStdErr: "\033[31mtest\033[0m\n",
			fn:         Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdOut := &__testWriter{}
			stdErr := &__testWriter{}

			Init(func(options *OutArgs) {
				options.Stdout = stdOut
				options.Stderr = stdErr
				options.Colorize = tt.colorize
			})

			tt.fn(tt.s)

			if tt.wantStdErr != stdErr.String() {
				t.Errorf("console() wantStdErr = %v, got %v", tt.wantStdErr, stdErr.String())
			}

			if tt.wantStdOut != stdOut.String() {
				t.Errorf("console() wantStdOut = %v, got %v", tt.wantStdOut, stdOut.String())
			}

		})
	}
}

// testWriter - write to buffer.
type __testWriter struct {
	buffer bytes.Buffer
}

// Write - write to buffer.
func (n *__testWriter) Write(p []byte) (int, error) {
	return n.buffer.Write(p)
}

// String - return buffer as string.
func (n *__testWriter) String() string {
	s := convert.B2S(n.buffer.Bytes())
	n.buffer.Reset()

	return s
}
