package ecrud

// Employee represents an employee record
type Employee struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	DateOfBirth string  `json:"dateOfBirth"`
	Email       string  `json:"email"`
	IsActive    *bool   `json:"isActive,omitempty"`
	Department  *string `json:"department,omitempty"`
	Role        *string `json:"role,omitempty"`
}

// EmployeeAttrs is used to create/update an employee record
// All fields are optional to avoid overwriting most recent change
// with values that are not specified by the user but is populated
// by an older read.
type EmployeeAttrs struct {
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	DateOfBirth *string `json:"dateOfBirth"`
	Email       *string `json:"email"`
	IsActive    *bool   `json:"isActive,omitempty"`
	Department  *string `json:"department,omitempty"`
	Role        *string `json:"role,omitempty"`
}
