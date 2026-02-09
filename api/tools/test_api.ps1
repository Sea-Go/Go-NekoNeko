# API 测试脚本 (Windows PowerShell)
# 使用方法: .\test_api.ps1

param(
    [string]$BaseUrl = "http://localhost:8888",
    [string]$JwtToken = ""
)

# 如果没有提供 JWT token，提示用户生成
if ([string]::IsNullOrEmpty($JwtToken)) {
    Write-Host "错误: 请提供有效的 JWT token" -ForegroundColor Red
    Write-Host "`n生成 JWT token 的方式:" -ForegroundColor Yellow
    Write-Host "  1. cd api\tools"
    Write-Host "  2. go run jwt_generator.go"
    Write-Host "  3. 将输出的 token 传递给此脚本"
    Write-Host "`n使用方法:" -ForegroundColor Yellow
    Write-Host "  .\test_api.ps1 -BaseUrl 'http://localhost:8888' -JwtToken 'Bearer <token>'"
    exit 1
}

# 设置默认头
$headers = @{
    "Content-Type"  = "application/json"
    "Authorization" = $JwtToken
}

Write-Host "=== 收藏夹 API 测试脚本 ===" -ForegroundColor Green
Write-Host "基础 URL: $BaseUrl`n" -ForegroundColor Gray

# 函数：发送请求并显示结果
function Test-ApiEndpoint {
    param(
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [string]$QueryString = ""
    )
    
    $url = "$BaseUrl$Endpoint"
    if (-not [string]::IsNullOrEmpty($QueryString)) {
        $url = "$url?$QueryString"
    }
    
    Write-Host "请求: $Method $Endpoint" -ForegroundColor Cyan
    if (-not [string]::IsNullOrEmpty($QueryString)) {
        Write-Host "查询字符串: $QueryString" -ForegroundColor Gray
    }
    if ($Body) {
        Write-Host "请求体: $(ConvertTo-Json $Body)" -ForegroundColor Gray
    }
    
    try {
        $params = @{
            Uri     = $url
            Method  = $Method
            Headers = $headers
        }
        
        if ($Body) {
            $params.Body = ConvertTo-Json $Body
        }
        
        $response = Invoke-WebRequest @params
        
        Write-Host "状态码: $($response.StatusCode)" -ForegroundColor Green
        Write-Host "响应体:" -ForegroundColor Gray
        Write-Host (ConvertFrom-Json $response.Content | ConvertTo-Json) -ForegroundColor Gray
    }
    catch {
        if ($_.Exception.Response) {
            $statusCode = [int]$_.Exception.Response.StatusCode
            Write-Host "状态码: $statusCode" -ForegroundColor Red
            
            try {
                $errorBody = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream()).ReadToEnd()
                Write-Host "错误响应:" -ForegroundColor Gray
                Write-Host $errorBody -ForegroundColor Red
            }
            catch {
                Write-Host "无法解析错误响应" -ForegroundColor Red
            }
        }
        else {
            Write-Host "请求失败: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host "`n" -ForegroundColor Gray
}

# 测试 1: 创建收藏项
Write-Host "测试 1: 创建收藏项" -ForegroundColor Yellow
Test-ApiEndpoint -Method "POST" -Endpoint "/favorite/v1/items" -Body @{
    folder_id  = 1
    object_type = "article"
    object_id  = "12345"
}

# 测试 2: 列表收藏项
Write-Host "测试 2: 列表收藏项" -ForegroundColor Yellow
Test-ApiEndpoint -Method "GET" -Endpoint "/favorite/v1/items" -QueryString "folder_id=1&page=1&page_size=10"

# 测试 3: 删除收藏项
Write-Host "测试 3: 删除收藏项" -ForegroundColor Yellow
Test-ApiEndpoint -Method "DELETE" -Endpoint "/favorite/v1/items" -Body @{
    object_type = "article"
    object_id   = "12345"
}

# 测试 4: 测试无效 token（应返回 401）
Write-Host "测试 4: 使用无效 token（应返回 401）" -ForegroundColor Yellow
$invalidHeaders = @{
    "Content-Type"  = "application/json"
    "Authorization" = "Bearer invalid-token-here"
}

try {
    $response = Invoke-WebRequest -Uri "$BaseUrl/favorite/v1/items?folder_id=1" `
        -Method GET `
        -Headers $invalidHeaders
}
catch {
    $statusCode = [int]$_.Exception.Response.StatusCode
    Write-Host "状态码: $statusCode (预期: 401)" -ForegroundColor $(if ($statusCode -eq 401) { "Green" } else { "Red" })
}

Write-Host "`n=== 测试完成 ===" -ForegroundColor Green
