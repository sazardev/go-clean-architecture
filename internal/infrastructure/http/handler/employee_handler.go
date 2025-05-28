package handler

import (
	"errors"

	"go-clean-architecture/internal/infrastructure/http/dto"
	"go-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// EmployeeHandler maneja las peticiones HTTP relacionadas con empleados
type EmployeeHandler struct {
	employeeUseCase *usecase.EmployeeUseCase
}

// NewEmployeeHandler crea una nueva instancia de EmployeeHandler
func NewEmployeeHandler(employeeUseCase *usecase.EmployeeUseCase) *EmployeeHandler {
	return &EmployeeHandler{
		employeeUseCase: employeeUseCase,
	}
}

// CreateEmployee maneja la creación de un nuevo empleado
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req dto.CreateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	employee, err := h.employeeUseCase.CreateEmployee(c.Context(), req.Name)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error:   "Invalid input",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse{
		Message: "Employee created successfully",
		Data:    dto.ToEmployeeResponse(employee),
	})
}

// GetEmployee maneja la obtención de un empleado por ID
func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
	}

	employee, err := h.employeeUseCase.GetEmployeeByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "Employee not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Employee retrieved successfully",
		Data:    dto.ToEmployeeResponse(employee),
	})
}

// GetAllEmployees maneja la obtención de todos los empleados
func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
	employees, err := h.employeeUseCase.GetAllEmployees(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Employees retrieved successfully",
		Data:    dto.ToEmployeeResponses(employees),
	})
}

// UpdateEmployee maneja la actualización de un empleado
func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
	}

	var req dto.UpdateEmployeeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	employee, err := h.employeeUseCase.UpdateEmployee(c.Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "Employee not found",
				Message: err.Error(),
			})
		}
		if errors.Is(err, usecase.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
				Error:   "Invalid input",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Employee updated successfully",
		Data:    dto.ToEmployeeResponse(employee),
	})
}

// DeleteEmployee maneja la eliminación de un empleado
func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error:   "Invalid employee ID",
			Message: "ID must be a valid UUID",
		})
	}

	err = h.employeeUseCase.DeleteEmployee(c.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrEmployeeNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
				Error:   "Employee not found",
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
	}

	return c.JSON(dto.SuccessResponse{
		Message: "Employee deleted successfully",
	})
}
