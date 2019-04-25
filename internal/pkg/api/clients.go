package api

import (
	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/heimdall/pkg/plugins"
)

func init() {
	addRegistrar(func(a *API) error {
		a.router.
			Methods("GET").
			Path("/api/v1/clients").
			Handler(girder.NewHandler(a.getSources)).
			Name("GET /v1/clients")

		a.router.
			Methods("GET").
			Path("/api/v1/client/{client}").
			Handler(girder.NewHandler(a.getSource)).
			Name("GET /v1/client/{client}")

		a.router.
			Methods("DELETE").
			Path("/api/v1/client/{client}").
			Handler(girder.NewHandler(a.removeSource)).
			Name("DELETE /v1/client/{client}")

		a.router.
			Methods("GET").
			Path("/api/v1/client/{client}/checks").
			Handler(girder.NewHandler(a.getSourceChecks)).
			Name("GET /v1/client/{client}checks")

		return nil
	})

}

func (a *API) getSources(c *girder.Context) (interface{}, error) {
	q := plugins.SourcesQuery{
		Tags: map[string]string{},
	}

	for k, v := range c.Request.URL.Query() {
		q.Tags[k] = v[0]
	}

	return a.store.GetSources(&q)
}

func (a *API) getSource(c *girder.Context) (interface{}, error) {
	return a.store.GetSource(c.Vars["client"])
}

func (a *API) removeSource(c *girder.Context) (interface{}, error) {
	return a.store.RemoveSource(c.Vars["client"])
}

func (a *API) getSourceChecks(c *girder.Context) (interface{}, error) {
	return a.store.GetSourceChecks(c.Vars["client"])
}
