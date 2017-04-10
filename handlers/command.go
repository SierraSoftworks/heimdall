package handlers

import (
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/utils"
)

type CommandHandler struct {
	URL   string
	Shell *utils.Shell

	cmd   *utils.ShellCommand
	execs chan string
	close chan struct{}
}

func NewCommandHandler(url string) *CommandHandler {
	return &CommandHandler{
		URL:   url,
		Shell: utils.NewDefaultShell(),

		close: make(chan struct{}),
	}
}

func (r *CommandHandler) Describe() string {
	return fmt.Sprintf("command://%s", r.URL)
}

func (r *CommandHandler) Start() error {
	cmd, err := r.Shell.NewCommand(r.URL, []string{})
	if err != nil {
		return err
	}

	r.cmd = cmd

	// TODO: Needs refactoring, likely to not work as expected
	go func() {
		for {
			select {
			case <-r.close:
				if cmd.Process != nil {
					cmd.Process.Kill()
				}
				cmd = nil
				return
			case <-cmd.WaitOn():
				cmd.Run()
				cmd.PushInput(r.execs)
			}
		}
	}()

	return cmd.PushInput(r.execs)
}

func (r *CommandHandler) Stop() error {
	r.close <- struct{}{}
	return nil
}

func (r *CommandHandler) OnExecution(exec *models.Execution) error {
	b := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(b).Encode(exec); err != nil {
		return err
	}

	r.execs <- b.String()
	return nil
}
