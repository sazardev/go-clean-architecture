# HR API - Clean Architecture en Go

Un ejemplo ultra-limpio de API REST para Recursos Humanos implementado con Clean Architecture en Go, utilizando las mejores prácticas de desarrollo modular, escalable y mantenible.

## 🏗️ Arquitectura

Este proyecto sigue los principios de **Clean Architecture** con las siguientes capas:

```
cmd/
└── server/              # Punto de entrada de la aplicación
    └── main.go

internal/
├── domain/              # Capa de Dominio (Entidades y Contratos)
│   ├── entity/          # Entidades de negocio
│   └── repository/      # Interfaces de repositorio
├── usecase/             # Capa de Casos de Uso (Lógica de Negocio)
├── infrastructure/      # Capa de Infraestructura
│   ├── config/          # Configuración de la aplicación
│   ├── container/       # Inyección de dependencias
│   ├── database/        # Implementación de persistencia
│   └── http/            # Capa de presentación HTTP
│       ├── dto/         # Data Transfer Objects
│       ├── handler/     # Manejadores HTTP
│       ├── middleware/  # Middlewares HTTP
│       └── router/      # Configuración de rutas
```

## 🚀 Tecnologías Utilizadas

- **Go 1.24+**
- **Fiber v2** - Framework web ultra-rápido
- **GORM** - ORM para Go
- **PostgreSQL** - Base de datos relacional
- **UUID** - Identificadores únicos
- **Godotenv** - Manejo de variables de entorno

## 📋 Características

- ✅ **Clean Architecture** - Separación clara de responsabilidades
- ✅ **Dependency Injection** - Inversión de dependencias
- ✅ **Repository Pattern** - Abstracción de persistencia
- ✅ **DTO Pattern** - Transferencia segura de datos
- ✅ **Error Handling** - Manejo robusto de errores
- ✅ **Middleware Support** - CORS, Logging, Recovery
- ✅ **Graceful Shutdown** - Cierre limpio del servidor
- ✅ **Environment Configuration** - Configuración flexible
- ✅ **Database Migrations** - Migración automática de esquemas

## 🔧 Configuración

1. **Clonar el repositorio**
   ```powershell
   git clone <repository-url>
   cd go-clean-architecture
   ```

2. **Configurar variables de entorno**
   ```powershell
   Copy-Item .env.example .env
   ```
   
   Editar el archivo `.env` con tus configuraciones:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=tu_password
   DB_NAME=hr_db
   DB_SSL_MODE=disable
   SERVER_PORT=8080
   ```

3. **Instalar dependencias**
   ```powershell
   go mod download
   ```

4. **Configurar PostgreSQL**
   ```powershell
   # Crear base de datos
   createdb hr_db
   ```

## 🏃‍♂️ Ejecución

```powershell
# Ejecutar la aplicación
go run cmd/server/main.go

# O compilar y ejecutar
go build -o hr-api.exe cmd/server/main.go
./hr-api.exe
```

## 📡 API Endpoints

### Health Check
- `GET /health` - Verificar estado del servidor

### Empleados
- `POST /api/v1/employees` - Crear empleado
- `GET /api/v1/employees` - Listar todos los empleados
- `GET /api/v1/employees/{id}` - Obtener empleado por ID
- `PUT /api/v1/employees/{id}` - Actualizar empleado
- `DELETE /api/v1/employees/{id}` - Eliminar empleado

### Ejemplos de Uso

#### Crear Empleado
```powershell
curl -X POST http://localhost:8080/api/v1/employees `
  -H "Content-Type: application/json" `
  -d '{"name": "Juan Pérez"}'
```

#### Obtener Todos los Empleados
```powershell
curl http://localhost:8080/api/v1/employees
```

#### Obtener Empleado por ID
```powershell
curl http://localhost:8080/api/v1/employees/{uuid}
```

#### Actualizar Empleado
```powershell
curl -X PUT http://localhost:8080/api/v1/employees/{uuid} `
  -H "Content-Type: application/json" `
  -d '{"name": "Juan Carlos Pérez"}'
```

#### Eliminar Empleado
```powershell
curl -X DELETE http://localhost:8080/api/v1/employees/{uuid}
```

## 🧪 Testing

```powershell
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con coverage
go test -cover ./...

# Ejecutar tests de un paquete específico
go test ./internal/usecase/
```

## 📁 Estructura de Archivos

- **`cmd/`** - Aplicaciones ejecutables
- **`internal/`** - Código interno de la aplicación
  - **`domain/`** - Lógica de negocio pura
  - **`usecase/`** - Casos de uso y reglas de negocio
  - **`infrastructure/`** - Implementaciones concretas
- **`go.mod`** - Definición de módulo y dependencias

## 🔄 Flujo de Datos

```
HTTP Request → Router → Middleware → Handler → UseCase → Repository → Database
                ↓
HTTP Response ← DTO ← Handler ← UseCase ← Repository ← Database
```

## 🎯 Principios Aplicados

- **Single Responsibility** - Cada componente tiene una única responsabilidad
- **Open/Closed** - Abierto para extensión, cerrado para modificación
- **Liskov Substitution** - Las interfaces pueden ser sustituidas sin romper el código
- **Interface Segregation** - Interfaces específicas y cohesivas
- **Dependency Inversion** - Dependencias hacia abstracciones, no concreciones

## 🚀 Despliegue

### Docker (Opcional)
```powershell
# Construir imagen
docker build -t hr-api .

# Ejecutar contenedor
docker run -p 8080:8080 --env-file .env hr-api
```

## 📚 Documentación Adicional

- [📐 Arquitectura Detallada](docs/ARCHITECTURE.md) - Explicación completa de la arquitectura
- [🧪 Ejemplos de Uso](examples/api-demo.ps1) - Scripts de demostración
- [⚙️ Configuración de Desarrollo](setup-dev.ps1) - Setup automático

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.
