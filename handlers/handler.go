package handlers

import (
	"github.com/SierraSoftworks/heimdall/models"
)

type Handler interface {
	Report(r *models.Report) error
	Shutdown()
}
