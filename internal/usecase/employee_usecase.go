package usecase

import (
	"context"
	"errors"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/domain/repository"

	"github.com/google/uuid"
)

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrInvalidInput     = errors.New("invalid input")
)

// EmployeeUseCase maneja la lógica de negocio de empleados
type EmployeeUseCase struct {
	employeeRepo repository.EmployeeRepository
}

// NewEmployeeUseCase crea una nueva instancia de EmployeeUseCase
func NewEmployeeUseCase(employeeRepo repository.EmployeeRepository) *EmployeeUseCase {
	return &EmployeeUseCase{
		employeeRepo: employeeRepo,
	}
}

// CreateEmployee crea un nuevo empleado
func (uc *EmployeeUseCase) CreateEmployee(ctx context.Context, name string) (*entity.Employee, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}

	employee := entity.NewEmployee(name)
	if err := uc.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// GetEmployeeByID obtiene un empleado por su ID
func (uc *EmployeeUseCase) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	employee, err := uc.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}

	return employee, nil
}

// GetAllEmployees obtiene todos los empleados
func (uc *EmployeeUseCase) GetAllEmployees(ctx context.Context) ([]*entity.Employee, error) {
	return uc.employeeRepo.FindAll(ctx)
}

// UpdateEmployee actualiza un empleado existente
func (uc *EmployeeUseCase) UpdateEmployee(ctx context.Context, id uuid.UUID, name string) (*entity.Employee, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}

	employee, err := uc.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrEmployeeNotFound
	}

	employee.Name = name
	if err := uc.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

// DeleteEmployee elimina un empleado
func (uc *EmployeeUseCase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	_, err := uc.employeeRepo.FindByID(ctx, id)
	if err != nil {
		return ErrEmployeeNotFound
	}

	return uc.employeeRepo.Delete(ctx, id)
}
