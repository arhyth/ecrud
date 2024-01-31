package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/arhyth/ecrud"
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

	stub := ecrud.NewServiceStub(records)
	svc := ecrud.NewServiceValidationMiddleware(stub)
	hndlr := ecrud.NewHTTPServer(svc)

	http.ListenAndServe(":3000", hndlr)
}
