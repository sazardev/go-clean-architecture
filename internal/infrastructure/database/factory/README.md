# factory/ - Factory de Repositorios

Implementa el patrón Factory para crear instancias de repositorios según la configuración.

## Responsabilidades

- Crear repositorios según el tipo de base de datos configurado
- Abstraer la lógica de creación de implementaciones
- Permitir cambio de base de datos sin modificar código cliente
- Gestionar inyección de dependencias de repositorios

## Patrón Factory

Permite cambiar la implementación de persistencia sin modificar el código que usa los repositorios:

```go
// Interfaz del factory
type RepositoryFactory interface {
    CreateEmployeeRepository() domain.EmployeeRepository
    CreateUserRepository() domain.UserRepository
    CreateDepartmentRepository() domain.DepartmentRepository
    CreateAuditRepository() domain.AuditRepository
}

// Implementaciones específicas
type PostgresRepositoryFactory struct {
    db *gorm.DB
}

type MongoRepositoryFactory struct {
    db *mongo.Database
}

type HybridRepositoryFactory struct {
    postgresDB *gorm.DB
    mongoDB    *mongo.Database
    redisClient *redis.Client
}
```

## Archivos Futuros

- **`repository_factory.go`** - Interface y factory base
- **`postgres_factory.go`** - Factory para PostgreSQL
- **`mongo_factory.go`** - Factory para MongoDB
- **`hybrid_factory.go`** - Factory híbrido
- **`config.go`** - Configuración del factory

## Implementación

### Interface Base
```go
type RepositoryFactory interface {
    // Core repositories
    CreateEmployeeRepository() domain.EmployeeRepository
    CreateUserRepository() domain.UserRepository
    CreateDepartmentRepository() domain.DepartmentRepository
    CreateRoleRepository() domain.RoleRepository
    
    // Specialized repositories
    CreateAuditRepository() domain.AuditRepository
    CreateCacheRepository() domain.CacheRepository
    CreateSessionRepository() domain.SessionRepository
}
```

### PostgreSQL Factory
```go
type PostgresRepositoryFactory struct {
    db    *gorm.DB
    cache *redis.Client
}

func NewPostgresRepositoryFactory(db *gorm.DB, cache *redis.Client) *PostgresRepositoryFactory {
    return &PostgresRepositoryFactory{
        db:    db,
        cache: cache,
    }
}

func (f *PostgresRepositoryFactory) CreateEmployeeRepository() domain.EmployeeRepository {
    return postgres.NewEmployeeRepository(f.db)
}

func (f *PostgresRepositoryFactory) CreateCacheRepository() domain.CacheRepository {
    return redis.NewCacheRepository(f.cache)
}
```

### Hybrid Factory
```go
type HybridRepositoryFactory struct {
    postgresDB  *gorm.DB
    mongoDB     *mongo.Database
    redisClient *redis.Client
}

func (f *HybridRepositoryFactory) CreateEmployeeRepository() domain.EmployeeRepository {
    // Datos estructurados en PostgreSQL
    return postgres.NewEmployeeRepository(f.postgresDB)
}

func (f *HybridRepositoryFactory) CreateAuditRepository() domain.AuditRepository {
    // Logs en MongoDB
    return mongodb.NewAuditRepository(f.mongoDB)
}

func (f *HybridRepositoryFactory) CreateCacheRepository() domain.CacheRepository {
    // Caché en Redis
    return redis.NewCacheRepository(f.redisClient)
}
```

## Configuración

```go
type DatabaseConfig struct {
    Type     string // "postgres", "mongodb", "hybrid"
    Postgres PostgresConfig
    MongoDB  MongoConfig
    Redis    RedisConfig
}

func CreateRepositoryFactory(config DatabaseConfig) (RepositoryFactory, error) {
    switch config.Type {
    case "postgres":
        db, err := setupPostgres(config.Postgres)
        if err != nil {
            return nil, err
        }
        return NewPostgresRepositoryFactory(db, nil), nil
        
    case "mongodb":
        db, err := setupMongoDB(config.MongoDB)
        if err != nil {
            return nil, err
        }
        return NewMongoRepositoryFactory(db), nil
        
    case "hybrid":
        return setupHybridFactory(config)
        
    default:
        return nil, fmt.Errorf("unsupported database type: %s", config.Type)
    }
}
```

## Ventajas

- **Flexibilidad**: Cambio fácil de implementación
- **Testabilidad**: Factories mock para testing
- **Configurabilidad**: Basado en variables de entorno
- **Mantenibilidad**: Cambios centralizados
