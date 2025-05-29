# infrastructure/ - Capa de Infraestructura

Implementaciones concretas de todas las interfaces definidas en las capas superiores.

## Responsabilidades

- Implementar interfaces de repositorio
- Manejar comunicación con sistemas externos
- Configurar frameworks y librerías
- Implementar mecanismos de caché y mensajería
- Proporcionar adaptadores para servicios externos

## Estructura

- **`auth/`** - Autenticación y autorización
- **`config/`** - Configuración de la aplicación
- **`container/`** - Inyección de dependencias
- **`database/`** - Persistencia de datos
- **`http/`** - API REST y presentación web
- **`grpc/`** - API gRPC (opcional)
- **`graphql/`** - API GraphQL (opcional)
- **`websocket/`** - Comunicación en tiempo real
- **`messaging/`** - Sistemas de mensajería
- **`eventbus/`** - Bus de eventos
- **`cache/`** - Sistemas de caché
- **`telemetry/`** - Observabilidad y métricas
- **`validation/`** - Validación centralizada

## Principios

- **Inversión de dependencias**: Implementa interfaces del dominio
- **Configurabilidad**: Puede cambiarse sin afectar el dominio
- **Testeable**: Implementaciones pueden ser mockeadas
- **Modular**: Cada componente es independiente

## Patrón Adapter

La infraestructura actúa como adaptador entre el dominio y el mundo exterior:

```
Domain Interface → Infrastructure Implementation → External System
```

## Configuración

Toda la configuración se centraliza y se inyecta a través del contenedor de dependencias.

```go
type Container struct {
    Config           *config.Config
    Database         *gorm.DB
    EmployeeRepo     domain.EmployeeRepository
    EmployeeUseCase  *usecase.EmployeeUseCase
    EventBus         EventBus
    Cache            Cache
}
```
