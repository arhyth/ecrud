package ecrud

import (
	"encoding/json"
	"errors"
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

type httpResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (hndlr *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	employees := hndlr.Svc.List()
	resp := httpResponse{
		Message: "list success",
		Data:    employees,
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *HTTPHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	resp := httpResponse{
		Message: "get success",
		Data:    employee,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
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
	resp := httpResponse{
		Message: "create success",
		Data: map[string]int{
			"id": id,
		},
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	resp := httpResponse{
		Message: "update success",
		Data:    attrs,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		WriteHTTPError(w, err)
	}
}

func (hndlr *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	resp := httpResponse{
		Message: "delete success",
		Data: map[string]int{
			"id": id,
		},
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
		resp := httpResponse{
			Message: errnf.Error(),
			Data:    errnf,
		}
		json.NewEncoder(w).Encode(resp)
	} else if errors.As(err, errbr) {
		w.WriteHeader(http.StatusBadRequest)
		resp := httpResponse{
			Message: errbr.Error(),
			Data:    errbr,
		}
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		resp := httpResponse{
			Message: "server error",
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func HTTPNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := httpResponse{
		Message: "path not found",
		Data: map[string]string{
			"path": r.URL.Path,
		},
	}
	json.NewEncoder(w).Encode(resp)
}
