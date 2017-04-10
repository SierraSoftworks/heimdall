package handlers

import (
	"fmt"

	"github.com/SierraSoftworks/heimdall/models"
)

// Handler represents a piece of pluggable logic which
// observes the information flowing through Heimdall
// for the purpose of integration with 3rd party services.
type Handler interface {
	Describe() string
	OnExecution(e *models.Execution) error
}

func GetHandler(driver string, url string) (Handler, error) {
	switch driver {
	case "command":
		return NewCommandHandler(url), nil
	default:
		return nil, fmt.Errorf("unknown handler driver")
	}
}
