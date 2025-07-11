# DailyTrackr Environment Setup Script
Write-Host "Setting up DailyTrackr Development Environment..." -ForegroundColor Green

# Set environment variable dengan force
Write-Host "Setting GO111MODULE environment variable..." -ForegroundColor Blue
[System.Environment]::SetEnvironmentVariable("GO111MODULE", "on", [System.EnvironmentVariableTarget]::User)
[System.Environment]::SetEnvironmentVariable("GO111MODULE", "on", [System.EnvironmentVariableTarget]::Process)

# Refresh environment variables
$env:GO111MODULE = "on"

# Verify
Write-Host "Verifying GO111MODULE setting..." -ForegroundColor Yellow
$goModuleStatus = go env GO111MODULE
Write-Host "GO111MODULE = $goModuleStatus" -ForegroundColor Green

# Check Go version
Write-Host "Go version info:" -ForegroundColor Blue
go version

# Check MySQL connection
Write-Host "Checking MySQL connection..." -ForegroundColor Blue
try {
    $mysql = Test-NetConnection -ComputerName localhost -Port 3306 -WarningAction SilentlyContinue
    if ($mysql.TcpTestSucceeded) {
        Write-Host "MySQL is running on localhost:3306" -ForegroundColor Green
    } else {
        Write-Host "MySQL is not running. Please start Laragon!" -ForegroundColor Red
    }
} catch {
    Write-Host "Could not test MySQL connection. Make sure Laragon is running." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Environment setup complete!" -ForegroundColor Magenta
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "   1. Make sure Laragon MySQL is running" -ForegroundColor White
Write-Host "   2. cd user-service" -ForegroundColor White
Write-Host "   3. go run main.go" -ForegroundColor White
Write-Host ""