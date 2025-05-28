package repository

import (
	"context"

	"go-clean-architecture/internal/domain/entity"

	"github.com/google/uuid"
)

// EmployeeRepository define el contrato para operaciones de persistencia de empleados
type EmployeeRepository interface {
	Create(ctx context.Context, employee *entity.Employee) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
	FindAll(ctx context.Context) ([]*entity.Employee, error)
	Update(ctx context.Context, employee *entity.Employee) error
	Delete(ctx context.Context, id uuid.UUID) error
}
