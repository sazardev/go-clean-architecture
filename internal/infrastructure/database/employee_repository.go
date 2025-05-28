package database

import (
	"context"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// employeeRepository implementa repository.EmployeeRepository
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository crea una nueva instancia de employeeRepository
func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

// Create crea un nuevo empleado en la base de datos
func (r *employeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	return r.db.WithContext(ctx).Create(employee).Error
}

// FindByID busca un empleado por su ID
func (r *employeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.WithContext(ctx).First(&employee, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindAll obtiene todos los empleados
func (r *employeeRepository) FindAll(ctx context.Context) ([]*entity.Employee, error) {
	var employees []*entity.Employee
	err := r.db.WithContext(ctx).Find(&employees).Error
	return employees, err
}

// Update actualiza un empleado existente
func (r *employeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	return r.db.WithContext(ctx).Save(employee).Error
}

// Delete elimina un empleado por su ID
func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Employee{}, "id = ?", id).Error
}
