# worker/ - Procesamiento Asíncrono

Aplicación separada para el procesamiento de tareas en segundo plano.

## Responsabilidades

- Procesamiento de trabajos en cola
- Manejo de eventos asincrónicos
- Tareas programadas (cron jobs)
- Procesamiento de notificaciones
- Limpieza de datos y mantenimiento

## Casos de Uso

- Envío de emails de notificación
- Generación de reportes pesados
- Sincronización con sistemas externos
- Procesamiento de archivos grandes
- Limpieza de datos temporales

## Configuración

Variables de entorno específicas:
- `WORKER_CONCURRENCY` - Número de workers concurrentes
- `QUEUE_CONNECTION` - Configuración de la cola (Redis, etc.)
- `WORKER_TIMEOUT` - Timeout para tareas

## Implementación Futura

```go
func main() {
    // 1. Configurar conexión a cola de mensajes
    // 2. Registrar handlers de tareas
    // 3. Iniciar workers con pool de goroutines
    // 4. Implementar graceful shutdown
}
```
