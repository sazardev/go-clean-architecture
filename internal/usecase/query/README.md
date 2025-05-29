# query/ - Casos de Uso de Lectura

Implementa todos los casos de uso de solo lectura optimizados para diferentes necesidades de consulta.

## Responsabilidades

- Implementar consultas específicas para la UI
- Optimizar rendimiento de lectura
- Formatear datos para presentación
- Implementar filtrado y paginación
- Cachear resultados cuando sea apropiado

## Características

- **Solo lectura**: No modifican estado del sistema
- **Optimizadas**: Pueden usar queries SQL específicas
- **Flexibles**: Soportan filtros, ordenamiento y paginación
- **Cacheables**: Resultados pueden cachearse
- **DTOs específicos**: Retornan estructuras optimizadas para cada caso

## Queries Futuras

### Employee Queries
- **`get_employee_by_id.go`** - Obtener empleado por ID
- **`list_employees.go`** - Listar empleados con filtros
- **`search_employees.go`** - Búsqueda de empleados
- **`get_employee_hierarchy.go`** - Jerarquía organizacional

### Department Queries
- **`list_departments.go`** - Listar departamentos
- **`get_department_stats.go`** - Estadísticas departamentales

### Reporting Queries
- **`get_employee_report.go`** - Reportes de empleados
- **`get_payroll_summary.go`** - Resumen de nómina

## Estructura de Query

```go
// Consulta de entrada
type ListEmployeesQuery struct {
    DepartmentID *uuid.UUID
    Position     *string
    IsActive     *bool
    Search       string
    Page         int
    PageSize     int
    SortBy       string
    SortOrder    string
}

// Resultado
type EmployeeListResult struct {
    Employees   []*EmployeeDTO
    TotalCount  int64
    Page        int
    PageSize    int
    TotalPages  int
}

// DTO específico para listado
type EmployeeDTO struct {
    ID           uuid.UUID
    Name         string
    Email        string
    Department   string
    Position     string
    IsActive     bool
    HiredAt      time.Time
}

// Caso de uso
type ListEmployeesUseCase struct {
    employeeRepo domain.EmployeeRepository
    cache        Cache
}

func (uc *ListEmployeesUseCase) Execute(
    ctx context.Context,
    query ListEmployeesQuery,
) (*EmployeeListResult, error) {
    // Implementación optimizada para lectura
}
```

## Optimizaciones

- Uso de índices específicos en base de datos
- Projections para evitar cargar datos innecesarios
- Caché de consultas frecuentes
- Paginación eficiente
- Read replicas para distribución de carga
