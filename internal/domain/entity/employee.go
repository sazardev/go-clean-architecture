package entity

import (
	"time"

	"github.com/google/uuid"
)

// Employee representa un empleado en el sistema de RH
type Employee struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (Employee) TableName() string {
	return "employees"
}

// NewEmployee crea una nueva instancia de Employee
func NewEmployee(name string) *Employee {
	return &Employee{
		ID:   uuid.New(),
		Name: name,
	}
}
