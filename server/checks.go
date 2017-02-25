package main

import (
	"strings"

	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/girder/errors"
	"github.com/SierraSoftworks/heimdall/models"
	"github.com/SierraSoftworks/heimdall/stores"
)

func init() {
	router.
		Methods("GET").
		Path("/api/v1/checks").
		Handler(girder.NewHandler(getChecks)).
		Name("GET /v1/checks")

	router.
		Methods("GET").
		Path("/api/v1/check/{check}").
		Handler(girder.NewHandler(getCheck)).
		Name("GET /v1/check/{check}")

	router.
		Methods("GET").
		Path("/api/v1/check/{check}/clients").
		Handler(girder.NewHandler(getCheckClients)).
		Name("GET /v1/check/{check}/clients")

	router.
		Methods("GET").
		Path("/api/v1/check/{check}/client/{client}/executions").
		Handler(girder.NewHandler(getCheckExecutions)).
		Name("GET /v1/check/{check}/client/{client}/executions")

	router.
		Methods("GET").
		Path("/api/v1/client/{client}/check/{check}/executions").
		Handler(girder.NewHandler(getCheckExecutions)).
		Name("GET /v1/client/{client}/check/{check}/executions")
}

func getChecks(c *girder.Context) (interface{}, error) {
	q := &stores.ChecksQuery{
		Status: []models.Status{},
	}

	if sq := c.Request.URL.Query().Get("status"); sq != "" {
		for _, s := range strings.Split(sq, ",") {
			q.Status = append(q.Status, models.ParseStatus(s))
		}
	}

	return stores.GetStore().GetChecks(q)
}

func getCheck(c *girder.Context) (interface{}, error) {
	ch, err := stores.GetStore().GetCheck(c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, errors.NotFound()
	}

	return ch, nil
}

func getCheckClients(c *girder.Context) (interface{}, error) {
	ch, err := stores.GetStore().GetCheckClients(c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, errors.NotFound()
	}

	return ch, nil
}

func getCheckExecutions(c *girder.Context) (interface{}, error) {
	es, err := stores.GetStore().GetCheckExecutions(c.Vars["client"], c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if es == nil {
		return nil, errors.NotFound()
	}

	return es, nil
}
