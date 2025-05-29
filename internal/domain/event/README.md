# event/ - Eventos de Dominio

Define los eventos que ocurren durante operaciones del dominio para comunicación asíncrona.

## Responsabilidades

- Definir eventos importantes del dominio
- Establecer contratos para event handlers
- Facilitar comunicación desacoplada entre contextos
- Habilitar arquitectura orientada a eventos

## Principios

- **Inmutables**: Los eventos no cambian una vez creados
- **Descriptivos**: Nombres claros que indican qué pasó
- **Ricos en información**: Contienen toda la data necesaria
- **Versioning**: Preparados para evolución

## Eventos Futuros

- **`employee_events.go`** - Eventos relacionados con empleados
- **`user_events.go`** - Eventos de usuarios
- **`department_events.go`** - Eventos departamentales
- **`dispatcher.go`** - Despachador de eventos

## Tipos de Eventos

### Eventos de Empleado
- `EmployeeCreated` - Cuando se crea un empleado
- `EmployeeUpdated` - Cuando se actualiza información
- `EmployeePromoted` - Cuando hay una promoción
- `EmployeeTerminated` - Cuando se termina el contrato

### Eventos de Usuario
- `UserRegistered` - Nuevo usuario registrado
- `UserLoggedIn` - Usuario inició sesión
- `PasswordChanged` - Cambio de contraseña

## Estructura de Evento

```go
type Event interface {
    EventName() string
    Payload() interface{}
    OccurredAt() time.Time
    AggregateID() uuid.UUID
}

type EmployeeCreatedEvent struct {
    Employee   *entity.Employee
    OccurredAt time.Time
    UserID     uuid.UUID
}

func (e EmployeeCreatedEvent) EventName() string {
    return "employee.created"
}
```

## Uso

Los eventos se disparan desde los casos de uso y son manejados por la infraestructura.
