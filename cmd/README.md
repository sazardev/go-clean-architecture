# cmd/ - Aplicaciones Ejecutables

Esta carpeta contiene los puntos de entrada de todas las aplicaciones ejecutables del proyecto.

## Estructura

- **`server/`** - Servidor HTTP principal de la API
- **`worker/`** - Procesamiento asíncrono de tareas en segundo plano
- **`migration/`** - Herramienta de línea de comandos para ejecutar migraciones de base de datos

## Principios

- Cada subcarpeta representa una aplicación independiente
- Los archivos `main.go` deben ser simples y delegar la lógica compleja a los paquetes internos
- Configuración y dependencias se inyectan desde la capa de infraestructura
- Cada aplicación puede tener su propia configuración específica

## Uso

```powershell
# Ejecutar servidor principal
go run cmd/server/main.go

# Ejecutar worker
go run cmd/worker/main.go

# Ejecutar migraciones
go run cmd/migration/main.go
```
