package console

import (
	"fmt"
	"strings"

	"github.com/klimby/version/internal/config"
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

// printLn - print line.
func printLn(clr col, s string) {
	if !viper.GetBool(config.Silent) {
		var b strings.Builder

		b.WriteString(string(clr))
		b.WriteString(s)
		b.WriteString(string(reset))

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

// Verbose - print verbose message.
func Verbose(s string) {
	if viper.GetBool(config.Verbose) {
		printLn(info, s)
	}
}
