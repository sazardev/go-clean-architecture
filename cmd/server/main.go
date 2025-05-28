package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-clean-architecture/internal/infrastructure/container"
	"go-clean-architecture/internal/infrastructure/http/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inicializar contenedor de dependencias
	container := container.NewContainer()
	defer func() {
		if err := container.Close(); err != nil {
			log.Printf("Error closing container: %v", err)
		}
	}()

	// Crear aplicación Fiber
	app := fiber.New(fiber.Config{
		AppName:      "HR API v1.0",
		ServerHeader: "HR-API",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		},
	})

	// Configurar rutas
	router.SetupRoutes(app, container.EmployeeHandler)

	// Configurar shutdown graceful
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	// Iniciar servidor
	port := fmt.Sprintf(":%s", container.Config.Server.Port)
	log.Printf("🚀 HR API Server starting on port %s", container.Config.Server.Port)
	log.Printf("📚 Health check available at: http://localhost%s/health", port)
	log.Printf("🔗 API documentation: http://localhost%s/api/v1/employees", port)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
