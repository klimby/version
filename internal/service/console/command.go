// Package console provides console service.
package console

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/klimby/version/internal/config/key"
	"github.com/spf13/viper"
)

// Cmd is a command runner.
type Cmd struct {
	commandFactory func(name string, arg ...string) runner
}

type runner interface {
	Run() error
}

// CmdArgs - command options.
type CmdArgs struct {
	CF func(name string, arg ...string) runner
}

// NewCmd creates new Cmd.
func NewCmd(args ...func(*CmdArgs)) *Cmd {
	options := CmdArgs{
		CF: func(name string, arg ...string) runner {
			cmd := exec.Command(name, arg...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			wd := viper.GetString(key.WorkDir)
			if wd != "." && wd != "" {
				cmd.Dir = wd
			}

			return cmd
		},
	}

	for _, opt := range args {
		opt(&options)
	}

	return &Cmd{
		commandFactory: options.CF,
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
