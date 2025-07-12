# DailyTrackr - Run All Services Script
Write-Host "🚀 Starting DailyTrackr Microservices..." -ForegroundColor Green

# Set environment variable
$env:GO111MODULE = "on"

# Function to start service in new window
function Start-Service {
    param(
        [string]$ServiceName,
        [string]$ServicePath,
        [string]$Port
    )

    Write-Host "🔧 Starting $ServiceName on port $Port..." -ForegroundColor Cyan

    # Check if service directory exists
    if (-not (Test-Path $ServicePath)) {
        Write-Host "   ❌ $ServicePath not found, skipping $ServiceName" -ForegroundColor Red
        return
    }

    # Start service in new PowerShell window
    $command = "cd '$ServicePath'; Write-Host '🚀 $ServiceName Starting...' -ForegroundColor Green; Write-Host 'Port: $Port' -ForegroundColor Yellow; go run main.go; Write-Host '❌ $ServiceName stopped. Press any key to close...' -ForegroundColor Red; `$null = `$Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')"

    Start-Process powershell -ArgumentList "-NoExit", "-Command", $command

    Start-Sleep -Seconds 2
}

# Check if we're in the right directory
if (-not (Test-Path "shared")) {
    Write-Host "❌ Please run this script from the project root directory" -ForegroundColor Red
    Write-Host "📁 Current directory: $(Get-Location)" -ForegroundColor Yellow
    Write-Host "💡 Expected structure: ./shared, ./gateway, ./user-service, etc." -ForegroundColor Blue
    exit 1
}

Write-Host "📁 Current directory: $(Get-Location)" -ForegroundColor Blue
Write-Host "🔧 Setting up environment..." -ForegroundColor Yellow

# Check MySQL connection
Write-Host "`n🗄️  Checking MySQL connection..." -ForegroundColor Magenta
try {
    $mysql = Test-NetConnection -ComputerName localhost -Port 3306 -WarningAction SilentlyContinue -ErrorAction SilentlyContinue
    if ($mysql.TcpTestSucceeded) {
        Write-Host "   ✅ MySQL is running on localhost:3306" -ForegroundColor Green
    } else {
        Write-Host "   ❌ MySQL is not running!" -ForegroundColor Red
        Write-Host "   💡 Please start Laragon or your MySQL server" -ForegroundColor Yellow
        Write-Host "   🔄 Continuing anyway..." -ForegroundColor Blue
    }
} catch {
    Write-Host "   ⚠️  Could not test MySQL connection" -ForegroundColor Yellow
}

# Start services in order
Write-Host "`n🏗️ Starting services..." -ForegroundColor Magenta

# 1. Gateway (Port 3000) - Start first as it routes to other services
Start-Service -ServiceName "Gateway" -ServicePath "$(Get-Location)\gateway" -Port "3000"

# 2. User Service (Port 3001) - Core authentication service
Start-Service -ServiceName "User Service" -ServicePath "$(Get-Location)\user-service" -Port "3001"

# 3. Activity Service (Port 3002) - Activity tracking
Start-Service -ServiceName "Activity Service" -ServicePath "$(Get-Location)\activity-service" -Port "3002"

# 4. Habit Service (Port 3003) - Habit management
Start-Service -ServiceName "Habit Service" -ServicePath "$(Get-Location)\habit-service" -Port "3003"

# Optional services (if they exist)
if (Test-Path "notification-service") {
    Start-Service -ServiceName "Notification Service" -ServicePath "$(Get-Location)\notification-service" -Port "3004"
}

if (Test-Path "stat-service") {
    Start-Service -ServiceName "Statistics Service" -ServicePath "$(Get-Location)\stat-service" -Port "3005"
}

if (Test-Path "ai-service") {
    Start-Service -ServiceName "AI Service" -ServicePath "$(Get-Location)\ai-service" -Port "3006"
}

Write-Host "`n✅ All services are starting up!" -ForegroundColor Green
Write-Host "🌐 Service URLs:" -ForegroundColor Cyan
Write-Host "   - Gateway: http://localhost:3000" -ForegroundColor White
Write-Host "   - User Service: http://localhost:3001" -ForegroundColor White
Write-Host "   - Activity Service: http://localhost:3002" -ForegroundColor White
Write-Host "   - Habit Service: http://localhost:3003" -ForegroundColor White

Write-Host "`n🧪 Test endpoints:" -ForegroundColor Yellow
Write-Host "   - Gateway Health: GET http://localhost:3000/" -ForegroundColor White
Write-Host "   - Register: POST http://localhost:3000/api/users/auth/register" -ForegroundColor White
Write-Host "   - Login: POST http://localhost:3000/api/users/auth/login" -ForegroundColor White
Write-Host "   - Activities: GET http://localhost:3000/api/activities/api/v1/activities (with JWT)" -ForegroundColor White
Write-Host "   - Habits: GET http://localhost:3000/api/habits/api/v1/habits (with JWT)" -ForegroundColor White

Write-Host "`n📋 Important Notes:" -ForegroundColor Magenta
Write-Host "   ⏱️  Wait 10-15 seconds for all services to start completely" -ForegroundColor Yellow
Write-Host "   📝 Check each service window for startup logs and errors" -ForegroundColor Blue
Write-Host "   🔄 If a service fails, check the error in its window" -ForegroundColor Blue
Write-Host "   🛑 Press Ctrl+C in any service window to stop that service" -ForegroundColor Red
Write-Host "   🧪 Run .\test_api_endpoints.ps1 to test all endpoints" -ForegroundColor Green

Write-Host "`n🎯 Press any key to close this launcher window..." -ForegroundColor Magenta
Read-Host "Press Enter to close"
