package api

import (
	"strings"

	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/girder/errors"
	"github.com/SierraSoftworks/heimdall/pkg/models"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
)

func init() {
	addRegistrar(func(a *API) error {
		a.router.
			Methods("GET").
			Path("/api/v1/checks").
			Handler(girder.NewHandler(a.getChecks)).
			Name("GET /v1/checks")

		a.router.
			Methods("GET").
			Path("/api/v1/check/{check}").
			Handler(girder.NewHandler(a.getCheck)).
			Name("GET /v1/check/{check}")

		a.router.
			Methods("GET").
			Path("/api/v1/check/{check}/clients").
			Handler(girder.NewHandler(a.getCheckSources)).
			Name("GET /v1/check/{check}/clients")

		a.router.
			Methods("GET").
			Path("/api/v1/check/{check}/client/{client}/executions").
			Handler(girder.NewHandler(a.getCheckExecutions)).
			Name("GET /v1/check/{check}/client/{client}/executions")

		a.router.
			Methods("GET").
			Path("/api/v1/client/{client}/check/{check}/executions").
			Handler(girder.NewHandler(a.getCheckExecutions)).
			Name("GET /v1/client/{client}/check/{check}/executions")

		return nil
	})
}

func (a *API) getChecks(c *girder.Context) (interface{}, error) {
	q := &plugins.ChecksQuery{
		Status: []models.Status{},
	}

	if sq := c.Request.URL.Query().Get("status"); sq != "" {
		for _, s := range strings.Split(sq, ",") {
			q.Status = append(q.Status, models.ParseStatus(s))
		}
	}

	return a.store.GetChecks(q)
}

func (a *API) getCheck(c *girder.Context) (interface{}, error) {
	ch, err := a.store.GetCheck(c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, errors.NotFound()
	}

	return ch, nil
}

func (a *API) getCheckSources(c *girder.Context) (interface{}, error) {
	ch, err := a.store.GetCheckSources(c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if ch == nil {
		return nil, errors.NotFound()
	}

	return ch, nil
}

func (a *API) getCheckExecutions(c *girder.Context) (interface{}, error) {
	es, err := a.store.GetCheckExecutions(c.Vars["client"], c.Vars["check"])
	if err != nil {
		return nil, err
	}

	if es == nil {
		return nil, errors.NotFound()
	}

	return es, nil
}
