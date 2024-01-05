package shell

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/klimby/version/internal/console"
	"github.com/klimby/version/pkg/convert"
)

// Cmd run cmd and return error.
func Cmd(name string, arg ...string) error {
	// split name and args.
	spl := strings.Split(name, " ")
	if len(spl) == 0 || spl[0] == "" {
		return fmt.Errorf("invalid command: %s", name)
	}

	n := spl[0]

	a := append(spl[1:], arg...)

	cmd := exec.Command(n, a...)

	// pipe the commands output to the applications.
	// standard output.
	// cmd.Stdout = stdOutOutput{}.
	cmd.Stdout = stdOutput{}
	cmd.Stderr = stdErrOutput{}

	if err := cmd.Run(); err != nil {
		n := commandString(name, arg...)
		return fmt.Errorf("could not run command %s: %w", n, err)
	}

	return nil
}

// commandString return command string.
func commandString(name string, arg ...string) string {
	var b strings.Builder

	b.WriteString(name)

	if len(arg) > 0 {
		for i := range arg {
			b.WriteString(" " + arg[i])
		}
	}

	return b.String()
}

type stdErrOutput struct{}

func (stdErrOutput) Write(p []byte) (int, error) {
	console.Error(convert.B2S(p))

	return len(p), nil
}

type stdOutput struct{}

func (stdOutput) Write(p []byte) (int, error) {
	console.Info(convert.B2S(p))

	return len(p), nil
}
