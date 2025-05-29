# messaging/ - Sistema de Mensajería

Sistema de colas de mensajes para procesamiento asíncrono de tareas y comunicación entre servicios.

## Responsabilidades

- Gestionar colas de tareas asíncronas
- Implementar workers para procesamiento en background
- Manejar retry logic y dead letter queues
- Priorización de tareas
- Distribución de carga entre workers

## Casos de Uso

### Tareas Asíncronas
- Envío de emails masivos
- Generación de reportes pesados
- Procesamiento de archivos
- Sincronización con sistemas externos
- Limpieza de datos

### Comunicación Entre Servicios
- Notificaciones entre microservicios
- Propagación de cambios de estado
- Integración con sistemas legacy

## Archivos Futuros

- **`queue_manager.go`** - Gestor principal de colas
- **`worker_pool.go`** - Pool de workers
- **`task_scheduler.go`** - Programador de tareas
- **`retry_strategy.go`** - Estrategias de reintento
- **`dead_letter_queue.go`** - Cola de mensajes fallidos

## Implementación

### Queue Manager
```go
type QueueManager interface {
    Enqueue(ctx context.Context, queueName string, task Task) error
    Dequeue(ctx context.Context, queueName string) (*Task, error)
    GetQueueSize(ctx context.Context, queueName string) (int64, error)
    PurgeQueue(ctx context.Context, queueName string) error
}

type Task struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Payload     map[string]interface{} `json:"payload"`
    Priority    int                    `json:"priority"`
    MaxRetries  int                    `json:"max_retries"`
    RetryCount  int                    `json:"retry_count"`
    ScheduledAt time.Time              `json:"scheduled_at"`
    CreatedAt   time.Time              `json:"created_at"`
}

type RedisQueueManager struct {
    client *redis.Client
    logger Logger
}

func (qm *RedisQueueManager) Enqueue(ctx context.Context, queueName string, task Task) error {
    task.ID = uuid.New().String()
    task.CreatedAt = time.Now()
    
    taskData, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    // Usar Redis sorted set para prioridad
    score := float64(task.Priority)*1000000 + float64(task.ScheduledAt.Unix())
    
    return qm.client.ZAdd(ctx, queueName, &redis.Z{
        Score:  score,
        Member: taskData,
    }).Err()
}

func (qm *RedisQueueManager) Dequeue(ctx context.Context, queueName string) (*Task, error) {
    // Obtener tarea con mayor prioridad que esté lista
    now := float64(time.Now().Unix())
    maxScore := fmt.Sprintf("(%f", now*1000000)
    
    result, err := qm.client.ZRangeByScore(ctx, queueName, &redis.ZRangeBy{
        Min:   "-inf",
        Max:   maxScore,
        Count: 1,
    }).Result()
    
    if err != nil {
        return nil, err
    }
    
    if len(result) == 0 {
        return nil, ErrQueueEmpty
    }
    
    // Remover de la cola
    qm.client.ZRem(ctx, queueName, result[0])
    
    var task Task
    err = json.Unmarshal([]byte(result[0]), &task)
    return &task, err
}
```

### Worker Pool
```go
type WorkerPool struct {
    queueManager QueueManager
    handlers     map[string]TaskHandler
    workers      int
    queues       []string
    stopCh       chan struct{}
    wg           sync.WaitGroup
    logger       Logger
}

type TaskHandler interface {
    Handle(ctx context.Context, task *Task) error
    CanHandle(taskType string) bool
}

func NewWorkerPool(queueManager QueueManager, workers int, logger Logger) *WorkerPool {
    return &WorkerPool{
        queueManager: queueManager,
        handlers:     make(map[string]TaskHandler),
        workers:      workers,
        queues:       make([]string, 0),
        stopCh:       make(chan struct{}),
        logger:       logger,
    }
}

func (wp *WorkerPool) RegisterHandler(taskType string, handler TaskHandler) {
    wp.handlers[taskType] = handler
}

func (wp *WorkerPool) AddQueue(queueName string) {
    wp.queues = append(wp.queues, queueName)
}

func (wp *WorkerPool) Start(ctx context.Context) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx, i)
    }
    
    wp.logger.Info("Worker pool started", "workers", wp.workers, "queues", wp.queues)
}

func (wp *WorkerPool) worker(ctx context.Context, workerID int) {
    defer wp.wg.Done()
    
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-wp.stopCh:
            return
        case <-ticker.C:
            wp.processQueues(ctx, workerID)
        }
    }
}

func (wp *WorkerPool) processQueues(ctx context.Context, workerID int) {
    for _, queueName := range wp.queues {
        task, err := wp.queueManager.Dequeue(ctx, queueName)
        if err != nil {
            if err != ErrQueueEmpty {
                wp.logger.Error("Failed to dequeue task", "error", err, "queue", queueName)
            }
            continue
        }
        
        wp.processTask(ctx, task, queueName, workerID)
    }
}

func (wp *WorkerPool) processTask(ctx context.Context, task *Task, queueName string, workerID int) {
    handler, exists := wp.handlers[task.Type]
    if !exists {
        wp.logger.Error("No handler for task type", "type", task.Type, "task_id", task.ID)
        return
    }
    
    wp.logger.Info("Processing task",
        "task_id", task.ID,
        "type", task.Type,
        "worker", workerID,
        "queue", queueName,
    )
    
    err := handler.Handle(ctx, task)
    if err != nil {
        wp.handleTaskError(ctx, task, queueName, err)
        return
    }
    
    wp.logger.Info("Task completed successfully", "task_id", task.ID, "worker", workerID)
}

func (wp *WorkerPool) handleTaskError(ctx context.Context, task *Task, queueName string, err error) {
    wp.logger.Error("Task failed", "error", err, "task_id", task.ID, "retry_count", task.RetryCount)
    
    task.RetryCount++
    
    if task.RetryCount >= task.MaxRetries {
        // Mover a dead letter queue
        dlqName := queueName + ":dlq"
        if dlqErr := wp.queueManager.Enqueue(ctx, dlqName, *task); dlqErr != nil {
            wp.logger.Error("Failed to enqueue to DLQ", "error", dlqErr, "task_id", task.ID)
        }
        return
    }
    
    // Reencolar con backoff exponencial
    backoffDelay := time.Duration(math.Pow(2, float64(task.RetryCount))) * time.Second
    task.ScheduledAt = time.Now().Add(backoffDelay)
    
    if retryErr := wp.queueManager.Enqueue(ctx, queueName, *task); retryErr != nil {
        wp.logger.Error("Failed to requeue task", "error", retryErr, "task_id", task.ID)
    }
}
```

## Task Handlers

### Email Handler
```go
type EmailTaskHandler struct {
    emailService EmailService
}

func (h *EmailTaskHandler) Handle(ctx context.Context, task *Task) error {
    switch task.Type {
    case "send_welcome_email":
        return h.sendWelcomeEmail(ctx, task.Payload)
    case "send_notification_email":
        return h.sendNotificationEmail(ctx, task.Payload)
    default:
        return fmt.Errorf("unknown email task type: %s", task.Type)
    }
}

func (h *EmailTaskHandler) CanHandle(taskType string) bool {
    return strings.HasPrefix(taskType, "send_") && strings.HasSuffix(taskType, "_email")
}
```

## Configuración

Variables de entorno:
- `QUEUE_PROVIDER` - Proveedor (redis, rabbitmq, sqs)
- `WORKER_COUNT` - Número de workers
- `QUEUE_POLL_INTERVAL` - Intervalo de polling
- `MAX_RETRY_ATTEMPTS` - Intentos máximos de retry
