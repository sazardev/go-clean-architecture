package container

import (
	"log"
	"time"

	"go-clean-architecture/internal/infrastructure/auth"
	"go-clean-architecture/internal/infrastructure/auth/jwt"
	"go-clean-architecture/internal/infrastructure/auth/middleware"
	"go-clean-architecture/internal/infrastructure/auth/rbac"
	"go-clean-architecture/internal/infrastructure/config"
	"go-clean-architecture/internal/infrastructure/database"
	"go-clean-architecture/internal/infrastructure/http/handler"
	"go-clean-architecture/internal/infrastructure/repository"
	"go-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Container mantiene todas las dependencias de la aplicación
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	// Auth components
	TokenService         *jwt.TokenService
	PolicyManager        *rbac.PolicyManager
	AuthService          *auth.AuthService
	AuthMiddleware       fiber.Handler
	PermissionMiddleware func(string, string) fiber.Handler

	// Handlers
	EmployeeHandler *handler.EmployeeHandler
	AuthHandler     *handler.AuthHandler

	// Use cases
	UserUseCase       *usecase.UserUseCase
	RoleUseCase       *usecase.RoleUseCase
	PermissionUseCase *usecase.PermissionUseCase
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
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// Inicializar servicios de autenticación
	tokenService := jwt.NewTokenService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.ExpirationHours)*time.Hour,
		cfg.JWT.Issuer,
	)
	// Inicializar policy manager
	enforcer, err := rbac.NewEnforcer(db, cfg.Casbin.ModelPath)
	if err != nil {
		log.Fatalf("Failed to create RBAC enforcer: %v", err)
	}
	policyManager := rbac.NewPolicyManager(enforcer)

	authService := auth.NewAuthService(userRepo, roleRepo, tokenService, policyManager)

	// Inicializar middlewares
	authMiddleware := middleware.AuthMiddleware(tokenService)
	permissionMiddleware := func(resource, action string) fiber.Handler {
		return middleware.RequirePermission(policyManager, resource, action)
	}

	// Inicializar casos de uso
	employeeUseCase := usecase.NewEmployeeUseCase(employeeRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, roleRepo, permissionRepo, authService, policyManager)
	roleUseCase := usecase.NewRoleUseCase(roleRepo, permissionRepo, userRepo, policyManager)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)

	// Inicializar handlers
	employeeHandler := handler.NewEmployeeHandler(employeeUseCase)
	authHandler := handler.NewAuthHandler(authService)

	return &Container{
		Config:               cfg,
		DB:                   db,
		TokenService:         tokenService,
		PolicyManager:        policyManager,
		AuthService:          authService,
		AuthMiddleware:       authMiddleware,
		PermissionMiddleware: permissionMiddleware,
		EmployeeHandler:      employeeHandler,
		AuthHandler:          authHandler,
		UserUseCase:          userUseCase,
		RoleUseCase:          roleUseCase,
		PermissionUseCase:    permissionUseCase,
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
