package console

import (
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
