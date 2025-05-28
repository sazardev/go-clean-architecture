# Arquitectura del Proyecto HR API

## 📁 Estructura Detallada

```
go-clean-architecture/
├── 📄 README.md                    # Documentación principal
├── 📄 go.mod                       # Definición del módulo Go
├── 📄 go.sum                       # Checksums de dependencias
├── 📄 Dockerfile                   # Configuración Docker
├── 📄 docker-compose.yml           # Orquestación de servicios
├── 📄 Makefile                     # Comandos automatizados
├── 📄 .env.example                 # Ejemplo de variables de entorno
├── 📄 .gitignore                   # Archivos ignorados por Git
├── 📄 setup-dev.ps1               # Script de configuración
│
├── 📂 cmd/                         # APLICACIONES EJECUTABLES
│   └── 📂 server/                  # Servidor principal
│       └── 📄 main.go              # Punto de entrada
│
├── 📂 internal/                    # CÓDIGO INTERNO
│   ├── 📂 domain/                  # CAPA DE DOMINIO
│   │   ├── 📂 entity/              # Entidades de negocio
│   │   │   └── 📄 employee.go      # Entidad Employee
│   │   └── 📂 repository/          # Contratos de repositorio
│   │       └── 📄 employee_repository.go # Interface EmployeeRepository
│   │
│   ├── 📂 usecase/                 # CAPA DE CASOS DE USO
│   │   ├── 📄 employee_usecase.go  # Lógica de negocio
│   │   └── 📄 employee_usecase_test.go # Tests unitarios
│   │
│   └── 📂 infrastructure/          # CAPA DE INFRAESTRUCTURA
│       ├── 📂 config/              # Configuración
│       │   └── 📄 config.go        # Manejo de configuración
│       │
│       ├── 📂 container/           # Inyección de dependencias
│       │   └── 📄 container.go     # Contenedor DI
│       │
│       ├── 📂 database/            # Persistencia
│       │   ├── 📄 connection.go    # Conexión a BD
│       │   └── 📄 employee_repository.go # Implementación repo
│       │
│       └── 📂 http/                # Capa HTTP
│           ├── 📂 dto/             # Data Transfer Objects
│           │   └── 📄 employee_dto.go # DTOs de Employee
│           │
│           ├── 📂 handler/         # Controladores HTTP
│           │   └── 📄 employee_handler.go # Handler de Employee
│           │
│           ├── 📂 middleware/      # Middlewares
│           │   └── 📄 middleware.go # CORS, Logging, etc.
│           │
│           └── 📂 router/          # Configuración de rutas
│               └── 📄 router.go    # Setup de rutas
│
└── 📂 examples/                    # EJEMPLOS Y DEMOS
    └── 📄 api-demo.ps1            # Demo de uso de API
```

## 🏗️ Capas de la Arquitectura

### 1. **Capa de Dominio** (`internal/domain/`)
- **Propósito**: Contiene la lógica de negocio pura
- **Características**:
  - Sin dependencias externas
  - Entidades de negocio
  - Interfaces de repositorio
  - Reglas de dominio

### 2. **Capa de Casos de Uso** (`internal/usecase/`)
- **Propósito**: Orquesta la lógica de aplicación
- **Características**:
  - Coordina entidades y repositorios
  - Implementa reglas de negocio específicas
  - Independiente de frameworks externos

### 3. **Capa de Infraestructura** (`internal/infrastructure/`)
- **Propósito**: Implementaciones concretas y detalles técnicos
- **Subcapas**:
  - **Config**: Manejo de configuración
  - **Container**: Inyección de dependencias
  - **Database**: Persistencia con GORM
  - **HTTP**: API REST con Fiber

## 🔄 Flujo de Dependencias

```
main.go → Container → Handler → UseCase → Repository → Database
   ↓         ↓          ↓         ↓          ↓          ↓
Config → Injector → HTTP → Business → Interface → GORM
```

## 📊 Principios Aplicados

### **SOLID**
- ✅ **S**ingle Responsibility: Cada componente tiene una responsabilidad
- ✅ **O**pen/Closed: Abierto para extensión, cerrado para modificación
- ✅ **L**iskov Substitution: Interfaces intercambiables
- ✅ **I**nterface Segregation: Interfaces específicas y cohesivas
- ✅ **D**ependency Inversion: Dependencias hacia abstracciones

### **Clean Architecture**
- ✅ **Independence**: Capas independientes
- ✅ **Testability**: Fácil testing con mocks
- ✅ **Framework Independence**: No dependiente de frameworks
- ✅ **Database Independence**: Intercambiable por interfaces

## 🔧 Patrones Implementados

### **Repository Pattern**
```go
type EmployeeRepository interface {
    Create(ctx context.Context, employee *entity.Employee) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Employee, error)
    // ... más métodos
}
```

### **Dependency Injection**
```go
type Container struct {
    Config          *config.Config
    DB              *gorm.DB
    EmployeeHandler *handler.EmployeeHandler
}
```

### **DTO Pattern**
```go
type CreateEmployeeRequest struct {
    Name string `json:"name" validate:"required,min=2,max=255"`
}
```

## 🧪 Testing Strategy

### **Niveles de Testing**
1. **Unit Tests**: Casos de uso con repositorios mock
2. **Integration Tests**: Base de datos real (opcional)
3. **E2E Tests**: API completa (opcional)

### **Ejemplo de Mock**
```go
type mockEmployeeRepository struct {
    employees map[uuid.UUID]*entity.Employee
}
```

## 📦 Tecnologías y Librerías

### **Core**
- **Go 1.24+**: Lenguaje principal
- **Fiber v2**: Framework web ultrarrápido
- **GORM**: ORM para Go

### **Database**
- **PostgreSQL**: Base de datos relacional
- **UUID**: Identificadores únicos

### **Utilities**
- **Godotenv**: Variables de entorno
- **Docker**: Containerización

## 🚀 Ventajas de esta Arquitectura

### **Mantenibilidad**
- Código organizado en capas claras
- Fácil localización de funcionalidades
- Separación de responsabilidades

### **Testabilidad**
- Mocks simples por interfaces
- Tests unitarios aislados
- Coverage alto posible

### **Escalabilidad**
- Nuevas funcionalidades sin impacto
- Microservicios ready
- Horizontal scaling friendly

### **Flexibilidad**
- Cambio de base de datos sin impacto
- Intercambio de frameworks
- Nuevos métodos de entrega (GraphQL, gRPC)

## 📈 Extensibilidad

### **Agregar Nueva Entidad**
1. Crear entidad en `domain/entity/`
2. Crear interface en `domain/repository/`
3. Implementar use case en `usecase/`
4. Implementar repositorio en `infrastructure/database/`
5. Crear DTOs en `infrastructure/http/dto/`
6. Crear handler en `infrastructure/http/handler/`
7. Agregar rutas en `infrastructure/http/router/`

### **Nuevos Casos de Uso**
- Agregar métodos a use case existente
- Mantener interfaces de repositorio
- Extender handlers para nuevos endpoints

Esta arquitectura garantiza un código limpio, mantenible y escalable siguiendo las mejores prácticas de desarrollo en Go.
