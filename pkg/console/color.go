package console

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	_Error   = color.New(color.FgRed)
	_Warning = color.New(color.FgYellow)
	_Success = color.New(color.FgGreen)
	_Notice  = color.New(color.FgCyan)
	_Info    = color.New(color.FgWhite)
	_Title   = color.New(color.FgCyan, color.Bold)
)

// Error - print error message.
func Error(s string) {
	fmt.Println(_Error.Println(s))
}

// Warning - print warning message.
func Warning(s string) {
	fmt.Println(_Warning.Println(s))
}

// Success - print success message.
func Success(s string) {
	fmt.Println(_Success.Println(s))
}

// Notice - print notice message.
func Notice(s string) {
	fmt.Println(_Notice.Println(s))
}

// Info - print info message.
func Info(s string) {
	fmt.Println(_Info.Println(s))
}

// Title - print title message.
func Title(s string) {
	fmt.Println()
	fmt.Println(_Title.Println(s))
}

// Colors - print all colors.
/*func Colors() {
	m := map[string]color.Attribute{
		"FgGreen":   color.FgGreen,
		"FgBlue":    color.FgBlue,
		"FgMagenta": color.FgMagenta,
		"FgCyan":    color.FgCyan,

		"FgHiBlack":   color.FgHiBlack,
		"FgHiRed":     color.FgHiRed,
		"FgHiGreen":   color.FgHiGreen,
		"FgHiYellow":  color.FgHiYellow,
		"FgHiBlue":    color.FgHiBlue,
		"FgHiMagenta": color.FgHiMagenta,
		"FgHiCyan":    color.FgHiCyan,
		"FgHiWhite":   color.FgHiWhite,
	}

	for k, v := range m {
		c := color.New(v)
		s := c.Sprintf(k)
		fmt.Println(s)
	}
}
*/
