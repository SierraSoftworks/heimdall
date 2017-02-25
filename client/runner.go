package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"strings"

	"github.com/SierraSoftworks/heimdall/models"
	log "github.com/Sirupsen/logrus"
)

type Runner struct {
	Executor     string
	ExecutorArgs []string
	Dir          string
}

func NewDefaultRunner() *Runner {
	switch runtime.GOOS {
	case "windows":
		return NewPowerShellRunner()
	case "linux":
		return NewBashRunner()
	default:
		return NewShRunner()
	}
}

func NewBashRunner() *Runner {
	return &Runner{
		Executor: "/bin/bash",
		ExecutorArgs: []string{
			"-c",
		},
		Dir: "/tmp",
	}
}

func NewShRunner() *Runner {
	return &Runner{
		Executor: "/bin/sh",
		ExecutorArgs: []string{
			"-c",
		},
		Dir: "/tmp",
	}
}

func NewPowerShellRunner() *Runner {
	return &Runner{
		Executor: "powershell.exe",
		ExecutorArgs: []string{
			"-ExecutionPolicy", "Unrestricted",
			"-NonInteractive",
			"-NoProfile",
			"-Command",
		},
		Dir: "$TEMP",
	}
}

func (r *Runner) Execute(c *models.Check) *models.Execution {
	logger := log.WithFields(log.Fields{
		"check":    c.Name,
		"executor": r.Executor,
	})

	ex := &models.Execution{
		Scheduled: time.Now(),
		Executed:  time.Now(),
		Duration:  0,
		Status:    models.StatusUnkn,
		Output:    "",
	}

	cmd := exec.Command(r.Executor, append(r.ExecutorArgs, c.Command)...)
	cmd.Dir = os.ExpandEnv(r.Dir)

	logger.Info("Running check")
	logger.
		WithField("cmd", strings.Join(append([]string{cmd.Path}, cmd.Args...), " ")).
		Debug("Running check command")

	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	cmd.Stderr = out

	ex.Executed = time.Now()
	err := cmd.Start()
	if err != nil {
		ex.Output = err.Error()
		ex.Status = models.StatusCrit
	} else {
		d := c.Timeout
		if d == 0 {
			d = time.Minute
		}

		t := time.NewTimer(d)

		select {
		case <-t.C:
			ex.Output = fmt.Sprintf("%s\nTimeout Expired!", out.String())
			ex.Status = models.StatusCrit
			err := cmd.Process.Kill()
			logger.
				WithField("timeout", d).
				Warn("Check execution timed out")

			if err != nil {
				logger.
					WithError(err).
					WithField("timeout", d).
					Error("Failed to kill check process after timeout")
			}
		case err := <-r.waitForProcess(cmd):
			t.Stop() // Reclaim the timer's handles

			ex.Output = out.String()

			if err != nil {
				ex.Output = fmt.Sprintf("%s\n%s", ex.Output, err.Error())
				ex.Status = models.StatusCrit
				break
			}

			type ExitStatus interface {
				ExitStatus() int
			}

			status, ok := cmd.ProcessState.Sys().(ExitStatus)
			if ok {
				ex.Status = models.Status(status.ExitStatus())
			} else if cmd.ProcessState.Success() {
				ex.Status = models.StatusOkay
			} else {
				ex.Status = models.StatusCrit
			}

			logger.
				WithField("status", ex.Status.String()).
				Debug("Check execution completed")
		}
	}

	ex.Duration = time.Now().Sub(ex.Executed)

	return ex
}

func (r *Runner) waitForProcess(cmd *exec.Cmd) <-chan error {
	c := make(chan error)

	go func() {
		defer close(c)
		c <- cmd.Wait()
	}()

	return c
}
