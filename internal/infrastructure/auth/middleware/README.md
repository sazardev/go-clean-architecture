# middleware/ - Middlewares de Autenticación

Middlewares de Fiber para manejar autenticación y autorización en las rutas HTTP.

## Responsabilidades

- Validar tokens JWT en requests
- Extraer información de usuario de tokens
- Verificar permisos específicos por ruta
- Manejar errores de autenticación
- Implementar rate limiting por usuario

## Middlewares Disponibles

### Authentication Middleware
- **`auth.go`** - Middleware principal de autenticación JWT
- **`optional_auth.go`** - Autenticación opcional
- **`refresh.go`** - Middleware para refresh tokens

### Authorization Middleware  
- **`permission.go`** - Verificación de permisos específicos
- **`role.go`** - Verificación de roles requeridos
- **`owner.go`** - Verificación de ownership de recursos

### Security Middleware
- **`rate_limit.go`** - Rate limiting por usuario
- **`cors.go`** - Configuración CORS
- **`security_headers.go`** - Headers de seguridad

## Implementación

### JWT Middleware
```go
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // 1. Extraer token del header Authorization
        token := extractToken(c)
        if token == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Authorization token required",
            })
        }
        
        // 2. Validar token
        claims, err := jwtService.ValidateToken(token)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        
        // 3. Almacenar en contexto
        c.Locals("user_id", claims.UserID)
        c.Locals("username", claims.Username)
        c.Locals("roles", claims.Roles)
        
        return c.Next()
    }
}
```

### Permission Middleware
```go
func RequirePermission(resource, action string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        roles := c.Locals("roles").([]string)
        
        if !rbacService.HasPermission(roles, resource, action) {
            return c.Status(403).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        
        return c.Next()
    }
}
```

## Uso en Rutas

```go
// Rutas que requieren autenticación
api := router.Group("/api/v1")
api.Use(middleware.JWTMiddleware())

// Rutas con permisos específicos
employees := api.Group("/employees")
employees.Post("/", 
    middleware.RequirePermission("employee", "create"),
    handler.CreateEmployee,
)
employees.Get("/", 
    middleware.RequirePermission("employee", "read"),
    handler.ListEmployees,
)
employees.Put("/:id", 
    middleware.RequirePermission("employee", "update"),
    handler.UpdateEmployee,
)
employees.Delete("/:id", 
    middleware.RequirePermission("employee", "delete"),
    handler.DeleteEmployee,
)

// Rutas con autenticación opcional
public := router.Group("/public")
public.Use(middleware.OptionalAuth())
public.Get("/info", handler.PublicInfo)
```

## Configuración

- Headers de autorización: `Authorization: Bearer <token>`
- Rate limiting configurable por endpoint
- CORS configurado para dominios específicos
- Headers de seguridad automaticos
