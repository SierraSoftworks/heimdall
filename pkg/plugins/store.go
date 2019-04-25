package plugins

import (
	"fmt"

	"github.com/SierraSoftworks/heimdall/pkg/driver"
	"github.com/SierraSoftworks/heimdall/pkg/models"
)

type Store interface {
	Driver() *driver.Driver
	Close() error

	Store(r *models.Report) error

	// AddReport(r *models.Report) error

	// GetSources(q *SourcesQuery) ([]models.Source, error)
	// GetSource(name string) (*models.Source, error)
	// GetSourceChecks(client string) ([]models.CheckDetails, error)
	// RemoveSource(name string) (*models.Source, error)

	// GetChecks(q *ChecksQuery) ([]models.Check, error)
	// GetCheck(name string) (*models.CheckDetails, error)
	// GetCheckSources(check string) ([]models.Source, error)
	// GetCheckExecutions(client, check string) ([]models.Execution, error)

	// GetAggregates() ([]models.Aggregate, error)
	// GetAggregate(name string) (*models.AggregateDetails, error)
	// GetAggregateChecks(name string) ([]models.Check, error)
	// GetAggregateSources(name string) ([]models.Source, error)
	// RemoveAggregate(name string) (*models.Aggregate, error)
}

type SourcesQuery struct {
	Tags map[string]string `json:"tags"`
}

type ChecksQuery struct {
	Status []models.Status `json:"status"`
}

var stores = map[string]func(*driver.Driver) (Store, error){}

func RegisterStore(name string, constructor func(*driver.Driver) (Store, error)) {
	stores[name] = constructor
}

func GetStore(driver *driver.Driver) (Store, error) {
	constructor, ok := stores[driver.Type]
	if !ok {
		return nil, fmt.Errorf("Unknown store driver '%s'", driver.Type)
	}

	return constructor(driver)
}
