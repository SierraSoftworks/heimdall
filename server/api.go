package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var router = mux.NewRouter()

func ListenAndServe(address string) error {
	mux := http.NewServeMux()
	mux.Handle("/api/", router)
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
