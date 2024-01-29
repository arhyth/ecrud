package ecrud

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// HTTPHandler implements net/http.HandlerFunc interfaces
// for each of the inner Service methods
type HTTPHandler struct {
	Svc Service
}

func NewHandler(svc Service) *HTTPHandler {
	return &HTTPHandler{
		Svc: svc,
	}
}

func (hndlr *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	employees := hndlr.Svc.List()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(employees)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "employeeID")
	// TODO: write a middleware to do this conversion and validate integer
	id, err := strconv.Atoi(idstr)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	employee, err := hndlr.Svc.Get(id)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(employee)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	if errors.Is(err, ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "employee record not found"}`))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
	}
}
