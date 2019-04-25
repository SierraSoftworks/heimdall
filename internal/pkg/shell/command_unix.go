// +build !windows

package shell

import (
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func (c *ShellCommand) configure() {
	c.Cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func (c *ShellCommand) Exit() error {
	if c.Cmd.ProcessState != nil && !c.Cmd.ProcessState.Exited() {
		logger := log.
			WithField("command", c.Cmd.Args).
			WithField("pgid", -c.Cmd.ProcessState.Pid())

		logger.Debug("Killing process group")

		err := syscall.Kill(-c.Cmd.ProcessState.Pid(), syscall.SIGKILL)
		if err != nil {
			logger.Warning("failed to shutdown process group")
			return errors.Wrap(err, "shell: failed to shutdown process group")
		}
	} else if c.Cmd.Process != nil && c.Cmd.Process.Pid != 0 {
		logger := log.
			WithField("command", c.Cmd.Args).
			WithField("pgid", -c.Cmd.Process.Pid)

		logger.Debug("Killing process group")
		err := syscall.Kill(-c.Cmd.Process.Pid, syscall.SIGKILL)
		if err != nil && err.Error() != "no such process" {
			logger.Warning("failed to shutdown process group")
			return errors.Wrap(err, "shell: failed to shutdown process group")
		}
	}

	if c.Cmd.ProcessState != nil && !c.Cmd.ProcessState.Exited() {
		logger := log.
			WithField("command", c.Cmd.Args).
			WithField("pid", c.Cmd.Process.Pid)

		err := c.Cmd.Process.Kill()
		if err != nil {
			logger.Warning("failed to shutdown process")
			return errors.Wrap(err, "runner: failed to shutdown process")
		}
	}

	return nil
}
