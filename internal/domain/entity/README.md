# entity/ - Entidades de Dominio

Contiene las entidades principales del negocio con sus propiedades y comportamientos.

## Responsabilidades

- Definir la estructura de datos de las entidades principales
- Implementar métodos de validación y lógica de negocio
- Mantener invariantes y reglas de consistencia
- Encapsular el estado y comportamiento del dominio

## Características de las Entidades

- **Identidad única**: Cada entidad tiene un identificador único (UUID)
- **Métodos de dominio**: Comportamientos que operan sobre la entidad
- **Validaciones**: Reglas que garantizan la consistencia
- **Inmutabilidad donde sea posible**: Estado protegido

## Entidades Actuales

- **`employee.go`** - Entidad principal para empleados

## Entidades Futuras

- **`user.go`** - Usuarios del sistema
- **`role.go`** - Roles y permisos
- **`department.go`** - Departamentos organizacionales
- **`position.go`** - Cargos/posiciones laborales

## Ejemplo de Entidad

```go
type Employee struct {
    ID          uuid.UUID
    Name        string
    Email       string
    Department  *Department
    Position    *Position
    HiredAt     time.Time
    IsActive    bool
}

func (e *Employee) Promote(newPosition *Position) error {
    // Lógica de negocio para promoción
}
```
