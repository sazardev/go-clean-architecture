package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// SetupMiddlewares configura todos los middlewares de la aplicación
func SetupMiddlewares(app *fiber.App) {
	// Middleware de recuperación de pánico
	app.Use(recover.New())

	// Middleware de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

	// Middleware de logging
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		TimeFormat: time.RFC3339,
		Output:     log.Writer(),
	}))

	// Middleware de validación de Content-Type para POST/PUT
	app.Use(ContentTypeMiddleware)
}

// ContentTypeMiddleware valida el Content-Type para operaciones que requieren JSON
func ContentTypeMiddleware(c *fiber.Ctx) error {
	if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
		if c.Get("Content-Type") != "application/json" && len(c.Body()) > 0 {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"error": "Content-Type must be application/json",
			})
		}
	}
	return c.Next()
}
