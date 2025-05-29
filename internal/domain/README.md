# domain/ - Capa de Dominio

La capa más interna de Clean Architecture. Contiene las reglas de negocio puras y entidades core.

## Responsabilidades

- Definir entidades de negocio y sus invariantes
- Establecer contratos (interfaces) para servicios
- Especificar eventos de dominio
- Definir objetos de valor inmutables
- Contener la lógica de negocio más crítica

## Estructura

- **`entity/`** - Entidades principales del dominio
- **`repository/`** - Interfaces para acceso a datos
- **`service/`** - Servicios de dominio para lógica compleja
- **`event/`** - Eventos que ocurren en el dominio
- **`valueobject/`** - Objetos de valor inmutables

## Principios

- **Sin dependencias externas**: No importa nada de infraestructura
- **Agnóstico a tecnología**: No conoce bases de datos, frameworks, etc.
- **Testeable**: Lógica pura que se puede probar fácilmente
- **Invariantes**: Las entidades mantienen su consistencia

## Ejemplo

```go
// Entidad Employee con reglas de negocio
type Employee struct {
    ID   uuid.UUID
    Name string
}

// Método de dominio que valida reglas de negocio
func (e *Employee) ValidateName() error {
    if len(strings.TrimSpace(e.Name)) == 0 {
        return errors.New("employee name cannot be empty")
    }
    return nil
}
```
