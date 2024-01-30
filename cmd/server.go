package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/arhyth/ecrud"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	f, err := os.Open("./seed.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("opening file failed")
	}

	var seed map[string][]ecrud.Employee
	err = json.NewDecoder(f).Decode(&seed)
	if err != nil {
		logger.Fatal().Err(err).Msg("decoding file failed")
	}

	records := map[int]ecrud.Employee{}
	for _, e := range seed["users"] {
		records[e.ID] = e
	}

	svc := ecrud.NewServiceStub(records)
	hndlr := ecrud.NewHandler(svc)
	mux := chi.NewRouter()
	mux.NotFound(ecrud.HTTPNotFound)
	mux.MethodFunc(http.MethodGet, "/employees", hndlr.List)
	mux.MethodFunc(http.MethodGet, "/employees/", hndlr.List)
	mux.Route("/employees/{employeeID:[0-9]+}", func(r chi.Router) {
		r.Get("/", hndlr.Get)
		r.Put("/", hndlr.Update)
		r.Delete("/", hndlr.Delete)
	})

	http.ListenAndServe(":3000", mux)
}
