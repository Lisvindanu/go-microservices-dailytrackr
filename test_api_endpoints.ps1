# DailyTrackr API Testing Script
Write-Host "🧪 Testing DailyTrackr API Endpoints..." -ForegroundColor Green

$baseUrl = "http://localhost:3000"
$token = ""

# Function to make HTTP requests
function Invoke-ApiTest {
    param(
        [string]$Method,
        [string]$Url,
        [hashtable]$Body = $null,
        [hashtable]$Headers = @{}
    )

    Write-Host "`n🔍 Testing: $Method $Url" -ForegroundColor Cyan

    try {
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            ContentType = "application/json"
        }

        if ($Body) {
            $params.Body = ($Body | ConvertTo-Json)
        }

        $response = Invoke-RestMethod @params
        Write-Host "✅ Success: $($response.message)" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $statusCode = $_.Exception.Response.StatusCode
            Write-Host "   Status Code: $statusCode" -ForegroundColor Yellow
        }
        return $null
    }
}

# 1. Test Gateway Health
Write-Host "`n🏥 Testing Gateway Health..." -ForegroundColor Magenta
$healthResponse = Invoke-ApiTest -Method "GET" -Url "$baseUrl/"

if ($healthResponse) {
    Write-Host "   Gateway: $($healthResponse.status)" -ForegroundColor White
    Write-Host "   Services Available:" -ForegroundColor White
    foreach ($service in $healthResponse.services.PSObject.Properties) {
        Write-Host "     - $($service.Name): $($service.Value)" -ForegroundColor Gray
    }
}

# 2. Test User Registration
Write-Host "`n👤 Testing User Registration..." -ForegroundColor Magenta
$registerData = @{
    username = "testuser_$(Get-Random -Minimum 1000 -Maximum 9999)"
    email = "testuser_$(Get-Random -Minimum 1000 -Maximum 9999)@example.com"
    password = "password123"
}

$registerResponse = Invoke-ApiTest -Method "POST" -Url "$baseUrl/api/users/auth/register" -Body $registerData

if ($registerResponse -and $registerResponse.data.token) {
    $token = $registerResponse.data.token
    $authHeaders = @{ "Authorization" = "Bearer $token" }
    Write-Host "   🔑 Token obtained: $($token.Substring(0, 20))..." -ForegroundColor Green
}

# 3. Test User Login
Write-Host "`n🔐 Testing User Login..." -ForegroundColor Magenta
$loginData = @{
    email = $registerData.email
    password = $registerData.password
}

$loginResponse = Invoke-ApiTest -Method "POST" -Url "$baseUrl/api/users/auth/login" -Body $loginData

# 4. Test Activity Service (if token available)
if ($token) {
    Write-Host "`n📝 Testing Activity Service..." -ForegroundColor Magenta

    # Create Activity
    $activityData = @{
        title = "Test Activity - API Test"
        start_time = "2025-07-12T10:00:00Z"
        duration_mins = 60
        cost = 25000
        note = "Testing activity creation via API"
    }

    $activityResponse = Invoke-ApiTest -Method "POST" -Url "$baseUrl/api/activities/api/v1/activities" -Body $activityData -Headers $authHeaders

    if ($activityResponse) {
        $activityId = $activityResponse.data.id
        Write-Host "   📋 Activity created with ID: $activityId" -ForegroundColor Green

        # Get Activities
        Invoke-ApiTest -Method "GET" -Url "$baseUrl/api/activities/api/v1/activities" -Headers $authHeaders

        # Get Single Activity
        if ($activityId) {
            Invoke-ApiTest -Method "GET" -Url "$baseUrl/api/activities/api/v1/activities/$activityId" -Headers $authHeaders
        }
    }
}

# 5. Test Habit Service (if token available)
if ($token) {
    Write-Host "`n🎯 Testing Habit Service..." -ForegroundColor Magenta

    # Create Habit
    $habitData = @{
        title = "Test Habit - API Test"
        start_date = "2025-07-12"
        end_date = "2025-08-11"
        reminder_time = "09:00"
    }

    $habitResponse = Invoke-ApiTest -Method "POST" -Url "$baseUrl/api/habits/api/v1/habits" -Body $habitData -Headers $authHeaders

    if ($habitResponse) {
        $habitId = $habitResponse.data.id
        Write-Host "   🎯 Habit created with ID: $habitId" -ForegroundColor Green

        # Get Habits
        Invoke-ApiTest -Method "GET" -Url "$baseUrl/api/habits/api/v1/habits" -Headers $authHeaders

        # Create Habit Log
        if ($habitId) {
            $habitLogData = @{
                habit_id = $habitId
                date = "2025-07-12"
                status = "DONE"
                note = "Completed habit for testing"
            }

            $habitLogResponse = Invoke-ApiTest -Method "POST" -Url "$baseUrl/api/habits/api/v1/habits/$habitId/logs" -Body $habitLogData -Headers $authHeaders

            # Get Habit Stats
            Invoke-ApiTest -Method "GET" -Url "$baseUrl/api/habits/api/v1/habits/$habitId/stats" -Headers $authHeaders
        }
    }
}

# Summary
Write-Host "`n📊 Test Summary:" -ForegroundColor Magenta
Write-Host "   🏥 Gateway Health: Tested" -ForegroundColor White
Write-Host "   👤 User Registration: Tested" -ForegroundColor White
Write-Host "   🔐 User Login: Tested" -ForegroundColor White
if ($token) {
    Write-Host "   📝 Activity Service: Tested" -ForegroundColor White
    Write-Host "   🎯 Habit Service: Tested" -ForegroundColor White
} else {
    Write-Host "   📝 Activity Service: Skipped (no token)" -ForegroundColor Yellow
    Write-Host "   🎯 Habit Service: Skipped (no token)" -ForegroundColor Yellow
}

Write-Host "`n✅ API Testing Complete!" -ForegroundColor Green
Write-Host "🔧 Check the logs above for any errors" -ForegroundColor Cyan