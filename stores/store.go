package stores

import (
	"github.com/SierraSoftworks/heimdall/models"
)

var activeStore Store

func GetStore() Store {
	if activeStore == nil {
		activeStore = NewMemory()
	}

	return activeStore
}

type Store interface {
	AddReport(r *models.Report) error

	GetClients(q *ClientsQuery) ([]models.Client, error)
	GetClient(name string) (*models.ClientDetails, error)
	GetClientChecks(client string) ([]models.CheckDetails, error)
	RemoveClient(name string) (*models.Client, error)

	GetChecks(q *ChecksQuery) ([]models.Check, error)
	GetCheck(name string) (*models.CheckDetails, error)
	GetCheckClients(check string) ([]models.ClientDetails, error)
	GetCheckExecutions(client, check string) ([]models.Execution, error)

	GetAggregates() ([]models.Aggregate, error)
	GetAggregate(name string) (*models.AggregateDetails, error)
	GetAggregateChecks(name string) ([]models.Check, error)
	GetAggregateClients(name string) ([]models.Client, error)
	RemoveAggregate(name string) (*models.Aggregate, error)
}

type ClientsQuery struct {
	Tags map[string]string `json:"tags"`
}

type ChecksQuery struct {
	Status []models.Status `json:"status"`
}
