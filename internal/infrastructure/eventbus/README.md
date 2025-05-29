# eventbus/ - Bus de Eventos

Sistema de eventos para comunicación desacoplada entre componentes del sistema.

## Responsabilidades

- Publicar eventos de dominio
- Suscribir handlers a eventos específicos
- Garantizar entrega confiable de eventos
- Manejar retry y dead letter queues
- Implementar múltiples transportes (memoria, Kafka, NATS)

## Estructura

- **`memory/`** - Implementación en memoria para desarrollo/testing
- **`kafka/`** - Implementación con Apache Kafka para producción
- **`nats/`** - Implementación con NATS para casos específicos

## Patrones Implementados

### Publisher-Subscriber
```go
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(eventType string, handler EventHandler) error
    Unsubscribe(eventType string, handler EventHandler) error
    Close() error
}

type Event interface {
    Type() string
    Payload() interface{}
    Metadata() map[string]string
    Timestamp() time.Time
}

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
    CanHandle(eventType string) bool
}
```

### Event Sourcing (Futuro)
```go
type EventStore interface {
    SaveEvents(ctx context.Context, aggregateID uuid.UUID, events []Event) error
    LoadEvents(ctx context.Context, aggregateID uuid.UUID) ([]Event, error)
    LoadEventsFromVersion(ctx context.Context, aggregateID uuid.UUID, version int) ([]Event, error)
}
```

## Eventos del Sistema

### Employee Events
- `employee.created` - Empleado creado
- `employee.updated` - Empleado actualizado
- `employee.promoted` - Empleado promovido
- `employee.terminated` - Empleado dado de baja

### User Events
- `user.registered` - Usuario registrado
- `user.logged_in` - Usuario inició sesión
- `user.password_changed` - Contraseña cambiada

### System Events
- `system.backup_completed` - Backup completado
- `system.maintenance_started` - Mantenimiento iniciado

## Implementación Base

```go
type BaseEvent struct {
    eventType string
    payload   interface{}
    metadata  map[string]string
    timestamp time.Time
    id        uuid.UUID
}

func NewEvent(eventType string, payload interface{}) *BaseEvent {
    return &BaseEvent{
        eventType: eventType,
        payload:   payload,
        metadata:  make(map[string]string),
        timestamp: time.Now(),
        id:        uuid.New(),
    }
}

func (e *BaseEvent) Type() string { return e.eventType }
func (e *BaseEvent) Payload() interface{} { return e.payload }
func (e *BaseEvent) Metadata() map[string]string { return e.metadata }
func (e *BaseEvent) Timestamp() time.Time { return e.timestamp }
```

## Handlers

```go
type EmployeeEventHandler struct {
    emailService EmailService
    auditService AuditService
}

func (h *EmployeeEventHandler) Handle(ctx context.Context, event Event) error {
    switch event.Type() {
    case "employee.created":
        return h.handleEmployeeCreated(ctx, event)
    case "employee.promoted":
        return h.handleEmployeePromoted(ctx, event)
    default:
        return fmt.Errorf("unknown event type: %s", event.Type())
    }
}

func (h *EmployeeEventHandler) handleEmployeeCreated(ctx context.Context, event Event) error {
    payload := event.Payload().(*EmployeeCreatedPayload)
    
    // Enviar email de bienvenida
    if err := h.emailService.SendWelcomeEmail(ctx, payload.Employee); err != nil {
        return err
    }
    
    // Registrar en auditoría
    return h.auditService.LogEvent(ctx, "employee_created", payload.Employee.ID)
}
```

## Configuración

Variables de entorno:
- `EVENTBUS_PROVIDER` - Proveedor (memory, kafka, nats)
- `EVENTBUS_RETRY_ATTEMPTS` - Intentos de retry
- `EVENTBUS_DEAD_LETTER_QUEUE` - Cola de eventos fallidos
- `EVENTBUS_BATCH_SIZE` - Tamaño de lote para procesamiento

## Uso en Casos de Uso

```go
func (uc *CreateEmployeeUseCase) Execute(ctx context.Context, cmd CreateEmployeeCommand) error {
    // ... lógica de creación
    
    // Publicar evento
    event := NewEvent("employee.created", &EmployeeCreatedPayload{
        Employee: employee,
        UserID:   cmd.UserID,
    })
    
    return uc.eventBus.Publish(ctx, event)
}
```
