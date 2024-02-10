package console

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/klimby/version/pkg/convert"
)

type col string

const (
	err     col = "\033[31m" // red
	warning col = "\033[33m" // yellow
	success col = "\033[32m" // green
	notice  col = "\033[36m" // cyan
	info    col = "\033[37m" // white
	reset   col = "\033[0m"  // reset color
)

var c = &out{
	stdout:   &nilWriter{},
	stderr:   &nilWriter{},
	colorize: false,
}

type out struct {
	stdout   io.Writer
	stderr   io.Writer
	colorize bool
}

// OutArgs is a console arguments.
type OutArgs struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Colorize bool
}

// InitTest - init console for testing.
func InitTest(colorize ...bool) func() {
	c.stderr = &testWriter{}
	c.stdout = &testWriter{}

	if len(colorize) > 0 {
		c.colorize = colorize[0]
	}

	return func() {
		c.stderr = &nilWriter{}
		c.stdout = &nilWriter{}
		c.colorize = false
	}
}

// Init - init console.
// WARNING: without args.
func Init(args ...func(*OutArgs)) {
	arg := OutArgs{
		Stdout:   &nilWriter{},
		Stderr:   &nilWriter{},
		Colorize: false,
	}

	for _, f := range args {
		f(&arg)
	}

	c.stdout = arg.Stdout
	c.stderr = arg.Stderr
	c.colorize = arg.Colorize
}

// Read - read stdout and stderr, if they are testWriter.
// For testing.
func Read() (stdout, stderr string) {
	// check if stdout is testWriter
	if w, ok := c.stdout.(*testWriter); ok {
		stdout = w.String()
	}

	// check if stderr is testWriter
	if w, ok := c.stderr.(*testWriter); ok {
		stderr = w.String()
	}

	return stdout, stderr
}

// printLn - print line.
func printLn(clr col, s string) {
	if s == "" {
		return
	}

	var b strings.Builder

	if c.colorize {
		b.Grow(len(clr) + len(s) + len(reset) + 1)

		b.WriteString(string(clr))
		b.WriteString(s)
		b.WriteString(string(reset))
		b.WriteString("\n")
	} else {
		b.Grow(len(s) + 1)

		b.WriteString(s)
		b.WriteString("\n")
	}

	w := c.stdout
	if clr == err {
		w = c.stderr
	}

	_, err := w.Write(convert.S2B(b.String()))
	if err != nil {
		//nolint:forbidigo
		fmt.Println(b.String())
	}
}

// Error - print error message.
func Error(s string) {
	printLn(err, s)
}

// Warn - print warning message.
func Warn(s string) {
	printLn(warning, s)
}

// Success - print success message.
func Success(s string) {
	printLn(success, s)
}

// Notice - print notice message.
func Notice(s string) {
	printLn(notice, s)
}

// Info - print info message.
func Info(s string) {
	printLn(info, s)
}

// nilWriter - write to nowhere.
type nilWriter struct{}

// Write - write to nowhere.
func (*nilWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

// testWriter - write to buffer.
type testWriter struct {
	buffer bytes.Buffer
}

// Write - write to buffer.
func (n *testWriter) Write(p []byte) (int, error) {
	return n.buffer.Write(p)
}

// String - return buffer as string.
func (n *testWriter) String() string {
	s := convert.B2S(n.buffer.Bytes())
	n.buffer.Reset()

	return s
}
