# postgres/ - Implementación PostgreSQL

Implementaciones específicas de repositorios usando PostgreSQL con GORM.

## Responsabilidades

- Implementar repositorios del dominio usando PostgreSQL
- Optimizar queries SQL específicas
- Manejar transacciones complejas
- Implementar full-text search
- Gestionar índices y performance

## Archivos Futuros

- **`employee_repository.go`** - Repositorio de empleados
- **`user_repository.go`** - Repositorio de usuarios
- **`department_repository.go`** - Repositorio de departamentos
- **`role_repository.go`** - Repositorio de roles
- **`models.go`** - Modelos específicos de GORM
- **`migrations.go`** - Helpers para migraciones

## Características PostgreSQL

### Ventajas
- ACID compliance completo
- Soporte para JSON/JSONB
- Full-text search nativo
- Extensiones avanzadas
- Excelente performance en queries complejas

### Funcionalidades Utilizadas
- **JSONB**: Para datos flexibles
- **UUID**: Como primary keys
- **Índices**: B-tree, GIN, GiST
- **Views**: Para queries complejas
- **Triggers**: Para auditoría automática

## Implementación

```go
type PostgresEmployeeRepository struct {
    db *gorm.DB
}

func NewPostgresEmployeeRepository(db *gorm.DB) *PostgresEmployeeRepository {
    return &PostgresEmployeeRepository{db: db}
}

func (r *PostgresEmployeeRepository) Create(ctx context.Context, employee *entity.Employee) error {
    model := &EmployeeModel{
        ID:        employee.ID,
        Name:      employee.Name,
        Email:     employee.Email,
        CreatedAt: employee.CreatedAt,
        UpdatedAt: employee.UpdatedAt,
    }
    
    return r.db.WithContext(ctx).Create(model).Error
}

func (r *PostgresEmployeeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error) {
    var model EmployeeModel
    err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrEmployeeNotFound
        }
        return nil, err
    }
    
    return model.ToEntity(), nil
}
```

## Modelos GORM

```go
type EmployeeModel struct {
    ID        uuid.UUID  `gorm:"type:uuid;primary_key"`
    Name      string     `gorm:"not null"`
    Email     string     `gorm:"unique;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `gorm:"index"`
}

func (EmployeeModel) TableName() string {
    return "employees"
}

func (m *EmployeeModel) ToEntity() *entity.Employee {
    return &entity.Employee{
        ID:        m.ID,
        Name:      m.Name,
        Email:     m.Email,
        CreatedAt: m.CreatedAt,
        UpdatedAt: m.UpdatedAt,
    }
}
```

## Migraciones

Archivos SQL en `migrations/postgres/`:
- `001_create_employees_table.sql`
- `002_create_users_table.sql`
- `003_create_roles_table.sql`
