# DailyTrackr - Run All Services (Simple Version)
Write-Host "🚀 Starting DailyTrackr Microservices..." -ForegroundColor Green

# Set environment variable
$env:GO111MODULE = "on"

# Check if we're in the right directory
if (-not (Test-Path "shared")) {
    Write-Host "❌ Please run this script from the project root directory" -ForegroundColor Red
    Write-Host "📁 Current directory: $(Get-Location)" -ForegroundColor Yellow
    exit 1
}

Write-Host "📁 Current directory: $(Get-Location)" -ForegroundColor Blue

# Check MySQL connection
Write-Host "`n🗄️ Checking MySQL connection..." -ForegroundColor Magenta
try {
    $mysql = Test-NetConnection -ComputerName localhost -Port 3306 -WarningAction SilentlyContinue -ErrorAction SilentlyContinue
    if ($mysql.TcpTestSucceeded) {
        Write-Host "   ✅ MySQL is running on localhost:3306" -ForegroundColor Green
    } else {
        Write-Host "   ❌ MySQL is not running!" -ForegroundColor Red
        Write-Host "   💡 Please start Laragon or your MySQL server" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ⚠️ Could not test MySQL connection" -ForegroundColor Yellow
}

# Start Gateway Service
Write-Host "`n🔧 Starting Gateway on port 3000..." -ForegroundColor Cyan
if (Test-Path "gateway") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd gateway; Write-Host 'Gateway Starting...' -ForegroundColor Green; go run main.go"
    Start-Sleep -Seconds 2
    Write-Host "   ✅ Gateway started" -ForegroundColor Green
} else {
    Write-Host "   ❌ Gateway directory not found" -ForegroundColor Red
}

# Start User Service
Write-Host "`n🔧 Starting User Service on port 3001..." -ForegroundColor Cyan
if (Test-Path "user-service") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd user-service; Write-Host 'User Service Starting...' -ForegroundColor Green; go run main.go"
    Start-Sleep -Seconds 2
    Write-Host "   ✅ User Service started" -ForegroundColor Green
} else {
    Write-Host "   ❌ User Service directory not found" -ForegroundColor Red
}

# Start Activity Service
Write-Host "`n🔧 Starting Activity Service on port 3002..." -ForegroundColor Cyan
if (Test-Path "activity-service") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd activity-service; Write-Host 'Activity Service Starting...' -ForegroundColor Green; go run main.go"
    Start-Sleep -Seconds 2
    Write-Host "   ✅ Activity Service started" -ForegroundColor Green
} else {
    Write-Host "   ❌ Activity Service directory not found" -ForegroundColor Red
}

# Start Habit Service
Write-Host "`n🔧 Starting Habit Service on port 3003..." -ForegroundColor Cyan
if (Test-Path "habit-service") {
    Start-Process powershell -ArgumentList "-NoExit", "-Command", "cd habit-service; Write-Host 'Habit Service Starting...' -ForegroundColor Green; go run main.go"
    Start-Sleep -Seconds 2
    Write-Host "   ✅ Habit Service started" -ForegroundColor Green
} else {
    Write-Host "   ❌ Habit Service directory not found" -ForegroundColor Red
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

Write-Host "`n📋 Important Notes:" -ForegroundColor Magenta
Write-Host "   ⏱️ Wait 10-15 seconds for all services to start completely" -ForegroundColor Yellow
Write-Host "   📝 Check each service window for startup logs and errors" -ForegroundColor Blue
Write-Host "   🛑 Press Ctrl+C in any service window to stop that service" -ForegroundColor Red

Write-Host "`n🎯 All services launched! Check individual windows for logs." -ForegroundColor Green