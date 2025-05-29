# cache/ - Sistema de Caché

Implementación de caché distribuido para mejorar el rendimiento de la aplicación.

## Responsabilidades

- Implementar estrategias de caché
- Gestionar invalidación de caché
- Optimizar tiempo de respuesta
- Reducir carga en base de datos
- Manejar caché distribuido

## Estrategias de Caché

### Cache-Aside (Lazy Loading)
```go
func GetEmployee(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
    // 1. Intentar obtener del caché
    cached, err := cache.Get(ctx, fmt.Sprintf("employee:%s", id))
    if err == nil {
        return cached, nil
    }
    
    // 2. Si no está en caché, obtener de BD
    employee, err := repository.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // 3. Guardar en caché
    cache.Set(ctx, fmt.Sprintf("employee:%s", id), employee, 1*time.Hour)
    
    return employee, nil
}
```

### Write-Through
```go
func UpdateEmployee(ctx context.Context, employee *entity.Employee) error {
    // 1. Actualizar en base de datos
    err := repository.Update(ctx, employee)
    if err != nil {
        return err
    }
    
    // 2. Actualizar en caché
    cache.Set(ctx, fmt.Sprintf("employee:%s", employee.ID), employee, 1*time.Hour)
    
    return nil
}
```

### Write-Behind (Write-Back)
```go
func UpdateEmployeeAsync(ctx context.Context, employee *entity.Employee) error {
    // 1. Actualizar en caché inmediatamente
    cache.Set(ctx, fmt.Sprintf("employee:%s", employee.ID), employee, 1*time.Hour)
    
    // 2. Marcar para escritura asíncrona
    queue.Push("employee_update", employee)
    
    return nil
}
```

## Archivos Futuros

- **`cache_service.go`** - Servicio principal de caché
- **`strategies.go`** - Estrategias de caché
- **`invalidation.go`** - Lógica de invalidación
- **`warmer.go`** - Pre-calentamiento de caché
- **`metrics.go`** - Métricas de hit/miss

## Implementación

```go
type CacheService struct {
    primary   Cache          // Redis principal
    secondary Cache          // Caché local (opcional)
    logger    Logger
    metrics   Metrics
}

func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
    // Intentar caché local primero
    if c.secondary != nil {
        if err := c.secondary.Get(ctx, key, dest); err == nil {
            c.metrics.IncrementHit("local")
            return nil
        }
    }
    
    // Intentar caché distribuido
    if err := c.primary.Get(ctx, key, dest); err == nil {
        c.metrics.IncrementHit("distributed")
        
        // Poblar caché local
        if c.secondary != nil {
            c.secondary.Set(ctx, key, dest, 5*time.Minute)
        }
        
        return nil
    }
    
    c.metrics.IncrementMiss()
    return domain.ErrCacheNotFound
}

func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    // Guardar en caché distribuido
    if err := c.primary.Set(ctx, key, value, expiration); err != nil {
        return err
    }
    
    // Guardar en caché local
    if c.secondary != nil {
        localExpiration := time.Minute * 5
        if expiration < localExpiration {
            localExpiration = expiration
        }
        c.secondary.Set(ctx, key, value, localExpiration)
    }
    
    return nil
}
```

## Invalidación

### Por Tags
```go
func (c *CacheService) InvalidateByTag(ctx context.Context, tag string) error {
    keys, err := c.primary.GetKeysByPattern(ctx, fmt.Sprintf("*:%s:*", tag))
    if err != nil {
        return err
    }
    
    for _, key := range keys {
        c.Delete(ctx, key)
    }
    
    return nil
}

// Uso: al actualizar un empleado, invalidar todos los cachés relacionados
c.InvalidateByTag(ctx, "employee:" + employeeID.String())
```

### Por TTL
```go
func (c *CacheService) SetWithTags(ctx context.Context, key string, value interface{}, expiration time.Duration, tags []string) error {
    // Guardar valor principal
    err := c.Set(ctx, key, value, expiration)
    if err != nil {
        return err
    }
    
    // Guardar referencias de tags
    for _, tag := range tags {
        tagKey := fmt.Sprintf("tag:%s", tag)
        c.primary.SAdd(ctx, tagKey, key)
        c.primary.Expire(ctx, tagKey, expiration)
    }
    
    return nil
}
```

## Configuración

Variables de entorno:
- `CACHE_PROVIDER` - Proveedor (redis, memory, hybrid)
- `CACHE_DEFAULT_TTL` - TTL por defecto
- `CACHE_MAX_SIZE` - Tamaño máximo del caché local
- `CACHE_COMPRESSION` - Habilitar compresión
