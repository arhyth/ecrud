package ecrud

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// NewHTTPServer returns an http.Handler
// that serves all the eCRUD endpoints
func NewHTTPServer(svc Service) http.Handler {
	hndlr := &httpHandler{Svc: svc}
	mux := chi.NewRouter()
	mux.NotFound(HTTPNotFound)
	mux.Route("/employees", func(r chi.Router) {
		r.Get("/", hndlr.List)
		r.Post("/", hndlr.Create)
		r.Route("/{employeeID:[0-9]+}", func(rr chi.Router) {
			rr.Get("/", hndlr.Get)
			rr.Put("/", hndlr.Update)
			rr.Delete("/", hndlr.Delete)
		})
	})

	return mux
}

// httpHandler implements net/http.HandlerFunc interfaces
// for each of the inner Service methods
type httpHandler struct {
	Svc Service
}

func (hndlr *httpHandler) List(w http.ResponseWriter, r *http.Request) {
	employees := hndlr.Svc.List()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(employees)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *httpHandler) Get(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "employeeID")
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
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(employee)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *httpHandler) Create(w http.ResponseWriter, r *http.Request) {
	var attrs EmployeeAttrs
	err := json.NewDecoder(r.Body).Decode(&attrs)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	id, err := hndlr.Svc.Create(attrs)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}

	resp := map[string]int{
		"id": id,
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *httpHandler) Update(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "employeeID")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	var attrs EmployeeAttrs
	err = json.NewDecoder(r.Body).Decode(&attrs)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	err = hndlr.Svc.Update(id, attrs)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(attrs)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *httpHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "employeeID")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	err = hndlr.Svc.Delete(id)
	if err != nil {
		WriteHTTPError(w, err)
		return
	}
	resp := map[string]int{
		"id": id,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	errnf := &ErrNotFound{}
	errbr := &ErrBadRequest{}
	if errors.As(err, errnf) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errnf)
	} else if errors.As(err, errbr) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errbr)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]string{
			"message": "server error",
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func HTTPNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"path": r.URL.Path,
	}
	json.NewEncoder(w).Encode(resp)
}
