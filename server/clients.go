package main

import (
	"github.com/SierraSoftworks/girder"
	"github.com/SierraSoftworks/heimdall/stores"
)

func init() {
	router.
		Methods("GET").
		Path("/api/v1/clients").
		Handler(girder.NewHandler(getClients)).
		Name("GET /v1/clients")

	router.
		Methods("GET").
		Path("/api/v1/client/{client}").
		Handler(girder.NewHandler(getClient)).
		Name("GET /v1/client/{client}")

	router.
		Methods("DELETE").
		Path("/api/v1/client/{client}").
		Handler(girder.NewHandler(removeClient)).
		Name("DELETE /v1/client/{client}")

	router.
		Methods("GET").
		Path("/api/v1/client/{client}/checks").
		Handler(girder.NewHandler(getClientChecks)).
		Name("GET /v1/client/{client}checks")

}

func getClients(c *girder.Context) (interface{}, error) {
	q := stores.ClientsQuery{
		Tags: map[string]string{},
	}

	for k, v := range c.Request.URL.Query() {
		q.Tags[k] = v[0]
	}

	return stores.GetStore().GetClients(&q)
}

func getClient(c *girder.Context) (interface{}, error) {
	return stores.GetStore().GetClient(c.Vars["client"])
}

func removeClient(c *girder.Context) (interface{}, error) {
	return stores.GetStore().RemoveClient(c.Vars["client"])
}

func getClientChecks(c *girder.Context) (interface{}, error) {
	return stores.GetStore().GetClientChecks(c.Vars["client"])
}
