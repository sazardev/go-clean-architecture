package router

import (
	"go-clean-architecture/internal/infrastructure/http/handler"
	"go-clean-architecture/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(app *fiber.App, employeeHandler *handler.EmployeeHandler) {
	// Configurar middlewares
	middleware.SetupMiddlewares(app)

	// Ruta de salud
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "HR API is running",
		})
	})

	// Grupo de rutas para la API
	api := app.Group("/api/v1")

	// Rutas de empleados
	employees := api.Group("/employees")
	employees.Post("/", employeeHandler.CreateEmployee)
	employees.Get("/", employeeHandler.GetAllEmployees)
	employees.Get("/:id", employeeHandler.GetEmployee)
	employees.Put("/:id", employeeHandler.UpdateEmployee)
	employees.Delete("/:id", employeeHandler.DeleteEmployee)
}
