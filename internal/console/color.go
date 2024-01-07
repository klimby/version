package console

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klimby/version/internal/config"
	"github.com/klimby/version/pkg/convert"
	"github.com/spf13/viper"
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

type out struct {
	stdout   io.Writer
	stderr   io.Writer
	colorize bool
}

var c = &out{
	stdout:   os.Stdout,
	stderr:   os.Stderr,
	colorize: true,
}

// OutArgs - arguments for SetOutput.
//
//goland:noinspection GoUnusedExportedFunction,GoUnnecessarilyExportedIdentifiers
type OutArgs struct {
	Stdout   io.Writer
	Stderr   io.Writer
	Colorize bool
}

// SetOutput - set stdout and stderr.
// For testing.
//
//goland:noinspection GoUnusedExportedFunction,GoUnnecessarilyExportedIdentifiers
func SetOutput(args ...func(arg *OutArgs)) {
	if len(args) == 0 {
		return
	}

	a := &OutArgs{
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		Colorize: true,
	}

	for _, arg := range args {
		arg(a)
	}

	c.stdout = a.Stdout
	c.stderr = a.Stderr
	c.colorize = a.Colorize
}

// printLn - print line.
func printLn(clr col, s string) {
	if !viper.GetBool(config.Silent) {
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

// Verbose - print verbose message.
func Verbose(s string) {
	if viper.GetBool(config.Verbose) {
		printLn(info, s)
	}
}
