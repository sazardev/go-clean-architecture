# HR API Authentication Test Script

$baseUrl = "http://localhost:8080/api/v1"
$headers = @{
    "Content-Type" = "application/json"
}

Write-Host "🚀 Testing HR API Authentication System" -ForegroundColor Green
Write-Host "=====================================`n" -ForegroundColor Green

# Test 1: Health Check
Write-Host "1️⃣  Testing Health Check..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get
    Write-Host "✅ Health Check: $($healthResponse.status) - $($healthResponse.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Health Check Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Register a new user
Write-Host "`n2️⃣  Testing User Registration..." -ForegroundColor Yellow
$registerData = @{
    email = "admin@company.com"
    password = "admin123"
    first_name = "Admin"
    last_name = "User"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method Post -Body $registerData -Headers $headers
    Write-Host "✅ User Registration successful" -ForegroundColor Green
} catch {
    Write-Host "⚠️  User Registration: $($_.Exception.Message)" -ForegroundColor Yellow
    Write-Host "   (This might fail if user already exists)" -ForegroundColor Gray
}

# Test 3: Login
Write-Host "`n3️⃣  Testing User Login..." -ForegroundColor Yellow
$loginData = @{
    email = "admin@company.com"
    password = "admin123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginData -Headers $headers
    $token = $loginResponse.access_token
    Write-Host "✅ Login successful - Token received" -ForegroundColor Green
    
    # Add authorization header for protected routes
    $authHeaders = @{
        "Content-Type" = "application/json"
        "Authorization" = "Bearer $token"
    }
} catch {
    Write-Host "❌ Login Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 4: Get Profile (Protected Route)
Write-Host "`n4️⃣  Testing Get Profile (Protected Route)..." -ForegroundColor Yellow
try {
    $profileResponse = Invoke-RestMethod -Uri "$baseUrl/profile" -Method Get -Headers $authHeaders
    Write-Host "✅ Profile access successful" -ForegroundColor Green
} catch {
    Write-Host "❌ Profile access failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 5: Test Permission-based Access
Write-Host "`n5️⃣  Testing Permission-based Access..." -ForegroundColor Yellow

# Test Users endpoint (requires permissions)
try {
    $usersResponse = Invoke-RestMethod -Uri "$baseUrl/users" -Method Get -Headers $authHeaders
    Write-Host "✅ Users endpoint accessible" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Users endpoint: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test Roles endpoint (requires permissions)
try {
    $rolesResponse = Invoke-RestMethod -Uri "$baseUrl/roles" -Method Get -Headers $authHeaders
    Write-Host "✅ Roles endpoint accessible" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Roles endpoint: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test Permissions endpoint (requires permissions)
try {
    $permissionsResponse = Invoke-RestMethod -Uri "$baseUrl/permissions" -Method Get -Headers $authHeaders
    Write-Host "✅ Permissions endpoint accessible" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Permissions endpoint: $($_.Exception.Message)" -ForegroundColor Yellow
}

# Test 6: Test Unauthorized Access
Write-Host "`n6️⃣  Testing Unauthorized Access..." -ForegroundColor Yellow
try {
    $unauthorizedResponse = Invoke-RestMethod -Uri "$baseUrl/profile" -Method Get -Headers $headers
    Write-Host "❌ Unauthorized access should have failed" -ForegroundColor Red
} catch {
    Write-Host "✅ Unauthorized access properly blocked" -ForegroundColor Green
}

# Test 7: Token Refresh
Write-Host "`n7️⃣  Testing Token Refresh..." -ForegroundColor Yellow
$refreshData = @{
    refresh_token = $token
} | ConvertTo-Json

try {
    $refreshResponse = Invoke-RestMethod -Uri "$baseUrl/auth/refresh" -Method Post -Body $refreshData -Headers $headers
    Write-Host "✅ Token refresh successful" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Token refresh: $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "`n🎉 Authentication System Test Complete!" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

Write-Host "`n📋 Available Endpoints:" -ForegroundColor Cyan
Write-Host "  • POST /api/v1/auth/register - User registration" -ForegroundColor White
Write-Host "  • POST /api/v1/auth/login - User login" -ForegroundColor White
Write-Host "  • POST /api/v1/auth/refresh - Token refresh" -ForegroundColor White
Write-Host "  • GET  /api/v1/profile - User profile (protected)" -ForegroundColor White
Write-Host "  • PUT  /api/v1/profile - Update profile (protected)" -ForegroundColor White
Write-Host "  • GET  /api/v1/users - List users (admin)" -ForegroundColor White
Write-Host "  • GET  /api/v1/roles - List roles (admin)" -ForegroundColor White
Write-Host "  • GET  /api/v1/permissions - List permissions (admin)" -ForegroundColor White
Write-Host "  • GET  /api/v1/employees - List employees (protected)" -ForegroundColor White
