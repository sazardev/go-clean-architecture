# kafka/ - Event Bus con Apache Kafka

Implementación del bus de eventos usando Apache Kafka para sistemas distribuidos de alta escala.

## Responsabilidades

- Publicar eventos a topics de Kafka
- Consumir eventos de manera confiable
- Manejar particionado y offset management
- Garantizar entrega al menos una vez
- Implementar dead letter queues

## Características

### Ventajas de Kafka
- **Alta throughput**: Millones de eventos por segundo
- **Durabilidad**: Persistencia configurable
- **Escalabilidad**: Particionado horizontal
- **Ordenamiento**: Por partición
- **Retención**: Configuración flexible de TTL

### Casos de Uso
- Sistemas de alta escala
- Event sourcing
- Streaming de datos
- Integración entre microservicios

## Archivos Futuros

- **`kafka_eventbus.go`** - Implementación principal
- **`producer.go`** - Productor de eventos
- **`consumer.go`** - Consumidor de eventos
- **`config.go`** - Configuración de Kafka
- **`partition_strategy.go`** - Estrategias de particionado

## Configuración de Topics

```go
type TopicConfig struct {
    Name              string
    Partitions        int
    ReplicationFactor int
    RetentionMs       int64
}

var EventTopics = map[string]TopicConfig{
    "employee.events": {
        Name:              "employee-events",
        Partitions:        6,
        ReplicationFactor: 3,
        RetentionMs:       7 * 24 * 60 * 60 * 1000, // 7 días
    },
    "user.events": {
        Name:              "user-events",
        Partitions:        3,
        ReplicationFactor: 3,
        RetentionMs:       30 * 24 * 60 * 60 * 1000, // 30 días
    },
    "audit.events": {
        Name:              "audit-events",
        Partitions:        12,
        ReplicationFactor: 3,
        RetentionMs:       365 * 24 * 60 * 60 * 1000, // 1 año
    },
}
```

## Implementación

### Producer
```go
type KafkaEventBus struct {
    producer kafka.Producer
    consumer kafka.Consumer
    config   *KafkaConfig
    logger   Logger
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

func NewKafkaEventBus(config *KafkaConfig, logger Logger) (*KafkaEventBus, error) {
    producer, err := kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers": config.Brokers,
        "client.id":         config.ClientID,
        "acks":             "all",
        "retries":          2147483647,
        "max.in.flight.requests.per.connection": 5,
        "enable.idempotence": true,
    })
    if err != nil {
        return nil, err
    }
    
    return &KafkaEventBus{
        producer: producer,
        config:   config,
        logger:   logger,
        handlers: make(map[string][]EventHandler),
    }, nil
}

func (bus *KafkaEventBus) Publish(ctx context.Context, event Event) error {
    topic := bus.getTopicForEvent(event.Type())
    key := bus.getPartitionKey(event)
    
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
    
    message := &kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &topic,
            Partition: kafka.PartitionAny,
        },
        Key:   []byte(key),
        Value: eventData,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte(event.Type())},
            {Key: "content-type", Value: []byte("application/json")},
        },
    }
    
    deliveryChan := make(chan kafka.Event)
    err = bus.producer.Produce(message, deliveryChan)
    if err != nil {
        return err
    }
    
    // Esperar confirmación
    e := <-deliveryChan
    m := e.(*kafka.Message)
    
    if m.TopicPartition.Error != nil {
        return m.TopicPartition.Error
    }
    
    bus.logger.Debug("Event published",
        "topic", topic,
        "partition", m.TopicPartition.Partition,
        "offset", m.TopicPartition.Offset,
        "event_type", event.Type(),
    )
    
    return nil
}
```

### Consumer
```go
func (bus *KafkaEventBus) StartConsumer(ctx context.Context) error {
    consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  bus.config.Brokers,
        "group.id":          bus.config.ConsumerGroup,
        "auto.offset.reset": "earliest",
        "enable.auto.commit": false,
    })
    if err != nil {
        return err
    }
    
    bus.consumer = consumer
    
    topics := bus.getSubscribedTopics()
    err = consumer.SubscribeTopics(topics, nil)
    if err != nil {
        return err
    }
    
    go bus.consumeLoop(ctx)
    return nil
}

func (bus *KafkaEventBus) consumeLoop(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            msg, err := bus.consumer.ReadMessage(100 * time.Millisecond)
            if err != nil {
                continue
            }
            
            if err := bus.processMessage(ctx, msg); err != nil {
                bus.logger.Error("Failed to process message",
                    "error", err,
                    "topic", *msg.TopicPartition.Topic,
                    "partition", msg.TopicPartition.Partition,
                    "offset", msg.TopicPartition.Offset,
                )
                continue
            }
            
            // Commit offset después del procesamiento exitoso
            bus.consumer.CommitMessage(msg)
        }
    }
}

func (bus *KafkaEventBus) processMessage(ctx context.Context, msg *kafka.Message) error {
    var envelope EventEnvelope
    if err := json.Unmarshal(msg.Value, &envelope); err != nil {
        return err
    }
    
    event := &BaseEvent{
        eventType: envelope.Type,
        payload:   envelope.Payload,
        metadata:  envelope.Metadata,
        timestamp: envelope.Timestamp,
    }
    
    bus.mutex.RLock()
    handlers, exists := bus.handlers[event.Type()]
    bus.mutex.RUnlock()
    
    if !exists {
        return nil
    }
    
    for _, handler := range handlers {
        if err := handler.Handle(ctx, event); err != nil {
            return err
        }
    }
    
    return nil
}
```

## Partición Strategy

```go
func (bus *KafkaEventBus) getPartitionKey(event Event) string {
    switch event.Type() {
    case "employee.created", "employee.updated", "employee.deleted":
        if payload, ok := event.Payload().(EmployeePayload); ok {
            return payload.EmployeeID.String()
        }
    case "user.logged_in", "user.password_changed":
        if payload, ok := event.Payload().(UserPayload); ok {
            return payload.UserID.String()
        }
    }
    
    // Fallback: distribución aleatoria
    return uuid.New().String()
}
```

## Configuración

Variables de entorno:
- `KAFKA_BROKERS` - Lista de brokers
- `KAFKA_CLIENT_ID` - ID del cliente
- `KAFKA_CONSUMER_GROUP` - Grupo de consumidores
- `KAFKA_SECURITY_PROTOCOL` - Protocolo de seguridad
- `KAFKA_SASL_*` - Configuración SASL
