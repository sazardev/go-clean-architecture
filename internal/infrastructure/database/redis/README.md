# redis/ - Implementación Redis

Implementación de caché y almacenamiento de sesiones usando Redis.

## Responsabilidades

- Implementar caché de aplicación
- Gestionar sesiones de usuario
- Manejar rate limiting
- Almacenar datos temporales
- Implementar pub/sub para eventos

## Casos de Uso

### Caché
- Resultados de queries costosas
- Datos frecuentemente accedidos
- Respuestas de APIs externas
- Configuraciones de aplicación

### Sesiones
- Tokens de autenticación
- Estado de usuario temporal
- Carritos de compra
- Preferencias temporales

### Rate Limiting
- Límites por usuario/IP
- Throttling de APIs
- Protección contra ataques

## Archivos Futuros

- **`cache_repository.go`** - Caché genérico
- **`session_repository.go`** - Gestión de sesiones
- **`rate_limiter.go`** - Rate limiting
- **`pubsub.go`** - Pub/Sub para eventos
- **`queue.go`** - Colas de tareas

## Implementación

### Cache Repository
```go
type RedisCacheRepository struct {
    client *redis.Client
}

func NewRedisCacheRepository(client *redis.Client) *RedisCacheRepository {
    return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return domain.ErrCacheNotFound
        }
        return err
    }
    
    return json.Unmarshal([]byte(data), dest)
}

func (r *RedisCacheRepository) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}

func (r *RedisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
    count, err := r.client.Exists(ctx, key).Result()
    return count > 0, err
}
```

### Session Repository
```go
type RedisSessionRepository struct {
    client *redis.Client
    prefix string
}

func (r *RedisSessionRepository) CreateSession(ctx context.Context, userID uuid.UUID, data map[string]interface{}) (string, error) {
    sessionID := uuid.New().String()
    key := r.prefix + sessionID
    
    sessionData := map[string]interface{}{
        "user_id":    userID.String(),
        "created_at": time.Now(),
        "data":       data,
    }
    
    err := r.Set(ctx, key, sessionData, 24*time.Hour)
    return sessionID, err
}

func (r *RedisSessionRepository) GetSession(ctx context.Context, sessionID string) (*SessionData, error) {
    key := r.prefix + sessionID
    var sessionData SessionData
    err := r.Get(ctx, key, &sessionData)
    return &sessionData, err
}
```

### Rate Limiter
```go
type RedisRateLimiter struct {
    client *redis.Client
}

func (r *RedisRateLimiter) IsAllowed(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    pipe := r.client.Pipeline()
    
    // Incrementar contador
    pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, window)
    
    results, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    count := results[0].(*redis.IntCmd).Val()
    return count <= int64(limit), nil
}
```

## Configuración

Variables de entorno:
- `REDIS_URL` - URL de conexión completa
- `REDIS_PASSWORD` - Contraseña (opcional)
- `REDIS_DB` - Número de base de datos
- `REDIS_MAX_RETRIES` - Reintentos máximos
- `REDIS_POOL_SIZE` - Tamaño del pool

## Patrones de Keys

- Cache: `cache:{type}:{id}`
- Sessions: `session:{session_id}`
- Rate limiting: `rate_limit:{user_id}:{endpoint}`
- Temporary data: `temp:{type}:{id}`
