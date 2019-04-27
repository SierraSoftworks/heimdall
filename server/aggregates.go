package main

import (
	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/girder/errors"
	"github.com/SierraSoftworks/heimdall/stores"
)

func init() {
	router.
		Methods("GET").
		Path("/api/v1/aggregates").
		Handler(girder.NewHandler(getAggregates)).
		Name("GET /v1/aggregates")

	router.
		Methods("GET").
		Path("/api/v1/aggregate/{aggregate}").
		Handler(girder.NewHandler(getAggregate)).
		Name("GET /v1/aggregate/{aggregate}")

	router.
		Methods("DELETE").
		Path("/api/v1/aggregate/{aggregate}").
		Handler(girder.NewHandler(removeAggregate)).
		Name("GET /v1/aggregate/{aggregate}")

	router.
		Methods("GET").
		Path("/api/v1/aggregate/{aggregate}/clients").
		Handler(girder.NewHandler(getAggregateClients)).
		Name("GET /v1/aggregate/{aggregate}/clients")

	router.
		Methods("GET").
		Path("/api/v1/aggregate/{aggregate}/checks").
		Handler(girder.NewHandler(getAggregateChecks)).
		Name("GET /v1/aggregate/{aggregate}/checks")
}

func getAggregates(c *girder.Context) (interface{}, error) {
	return stores.GetStore().GetAggregates()
}

func getAggregate(c *girder.Context) (interface{}, error) {
	a, err := stores.GetStore().GetAggregate(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if a == nil {
		return nil, errors.NotFound()
	}

	return a, nil
}

func removeAggregate(c *girder.Context) (interface{}, error) {
	a, err := stores.GetStore().RemoveAggregate(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if a == nil {
		return nil, errors.NotFound()
	}

	return a, nil
}

func getAggregateClients(c *girder.Context) (interface{}, error) {
	ac, err := stores.GetStore().GetAggregateClients(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ac == nil {
		return nil, errors.NotFound()
	}

	return ac, nil
}

func getAggregateChecks(c *girder.Context) (interface{}, error) {
	ac, err := stores.GetStore().GetAggregateChecks(c.Vars["aggregate"])
	if err != nil {
		return nil, err
	}

	if ac == nil {
		return nil, errors.NotFound()
	}

	return ac, nil
}
