package ecrud

import (
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("employee not found")
)

type Service interface {
	List() []Employee
	Get(int) (Employee, error)
	Create(EmployeeAttrs) (int, error)
	Update(int, EmployeeAttrs) error
	Delete(int) error
}

type ServiceStub struct {
	mtx     *sync.RWMutex
	records map[int]Employee
	seq     int
}

var _ Service = (*ServiceStub)(nil)

func NewServiceStub(records map[int]Employee) *ServiceStub {
	return &ServiceStub{
		mtx:     &sync.RWMutex{},
		records: records,
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
		return e, ErrNotFound
	}

	return e, nil
}

func (stub *ServiceStub) Create(attrs EmployeeAttrs) (int, error) {
	stub.mtx.Lock()
	defer stub.mtx.Unlock()

	seq := stub.seq + 1
	stub.records[seq] = Employee{
		ID:          seq,
		FirstName:   *attrs.FirstName,
		LastName:    *attrs.LastName,
		DateOfBirth: *attrs.DateOfBirth,
		Email:       *attrs.Email,
		IsActive:    attrs.IsActive,
		Department:  attrs.Department,
		Role:        attrs.Role,
	}

	return seq, nil
}

func (stub *ServiceStub) Update(id int, attrs EmployeeAttrs) error {
	stub.mtx.Lock()
	defer stub.mtx.Unlock()

	e, found := stub.records[id]
	if !found {
		return ErrNotFound
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
		return ErrNotFound
	}

	delete(stub.records, id)

	return nil
}
