package container

import (
	"log"

	"go-clean-architecture/internal/infrastructure/config"
	"go-clean-architecture/internal/infrastructure/database"
	"go-clean-architecture/internal/infrastructure/http/handler"
	"go-clean-architecture/internal/usecase"

	"gorm.io/gorm"
)

// Container mantiene todas las dependencias de la aplicación
type Container struct {
	Config          *config.Config
	DB              *gorm.DB
	EmployeeHandler *handler.EmployeeHandler
}

// NewContainer crea e inicializa todas las dependencias
func NewContainer() *Container {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Establecer conexión a la base de datos
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Inicializar repositorios
	employeeRepo := database.NewEmployeeRepository(db)

	// Inicializar casos de uso
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepo)

	// Inicializar handlers
	employeeHandler := handler.NewEmployeeHandler(employeeUseCase)

	return &Container{
		Config:          cfg,
		DB:              db,
		EmployeeHandler: employeeHandler,
	}
}

// Close cierra todas las conexiones del contenedor
func (c *Container) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
