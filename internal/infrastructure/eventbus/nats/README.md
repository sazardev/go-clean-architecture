# nats/ - Event Bus con NATS

Implementación del bus de eventos usando NATS para comunicación de baja latencia y alta performance.

## Responsabilidades

- Publicar eventos a subjects de NATS
- Suscribirse a eventos específicos
- Manejar request-reply patterns
- Implementar streaming con NATS JetStream
- Gestionar durabilidad y acknowledgments

## Características

### Ventajas de NATS
- **Ultra baja latencia**: < 1ms en redes locales
- **Simplicidad**: Protocolo lightweight
- **Patrrones múltiples**: Pub/Sub, Request/Reply, Queuing
- **Clustering**: Tolerancia a fallos automática
- **JetStream**: Persistencia y exactly-once delivery

### Casos de Uso
- Microservicios comunicándose
- Sistemas de tiempo real
- Notificaciones inmediatas
- Command/Query buses

## Archivos Futuros

- **`nats_eventbus.go`** - Implementación principal
- **`subjects.go`** - Definición de subjects
- **`jetstream.go`** - Configuración de JetStream
- **`request_reply.go`** - Patrón request-reply

## Subjects Structure

```go
const (
    // Core events
    SubjectEmployeeCreated   = "hr.employee.created"
    SubjectEmployeeUpdated   = "hr.employee.updated"
    SubjectEmployeePromoted  = "hr.employee.promoted"
    SubjectEmployeeDeleted   = "hr.employee.deleted"
    
    // User events
    SubjectUserRegistered    = "hr.user.registered"
    SubjectUserLoggedIn      = "hr.user.logged_in"
    SubjectPasswordChanged   = "hr.user.password_changed"
    
    // System events
    SubjectSystemHealth      = "hr.system.health"
    SubjectSystemMaintenance = "hr.system.maintenance"
    
    // Audit events
    SubjectAuditLog          = "hr.audit.*"
)

// Wildcard subjects para suscripciones
const (
    SubjectAllEmployeeEvents = "hr.employee.*"
    SubjectAllUserEvents     = "hr.user.*"
    SubjectAllSystemEvents   = "hr.system.*"
)
```

## Implementación

### Basic NATS EventBus
```go
type NATSEventBus struct {
    conn     *nats.Conn
    js       nats.JetStreamContext
    logger   Logger
    config   *NATSConfig
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
    subs     []*nats.Subscription
}

func NewNATSEventBus(config *NATSConfig, logger Logger) (*NATSEventBus, error) {
    // Conectar a NATS
    conn, err := nats.Connect(config.URL,
        nats.Name(config.ClientID),
        nats.MaxReconnects(config.MaxReconnects),
        nats.ReconnectWait(config.ReconnectWait),
    )
    if err != nil {
        return nil, err
    }
    
    // Configurar JetStream (opcional)
    var js nats.JetStreamContext
    if config.EnableJetStream {
        js, err = conn.JetStream()
        if err != nil {
            return nil, err
        }
    }
    
    return &NATSEventBus{
        conn:     conn,
        js:       js,
        logger:   logger,
        config:   config,
        handlers: make(map[string][]EventHandler),
        subs:     make([]*nats.Subscription, 0),
    }, nil
}

func (bus *NATSEventBus) Publish(ctx context.Context, event Event) error {
    subject := bus.getSubjectForEvent(event.Type())
    
    eventData, err := json.Marshal(EventEnvelope{
        ID:        uuid.New().String(),
        Type:      event.Type(),
        Payload:   event.Payload(),
        Metadata:  event.Metadata(),
        Timestamp: event.Timestamp(),
    })
    if err != nil {
        return err
    }
    
    if bus.js != nil && bus.shouldUsePersistence(event.Type()) {
        // Usar JetStream para eventos que requieren persistencia
        _, err = bus.js.Publish(subject, eventData)
    } else {
        // Usar NATS core para eventos efímeros
        err = bus.conn.Publish(subject, eventData)
    }
    
    if err != nil {
        return err
    }
    
    bus.logger.Debug("Event published",
        "subject", subject,
        "event_type", event.Type(),
        "persistent", bus.shouldUsePersistence(event.Type()),
    )
    
    return nil
}
```

### Suscripciones
```go
func (bus *NATSEventBus) Subscribe(eventType string, handler EventHandler) error {
    subject := bus.getSubjectForEvent(eventType)
    
    var sub *nats.Subscription
    var err error
    
    if bus.js != nil && bus.shouldUsePersistence(eventType) {
        // Suscripción durable con JetStream
        sub, err = bus.js.QueueSubscribe(subject, bus.config.QueueGroup, bus.createHandler(handler))
    } else {
        // Suscripción simple con NATS core
        sub, err = bus.conn.QueueSubscribe(subject, bus.config.QueueGroup, bus.createHandler(handler))
    }
    
    if err != nil {
        return err
    }
    
    bus.mutex.Lock()
    bus.subs = append(bus.subs, sub)
    if bus.handlers[eventType] == nil {
        bus.handlers[eventType] = make([]EventHandler, 0)
    }
    bus.handlers[eventType] = append(bus.handlers[eventType], handler)
    bus.mutex.Unlock()
    
    bus.logger.Info("Handler subscribed", "subject", subject, "event_type", eventType)
    return nil
}

func (bus *NATSEventBus) createHandler(handler EventHandler) nats.MsgHandler {
    return func(msg *nats.Msg) {
        var envelope EventEnvelope
        if err := json.Unmarshal(msg.Data, &envelope); err != nil {
            bus.logger.Error("Failed to unmarshal event", "error", err)
            return
        }
        
        event := &BaseEvent{
            eventType: envelope.Type,
            payload:   envelope.Payload,
            metadata:  envelope.Metadata,
            timestamp: envelope.Timestamp,
        }
        
        ctx := context.Background()
        if err := handler.Handle(ctx, event); err != nil {
            bus.logger.Error("Handler failed", "error", err, "event_type", event.Type())
            
            // En JetStream, no hacer ACK para reintento
            if msg.Reply == "" {
                return
            }
        }
        
        // ACK el mensaje si es de JetStream
        if msg.Reply != "" {
            msg.Ack()
        }
    }
}
```

### Request-Reply Pattern
```go
func (bus *NATSEventBus) Request(ctx context.Context, subject string, data interface{}) (interface{}, error) {
    requestData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    
    msg, err := bus.conn.RequestWithContext(ctx, subject, requestData)
    if err != nil {
        return nil, err
    }
    
    var response interface{}
    if err := json.Unmarshal(msg.Data, &response); err != nil {
        return nil, err
    }
    
    return response, nil
}

func (bus *NATSEventBus) RegisterRequestHandler(subject string, handler func([]byte) (interface{}, error)) error {
    sub, err := bus.conn.Subscribe(subject, func(msg *nats.Msg) {
        response, err := handler(msg.Data)
        if err != nil {
            bus.logger.Error("Request handler failed", "error", err, "subject", subject)
            msg.Respond([]byte(`{"error":"` + err.Error() + `"}`))
            return
        }
        
        responseData, err := json.Marshal(response)
        if err != nil {
            bus.logger.Error("Failed to marshal response", "error", err)
            msg.Respond([]byte(`{"error":"internal error"}`))
            return
        }
        
        msg.Respond(responseData)
    })
    
    if err != nil {
        return err
    }
    
    bus.mutex.Lock()
    bus.subs = append(bus.subs, sub)
    bus.mutex.Unlock()
    
    return nil
}
```

## JetStream Configuration

```go
func (bus *NATSEventBus) setupJetStreamStreams() error {
    streams := []nats.StreamConfig{
        {
            Name:        "HR_EVENTS",
            Subjects:    []string{"hr.employee.*", "hr.user.*"},
            Storage:     nats.FileStorage,
            MaxAge:      30 * 24 * time.Hour, // 30 días
            MaxBytes:    1024 * 1024 * 1024,  // 1GB
            Retention:   nats.LimitsPolicy,
            Replicas:    3,
        },
        {
            Name:        "HR_AUDIT",
            Subjects:    []string{"hr.audit.*"},
            Storage:     nats.FileStorage,
            MaxAge:      365 * 24 * time.Hour, // 1 año
            MaxBytes:    10 * 1024 * 1024 * 1024, // 10GB
            Retention:   nats.LimitsPolicy,
            Replicas:    3,
        },
    }
    
    for _, streamConfig := range streams {
        _, err := bus.js.AddStream(&streamConfig)
        if err != nil && !strings.Contains(err.Error(), "already exists") {
            return err
        }
    }
    
    return nil
}
```

## Configuración

Variables de entorno:
- `NATS_URL` - URL del servidor NATS
- `NATS_CLIENT_ID` - ID del cliente
- `NATS_QUEUE_GROUP` - Grupo de cola
- `NATS_ENABLE_JETSTREAM` - Habilitar JetStream
- `NATS_MAX_RECONNECTS` - Reconexiones máximas
