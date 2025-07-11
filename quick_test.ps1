# Quick Test Script untuk DailyTrackr
Write-Host "🚀 Quick Test - DailyTrackr User Service" -ForegroundColor Green

# Set environment variable dengan force
Write-Host "🔧 Setting GO111MODULE..." -ForegroundColor Blue
[System.Environment]::SetEnvironmentVariable("GO111MODULE", "on", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GO111MODULE", "on", [System.EnvironmentVariableTarget]::Process)

# Refresh environment variables
$env:GO111MODULE = "on"

# Verify
Write-Host "✅ GO111MODULE = $(go env GO111MODULE)" -ForegroundColor Green

# Navigate to user-service and run
Write-Host "📁 Navigating to user-service..." -ForegroundColor Blue
Set-Location "user-service"

Write-Host "🏃 Starting user-service with environment override..." -ForegroundColor Green
Write-Host "🌐 Service will be available at: http://localhost:3001" -ForegroundColor Cyan
Write-Host ""

# Run with environment variable override
$env:GO111MODULE="on"; go run main.go