package api

import (
	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/girder/errors"
)

func init() {
	addRegistrar(func(a *API) error {
		a.router.
			Methods("GET").
			Path("/api/v1/aggregates").
			Handler(girder.NewHandler(a.getAggregates)).
			Name("GET /v1/aggregates")

		a.router.
			Methods("GET").
			Path("/api/v1/aggregate/{aggregate}").
			Handler(girder.NewHandler(a.getAggregate)).
			Name("GET /v1/aggregate/{aggregate}")

		a.router.
			Methods("DELETE").
			Path("/api/v1/aggregate/{aggregate}").
			Handler(girder.NewHandler(a.removeAggregate)).
			Name("GET /v1/aggregate/{aggregate}")

		a.router.
			Methods("GET").
			Path("/api/v1/aggregate/{aggregate}/clients").
			Handler(girder.NewHandler(a.getAggregateSources)).
			Name("GET /v1/aggregate/{aggregate}/clients")

		a.router.
			Methods("GET").
			Path("/api/v1/aggregate/{aggregate}/checks").
			Handler(girder.NewHandler(a.getAggregateChecks)).
			Name("GET /v1/aggregate/{aggregate}/checks")

		return nil
	})
}

func (a *API) getAggregates(c *girder.Context) (interface{}, error) {
	return a.store.GetAggregates()
}

func (a *API) getAggregate(c *girder.Context) (interface{}, error) {
	ag, err := a.store.GetAggregate(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ag == nil {
		return nil, errors.NotFound()
	}

	return ag, nil
}

func (a *API) removeAggregate(c *girder.Context) (interface{}, error) {
	ag, err := a.store.RemoveAggregate(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ag == nil {
		return nil, errors.NotFound()
	}

	return ag, nil
}

func (a *API) getAggregateSources(c *girder.Context) (interface{}, error) {
	ac, err := a.store.GetAggregateSources(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ac == nil {
		return nil, errors.NotFound()
	}

	return ac, nil
}

func (a *API) getAggregateChecks(c *girder.Context) (interface{}, error) {
	ac, err := a.store.GetAggregateChecks(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ac == nil {
		return nil, errors.NotFound()
	}

	return ac, nil
}
