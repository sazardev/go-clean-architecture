# database/ - Capa de Persistencia

Implementaciones concretas de repositorios y manejo de múltiples bases de datos.

## Responsabilidades

- Implementar interfaces de repositorio del dominio
- Manejar conexiones a diferentes bases de datos
- Gestionar transacciones y consistencia
- Optimizar queries y performance
- Manejar migraciones y esquemas

## Estructura

- **`postgres/`** - Implementación para PostgreSQL
- **`mongodb/`** - Implementación para MongoDB  
- **`redis/`** - Implementación para Redis (caché/sesiones)
- **`factory/`** - Factory pattern para crear repositorios
- **`connection.go`** - Gestión de conexiones
- **`employee_repository.go`** - Implementación actual (se moverá a postgres/)

## Patrón Repository

Cada base de datos implementa las mismas interfaces del dominio:

```go
// Dominio define la interfaz
type EmployeeRepository interface {
    Create(ctx context.Context, employee *entity.Employee) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
    // ... más métodos
}

// PostgreSQL implementa la interfaz
type PostgresEmployeeRepository struct {
    db *gorm.DB
}

func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
    return r.db.WithContext(ctx).Create(employee).Error
}
```

## Factory Pattern

```go
type RepositoryFactory interface {
    CreateEmployeeRepository() domain.EmployeeRepository
    CreateUserRepository() domain.UserRepository
}

type PostgresRepositoryFactory struct {
    db *gorm.DB
}

func (f *PostgresRepositoryFactory) CreateEmployeeRepository() domain.EmployeeRepository {
    return &PostgresEmployeeRepository{db: f.db}
}
```

## Configuración

Variables de entorno:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - PostgreSQL
- `MONGODB_URI` - MongoDB connection string
- `REDIS_URL` - Redis connection string
- `DB_MAX_CONNECTIONS` - Pool de conexiones
- `DB_SSL_MODE` - Modo SSL

## Características

- **Multi-database**: Soporte para diferentes motores
- **Connection pooling**: Gestión eficiente de conexiones
- **Transactions**: Soporte completo para transacciones
- **Migrations**: Sistema de migraciones automáticas
- **Health checks**: Verificación de estado de conexiones
