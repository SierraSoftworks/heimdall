package handlers

import (
	"fmt"

	"github.com/SierraSoftworks/heimdall/models"
)

type StdinHandler struct {
	Command string
}

func (h *StdinHandler) Report(r *models.Report) error {
	// TODO: Consider using circuit breaker
	// TODO: Start process if it isn't running
	// TODO: Write report to stdin of running process
	// TODO: If write fails, queue process for killing and spawn another
	return fmt.Errorf("not yet implemented")
}

func (h *StdinHandler) Shutdown() {
	// TODO: Implement logic
}
