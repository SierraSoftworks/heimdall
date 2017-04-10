package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/utils"
	log "github.com/Sirupsen/logrus"
)

type Runner struct {
	Shell *utils.Shell
}

func NewDefaultRunner() *Runner {
	return &Runner{
		Shell: utils.NewDefaultShell(),
	}
}

func (r *Runner) ExecuteCheck(c *models.Check) *models.Execution {
	logger := log.WithFields(log.Fields{
		"check":    c.Name,
		"executor": r.Shell.Describe(),
	})

	ex := &models.Execution{
		Scheduled: time.Now(),
		Executed:  time.Now(),
		Duration:  0,
		Status:    models.StatusUnkn,
		Output:    "",
	}

	cmd, err := r.Shell.NewCommand(c.Command, []string{})
	if err == nil {
		ex.Output = err.Error()
		ex.Status = models.StatusCrit
	}

	logger.Info("Running check")
	logger.
		WithField("cmd", cmd.Describe()).
		Debug("Running check command")

	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	cmd.Stderr = out

	ex.Executed = time.Now()
	err = cmd.Start()

	if err != nil {
		ex.Output = err.Error()
		ex.Status = models.StatusCrit

		ex.Duration = time.Now().Sub(ex.Executed)

		return ex
	}

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
	case err := <-cmd.WaitOn():
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

	ex.Duration = time.Now().Sub(ex.Executed)

	return ex
}
