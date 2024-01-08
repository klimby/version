package console

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/klimby/version/pkg/convert"
)

// Cmd is a command runner.
type Cmd struct {
	commandFactory func(name string, arg ...string) runner
	stdout         io.Writer
	stderr         errorWriter
}

type errorWriter interface {
	io.Writer
	isError() bool
}

type runner interface {
	Run() error
}

// NewCmd creates new Cmd.
func NewCmd() *Cmd {
	sO := &stdOutput{}
	eO := &stdErrOutput{}
	return &Cmd{
		commandFactory: func(name string, arg ...string) runner {
			cmd := exec.Command(name, arg...)
			cmd.Stdout = sO
			cmd.Stderr = eO

			return cmd
		},
		stdout: sO,
		stderr: eO,
	}
}

// Run runs command.
func (c *Cmd) Run(name string, arg ...string) error {
	nm, a, err := normalizeArgs(name, arg...)
	if err != nil {
		return err
	}

	cmd := c.commandFactory(nm, a...)

	n := commandString(name, arg...)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not run command %s: %w", n, err)
	}

	if c.stderr.isError() {
		return fmt.Errorf("command %s run error", n)
	}

	return nil
}

// normalizeArgs transform arguments.
func normalizeArgs(name string, arg ...string) (string, []string, error) {
	// split name and args.
	spl := strings.Split(name, " ")
	if len(spl) == 0 || spl[0] == "" {
		return "", nil, fmt.Errorf("invalid command: %s", name)
	}

	var a []string

	if len(arg) > 0 {
		a = append(a, arg...)
	}

	if len(spl) > 1 {
		a = append(spl[1:], a...)
	}

	return spl[0], a, nil
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
	isErr bool
}

func (s *stdErrOutput) isError() bool {
	return s.isErr
}

func (s *stdErrOutput) Write(p []byte) (int, error) {
	Error(convert.B2S(p))
	s.isErr = true

	return len(p), nil
}

type stdOutput struct{}

func (stdOutput) Write(p []byte) (int, error) {
	Info(convert.B2S(p))

	return len(p), nil
}
