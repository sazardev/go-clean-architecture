# migration/ - Herramienta de Migraciones

Utilidad de línea de comandos para gestionar migraciones de base de datos.

## Responsabilidades

- Ejecutar migraciones hacia adelante (up)
- Revertir migraciones (down)
- Crear nuevas migraciones
- Verificar estado de migraciones
- Rollback a versiones específicas

## Comandos

```powershell
# Ejecutar todas las migraciones pendientes
go run cmd/migration/main.go up

# Revertir la última migración
go run cmd/migration/main.go down

# Crear nueva migración
go run cmd/migration/main.go create nombre_migracion

# Ver estado de migraciones
go run cmd/migration/main.go status

# Rollback a versión específica
go run cmd/migration/main.go rollback 001
```

## Estructura de Migraciones

Las migraciones se almacenarán en `migrations/` con formato:
- `001_create_users_table.up.sql`
- `001_create_users_table.down.sql`

## Configuración

Utiliza las mismas variables de base de datos que el servidor principal.
