package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type ShellCommand struct {
	*exec.Cmd
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
