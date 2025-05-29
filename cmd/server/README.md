# server/ - Servidor HTTP Principal

Punto de entrada del servidor HTTP que expone la API REST.

## Responsabilidades

- Inicialización del servidor HTTP
- Configuración de rutas y middlewares
- Inyección de dependencias
- Manejo de señales del sistema para shutdown graceful
- Configuración de CORS y otros middlewares globales

## Estructura del main.go

```go
func main() {
    // 1. Cargar configuración
    // 2. Inicializar base de datos
    // 3. Configurar contenedor de dependencias
    // 4. Inicializar rutas y middlewares
    // 5. Iniciar servidor con graceful shutdown
}
```

## Configuración

El servidor lee la configuración desde variables de entorno y archivo `.env`.

Variables principales:
- `SERVER_PORT` - Puerto del servidor (default: 8080)
- `DB_*` - Configuración de base de datos
- `JWT_SECRET` - Secreto para tokens JWT
