package console

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/klimby/version/internal/config"
	"github.com/spf13/viper"
)

var (
	_Error   = color.New(color.FgRed)
	_Warning = color.New(color.FgYellow)
	_Success = color.New(color.FgGreen)
	_Notice  = color.New(color.FgCyan)
	_Info    = color.New(color.FgWhite)
)

// printLn - print line.
func printLn(clr *color.Color, s string) {
	if !viper.GetBool(config.Silent) {
		if _, err := clr.Println(s); err != nil {
			//nolint:forbidigo
			fmt.Println(s)
		}
	}
}

// Error - print error message.
func Error(s string) {
	printLn(_Error, s)
}

// Warn - print warning message.
func Warn(s string) {
	printLn(_Warning, s)
}

// Success - print success message.
func Success(s string) {
	printLn(_Success, s)
}

// Notice - print notice message.
func Notice(s string) {
	printLn(_Notice, s)
}

// Info - print info message.
func Info(s string) {
	printLn(_Info, s)
}

// Verbose - print verbose message.
func Verbose(s string) {
	if viper.GetBool(config.Verbose) {
		printLn(_Info, s)
	}
}
