package usecase_test

import (
	"context"
	"errors"
	"testing"

	"go-clean-architecture/internal/domain/entity"
	"go-clean-architecture/internal/usecase"

	"github.com/google/uuid"
)

// mockEmployeeRepository es un mock del repositorio de empleados para testing
type mockEmployeeRepository struct {
	employees map[uuid.UUID]*entity.Employee
	createErr error
	findErr   error
	updateErr error
	deleteErr error
}

func newMockEmployeeRepository() *mockEmployeeRepository {
	return &mockEmployeeRepository{
		employees: make(map[uuid.UUID]*entity.Employee),
	}
}

func (m *mockEmployeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.employees[employee.ID] = employee
	return nil
}

func (m *mockEmployeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	employee, exists := m.employees[id]
	if !exists {
		return nil, errors.New("employee not found")
	}
	return employee, nil
}

func (m *mockEmployeeRepository) FindAll(ctx context.Context) ([]*entity.Employee, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	employees := make([]*entity.Employee, 0, len(m.employees))
	for _, employee := range m.employees {
		employees = append(employees, employee)
	}
	return employees, nil
}

func (m *mockEmployeeRepository) Update(ctx context.Context, employee *entity.Employee) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.employees[employee.ID] = employee
	return nil
}

func (m *mockEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.employees, id)
	return nil
}

func TestEmployeeUseCase_CreateEmployee(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		createErr   error
		expectError bool
		errorType   error
	}{
		{
			name:        "successful creation",
			inputName:   "John Doe",
			expectError: false,
		},
		{
			name:        "empty name should return error",
			inputName:   "",
			expectError: true,
			errorType:   usecase.ErrInvalidInput,
		},
		{
			name:        "repository error",
			inputName:   "John Doe",
			createErr:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := newMockEmployeeRepository()
			mockRepo.createErr = tt.createErr
			uc := usecase.NewEmployeeUseCase(mockRepo)

			employee, err := uc.CreateEmployee(context.Background(), tt.inputName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if employee.Name != tt.inputName {
				t.Errorf("expected name %s, got %s", tt.inputName, employee.Name)
			}
		})
	}
}

func TestEmployeeUseCase_GetEmployeeByID(t *testing.T) {
	mockRepo := newMockEmployeeRepository()
	uc := usecase.NewEmployeeUseCase(mockRepo)

	// Crear un empleado de prueba
	employee := entity.NewEmployee("John Doe")
	mockRepo.employees[employee.ID] = employee

	tests := []struct {
		name        string
		id          uuid.UUID
		findErr     error
		expectError bool
	}{
		{
			name:        "successful retrieval",
			id:          employee.ID,
			expectError: false,
		},
		{
			name:        "employee not found",
			id:          uuid.New(),
			expectError: true,
		},
		{
			name:        "repository error",
			id:          employee.ID,
			findErr:     errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.findErr = tt.findErr
			result, err := uc.GetEmployeeByID(context.Background(), tt.id)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ID != tt.id {
				t.Errorf("expected ID %v, got %v", tt.id, result.ID)
			}
		})
	}
}

func TestEmployeeUseCase_UpdateEmployee(t *testing.T) {
	mockRepo := newMockEmployeeRepository()
	uc := usecase.NewEmployeeUseCase(mockRepo)

	// Crear un empleado de prueba
	employee := entity.NewEmployee("John Doe")
	mockRepo.employees[employee.ID] = employee

	tests := []struct {
		name        string
		id          uuid.UUID
		newName     string
		expectError bool
		errorType   error
	}{
		{
			name:        "successful update",
			id:          employee.ID,
			newName:     "Jane Doe",
			expectError: false,
		},
		{
			name:        "empty name should return error",
			id:          employee.ID,
			newName:     "",
			expectError: true,
			errorType:   usecase.ErrInvalidInput,
		},
		{
			name:        "employee not found",
			id:          uuid.New(),
			newName:     "Jane Doe",
			expectError: true,
			errorType:   usecase.ErrEmployeeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := uc.UpdateEmployee(context.Background(), tt.id, tt.newName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !errors.Is(err, tt.errorType) {
					t.Errorf("expected error %v, got %v", tt.errorType, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Name != tt.newName {
				t.Errorf("expected name %s, got %s", tt.newName, result.Name)
			}
		})
	}
}
