# repository/ - Interfaces de Repositorio

Define los contratos para el acceso a datos sin especificar la implementación.

## Responsabilidades

- Definir métodos para operaciones CRUD
- Especificar queries complejas necesarias
- Establecer contratos para búsquedas y filtros
- Mantener la abstracción sobre la persistencia

## Principios

- **Solo interfaces**: No contiene implementaciones
- **Agnóstico a tecnología**: No especifica base de datos
- **Orientado al dominio**: Métodos expresados en términos de negocio
- **Testeable**: Fácil de mockear para testing

## Interfaces Actuales

- **`employee_repository.go`** - Contrato para persistencia de empleados

## Interfaces Futuras

- **`user_repository.go`** - Gestión de usuarios
- **`role_repository.go`** - Gestión de roles
- **`department_repository.go`** - Gestión de departamentos

## Patrón de Diseño

```go
type EmployeeRepository interface {
    Create(ctx context.Context, employee *entity.Employee) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
    FindAll(ctx context.Context) ([]*entity.Employee, error)
    Update(ctx context.Context, employee *entity.Employee) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // Métodos específicos del dominio
    FindByDepartment(ctx context.Context, deptID uuid.UUID) ([]*entity.Employee, error)
    FindActiveEmployees(ctx context.Context) ([]*entity.Employee, error)
}
```

La implementación concreta estará en `infrastructure/database/`.
