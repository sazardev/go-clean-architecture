package dto

import (
	"time"

	"go-clean-architecture/internal/domain/entity"

	"github.com/google/uuid"
)

// CreateEmployeeRequest representa la petición para crear un empleado
type CreateEmployeeRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// UpdateEmployeeRequest representa la petición para actualizar un empleado
type UpdateEmployeeRequest struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// EmployeeResponse representa la respuesta de un empleado
type EmployeeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse representa una respuesta exitosa genérica
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToEmployeeResponse convierte una entidad Employee a EmployeeResponse
func ToEmployeeResponse(employee *entity.Employee) *EmployeeResponse {
	return &EmployeeResponse{
		ID:        employee.ID,
		Name:      employee.Name,
		CreatedAt: employee.CreatedAt,
		UpdatedAt: employee.UpdatedAt,
	}
}

// ToEmployeeResponses convierte una slice de entidades Employee a EmployeeResponse
func ToEmployeeResponses(employees []*entity.Employee) []*EmployeeResponse {
	responses := make([]*EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = ToEmployeeResponse(employee)
	}
	return responses
}
