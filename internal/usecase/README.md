# usecase/ - Capa de Casos de Uso

Contiene la lógica de aplicación que orquesta el dominio para cumplir requerimientos específicos.

## Responsabilidades

- Implementar casos de uso específicos de la aplicación
- Orquestar entidades y servicios de dominio
- Manejar transacciones y coordinación
- Implementar reglas de autorización de aplicación
- Convertir entre DTOs y entidades de dominio

## Estructura

- **`command/`** - Casos de uso de escritura (CQS - Command)
- **`query/`** - Casos de uso de lectura (CQS - Query)  
- **`common/`** - Utilidades y código compartido

## Principios

- **Single Responsibility**: Cada caso de uso tiene una responsabilidad específica
- **Dependencias hacia adentro**: Solo conoce el dominio
- **Orquestación**: Coordina pero no contiene lógica de negocio
- **Transaccional**: Maneja la consistencia de las operaciones

## Patrón CQRS

Separamos las operaciones de lectura y escritura:

### Commands (Escritura)
- Modifican estado del sistema
- Retornan confirmación de éxito/fallo
- Pueden disparar eventos
- Validaciones de negocio

### Queries (Lectura)
- Solo leen datos
- Optimizadas para casos específicos
- Pueden retornar DTOs especializados
- No modifican estado

## Estructura de Caso de Uso

```go
type CreateEmployeeUseCase struct {
    employeeRepo domain.EmployeeRepository
    eventBus     EventBus
    logger       Logger
}

func (uc *CreateEmployeeUseCase) Execute(
    ctx context.Context, 
    cmd CreateEmployeeCommand,
) (*CreateEmployeeResult, error) {
    // 1. Validar entrada
    // 2. Crear entidad de dominio
    // 3. Persistir usando repositorio
    // 4. Publicar eventos
    // 5. Retornar resultado
}
```
