package api

import (
	"net/http"

	"github.com/SierraSoftworks/heimdall/pkg/plugins"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var registrars = []func(*API) error{}

func addRegistrar(r func(*API) error) {
	registrars = append(registrars, r)
}

type API struct {
	store  plugins.Store
	router *mux.Router
}

func NewAPI(store plugins.Store) (*API, error) {
	api := &API{
		store:  store,
		router: mux.NewRouter(),
	}

	for _, registrar := range registrars {
		if err := registrar(api); err != nil {
			return nil, err
		}
	}

	return api, nil
}

func (a *API) Store() plugins.Store {
	return a.store
}

func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) ListenAndServe(address string) error {
	mux := http.NewServeMux()
	mux.Handle("/api/", a.router)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		w.Write([]byte(`{"code": 404, "error": "Not Found", "message": "The method you attempted to make use of could not be found on our system."}`))
	})

	log.WithField("address", address).Info("Starting server")
	return http.ListenAndServe(address, cors.New(cors.Options{
		Debug: false,
	}).Handler(mux))
}
