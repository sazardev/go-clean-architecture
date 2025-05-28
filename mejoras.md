# Arquitectura "Todo Terreno" para HR API: Implementación Detallada

## 1. Estructura de Proyecto Ampliada

Para soportar todas las capacidades requeridas, propongo esta estructura expandida:

```
go-clean-architecture/
├── cmd/
│   ├── server/             # Servidor HTTP principal
│   ├── worker/             # Procesamiento asíncrono de tareas
│   └── migration/          # Herramienta de migraciones
├── internal/
│   ├── domain/
│   │   ├── entity/         # Entidades core del dominio
│   │   ├── repository/     # Interfaces de repositorio
│   │   ├── service/        # Servicios de dominio
│   │   ├── event/          # Eventos de dominio
│   │   └── valueobject/    # Objetos de valor
│   ├── usecase/
│   │   ├── command/        # Comandos (escritura)
│   │   ├── query/          # Consultas (lectura)
│   │   └── common/         # Código compartido
│   ├── infrastructure/
│   │   ├── auth/           # Autenticación/autorización
│   │   ├── config/         # Configuración
│   │   ├── container/      # Inyección de dependencias
│   │   ├── database/       # Persistencia
│   │   │   ├── postgres/   # Implementación PostgreSQL
│   │   │   ├── mongodb/    # Implementación MongoDB
│   │   │   ├── redis/      # Caché/mensajería
│   │   │   └── factory/    # Fábrica de repositorios
│   │   ├── http/           # API REST
│   │   ├── grpc/           # API gRPC (opcional)
│   │   ├── graphql/        # API GraphQL (opcional)
│   │   ├── websocket/      # Comunicación en tiempo real
│   │   ├── messaging/      # Sistema de mensajería
│   │   ├── eventbus/       # Bus de eventos
│   │   ├── cache/          # Capa de caché
│   │   ├── telemetry/      # Observabilidad
│   │   └── validation/     # Validación centralizada
│   └── pkg/                # Utilidades compartidas
└── pkg/                    # Bibliotecas públicas (opcional)
```

## 2. Autenticación y Autorización

### Estructura Detallada
```
internal/domain/entity/
├── user.go               # Entidad de usuario
├── role.go               # Entidad de rol
└── permission.go         # Entidad de permiso

internal/domain/repository/
├── user_repository.go    # Interfaz de repositorio de usuarios
└── role_repository.go    # Interfaz de repositorio de roles

internal/infrastructure/auth/
├── jwt/                  # Implementación JWT
│   ├── token.go          # Generación/validación de tokens
│   └── claims.go         # Claims personalizados
├── rbac/                 # Control de acceso basado en roles
│   ├── enforcer.go       # Verificador de permisos
│   ├── policy.go         # Definición de políticas
│   └── adapter.go        # Adaptador para Casbin
├── middleware/           # Middlewares de autenticación
│   ├── auth.go           # Middleware principal de autenticación
│   └── permission.go     # Middleware de verificación de permisos
└── service.go            # Servicio de autenticación
```

### Implementación Recomendada

1. **Entidad de Usuario**
```go
// internal/domain/entity/user.go
type User struct {
    ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
    Username    string       `json:"username" gorm:"unique;not null"`
    Email       string       `json:"email" gorm:"unique;not null"`
    Password    string       `json:"-" gorm:"not null"`  // Contraseña hasheada
    Roles       []Role       `json:"roles" gorm:"many2many:user_roles"`
    IsActive    bool         `json:"is_active" gorm:"default:true"`
    LastLoginAt *time.Time   `json:"last_login_at"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
}
```

2. **Middleware JWT**
```go
// internal/infrastructure/auth/middleware/auth.go
func JWTMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        token := extractToken(c)
        if token == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Authorization token required",
            })
        }
        
        claims, err := validateToken(token)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid or expired token",
            })
        }
        
        // Almacenar claims en el contexto
        c.Locals("user", claims.User)
        c.Locals("roles", claims.Roles)
        
        return c.Next()
    }
}
```

3. **Verificación de Permisos**
```go
// internal/infrastructure/auth/middleware/permission.go
func RequirePermission(permission string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(domain.User)
        roles := c.Locals("roles").([]string)
        
        enforcer := rbac.GetEnforcer()
        hasPermission := false
        
        for _, role := range roles {
            if enforcer.Enforce(role, permission) {
                hasPermission = true
                break
            }
        }
        
        if !hasPermission {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        
        return c.Next()
    }
}
```

4. **Integración con Casbin**
```go
// internal/infrastructure/auth/rbac/enforcer.go
type Enforcer struct {
    enforcer *casbin.Enforcer
}

func NewEnforcer() (*Enforcer, error) {
    adapter, err := gormadapter.NewAdapter("postgres", dbConnectionString)
    if err != nil {
        return nil, err
    }
    
    enforcer, err := casbin.NewEnforcer("configs/rbac_model.conf", adapter)
    if err != nil {
        return nil, err
    }
    
    return &Enforcer{enforcer: enforcer}, nil
}

func (e *Enforcer) Enforce(role, permission string) bool {
    return e.enforcer.Enforce(role, permission)
}
```

## 3. Capacidades REST Avanzadas

### Estructura Detallada
```
internal/infrastructure/http/query/
├── options.go         # Definición de opciones de consulta
├── pagination.go      # Paginación
├── sorting.go         # Ordenamiento
├── filtering.go       # Filtrado
├── fields.go          # Selección de campos
├── expanding.go       # Expansión de relaciones
└── parser.go          # Parsing de query parameters

internal/domain/repository/
└── options.go         # Opciones de repositorio genéricas
```

### Implementación Recomendada

1. **Opciones de Consulta**
```go
// internal/infrastructure/http/query/options.go
type QueryOptions struct {
    Pagination *PaginationOptions `json:"pagination,omitempty"`
    Sorting    []*SortOption      `json:"sorting,omitempty"`
    Filtering  []*FilterOption    `json:"filtering,omitempty"`
    Fields     []string           `json:"fields,omitempty"`
    Expands    []string           `json:"expands,omitempty"`
    Search     string             `json:"search,omitempty"`
}

type PaginationOptions struct {
    Page       int `json:"page" query:"page"`
    PageSize   int `json:"page_size" query:"limit"`
    Offset     int `json:"offset,omitempty"`
}

type SortOption struct {
    Field     string `json:"field"`
    Direction string `json:"direction"` // asc, desc
}

type FilterOption struct {
    Field    string      `json:"field"`
    Operator string      `json:"operator"` // eq, neq, gt, gte, lt, lte, in, nin, like
    Value    interface{} `json:"value"`
}
```

2. **Parser de Query Parameters**
```go
// internal/infrastructure/http/query/parser.go
func ParseQueryOptions(c *fiber.Ctx) (*QueryOptions, error) {
    opts := &QueryOptions{}
    
    // Paginación
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "20"))
    opts.Pagination = &PaginationOptions{
        Page:     page,
        PageSize: limit,
        Offset:   (page - 1) * limit,
    }
    
    // Ordenamiento
    if sort := c.Query("sort"); sort != "" {
        parts := strings.Split(sort, ",")
        opts.Sorting = make([]*SortOption, 0, len(parts))
        
        for _, part := range parts {
            sortParts := strings.Split(strings.TrimSpace(part), ":")
            field := sortParts[0]
            direction := "asc"
            if len(sortParts) > 1 {
                direction = strings.ToLower(sortParts[1])
            }
            
            opts.Sorting = append(opts.Sorting, &SortOption{
                Field:     field,
                Direction: direction,
            })
        }
    }
    
    // Filtrado
    if filter := c.Query("filter"); filter != "" {
        // Implementar lógica de parsing de filtros
        // Ejemplo: filter=name:eq:John,age:gt:25
    }
    
    // Selección de campos
    if fields := c.Query("fields"); fields != "" {
        opts.Fields = strings.Split(fields, ",")
    }
    
    // Expansión de relaciones
    if expand := c.Query("expand"); expand != "" {
        opts.Expands = strings.Split(expand, ",")
    }
    
    // Búsqueda
    opts.Search = c.Query("search")
    
    return opts, nil
}
```

3. **Implementación en Repositorio**
```go
// internal/domain/repository/employee_repository.go
type EmployeeRepository interface {
    // Métodos existentes
    Create(ctx context.Context, employee *entity.Employee) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
    FindAll(ctx context.Context) ([]*entity.Employee, error)
    Update(ctx context.Context, employee *entity.Employee) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // Nuevo método con opciones avanzadas
    FindWithOptions(ctx context.Context, opts *query.QueryOptions) ([]*entity.Employee, int64, error)
}
```

4. **Implementación GORM**
```go
// internal/infrastructure/database/postgres/employee_repository.go
func (r *employeeRepository) FindWithOptions(ctx context.Context, opts *query.QueryOptions) ([]*entity.Employee, int64, error) {
    var employees []*entity.Employee
    var count int64
    
    db := r.db.WithContext(ctx)
    
    // Construcción de la consulta base
    q := db.Model(&entity.Employee{})
    
    // Aplicar filtros
    if opts.Filtering != nil {
        for _, filter := range opts.Filtering {
            switch filter.Operator {
            case "eq":
                q = q.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
            case "neq":
                q = q.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
            case "gt":
                q = q.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
            // Más operadores...
            }
        }
    }
    
    // Aplicar búsqueda
    if opts.Search != "" {
        q = q.Where("name ILIKE ?", "%"+opts.Search+"%")
    }
    
    // Contar resultados totales
    q.Count(&count)
    
    // Aplicar ordenamiento
    if opts.Sorting != nil {
        for _, sort := range opts.Sorting {
            direction := sort.Direction
            if direction != "asc" && direction != "desc" {
                direction = "asc"
            }
            q = q.Order(fmt.Sprintf("%s %s", sort.Field, direction))
        }
    }
    
    // Aplicar paginación
    if opts.Pagination != nil {
        q = q.Offset(opts.Pagination.Offset).Limit(opts.Pagination.PageSize)
    }
    
    // Aplicar expansión de relaciones
    if opts.Expands != nil {
        for _, expand := range opts.Expands {
            q = q.Preload(strings.Title(expand))
        }
    }
    
    // Ejecutar consulta
    result := q.Find(&employees)
    if result.Error != nil {
        return nil, 0, result.Error
    }
    
    return employees, count, nil
}
```

## 4. Comunicación en Tiempo Real (WebSockets)

### Estructura Detallada
```
internal/infrastructure/websocket/
├── hub.go        # Centro de gestión de conexiones
├── client.go     # Gestión de clientes individuales
├── message.go    # Definición de mensajes
├── handler.go    # Manejadores HTTP de WebSocket
└── subscription/ # Gestión de suscripciones
    ├── manager.go    # Administrador de suscripciones
    ├── topic.go      # Definición de tópicos
    └── subscriber.go # Suscriptores

internal/domain/event/
├── employee_events.go    # Eventos específicos de empleados
└── dispatcher.go         # Despachador de eventos
```

### Implementación Recomendada

1. **Hub Central**
```go
// internal/infrastructure/websocket/hub.go
type Hub struct {
    clients       map[*Client]bool
    register      chan *Client
    unregister    chan *Client
    broadcast     chan *Message
    topics        map[string]map[*Client]bool
    subscribe     chan *Subscription
    unsubscribe   chan *Subscription
    mutex         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:      make(map[*Client]bool),
        register:     make(chan *Client),
        unregister:   make(chan *Client),
        broadcast:    make(chan *Message),
        topics:       make(map[string]map[*Client]bool),
        subscribe:    make(chan *Subscription),
        unsubscribe:  make(chan *Subscription),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                
                // Eliminar de tópicos
                for topic, clients := range h.topics {
                    if _, ok := clients[client]; ok {
                        delete(h.topics[topic], client)
                    }
                }
            }
            h.mutex.Unlock()
            
        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
            
        case subscription := <-h.subscribe:
            h.mutex.Lock()
            if _, ok := h.topics[subscription.Topic]; !ok {
                h.topics[subscription.Topic] = make(map[*Client]bool)
            }
            h.topics[subscription.Topic][subscription.Client] = true
            h.mutex.Unlock()
            
        case subscription := <-h.unsubscribe:
            h.mutex.Lock()
            if _, ok := h.topics[subscription.Topic]; ok {
                delete(h.topics[subscription.Topic], subscription.Client)
            }
            h.mutex.Unlock()
        }
    }
}

func (h *Hub) PublishToTopic(topic string, message *Message) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    
    if clients, ok := h.topics[topic]; ok {
        for client := range clients {
            select {
            case client.send <- message:
            default:
                close(client.send)
                delete(h.clients, client)
                delete(clients, client)
            }
        }
    }
}
```

2. **Cliente WebSocket**
```go
// internal/infrastructure/websocket/client.go
type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan *Message
    userID   uuid.UUID
    metadata map[string]interface{}
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { 
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil 
    })
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }
        
        // Procesar mensaje (comando, suscripción, etc.)
        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("error parsing message: %v", err)
            continue
        }
        
        // Manejar comando
        if msg.Type == "subscribe" {
            c.hub.subscribe <- &Subscription{
                Client: c,
                Topic:  msg.Topic,
            }
        }
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            
            payload, _ := json.Marshal(message)
            w.Write(payload)
            
            if err := w.Close(); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

3. **Integración con Eventos de Dominio**
```go
// internal/domain/event/dispatcher.go
type EventDispatcher struct {
    listeners map[string][]EventListener
    mu        sync.RWMutex
}

type EventListener func(ctx context.Context, event Event) error

type Event interface {
    EventName() string
    Payload() interface{}
    OccurredAt() time.Time
}

func (d *EventDispatcher) Register(eventName string, listener EventListener) {
    d.mu.Lock()
    defer d.mu.Unlock()
    
    if d.listeners == nil {
        d.listeners = make(map[string][]EventListener)
    }
    
    d.listeners[eventName] = append(d.listeners[eventName], listener)
}

func (d *EventDispatcher) Dispatch(ctx context.Context, event Event) error {
    d.mu.RLock()
    defer d.mu.RUnlock()
    
    if listeners, ok := d.listeners[event.EventName()]; ok {
        for _, listener := range listeners {
            if err := listener(ctx, event); err != nil {
                return err
            }
        }
    }
    
    return nil
}
```

4. **Emisión de Eventos a WebSockets**
```go
// internal/infrastructure/websocket/handler.go
func RegisterEventHandlers(eventDispatcher *event.EventDispatcher, hub *Hub) {
    // Cuando se crea un empleado, notificar a los suscriptores
    eventDispatcher.Register("employee.created", func(ctx context.Context, evt event.Event) error {
        employee := evt.Payload().(*entity.Employee)
        
        // Crear mensaje para WebSocket
        message := &Message{
            Type:    "event",
            Topic:   "employees",
            Action:  "created",
            Payload: employee,
        }
        
        // Publicar a todos los suscriptores del tópico "employees"
        hub.PublishToTopic("employees", message)
        
        return nil
    })
}
```

## 5. Arquitectura Orientada a Eventos

### Estructura Detallada
```
internal/domain/event/
├── events.go           # Definiciones básicas de eventos
├── employee_events.go  # Eventos específicos de empleados
└── dispatcher.go       # Dispatcher de eventos

internal/infrastructure/eventbus/
├── bus.go              # Bus de eventos central
├── memory/             # Implementación en memoria
├── kafka/              # Implementación con Kafka (opcional)
└── nats/               # Implementación con NATS (opcional)

internal/infrastructure/messaging/
├── producer.go         # Productor de mensajes
├── consumer.go         # Consumidor de mensajes
└── handler.go          # Manejadores de mensajes
```

### Implementación Recomendada

1. **Eventos de Dominio**
```go
// internal/domain/event/employee_events.go
type EmployeeCreatedEvent struct {
    Employee   *entity.Employee
    OccurredAt time.Time
}

func (e EmployeeCreatedEvent) EventName() string {
    return "employee.created"
}

func (e EmployeeCreatedEvent) Payload() interface{} {
    return e.Employee
}

func (e EmployeeCreatedEvent) OccurredAt() time.Time {
    return e.OccurredAt
}

type EmployeeUpdatedEvent struct {
    Employee   *entity.Employee
    OldValues  map[string]interface{}
    OccurredAt time.Time
}

func (e EmployeeUpdatedEvent) EventName() string {
    return "employee.updated"
}

func (e EmployeeUpdatedEvent) Payload() interface{} {
    return map[string]interface{}{
        "employee":   e.Employee,
        "old_values": e.OldValues,
    }
}

func (e EmployeeUpdatedEvent) OccurredAt() time.Time {
    return e.OccurredAt
}
```

2. **Bus de Eventos**
```go
// internal/infrastructure/eventbus/bus.go
type EventBus interface {
    Publish(ctx context.Context, event event.Event) error
    Subscribe(eventName string, handler event.EventListener) error
    Unsubscribe(eventName string, handler event.EventListener) error
}

// internal/infrastructure/eventbus/memory/bus.go
type InMemoryEventBus struct {
    dispatcher *event.EventDispatcher
}

func NewInMemoryEventBus() *InMemoryEventBus {
    return &InMemoryEventBus{
        dispatcher: &event.EventDispatcher{},
    }
}

func (b *InMemoryEventBus) Publish(ctx context.Context, event event.Event) error {
    return b.dispatcher.Dispatch(ctx, event)
}

func (b *InMemoryEventBus) Subscribe(eventName string, handler event.EventListener) error {
    b.dispatcher.Register(eventName, handler)
    return nil
}

func (b *InMemoryEventBus) Unsubscribe(eventName string, handler event.EventListener) error {
    // Implementar lógica de anulación de suscripción
    return nil
}
```

3. **Integración en Casos de Uso**
```go
// internal/usecase/command/create_employee.go
type CreateEmployeeCommand struct {
    Name string
}

type CreateEmployeeHandler struct {
    employeeRepo repository.EmployeeRepository
    eventBus     eventbus.EventBus
}

func NewCreateEmployeeHandler(repo repository.EmployeeRepository, bus eventbus.EventBus) *CreateEmployeeHandler {
    return &CreateEmployeeHandler{
        employeeRepo: repo,
        eventBus:     bus,
    }
}

func (h *CreateEmployeeHandler) Handle(ctx context.Context, cmd CreateEmployeeCommand) (*entity.Employee, error) {
    if cmd.Name == "" {
        return nil, usecase.ErrInvalidInput
    }
    
    employee := entity.NewEmployee(cmd.Name)
    
    if err := h.employeeRepo.Create(ctx, employee); err != nil {
        return nil, err
    }
    
    // Publicar evento
    event := event.EmployeeCreatedEvent{
        Employee:   employee,
        OccurredAt: time.Now(),
    }
    
    if err := h.eventBus.Publish(ctx, event); err != nil {
        log.Printf("Error publishing event: %v", err)
    }
    
    return employee, nil
}
```

4. **Manejadores de Eventos**
```go
// internal/infrastructure/messaging/handler.go
func RegisterEventHandlers(eventBus eventbus.EventBus) {
    // Registrar handler para enviar notificaciones
    eventBus.Subscribe("employee.created", func(ctx context.Context, evt event.Event) error {
        employee := evt.Payload().(*entity.Employee)
        
        // Enviar notificación por email, SMS, etc.
        log.Printf("Employee created: %s", employee.Name)
        
        return nil
    })
    
    // Registrar handler para actualizar caché
    eventBus.Subscribe("employee.created", func(ctx context.Context, evt event.Event) error {
        employee := evt.Payload().(*entity.Employee)
        
        // Actualizar caché
        cache.Set(fmt.Sprintf("employee:%s", employee.ID), employee, 1*time.Hour)
        
        return nil
    })
}
```

## 6. CQRS (Command Query Responsibility Segregation)

### Estructura Detallada
```
internal/usecase/
├── command/                  # Comandos (escritura)
│   ├── create_employee.go    # Comando para crear empleado
│   ├── update_employee.go    # Comando para actualizar empleado
│   └── delete_employee.go    # Comando para eliminar empleado
└── query/                    # Consultas (lectura)
    ├── get_employee.go       # Consulta para obtener empleado
    ├── list_employees.go     # Consulta para listar empleados
    └── search_employees.go   # Consulta para buscar empleados
```

### Implementación Recomendada

1. **Comandos (Escritura)**
```go
// internal/usecase/command/create_employee.go
type CreateEmployeeCommand struct {
    Name string
}

type CreateEmployeeHandler struct {
    employeeRepo repository.EmployeeRepository
    eventBus     eventbus.EventBus
}

func (h *CreateEmployeeHandler) Handle(ctx context.Context, cmd CreateEmployeeCommand) (*entity.Employee, error) {
    // Implementación...
}

// internal/usecase/command/update_employee.go
type UpdateEmployeeCommand struct {
    ID   uuid.UUID
    Name string
}

type UpdateEmployeeHandler struct {
    employeeRepo repository.EmployeeRepository
    eventBus     eventbus.EventBus
}

func (h *UpdateEmployeeHandler) Handle(ctx context.Context, cmd UpdateEmployeeCommand) (*entity.Employee, error) {
    // Implementación...
}
```

2. **Queries (Lectura)**
```go
// internal/usecase/query/get_employee.go
type GetEmployeeQuery struct {
    ID uuid.UUID
}

type GetEmployeeHandler struct {
    employeeRepo repository.EmployeeRepository
    cache        cache.Cache
}

func (h *GetEmployeeHandler) Handle(ctx context.Context, query GetEmployeeQuery) (*entity.Employee, error) {
    // Intentar obtener de caché primero
    cacheKey := fmt.Sprintf("employee:%s", query.ID)
    if cached, found := h.cache.Get(cacheKey); found {
        return cached.(*entity.Employee), nil
    }
    
    // Si no está en caché, obtener de base de datos
    employee, err := h.employeeRepo.FindByID(ctx, query.ID)
    if err != nil {
        return nil, err
    }
    
    // Guardar en caché
    h.cache.Set(cacheKey, employee, 1*time.Hour)
    
    return employee, nil
}

// internal/usecase/query/list_employees.go
type ListEmployeesQuery struct {
    Options *query.QueryOptions
}

type ListEmployeesHandler struct {
    employeeRepo repository.EmployeeRepository
}

func (h *ListEmployeesHandler) Handle(ctx context.Context, query ListEmployeesQuery) ([]*entity.Employee, int64, error) {
    return h.employeeRepo.FindWithOptions(ctx, query.Options)
}
```

3. **Integración en Handlers HTTP**
```go
// internal/infrastructure/http/handler/employee_handler.go
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
    var req dto.CreateEmployeeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
            Error:   "Invalid request body",
            Message: err.Error(),
        })
    }
    
    cmd := command.CreateEmployeeCommand{
        Name: req.Name,
    }
    
    handler := h.container.GetCreateEmployeeHandler()
    employee, err := handler.Handle(c.Context(), cmd)
    if err != nil {
        // Manejar error...
    }
    
    return c.Status(fiber.StatusCreated).JSON(dto.SuccessResponse{
        Message: "Employee created successfully",
        Data:    dto.ToEmployeeResponse(employee),
    })
}

func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
    opts, err := query.ParseQueryOptions(c)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
            Error:   "Invalid query parameters",
            Message: err.Error(),
        })
    }
    
    q := query.ListEmployeesQuery{
        Options: opts,
    }
    
    handler := h.container.GetListEmployeesHandler()
    employees, total, err := handler.Handle(c.Context(), q)
    if err != nil {
        // Manejar error...
    }
    
    return c.JSON(dto.PaginatedResponse{
        Message: "Employees retrieved successfully",
        Data:    dto.ToEmployeeResponses(employees),
        Meta: dto.PaginationMeta{
            Page:       opts.Pagination.Page,
            PageSize:   opts.Pagination.PageSize,
            TotalItems: total,
            TotalPages: int(math.Ceil(float64(total) / float64(opts.Pagination.PageSize))),
        },
    })
}
```

## 7. Validación Avanzada

### Estructura Detallada
```
internal/infrastructure/validation/
├── validator.go         # Validador central
├── rules.go             # Reglas de validación personalizadas
└── employee_rules.go    # Reglas específicas para empleados
```

### Implementación Recomendada

1. **Validador Central**
```go
// internal/infrastructure/validation/validator.go
type Validator struct {
    validate *validator.Validate
}

func NewValidator() *Validator {
    v := validator.New()
    
    // Registrar validadores personalizados
    v.RegisterValidation("unique_email", validateUniqueEmail)
    
    return &Validator{validate: v}
}

func (v *Validator) Struct(s interface{}) error {
    return v.validate.Struct(s)
}

func (v *Validator) Var(field interface{}, tag string) error {
    return v.validate.Var(field, tag)
}

// Función para validar email único
func validateUniqueEmail(fl validator.FieldLevel) bool {
    email := fl.Field().String()
    
    // Lógica para verificar unicidad en la base de datos
    // ...
    
    return true // Cambiar según el resultado de la validación
}
```

2. **Reglas de Validación para DTOs**
```go
// internal/infrastructure/http/dto/employee_dto.go
type CreateEmployeeRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=255"`
    Email string `json:"email" validate:"required,email,unique_email"`
    Role  string `json:"role" validate:"required,oneof=admin user manager"`
}

// internal/infrastructure/http/handler/employee_handler.go
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
    var req dto.CreateEmployeeRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
            Error:   "Invalid request body",
            Message: err.Error(),
        })
    }
    
    // Validar request
    if err := h.validator.Struct(req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(dto.ValidationErrorResponse{
            Error:   "Validation failed",
            Details: formatValidationErrors(err),
        })
    }
    
    // Continuar con la lógica...
}

// Función para formatear errores de validación
func formatValidationErrors(err error) map[string]string {
    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        errMap := make(map[string]string)
        
        for _, e := range validationErrs {
            field := strings.ToLower(e.Field())
            errMap[field] = getErrorMessage(e)
        }
        
        return errMap
    }
    
    return map[string]string{"general": err.Error()}
}

func getErrorMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required":
        return "This field is required"
    case "email":
        return "Invalid email format"
    case "min":
        return fmt.Sprintf("Must be at least %s characters long", e.Param())
    case "max":
        return fmt.Sprintf("Must not be longer than %s characters", e.Param())
    case "unique_email":
        return "Email already exists"
    default:
        return fmt.Sprintf("Failed validation on %s", e.Tag())
    }
}
```

## 8. Cache y Rendimiento

### Estructura Detallada
```
internal/infrastructure/cache/
├── cache.go       # Interfaz de caché
├── memory.go      # Implementación en memoria
└── redis.go       # Implementación con Redis

internal/usecase/query/
└── base_query.go  # Base para consultas con caché
```

### Implementación Recomendada

1. **Interfaz de Caché**
```go
// internal/infrastructure/cache/cache.go
type Cache interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, expiration time.Duration)
    Delete(key string)
    Clear()
}
```

2. **Implementación Redis**
```go
// internal/infrastructure/cache/redis.go
type RedisCache struct {
    client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    
    return &RedisCache{client: client}
}

func (c *RedisCache) Get(key string) (interface{}, bool) {
    val, err := c.client.Get(context.Background(), key).Result()
    if err != nil {
        return nil, false
    }
    
    var result interface{}
    if err := json.Unmarshal([]byte(val), &result); err != nil {
        return nil, false
    }
    
    return result, true
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) {
    data, err := json.Marshal(value)
    if err != nil {
        return
    }
    
    c.client.Set(context.Background(), key, data, expiration)
}

func (c *RedisCache) Delete(key string) {
    c.client.Del(context.Background(), key)
}

func (c *RedisCache) Clear() {
    c.client.FlushAll(context.Background())
}
```

3. **Consulta con Caché**
```go
// internal/usecase/query/base_query.go
type CachedQueryHandler struct {
    cache cache.Cache
}

func (h *CachedQueryHandler) WithCache(key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
    // Intentar obtener de caché
    if cached, found := h.cache.Get(key); found {
        return cached, nil
    }
    
    // Ejecutar función si no está en caché
    result, err := fn()
    if err != nil {
        return nil, err
    }
    
    // Guardar en caché
    h.cache.Set(key, result, expiration)
    
    return result, nil
}
```

4. **Uso en Query Handler**
```go
// internal/usecase/query/get_employee.go
type GetEmployeeHandler struct {
    employeeRepo repository.EmployeeRepository
    cacheHandler *CachedQueryHandler
}

func (h *GetEmployeeHandler) Handle(ctx context.Context, query GetEmployeeQuery) (*entity.Employee, error) {
    cacheKey := fmt.Sprintf("employee:%s", query.ID)
    
    result, err := h.cacheHandler.WithCache(cacheKey, 1*time.Hour, func() (interface{}, error) {
        return h.employeeRepo.FindByID(ctx, query.ID)
    })
    if err != nil {
        return nil, err
    }
    
    return result.(*entity.Employee), nil
}
```

## 9. Observabilidad y Telemetría

### Estructura Detallada
```
internal/infrastructure/telemetry/
├── logger/           # Logging estructurado
│   ├── logger.go     # Interfaz de logger
│   └── zap.go        # Implementación con Zap
├── metrics/          # Métricas
│   ├── metrics.go    # Interfaz de métricas
│   └── prometheus.go # Implementación con Prometheus
└── tracing/          # Tracing distribuido
    ├── tracer.go     # Interfaz de tracer
    └── opentelemetry.go # Implementación con OpenTelemetry
```

### Implementación Recomendada

1. **Logger Estructurado**
```go
// internal/infrastructure/telemetry/logger/logger.go
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
    With(fields ...Field) Logger
}

type Field struct {
    Key   string
    Value interface{}
}

func String(key string, value string) Field {
    return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
    return Field{Key: key, Value: value}
}

func Any(key string, value interface{}) Field {
    return Field{Key: key, Value: value}
}
```

2. **Implementación con Zap**
```go
// internal/infrastructure/telemetry/logger/zap.go
type ZapLogger struct {
    logger *zap.Logger
}

func NewZapLogger(development bool) (*ZapLogger, error) {
    var logger *zap.Logger
    var err error
    
    if development {
        logger, err = zap.NewDevelopment()
    } else {
        logger, err = zap.NewProduction()
    }
    
    if err != nil {
        return nil, err
    }
    
    return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
    l.logger.Debug(msg, convertFields(fields)...)
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
    l.logger.Info(msg, convertFields(fields)...)
}

// ... más métodos

func (l *ZapLogger) With(fields ...Field) Logger {
    return &ZapLogger{
        logger: l.logger.With(convertFields(fields)...),
    }
}

func convertFields(fields []Field) []zap.Field {
    zapFields := make([]zap.Field, len(fields))
    for i, f := range fields {
        zapFields[i] = zap.Any(f.Key, f.Value)
    }
    return zapFields
}
```

3. **Middleware de Logging**
```go
// internal/infrastructure/http/middleware/logging.go
func LoggingMiddleware(logger telemetry.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        
        // Extraer información de la petición
        method := c.Method()
        path := c.Path()
        ip := c.IP()
        
        // Procesar la petición
        err := c.Next()
        
        // Registrar la respuesta
        duration := time.Since(start)
        status := c.Response().StatusCode()
        
        logLevel := "info"
        if status >= 500 {
            logLevel = "error"
        } else if status >= 400 {
            logLevel = "warn"
        }
        
        fields := []telemetry.Field{
            telemetry.String("method", method),
            telemetry.String("path", path),
            telemetry.String("ip", ip),
            telemetry.Int("status", status),
            telemetry.String("duration", duration.String()),
        }
        
        switch logLevel {
        case "error":
            logger.Error("request failed", fields...)
        case "warn":
            logger.Warn("request warning", fields...)
        default:
            logger.Info("request completed", fields...)
        }
        
        return err
    }
}
```

## 10. Conclusión de la Arquitectura Todo Terreno

La arquitectura propuesta ofrece:

1. **Alta Modularidad**: Cada componente tiene una responsabilidad única y está claramente separado.
2. **Flexibilidad**: Facilidad para cambiar implementaciones o agregar nuevas funcionalidades.
3. **Testabilidad**: Interfaces claras que facilitan el uso de mocks.
4. **Escalabilidad**: Preparada para crecer horizontal y verticalmente.
5. **Preparación para Microservicios**: Puede evolucionar fácilmente a una arquitectura distribuida.
6. **Desacoplamiento**: Dependencias hacia abstracciones, no implementaciones concretas.
7. **Experiencia de Desarrollo**: Estructura predecible y consistente.

Para implementar esta arquitectura, recomiendo comenzar por:

1. Autenticación y autorización (seguridad primero)
2. Capacidades REST avanzadas (paginación, filtrado, ordenamiento)
3. Validación robusta
4. Arquitectura CQRS
5. Agregar WebSockets y eventos

Estas mejoras convertirán tu actual proyecto en una plataforma verdaderamente "todo terreno" lista para cualquier requerimiento futuro.

Similar code found with 5 license types