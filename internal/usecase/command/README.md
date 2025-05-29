# command/ - Casos de Uso de Escritura

Implementa todos los casos de uso que modifican el estado del sistema siguiendo el patrón Command.

## Responsabilidades

- Validar comandos de entrada
- Orquestar operaciones de escritura
- Mantener consistencia transaccional
- Publicar eventos de dominio
- Manejar errores y rollbacks

## Características

- **Idempotentes cuando sea posible**: Mismo resultado en múltiples ejecuciones
- **Transaccionales**: Operaciones atómicas
- **Event-driven**: Publican eventos para comunicación asíncrona
- **Validación robusta**: Verifican precondiciones

## Comandos Actuales

- Casos de uso de Employee existentes se moverán aquí

## Comandos Futuros

### Employee Commands
- **`create_employee.go`** - Crear nuevo empleado
- **`update_employee.go`** - Actualizar información
- **`promote_employee.go`** - Promocionar empleado
- **`terminate_employee.go`** - Terminar contrato

### User Commands  
- **`register_user.go`** - Registrar nuevo usuario
- **`change_password.go`** - Cambiar contraseña
- **`assign_role.go`** - Asignar rol a usuario

### Department Commands
- **`create_department.go`** - Crear departamento
- **`assign_manager.go`** - Asignar gerente

## Estructura de Command

```go
// Comando de entrada
type CreateEmployeeCommand struct {
    Name       string
    Email      string
    Department string
    Position   string
    UserID     uuid.UUID // Quien ejecuta el comando
}

// Resultado
type CreateEmployeeResult struct {
    EmployeeID uuid.UUID
    CreatedAt  time.Time
}

// Caso de uso
type CreateEmployeeUseCase struct {
    employeeRepo domain.EmployeeRepository
    deptRepo     domain.DepartmentRepository
    eventBus     EventBus
}

func (uc *CreateEmployeeUseCase) Execute(
    ctx context.Context,
    cmd CreateEmployeeCommand,
) (*CreateEmployeeResult, error) {
    // Implementación del caso de uso
}
```
