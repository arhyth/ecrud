package ecrud_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/arhyth/ecrud"
)

func TestServiceStub(t *testing.T) {
	svc := ecrud.NewServiceStub(map[int]ecrud.Employee{
		3: {
			FirstName:   "David",
			LastName:    "Ebreo",
			DateOfBirth: "2001-08-15",
			Email:       "hire@me.com",
		},
	})

	t.Run("`Create` increments id", func(tt *testing.T) {
		as := assert.New(tt)
		fn, ln, dob, em, ro := "Steve", "Jobs", "1955-02-24", "steve@apple.com", "CEO"
		attrs := ecrud.EmployeeAttrs{
			FirstName:   &fn,
			LastName:    &ln,
			DateOfBirth: &dob,
			Email:       &em,
			Role:        &ro,
		}
		id, err := svc.Create(attrs)
		as.NoError(err)
		as.Greater(id, 1)
	})

	t.Run("`Create` returns error on existing email", func(tt *testing.T) {
		as := assert.New(tt)
		fn, ln, dob, em := "Linus", "Torvalds", "1969-12-28", "hire@me.com"
		attrs := ecrud.EmployeeAttrs{
			FirstName:   &fn,
			LastName:    &ln,
			DateOfBirth: &dob,
			Email:       &em,
		}
		_, err := svc.Create(attrs)
		concrete := ecrud.ErrBadRequest{}
		as.ErrorAs(err, &concrete)
		as.Contains(concrete.Fields, "email")
	})

	t.Run("`List` returns list of employees", func(tt *testing.T) {
		as := assert.New(tt)
		employees := svc.List()
		as.NotNil(employees)
		as.NotEmpty(employees)
	})

	t.Run("`Get` returns not found on non-existent record", func(tt *testing.T) {
		as := assert.New(tt)
		_, err := svc.Get(99)
		var enf ecrud.ErrNotFound
		as.ErrorAs(err, &enf)
	})
}

func TestServiceMiddleware(t *testing.T) {
	stub := ecrud.NewServiceStub(map[int]ecrud.Employee{
		1: {
			FirstName:   "David",
			LastName:    "Ebreo",
			DateOfBirth: "2001-04-15",
			Email:       "hire@me.com",
		},
	})
	svc := ecrud.NewServiceValidationMiddleware(stub)

	t.Run("validates `Create` params", func(tt *testing.T) {
		as := assert.New(tt)
		fn, ln, dob, em := "John", "Smith", "16001020", "notavalid-email"
		attrs := ecrud.EmployeeAttrs{
			FirstName:   &fn,
			LastName:    &ln,
			DateOfBirth: &dob,
			Email:       &em,
		}
		_, err := svc.Create(attrs)
		var ebr ecrud.ErrBadRequest
		as.ErrorAs(err, &ebr)
		as.Contains(ebr.Fields, "dateOfBirth")
		as.Contains(ebr.Fields, "email")
	})

	t.Run("validates `Update` params", func(tt *testing.T) {
		as := assert.New(tt)
		dob, em := "16001020", "notavalid-email"
		attrs := ecrud.EmployeeAttrs{
			DateOfBirth: &dob,
			Email:       &em,
		}
		err := svc.Update(1, attrs)
		var ebr ecrud.ErrBadRequest
		as.ErrorAs(err, &ebr)
		as.Contains(ebr.Fields, "dateOfBirth")
		as.Contains(ebr.Fields, "email")
	})
}
