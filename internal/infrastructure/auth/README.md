# auth/ - Autenticación y Autorización

Sistema completo de autenticación y control de acceso basado en roles (RBAC).

## Responsabilidades

- Generar y validar tokens JWT
- Implementar middleware de autenticación
- Controlar acceso basado en roles y permisos
- Gestionar sesiones de usuario
- Hashear y verificar contraseñas

## Estructura

- **`jwt/`** - Manejo de tokens JWT
- **`rbac/`** - Control de acceso basado en roles
- **`middleware/`** - Middlewares de autenticación

## Componentes

### JWT (JSON Web Tokens)
- Generación de tokens seguros
- Validación y parsing de tokens
- Refresh tokens para renovación
- Claims personalizados

### RBAC (Role-Based Access Control)
- Definición de roles y permisos
- Verificación de permisos
- Políticas de acceso flexibles
- Integración con Casbin (opcional)

### Middlewares
- Autenticación obligatoria/opcional
- Verificación de permisos específicos
- Rate limiting por usuario
- Logging de accesos

## Flujo de Autenticación

```
1. Usuario → Login (email/password)
2. Verificar credenciales
3. Generar JWT token
4. Cliente → Request con token en header
5. Middleware → Validar token
6. Extraer usuario y roles
7. Verificar permisos
8. Continuar o rechazar
```

## Configuración

Variables de entorno necesarias:
- `JWT_SECRET` - Secreto para firmar tokens
- `JWT_EXPIRATION` - Tiempo de expiración
- `JWT_REFRESH_EXPIRATION` - Tiempo de refresh token
- `BCRYPT_COST` - Costo de hasheo bcrypt

## Uso

```go
// En rutas protegidas
router.Use(auth.JWTMiddleware())
router.Use(auth.RequirePermission("employee:read"))

// En handlers
user := c.Locals("user").(domain.User)
roles := c.Locals("roles").([]string)
```
