package ecrud_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/arhyth/ecrud"
)

func TestHandler(t *testing.T) {
	log := zerolog.Nop()
	seed := map[int]ecrud.Employee{
		1: {
			FirstName:   "David",
			LastName:    "Ebreo",
			DateOfBirth: "2001-04-15",
			Email:       "hire@me.com",
		},
	}
	stub := ecrud.NewServiceStub(seed, &log)
	svc := ecrud.NewServiceValidationMiddleware(stub, &log)
	hndlr := ecrud.NewHTTPServer(svc, &log)

	t.Run("`List` returns employee records", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/employees", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)
		resp := []ecrud.Employee{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.NotEmpty(resp)
	})

	t.Run("`Get` returns an employee record", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)
		resp := ecrud.Employee{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.NotEmpty(resp.FirstName)
		as.NotEmpty(resp.LastName)
		as.NotEmpty(resp.DateOfBirth)
		as.NotEmpty(resp.Email)
	})

	t.Run("`Get` returns 404 on non-existent employee record", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/employees/999", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusNotFound, w.Result().StatusCode)
		resp := ecrud.ErrNotFound{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.Equal(resp.ID, 999)
	})

	t.Run("`Create` creates an employee record", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(`{
			"firstName": "Steve",
			"lastName": "Jobs",
			"dateOfBirth": "1955-02-24",
			"email": "spj@apple.com"
		}`))
		r := httptest.NewRequest(http.MethodPost, "/employees", body)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusCreated, w.Result().StatusCode)
		resp := map[string]int{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.Greater(resp["id"], 1)
	})

	t.Run("`Create` returns an error on existing email", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		body := bytes.NewBuffer([]byte(`{
			"firstName": "Steve",
			"lastName": "Jobs",
			"dateOfBirth": "1955-02-24",
			"email": "hire@me.com"
		}`))
		r := httptest.NewRequest(http.MethodPost, "/employees", body)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusBadRequest, w.Result().StatusCode)
		resp := ecrud.ErrBadRequest{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.Contains(resp.Fields, "email")
	})

	t.Run("`Update` updates an employee record", func(tt *testing.T) {
		as := assert.New(tt)
		fn, ln, dob, em, ro := "Saul", "Goodman", "1960-10-15", "saul@good.man", "CEO"
		update := ecrud.EmployeeAttrs{
			FirstName:   &fn,
			LastName:    &ln,
			DateOfBirth: &dob,
			Email:       &em,
			Role:        &ro,
		}
		// check initial state
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)
		resp := ecrud.Employee{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.NotEqual(fn, resp.FirstName)
		as.NotEqual(ln, resp.LastName)
		as.NotEqual(dob, resp.DateOfBirth)
		as.NotEqual(em, resp.Email)

		// do update
		w = httptest.NewRecorder()
		buf := bytes.Buffer{}
		err = json.NewEncoder(&buf).Encode(update)
		as.NoError(err)
		r = httptest.NewRequest(http.MethodPut, "/employees/1", &buf)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)

		// assert updated
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)
		resp = ecrud.Employee{}
		err = json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.Equal(fn, resp.FirstName)
		as.Equal(ln, resp.LastName)
		as.Equal(dob, resp.DateOfBirth)
		as.Equal(em, resp.Email)
		as.Equal(ro, *resp.Role)
	})

	t.Run("`Update` returns 404 on non-existent employee record", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		buf := bytes.Buffer{}
		fn, ln := "Mickey", "Mouse"
		update := ecrud.EmployeeAttrs{
			FirstName: &fn,
			LastName:  &ln,
		}
		err := json.NewEncoder(&buf).Encode(update)
		as.NoError(err)
		r := httptest.NewRequest(http.MethodPut, "/employees/999", &buf)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("`Delete` deletes an employee record", func(tt *testing.T) {
		as := assert.New(tt)
		// check initial state
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)
		resp := ecrud.Employee{}
		err := json.NewDecoder(w.Result().Body).Decode(&resp)
		as.NoError(err)
		as.NotEmpty(resp.FirstName)
		as.NotEmpty(resp.LastName)

		// do delete
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodDelete, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusOK, w.Result().StatusCode)

		// assert deleted
		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusNotFound, w.Result().StatusCode)
		notfound := ecrud.ErrNotFound{}
		err = json.NewDecoder(w.Result().Body).Decode(&notfound)
		as.NoError(err)
		as.Equal(1, notfound.ID)
	})

	t.Run("`Delete` returns 404 on non-existent employee record", func(tt *testing.T) {
		as := assert.New(tt)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/employees/1", nil)
		hndlr.ServeHTTP(w, r)
		as.Equal(http.StatusNotFound, w.Result().StatusCode)
	})
}
