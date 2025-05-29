# memory/ - Event Bus en Memoria

Implementación del bus de eventos en memoria para desarrollo, testing y casos simples.

## Responsabilidades

- Proporcionar un bus de eventos simple y rápido
- Facilitar testing unitario y de integración
- Servir como fallback cuando no hay infraestructura externa
- Demostrar la interface del event bus

## Características

- **Síncrono**: Los eventos se procesan inmediatamente
- **Sin persistencia**: Los eventos no sobreviven al reinicio
- **Thread-safe**: Seguro para uso concurrente
- **Simple**: Implementación minimalista

## Archivos Futuros

- **`memory_eventbus.go`** - Implementación principal
- **`subscriber_registry.go`** - Registro de suscriptores
- **`event_queue.go`** - Cola interna de eventos
- **`metrics.go`** - Métricas básicas

## Implementación

```go
type MemoryEventBus struct {
    subscribers map[string][]EventHandler
    mutex       sync.RWMutex
    logger      Logger
    metrics     *EventBusMetrics
}

func NewMemoryEventBus(logger Logger) *MemoryEventBus {
    return &MemoryEventBus{
        subscribers: make(map[string][]EventHandler),
        logger:      logger,
        metrics:     NewEventBusMetrics(),
    }
}

func (bus *MemoryEventBus) Publish(ctx context.Context, event Event) error {
    bus.mutex.RLock()
    handlers, exists := bus.subscribers[event.Type()]
    bus.mutex.RUnlock()
    
    if !exists {
        bus.logger.Debug("No handlers found for event type", "type", event.Type())
        return nil
    }
    
    bus.metrics.IncrementPublished(event.Type())
    
    // Ejecutar handlers secuencialmente
    for _, handler := range handlers {
        if err := bus.executeHandler(ctx, handler, event); err != nil {
            bus.metrics.IncrementFailed(event.Type())
            bus.logger.Error("Handler failed", "error", err, "event_type", event.Type())
            // Continuar con otros handlers (fire-and-forget)
        } else {
            bus.metrics.IncrementProcessed(event.Type())
        }
    }
    
    return nil
}

func (bus *MemoryEventBus) Subscribe(eventType string, handler EventHandler) error {
    bus.mutex.Lock()
    defer bus.mutex.Unlock()
    
    if bus.subscribers[eventType] == nil {
        bus.subscribers[eventType] = make([]EventHandler, 0)
    }
    
    bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
    bus.logger.Info("Handler subscribed", "event_type", eventType)
    
    return nil
}

func (bus *MemoryEventBus) Unsubscribe(eventType string, handler EventHandler) error {
    bus.mutex.Lock()
    defer bus.mutex.Unlock()
    
    handlers, exists := bus.subscribers[eventType]
    if !exists {
        return nil
    }
    
    // Filtrar el handler a remover
    filtered := make([]EventHandler, 0, len(handlers))
    for _, h := range handlers {
        if h != handler {
            filtered = append(filtered, h)
        }
    }
    
    bus.subscribers[eventType] = filtered
    return nil
}

func (bus *MemoryEventBus) executeHandler(ctx context.Context, handler EventHandler, event Event) error {
    defer func() {
        if r := recover(); r != nil {
            bus.logger.Error("Handler panicked", "panic", r, "event_type", event.Type())
        }
    }()
    
    return handler.Handle(ctx, event)
}

func (bus *MemoryEventBus) Close() error {
    bus.mutex.Lock()
    defer bus.mutex.Unlock()
    
    bus.subscribers = make(map[string][]EventHandler)
    bus.logger.Info("Memory event bus closed")
    
    return nil
}
```

## Modo Asíncrono (Opcional)

```go
type AsyncMemoryEventBus struct {
    *MemoryEventBus
    workers    int
    eventQueue chan EventWithHandlers
    stopCh     chan struct{}
    wg         sync.WaitGroup
}

type EventWithHandlers struct {
    Event    Event
    Handlers []EventHandler
}

func NewAsyncMemoryEventBus(workers int, logger Logger) *AsyncMemoryEventBus {
    bus := &AsyncMemoryEventBus{
        MemoryEventBus: NewMemoryEventBus(logger),
        workers:        workers,
        eventQueue:     make(chan EventWithHandlers, 1000),
        stopCh:         make(chan struct{}),
    }
    
    // Iniciar workers
    for i := 0; i < workers; i++ {
        bus.wg.Add(1)
        go bus.worker()
    }
    
    return bus
}

func (bus *AsyncMemoryEventBus) Publish(ctx context.Context, event Event) error {
    bus.mutex.RLock()
    handlers, exists := bus.subscribers[event.Type()]
    bus.mutex.RUnlock()
    
    if !exists {
        return nil
    }
    
    select {
    case bus.eventQueue <- EventWithHandlers{Event: event, Handlers: handlers}:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (bus *AsyncMemoryEventBus) worker() {
    defer bus.wg.Done()
    
    for {
        select {
        case eventWithHandlers := <-bus.eventQueue:
            for _, handler := range eventWithHandlers.Handlers {
                bus.executeHandler(context.Background(), handler, eventWithHandlers.Event)
            }
        case <-bus.stopCh:
            return
        }
    }
}
```

## Ventajas

- **Simplicidad**: Sin dependencias externas
- **Velocidad**: Procesamiento inmediato en memoria
- **Testing**: Ideal para pruebas unitarias
- **Desarrollo**: Sin setup de infraestructura

## Limitaciones

- **No persistente**: Eventos se pierden al reiniciar
- **No distribuido**: Solo funciona en una instancia
- **Memory bound**: Limitado por la memoria disponible
- **No reliable**: Sin garantías de entrega
