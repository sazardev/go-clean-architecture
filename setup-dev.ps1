# HR API - Configuración de Desarrollo
# Script de PowerShell para configurar el entorno de desarrollo

Write-Host "🚀 Configurando entorno de desarrollo para HR API..." -ForegroundColor Green

# Verificar que Go esté instalado
try {
    $goVersion = go version
    Write-Host "✅ Go encontrado: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go no está instalado. Por favor instala Go primero." -ForegroundColor Red
    exit 1
}

# Verificar que Docker esté instalado (opcional)
try {
    $dockerVersion = docker --version
    Write-Host "✅ Docker encontrado: $dockerVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Docker no encontrado. Puedes usar PostgreSQL local." -ForegroundColor Yellow
}

# Crear archivo .env si no existe
if (!(Test-Path ".env")) {
    Write-Host "📝 Creando archivo .env..." -ForegroundColor Blue
    Copy-Item ".env.example" ".env"
    Write-Host "✅ Archivo .env creado. Configura tus variables de entorno." -ForegroundColor Green
} else {
    Write-Host "✅ Archivo .env ya existe." -ForegroundColor Green
}

# Descargar dependencias
Write-Host "📦 Descargando dependencias de Go..." -ForegroundColor Blue
go mod download
go mod tidy

# Compilar la aplicación
Write-Host "🔨 Compilando aplicación..." -ForegroundColor Blue
go build -o hr-api.exe cmd/server/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Aplicación compilada exitosamente." -ForegroundColor Green
} else {
    Write-Host "❌ Error compilando la aplicación." -ForegroundColor Red
    exit 1
}

# Ejecutar tests
Write-Host "🧪 Ejecutando tests..." -ForegroundColor Blue
go test ./...

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Todos los tests pasaron." -ForegroundColor Green
} else {
    Write-Host "❌ Algunos tests fallaron." -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 Configuración completada!" -ForegroundColor Green
Write-Host ""
Write-Host "Próximos pasos:" -ForegroundColor Yellow
Write-Host "1. Configura tu base de datos PostgreSQL" -ForegroundColor White
Write-Host "2. Edita el archivo .env con tus configuraciones" -ForegroundColor White
Write-Host "3. Ejecuta: ./hr-api.exe o go run cmd/server/main.go" -ForegroundColor White
Write-Host "4. Visita: http://localhost:8080/health" -ForegroundColor White
Write-Host ""
Write-Host "Con Docker:" -ForegroundColor Yellow
Write-Host "docker-compose up -d postgres  # Solo PostgreSQL" -ForegroundColor White
Write-Host "docker-compose up              # Toda la aplicación" -ForegroundColor White
