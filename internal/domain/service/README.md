# service/ - Servicios de Dominio

Contiene lógica de negocio que no pertenece a una entidad específica pero opera sobre el dominio.

## Responsabilidades

- Implementar reglas de negocio complejas que involucran múltiples entidades
- Coordinar operaciones entre diferentes entidades
- Encapsular algoritmos de dominio específicos
- Mantener la consistencia del dominio

## Cuándo Usar Servicios de Dominio

- Lógica que no pertenece naturalmente a una entidad
- Operaciones que requieren múltiples entidades
- Cálculos complejos específicos del dominio
- Reglas de negocio que cambian frecuentemente

## Servicios Futuros

- **`employee_domain_service.go`** - Lógica compleja de empleados
- **`payroll_service.go`** - Cálculos de nómina
- **`promotion_service.go`** - Reglas de promoción
- **`department_service.go`** - Gestión de departamentos

## Principios

- **Stateless**: No mantienen estado entre llamadas
- **Puros**: Solo dependen del dominio
- **Específicos**: Cada servicio tiene una responsabilidad clara
- **Testeable**: Lógica pura sin dependencias externas

## Ejemplo

```go
type EmployeeDomainService struct {
    // Sin dependencias externas, solo del dominio
}

func (s *EmployeeDomainService) CalculateSalaryIncrease(
    employee *entity.Employee,
    performanceRating PerformanceRating,
) (amount Money, err error) {
    // Lógica compleja de negocio para aumento salarial
    // basada en reglas del dominio
}

func (s *EmployeeDomainService) ValidatePromotion(
    employee *entity.Employee,
    newPosition *entity.Position,
) error {
    // Validar reglas de promoción
}
```
