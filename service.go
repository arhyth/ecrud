package ecrud

import (
	"net/mail"
	"sync"
	"time"
)

// Service is complete domain interface of eCRUD
type Service interface {
	List() []Employee
	Get(int) (Employee, error)
	Create(EmployeeAttrs) (int, error)
	Update(int, EmployeeAttrs) error
	Delete(int) error
}

// ServiceStub is a "stub" implementation of Service
type ServiceStub struct {
	mtx     *sync.RWMutex
	records map[int]Employee
	seq     int
}

var _ Service = (*ServiceStub)(nil)

func NewServiceStub(records map[int]Employee) *ServiceStub {
	var seq int
	for id := range records {
		if id > seq {
			seq = id
		}
	}
	return &ServiceStub{
		mtx:     &sync.RWMutex{},
		records: records,
		seq:     seq,
	}
}

func (stub *ServiceStub) List() (employees []Employee) {
	stub.mtx.RLock()
	defer stub.mtx.RUnlock()

	for _, e := range stub.records {
		employees = append(employees, e)
	}

	return employees
}

func (stub *ServiceStub) Get(id int) (Employee, error) {
	stub.mtx.RLock()
	defer stub.mtx.RUnlock()

	e, found := stub.records[id]
	if !found {
		return e, ErrNotFound{ID: id}
	}

	return e, nil
}

func (stub *ServiceStub) Create(attrs EmployeeAttrs) (int, error) {
	stub.mtx.Lock()
	defer stub.mtx.Unlock()

	stub.seq += 1
	stub.records[stub.seq] = Employee{
		ID:          stub.seq,
		FirstName:   *attrs.FirstName,
		LastName:    *attrs.LastName,
		DateOfBirth: *attrs.DateOfBirth,
		Email:       *attrs.Email,
		IsActive:    attrs.IsActive,
		Department:  attrs.Department,
		Role:        attrs.Role,
	}

	return stub.seq, nil
}

func (stub *ServiceStub) Update(id int, attrs EmployeeAttrs) error {
	stub.mtx.Lock()
	defer stub.mtx.Unlock()

	e, found := stub.records[id]
	if !found {
		return ErrNotFound{ID: id}
	}

	if attrs.FirstName != nil {
		e.FirstName = *attrs.FirstName
	}
	if attrs.LastName != nil {
		e.LastName = *attrs.LastName
	}
	if attrs.DateOfBirth != nil {
		e.DateOfBirth = *attrs.DateOfBirth
	}
	if attrs.Email != nil {
		e.Email = *attrs.Email
	}
	if attrs.IsActive != nil {
		e.IsActive = attrs.IsActive
	}
	if attrs.Department != nil {
		e.Department = attrs.Department
	}
	if attrs.Role != nil {
		e.Role = attrs.Role
	}

	stub.records[id] = e

	return nil
}

func (stub *ServiceStub) Delete(id int) error {
	stub.mtx.Lock()
	defer stub.mtx.Unlock()

	_, found := stub.records[id]
	if !found {
		return ErrNotFound{ID: id}
	}

	delete(stub.records, id)

	return nil
}

// ServiceValidationMiddleware is a middleware that validates request parameters
// at the domain layer. This avoids having to duplicate decoding when done at
// the protocol (HTTP) layer.
type ServiceValidationMiddleware struct {
	inner Service
}

var _ Service = (*ServiceValidationMiddleware)(nil)

func NewServiceValidationMiddleware(svc Service) *ServiceValidationMiddleware {
	return &ServiceValidationMiddleware{inner: svc}
}

func (mw *ServiceValidationMiddleware) List() (employees []Employee) {
	return mw.inner.List()
}

func (mw *ServiceValidationMiddleware) Get(id int) (Employee, error) {
	return mw.Get(id)
}

func (mw *ServiceValidationMiddleware) Create(attrs EmployeeAttrs) (int, error) {
	var witherrors []string
	if attrs.FirstName == nil || len(*attrs.FirstName) <= 1 {
		witherrors = append(witherrors, "firstName")
	}
	if attrs.LastName == nil || len(*attrs.LastName) <= 1 {
		witherrors = append(witherrors, "lastName")
	}
	if attrs.DateOfBirth == nil {
		witherrors = append(witherrors, "dateOfBirth")
	} else if _, err := time.Parse(time.DateOnly, *attrs.DateOfBirth); err != nil {
		witherrors = append(witherrors, "dateOfBirth")
	}
	if attrs.Email == nil {
		witherrors = append(witherrors, "dateOfBirth")
	} else if _, err := mail.ParseAddress(*attrs.Email); err != nil {
		witherrors = append(witherrors, "email")
	}

	if attrs.Department != nil && len(*attrs.Department) <= 1 {
		witherrors = append(witherrors, "department")
	}
	if attrs.Role != nil && len(*attrs.Department) <= 1 {
		witherrors = append(witherrors, "role")
	}

	if witherrors != nil {
		return 0, ErrBadRequest{
			Fields: witherrors,
		}
	}

	return mw.inner.Create(attrs)
}

func (mw *ServiceValidationMiddleware) Update(id int, attrs EmployeeAttrs) error {
	var witherrors []string
	if attrs.FirstName != nil && len(*attrs.FirstName) <= 1 {
		witherrors = append(witherrors, "firstName")
	}
	if attrs.LastName != nil && len(*attrs.LastName) <= 1 {
		witherrors = append(witherrors, "lastName")
	}
	if attrs.DateOfBirth != nil {
		if _, err := time.Parse(time.DateOnly, *attrs.DateOfBirth); err != nil {
			witherrors = append(witherrors, "dateOfBirth")
		}
	}
	if attrs.Email != nil {
		if _, err := mail.ParseAddress(*attrs.Email); err != nil {
			witherrors = append(witherrors, "email")
		}
	}
	if attrs.Department != nil && len(*attrs.Department) <= 1 {
		witherrors = append(witherrors, "department")
	}
	if attrs.Role != nil && len(*attrs.Department) <= 1 {
		witherrors = append(witherrors, "role")
	}

	if witherrors != nil {
		return ErrBadRequest{
			Fields: witherrors,
		}
	}

	return mw.inner.Update(id, attrs)
}

func (mw *ServiceValidationMiddleware) Delete(id int) error {
	return mw.Delete(id)
}
