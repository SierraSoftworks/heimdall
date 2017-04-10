package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Shell struct {
	Command     string
	Args        []string
	Directory   string
	Environment []string
}

type ShellCommand struct {
	*exec.Cmd
}

func NewDefaultShell() *Shell {
	switch runtime.GOOS {
	case "windows":
		return NewPowerShell()
	case "linux":
		return NewBashShell()
	default:
		return NewShShell()
	}
}

func NewShShell() *Shell {
	return &Shell{
		Command: "/bin/sh",
		Args: []string{
			"-c",
		},
		Directory:   "/tmp",
		Environment: []string{},
	}
}

func NewBashShell() *Shell {
	return &Shell{
		Command: "/bin/bash",
		Args: []string{
			"-c",
		},
		Directory:   "/tmp",
		Environment: []string{},
	}
}

func NewPowerShell() *Shell {
	return &Shell{
		Command: "powershell.exe",
		Args: []string{
			"-ExecutionPolicy", "Unrestricted",
			"-NonInteractive",
			"-NoProfile",
			"-Command",
		},
		Directory:   "$TEMP",
		Environment: []string{},
	}
}

func (s *Shell) Describe() string {
	return s.Command
}

func (s *Shell) NewCommand(command string, env []string) (*ShellCommand, error) {
	cmd := exec.Command(s.Command, append(s.Args, command)...)
	cmd.Dir = os.ExpandEnv(s.Directory)
	cmd.Env = append(s.Environment, env...)

	return &ShellCommand{
		cmd,
	}, nil
}

func (c *ShellCommand) WaitOn() <-chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)
		ch <- c.Wait()
	}()

	return ch
}

func (c *ShellCommand) Describe() string {
	return strings.Join(append([]string{c.Path}, c.Args...), "")
}

func (c *ShellCommand) PushInput(input <-chan string) error {
	in, err := c.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		for item := range input {
			_, err := fmt.Fprintln(in, item)
			if err != nil {
				log.
					WithFields(log.Fields{
						"command": c.Describe(),
						"line":    item,
					}).
					WithError(err).
					Errorln("failed to write input to command stdin")
				break
			}
		}
	}()

	return nil
}

func (c *ShellCommand) ScanOutput() (<-chan string, error) {
	out, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	return c.scanPipe(out), nil
}

func (c *ShellCommand) ScanError() (<-chan string, error) {
	out, err := c.StderrPipe()
	if err != nil {
		return nil, err
	}

	return c.scanPipe(out), nil
}

func (c *ShellCommand) scanPipe(p io.ReadCloser) <-chan string {
	ch := make(chan string)

	scanner := bufio.NewScanner(p)
	go func() {
		defer close(ch)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()

	return ch
}
