# common/ - Código Compartido de Casos de Uso

Contiene utilidades, interfaces y código común compartido entre commands y queries.

## Responsabilidades

- Definir interfaces comunes para casos de uso
- Implementar utilidades de validación
- Manejar errores comunes
- Proporcionar helpers para paginación y filtrado
- Definir estructuras base compartidas

## Contenido

### Interfaces Base
- **`use_case.go`** - Interface base para casos de uso
- **`command_handler.go`** - Interface para command handlers
- **`query_handler.go`** - Interface para query handlers

### Utilidades
- **`validation.go`** - Validaciones comunes
- **`errors.go`** - Tipos de error específicos
- **`pagination.go`** - Helpers para paginación
- **`filtering.go`** - Helpers para filtrado

### DTOs Base
- **`base_dto.go`** - Estructuras base para DTOs
- **`result.go`** - Tipos de resultado comunes

## Ejemplos

```go
// Interface base para casos de uso
type UseCase[TRequest any, TResponse any] interface {
    Execute(ctx context.Context, request TRequest) (TResponse, error)
}

// Command handler base
type CommandHandler[TCommand any, TResult any] interface {
    Handle(ctx context.Context, command TCommand) (TResult, error)
}

// Query handler base  
type QueryHandler[TQuery any, TResult any] interface {
    Handle(ctx context.Context, query TQuery) (TResult, error)
}

// Resultado paginado común
type PaginatedResult[T any] struct {
    Items      []T   `json:"items"`
    TotalCount int64 `json:"total_count"`
    Page       int   `json:"page"`
    PageSize   int   `json:"page_size"`
    TotalPages int   `json:"total_pages"`
}

// Errores específicos de casos de uso
var (
    ErrEmployeeNotFound    = errors.New("employee not found")
    ErrInvalidEmployeeData = errors.New("invalid employee data")
    ErrUnauthorized        = errors.New("unauthorized operation")
)

// Validaciones comunes
func ValidateEmail(email string) error {
    // Implementación de validación
}

func ValidatePagination(page, pageSize int) error {
    // Implementación de validación
}
```

## Principios

- **DRY**: Evita duplicación entre commands y queries
- **Consistencia**: Patrones uniformes en toda la aplicación
- **Reutilización**: Código que puede usarse en múltiples casos de uso
- **Testeable**: Utilidades fáciles de probar
