# Script de PowerShell para crear la estructura completa del proyecto
# Ejecutar desde la raíz del proyecto: .\create-structure.ps1

param(
    [switch]$Force = $false
)

Write-Host "🚀 Creando estructura completa del proyecto HR API - Clean Architecture" -ForegroundColor Green
Write-Host "=" * 80 -ForegroundColor Gray

$baseDir = $PWD.Path

# Función para crear directorio si no existe
function New-DirectoryIfNotExists {
    param([string]$Path)
    
    if (-not (Test-Path $Path)) {
        New-Item -ItemType Directory -Path $Path -Force | Out-Null
        Write-Host "✅ Creado: $($Path.Replace($baseDir, '.'))" -ForegroundColor Green
    } else {
        Write-Host "⚡ Existe: $($Path.Replace($baseDir, '.'))" -ForegroundColor Yellow
    }
}

# Función para crear archivo README si no existe
function New-ReadmeIfNotExists {
    param(
        [string]$Path,
        [string]$Title,
        [string]$Description
    )
    
    $readmePath = Join-Path $Path "README.md"
    
    if (-not (Test-Path $readmePath) -or $Force) {
        $content = @"
# $Title

$Description

## Responsabilidades

- [Definir responsabilidades específicas]

## Estructura

- [Describir estructura interna]

## Implementación

```go
// Código de ejemplo aquí
```

## Configuración

Variables de entorno:
- `VARIABLE_NAME` - Descripción

## Uso

```go
// Ejemplo de uso
```
"@
        
        Set-Content -Path $readmePath -Value $content -Encoding UTF8
        Write-Host "📝 README creado: $($readmePath.Replace($baseDir, '.'))" -ForegroundColor Cyan
    }
}

Write-Host "`n📁 Creando estructura de directorios..." -ForegroundColor Blue

# Crear todas las carpetas principales
$directories = @(
    # CMD directories
    "cmd\migration",
    "cmd\worker",
    
    # Internal domain
    "internal\domain\entity",
    "internal\domain\repository", 
    "internal\domain\service",
    "internal\domain\event",
    "internal\domain\valueobject",
    
    # Internal usecase
    "internal\usecase\command",
    "internal\usecase\query",
    "internal\usecase\common",
    
    # Internal infrastructure - auth
    "internal\infrastructure\auth\jwt",
    "internal\infrastructure\auth\rbac", 
    "internal\infrastructure\auth\middleware",
    
    # Internal infrastructure - database
    "internal\infrastructure\database\postgres",
    "internal\infrastructure\database\mongodb",
    "internal\infrastructure\database\redis",
    "internal\infrastructure\database\factory",
    
    # Internal infrastructure - communication
    "internal\infrastructure\grpc",
    "internal\infrastructure\graphql",
    "internal\infrastructure\websocket\subscription",
    "internal\infrastructure\http\query",
    
    # Internal infrastructure - messaging & events
    "internal\infrastructure\messaging",
    "internal\infrastructure\eventbus\memory",
    "internal\infrastructure\eventbus\kafka", 
    "internal\infrastructure\eventbus\nats",
    
    # Internal infrastructure - other
    "internal\infrastructure\cache",
    "internal\infrastructure\telemetry",
    "internal\infrastructure\validation",
    
    # Internal pkg
    "internal\pkg",
    
    # Public pkg
    "pkg\logger",
    "pkg\validator",
    "pkg\crypto",
    "pkg\utils",
    
    # Additional directories
    "migrations\postgres",
    "migrations\mongodb",
    "deployments\docker",
    "deployments\kubernetes",
    "scripts\dev",
    "scripts\prod",
    "tests\integration",
    "tests\e2e",
    "api\openapi"
)

foreach ($dir in $directories) {
    $fullPath = Join-Path $baseDir $dir
    New-DirectoryIfNotExists $fullPath
}

Write-Host "`n📚 Creando archivos README.md..." -ForegroundColor Blue

# Crear READMEs para nuevas carpetas que no los tienen
$readmeConfig = @{
    "migrations" = @{
        "title" = "migrations/ - Migraciones de Base de Datos"
        "desc" = "Archivos de migración para diferentes bases de datos del sistema."
    }
    "migrations\postgres" = @{
        "title" = "postgres/ - Migraciones PostgreSQL"  
        "desc" = "Scripts SQL para migrar esquemas de PostgreSQL."
    }
    "migrations\mongodb" = @{
        "title" = "mongodb/ - Migraciones MongoDB"
        "desc" = "Scripts para migrar colecciones y índices de MongoDB."
    }
    "pkg" = @{
        "title" = "pkg/ - Bibliotecas Públicas"
        "desc" = "Paquetes reutilizables que pueden ser importados por otros proyectos."
    }
    "pkg\logger" = @{
        "title" = "logger/ - Sistema de Logging"
        "desc" = "Implementación de logging estructurado para toda la aplicación."
    }
    "pkg\validator" = @{
        "title" = "validator/ - Validaciones"
        "desc" = "Utilidades de validación reutilizables."
    }
    "pkg\crypto" = @{
        "title" = "crypto/ - Criptografía"
        "desc" = "Funciones criptográficas seguras para hashing, encryption, etc."
    }
    "pkg\utils" = @{
        "title" = "utils/ - Utilidades"
        "desc" = "Funciones de utilidad general."
    }
    "deployments" = @{
        "title" = "deployments/ - Configuración de Despliegue"
        "desc" = "Archivos de configuración para diferentes entornos de despliegue."
    }
    "deployments\docker" = @{
        "title" = "docker/ - Configuración Docker"
        "desc" = "Dockerfiles y docker-compose para diferentes entornos."
    }
    "deployments\kubernetes" = @{
        "title" = "kubernetes/ - Configuración K8s"
        "desc" = "Manifiestos de Kubernetes para despliegue en cluster."
    }
    "scripts" = @{
        "title" = "scripts/ - Scripts de Automatización"
        "desc" = "Scripts para automatizar tareas de desarrollo y despliegue."
    }
    "scripts\dev" = @{
        "title" = "dev/ - Scripts de Desarrollo"
        "desc" = "Scripts para facilitar el desarrollo local."
    }
    "scripts\prod" = @{
        "title" = "prod/ - Scripts de Producción"
        "desc" = "Scripts para despliegue y mantenimiento en producción."
    }
    "tests" = @{
        "title" = "tests/ - Tests"
        "desc" = "Tests de integración y end-to-end."
    }
    "tests\integration" = @{
        "title" = "integration/ - Tests de Integración"
        "desc" = "Tests que verifican la integración entre componentes."
    }
    "tests\e2e" = @{
        "title" = "e2e/ - Tests End-to-End"
        "desc" = "Tests completos del flujo de usuario."
    }
    "api" = @{
        "title" = "api/ - Documentación de API"
        "desc" = "Especificaciones y documentación de las APIs."
    }
    "api\openapi" = @{
        "title" = "openapi/ - Especificaciones OpenAPI"
        "desc" = "Archivos OpenAPI/Swagger para documentar la API REST."
    }
    "internal\infrastructure\telemetry" = @{
        "title" = "telemetry/ - Observabilidad"
        "desc" = "Métricas, tracing y logging para observabilidad del sistema."
    }
    "internal\infrastructure\validation" = @{
        "title" = "validation/ - Validación Centralizada"
        "desc" = "Sistema centralizado de validación de datos de entrada."
    }
    "internal\infrastructure\grpc" = @{
        "title" = "grpc/ - API gRPC"
        "desc" = "Implementación de servicios gRPC para comunicación eficiente."
    }
    "internal\infrastructure\graphql" = @{
        "title" = "graphql/ - API GraphQL"
        "desc" = "Implementación de API GraphQL para consultas flexibles."
    }
    "internal\infrastructure\websocket" = @{
        "title" = "websocket/ - Comunicación en Tiempo Real"
        "desc" = "WebSockets para comunicación bidireccional en tiempo real."
    }
    "internal\infrastructure\websocket\subscription" = @{
        "title" = "subscription/ - Suscripciones WebSocket"
        "desc" = "Gestión de suscripciones y notificaciones en tiempo real."
    }
    "internal\pkg" = @{
        "title" = "pkg/ - Utilidades Internas"
        "desc" = "Paquetes de utilidad compartidos dentro de la aplicación."
    }
}

foreach ($path in $readmeConfig.Keys) {
    $fullPath = Join-Path $baseDir $path
    if (Test-Path $fullPath) {
        $config = $readmeConfig[$path]
        New-ReadmeIfNotExists $fullPath $config.title $config.desc
    }
}

Write-Host "`n🎯 Creando archivos de configuración adicionales..." -ForegroundColor Blue

# .gitkeep para directorios que deben existir pero pueden estar vacíos
$gitkeepDirs = @(
    "migrations\postgres",
    "migrations\mongodb", 
    "tests\integration",
    "tests\e2e",
    "api\openapi"
)

foreach ($dir in $gitkeepDirs) {
    $gitkeepPath = Join-Path $baseDir $dir ".gitkeep"
    if (-not (Test-Path $gitkeepPath)) {
        New-Item -ItemType File -Path $gitkeepPath -Force | Out-Null
        Write-Host "📌 .gitkeep creado: $($gitkeepPath.Replace($baseDir, '.'))" -ForegroundColor Gray
    }
}

Write-Host "`n📊 Resumen de la estructura creada:" -ForegroundColor Blue
Write-Host "=" * 50 -ForegroundColor Gray

$totalDirs = (Get-ChildItem -Path $baseDir -Recurse -Directory | Measure-Object).Count
$totalFiles = (Get-ChildItem -Path $baseDir -Recurse -File | Measure-Object).Count
$readmeFiles = (Get-ChildItem -Path $baseDir -Recurse -Name "README.md" | Measure-Object).Count

Write-Host "📁 Total de directorios: $totalDirs" -ForegroundColor Cyan
Write-Host "📄 Total de archivos: $totalFiles" -ForegroundColor Cyan  
Write-Host "📚 Archivos README.md: $readmeFiles" -ForegroundColor Cyan

Write-Host "`n✅ Estructura del proyecto creada exitosamente!" -ForegroundColor Green
Write-Host "🚀 Próximos pasos:" -ForegroundColor Yellow
Write-Host "   1. Revisar y personalizar los archivos README.md creados" -ForegroundColor White
Write-Host "   2. Implementar las interfaces y estructuras definidas" -ForegroundColor White
Write-Host "   3. Configurar las variables de entorno necesarias" -ForegroundColor White
Write-Host "   4. Ejecutar las migraciones de base de datos" -ForegroundColor White

Write-Host "`n🎉 ¡Proyecto listo para desarrollo!" -ForegroundColor Green
