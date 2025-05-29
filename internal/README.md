# internal/ - Código Interno de la Aplicación

Contiene toda la lógica de negocio y componentes internos del proyecto siguiendo Clean Architecture.

## Estructura

- **`domain/`** - Capa de Dominio: Entidades, interfaces y reglas de negocio puras
- **`usecase/`** - Capa de Casos de Uso: Lógica de aplicación y orquestación
- **`infrastructure/`** - Capa de Infraestructura: Implementaciones concretas
- **`pkg/`** - Utilidades y herramientas compartidas internamente

## Principios de Clean Architecture

1. **Independencia de Frameworks**: El núcleo no depende de frameworks externos
2. **Testabilidad**: Cada capa puede ser probada independientemente
3. **Independencia de UI**: La lógica de negocio no conoce la presentación
4. **Independencia de Base de Datos**: El dominio no conoce detalles de persistencia
5. **Independencia de Agentes Externos**: Sin dependencias externas en el core

## Flujo de Dependencias

```
External → Infrastructure → UseCase → Domain
```

Las dependencias apuntan hacia adentro, nunca hacia afuera.
