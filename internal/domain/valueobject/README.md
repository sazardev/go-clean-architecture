# valueobject/ - Objetos de Valor

Contiene objetos inmutables que representan conceptos del dominio sin identidad.

## Responsabilidades

- Encapsular datos que van juntos conceptualmente
- Implementar validaciones de formato y consistencia
- Proporcionar métodos de comparación y conversión
- Mantener inmutabilidad y validez

## Características

- **Sin identidad**: Se comparan por valor, no por referencia
- **Inmutables**: No cambian después de la creación
- **Autovalidantes**: Se validan en la construcción
- **Reutilizables**: Pueden usarse en múltiples entidades

## Objetos de Valor Futuros

- **`email.go`** - Dirección de email válida
- **`phone.go`** - Número de teléfono
- **`address.go`** - Dirección postal completa
- **`money.go`** - Valor monetario con moneda
- **`date_range.go`** - Rango de fechas
- **`username.go`** - Nombre de usuario válido

## Ejemplos

```go
// Email como Value Object
type Email struct {
    value string
}

func NewEmail(email string) (Email, error) {
    if !isValidEmail(email) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{value: strings.ToLower(email)}, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

// Money como Value Object
type Money struct {
    amount   decimal.Decimal
    currency string
}

func NewMoney(amount decimal.Decimal, currency string) (Money, error) {
    if amount.IsNegative() {
        return Money{}, errors.New("amount cannot be negative")
    }
    return Money{amount: amount, currency: currency}, nil
}
```

## Ventajas

- Eliminan duplicación de validaciones
- Hacen el código más expresivo
- Previenen errores de tipo en tiempo de compilación
- Facilitan testing y refactoring
