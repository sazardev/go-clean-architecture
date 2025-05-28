# Ejemplos de uso de la HR API
# Script de PowerShell para demostrar funcionalidades

$baseUrl = "http://localhost:8080/api/v1"
$headers = @{ "Content-Type" = "application/json" }

Write-Host "🧪 Demostrando HR API..." -ForegroundColor Green
Write-Host ""

# 1. Health Check
Write-Host "1. Verificando estado del servidor..." -ForegroundColor Blue
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    Write-Host "✅ Servidor funcionando: $($healthResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Error: El servidor no está ejecutándose. Ejecuta primero: go run cmd/server/main.go" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 2. Crear empleados
Write-Host "2. Creando empleados..." -ForegroundColor Blue

$employees = @(
    @{ name = "Juan Pérez" },
    @{ name = "María García" },
    @{ name = "Carlos López" },
    @{ name = "Ana Martínez" }
)

$createdEmployees = @()

foreach ($employee in $employees) {
    try {
        $body = $employee | ConvertTo-Json
        $response = Invoke-RestMethod -Uri "$baseUrl/employees" -Method Post -Body $body -Headers $headers
        $createdEmployees += $response.data
        Write-Host "✅ Empleado creado: $($response.data.name) (ID: $($response.data.id))" -ForegroundColor Green
    } catch {
        Write-Host "❌ Error creando empleado: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 1
Write-Host ""

# 3. Listar todos los empleados
Write-Host "3. Listando todos los empleados..." -ForegroundColor Blue
try {
    $allEmployees = Invoke-RestMethod -Uri "$baseUrl/employees" -Method Get
    Write-Host "✅ Total de empleados: $($allEmployees.data.Count)" -ForegroundColor Green
    foreach ($emp in $allEmployees.data) {
        Write-Host "   - $($emp.name) (ID: $($emp.id))" -ForegroundColor White
    }
} catch {
    Write-Host "❌ Error obteniendo empleados: $($_.Exception.Message)" -ForegroundColor Red
}

Start-Sleep -Seconds 1
Write-Host ""

# 4. Obtener empleado por ID
if ($createdEmployees.Count -gt 0) {
    $firstEmployee = $createdEmployees[0]
    Write-Host "4. Obteniendo empleado por ID..." -ForegroundColor Blue
    try {
        $employee = Invoke-RestMethod -Uri "$baseUrl/employees/$($firstEmployee.id)" -Method Get
        Write-Host "✅ Empleado encontrado: $($employee.data.name)" -ForegroundColor Green
        Write-Host "   - ID: $($employee.data.id)" -ForegroundColor White
        Write-Host "   - Nombre: $($employee.data.name)" -ForegroundColor White
        Write-Host "   - Creado: $($employee.data.created_at)" -ForegroundColor White
    } catch {
        Write-Host "❌ Error obteniendo empleado: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 1
Write-Host ""

# 5. Actualizar empleado
if ($createdEmployees.Count -gt 0) {
    $employeeToUpdate = $createdEmployees[0]
    Write-Host "5. Actualizando empleado..." -ForegroundColor Blue
    try {
        $updateBody = @{ name = "$($employeeToUpdate.name) - Actualizado" } | ConvertTo-Json
        $updatedEmployee = Invoke-RestMethod -Uri "$baseUrl/employees/$($employeeToUpdate.id)" -Method Put -Body $updateBody -Headers $headers
        Write-Host "✅ Empleado actualizado: $($updatedEmployee.data.name)" -ForegroundColor Green
    } catch {
        Write-Host "❌ Error actualizando empleado: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 1
Write-Host ""

# 6. Eliminar empleado
if ($createdEmployees.Count -gt 1) {
    $employeeToDelete = $createdEmployees[-1]  # Último empleado
    Write-Host "6. Eliminando empleado..." -ForegroundColor Blue
    try {
        $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/employees/$($employeeToDelete.id)" -Method Delete
        Write-Host "✅ Empleado eliminado: $($deleteResponse.message)" -ForegroundColor Green
    } catch {
        Write-Host "❌ Error eliminando empleado: $($_.Exception.Message)" -ForegroundColor Red
    }
}

Start-Sleep -Seconds 1
Write-Host ""

# 7. Verificar estado final
Write-Host "7. Estado final - Listando empleados restantes..." -ForegroundColor Blue
try {
    $finalEmployees = Invoke-RestMethod -Uri "$baseUrl/employees" -Method Get
    Write-Host "✅ Empleados restantes: $($finalEmployees.data.Count)" -ForegroundColor Green
    foreach ($emp in $finalEmployees.data) {
        Write-Host "   - $($emp.name) (ID: $($emp.id))" -ForegroundColor White
    }
} catch {
    Write-Host "❌ Error obteniendo empleados finales: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎉 Demo completada!" -ForegroundColor Green
Write-Host "Puedes usar estos mismos comandos para interactuar con la API." -ForegroundColor Yellow
