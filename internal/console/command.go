package console

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/klimby/version/pkg/convert"
)

// Cmd is a command runner.
type Cmd struct{}

// NewCmd creates new Cmd.
func NewCmd() *Cmd {
	return &Cmd{}
}

// Run runs command.
func (c *Cmd) Run(name string, arg ...string) error {
	// split name and args.
	spl := strings.Split(name, " ")
	if len(spl) == 0 || spl[0] == "" {
		return fmt.Errorf("invalid command: %s", name)
	}

	var a []string

	copy(a, arg)

	if len(spl) > 1 {
		a = append(spl[1:], a...)
	}

	cmd := exec.Command(spl[0], a...)

	errOut := &stdErrOutput{}

	n := commandString(name, arg...)

	// pipe the commands output to the applications.
	// standard output.
	// cmd.Stdout = stdOutOutput{}.
	cmd.Stdout = stdOutput{}
	cmd.Stderr = errOut

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run command %s: %w", n, err)
	}

	if errOut.isError {
		return fmt.Errorf("command %s run error", n)
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

type stdErrOutput struct {
	isError bool
}

func (s *stdErrOutput) Write(p []byte) (int, error) {
	Error(convert.B2S(p))
	s.isError = true

	return len(p), nil
}

type stdOutput struct{}

func (stdOutput) Write(p []byte) (int, error) {
	Info(convert.B2S(p))

	return len(p), nil
}
