package router

import (
	"go-clean-architecture/internal/infrastructure/http/handler"
	httpMiddleware "go-clean-architecture/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(app *fiber.App, employeeHandler *handler.EmployeeHandler, authHandler *handler.AuthHandler, authMiddleware fiber.Handler, permissionMiddleware func(string, string) fiber.Handler) {
	// Configurar middlewares generales
	httpMiddleware.SetupMiddlewares(app)

	// Ruta de salud
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "HR API is running",
		})
	})

	// Grupo de rutas para la API
	api := app.Group("/api/v1")

	// Rutas de autenticación (públicas)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Rutas protegidas
	protected := api.Group("/", authMiddleware)

	// Rutas de perfil de usuario (requiere autenticación)
	profile := protected.Group("/profile")
	profile.Get("/", authHandler.GetProfile)
	profile.Put("/", authHandler.UpdateProfile)
	profile.Put("/password", authHandler.ChangePassword)

	// Rutas de empleados (requiere autenticación)
	employees := protected.Group("/employees")
	employees.Post("/", permissionMiddleware("users", "create"), employeeHandler.CreateEmployee)
	employees.Get("/", permissionMiddleware("users", "list"), employeeHandler.GetAllEmployees)
	employees.Get("/:id", permissionMiddleware("users", "read"), employeeHandler.GetEmployee)
	employees.Put("/:id", permissionMiddleware("users", "update"), employeeHandler.UpdateEmployee)
	employees.Delete("/:id", permissionMiddleware("users", "delete"), employeeHandler.DeleteEmployee)

	// Rutas de administración de usuarios (requiere permisos especiales)
	users := protected.Group("/users", permissionMiddleware("users", "read"))
	users.Get("/", permissionMiddleware("users", "list"), authHandler.GetUsers)
	users.Get("/:id", authHandler.GetUser)
	users.Put("/:id", permissionMiddleware("users", "update"), authHandler.UpdateUser)
	users.Delete("/:id", permissionMiddleware("users", "delete"), authHandler.DeleteUser)
	users.Post("/:id/roles", permissionMiddleware("roles", "assign"), authHandler.AssignRole)
	users.Delete("/:id/roles/:roleId", permissionMiddleware("roles", "assign"), authHandler.RemoveRole)

	// Rutas de administración de roles (requiere permisos de administrador)
	roles := protected.Group("/roles", permissionMiddleware("roles", "read"))
	roles.Get("/", permissionMiddleware("roles", "list"), authHandler.GetRoles)
	roles.Post("/", permissionMiddleware("roles", "create"), authHandler.CreateRole)
	roles.Get("/:id", authHandler.GetRole)
	roles.Put("/:id", permissionMiddleware("roles", "update"), authHandler.UpdateRole)
	roles.Delete("/:id", permissionMiddleware("roles", "delete"), authHandler.DeleteRole)

	// Rutas de administración de permisos (requiere permisos de administrador)
	permissions := protected.Group("/permissions", permissionMiddleware("permissions", "read"))
	permissions.Get("/", permissionMiddleware("permissions", "list"), authHandler.GetPermissions)
	permissions.Post("/", permissionMiddleware("permissions", "create"), authHandler.CreatePermission)
	permissions.Get("/:id", authHandler.GetPermission)
	permissions.Put("/:id", permissionMiddleware("permissions", "update"), authHandler.UpdatePermission)
	permissions.Delete("/:id", permissionMiddleware("permissions", "delete"), authHandler.DeletePermission)
}
