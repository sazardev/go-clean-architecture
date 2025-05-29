# jwt/ - Gestión de Tokens JWT

Implementación completa para generación, validación y manejo de tokens JWT.

## Responsabilidades

- Generar tokens JWT seguros
- Validar y parsear tokens
- Manejar refresh tokens
- Configurar claims personalizados
- Gestionar expiración y renovación

## Archivos

- **`token.go`** - Generación y validación de tokens
- **`claims.go`** - Definición de claims personalizados
- **`service.go`** - Servicio principal de JWT
- **`config.go`** - Configuración de JWT

## Claims Personalizados

```go
type CustomClaims struct {
    UserID   uuid.UUID `json:"user_id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Roles    []string  `json:"roles"`
    jwt.RegisteredClaims
}
```

## Funcionalidades

### Generación de Tokens
```go
func GenerateToken(user *domain.User) (string, error)
func GenerateRefreshToken(userID uuid.UUID) (string, error)
```

### Validación
```go
func ValidateToken(tokenString string) (*CustomClaims, error)
func ValidateRefreshToken(tokenString string) (uuid.UUID, error)
```

### Configuración
- Algoritmo de firma: HS256
- Tiempo de vida configurable
- Secreto seguro desde variables de entorno
- Issuer y audience personalizables

## Seguridad

- Secretos aleatorios y seguros
- Tokens con tiempo de vida limitado
- Refresh tokens para renovación segura
- Validación estricta de claims
- Protección contra ataques de timing

## Uso

```go
// Generar token en login
token, err := jwtService.GenerateToken(user)

// Validar en middleware
claims, err := jwtService.ValidateToken(tokenString)

// Renovar token
newToken, err := jwtService.RefreshToken(refreshToken)
```
