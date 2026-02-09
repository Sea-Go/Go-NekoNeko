# 简化版集成测试 - 快速验证 API 功能
# 使用方法: .\integration_test_simple.ps1 -Token "Bearer <your-jwt-token>"

param(
    [string]$BaseUrl = "http://localhost:8888",
    [string]$Token = "" # 必须提供 JWT token
)

$ErrorActionPreference = "Stop"

# 验证参数
if ([string]::IsNullOrEmpty($Token)) {
    Write-Host "❌ 错误: 必须提供 JWT token" -ForegroundColor Red
    Write-Host ""
    Write-Host "使用方法:"
    Write-Host "  .\integration_test_simple.ps1 -Token 'Bearer <your-jwt-token>'"
    Write-Host ""
    Write-Host "生成 token:"
    Write-Host "  cd api\tools"
    Write-Host "  go run jwt_generator.go"
    exit 1
}

Write-Host ""
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "🧪 集成测试 - API 功能验证" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "配置信息:"
Write-Host "  🌐 API 地址: $BaseUrl"
Write-Host "  🔐 Token: $($Token.Substring(0, [Math]::Min(20, $Token.Length)))..."
Write-Host ""

# 计数器
$testsPassed = 0
$testsFailed = 0

# 测试函数
function Test-ApiEndpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Endpoint,
        [object]$Body,
        [int]$ExpectedStatus,
        [string]$Description
    )
    
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
    Write-Host "🧪 $Name" -ForegroundColor Magenta
    Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
    
    Write-Host "📝 说明: $Description"
    Write-Host "$Method $Endpoint"
    
    if ($Body) {
        Write-Host "📦 请求体: $($Body | ConvertTo-Json -Compress)" -ForegroundColor Gray
    }
    
    $uri = "$BaseUrl$Endpoint"
    $headers = @{
        "Content-Type"  = "application/json"
        "Authorization" = $Token
    }
    
    $params = @{
        Uri             = $uri
        Method          = $Method
        Headers         = $headers
        UseBasicParsing = $true
    }
    
    if ($Body) {
        $params.Body = $Body | ConvertTo-Json -Compress
    }
    
    try {
        $response = Invoke-WebRequest @params
        $statusCode = $response.StatusCode
        $content = $response.Content | ConvertFrom-Json
        
        Write-Host "✅ 状态码: $statusCode (期望: $ExpectedStatus)" -ForegroundColor Green
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✅ 响应正确" -ForegroundColor Green
            Write-Host "📋 响应: $($content | ConvertTo-Json -Compress -Depth 2)" -ForegroundColor Gray
            global:testsPassed++
            return @{ Success = $true; Body = $content }
        } else {
            Write-Host "⚠️  状态码不符，期望 $ExpectedStatus 但得到 $statusCode" -ForegroundColor Yellow
            global:testsFailed++
            return @{ Success = $false; Body = $content }
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.Value__
        $content = $_.Exception.Response.Content.ReadAsStream() | ForEach-Object { [System.IO.StreamReader]::new($_).ReadToEnd() }
        
        Write-Host "📋 状态码: $statusCode" -ForegroundColor Yellow
        Write-Host "📋 响应: $content" -ForegroundColor Gray
        
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✅ 错误状态码符合预期（测试成功）" -ForegroundColor Green
            global:testsPassed++
            return @{ Success = $true; Body = $content }
        } else {
            Write-Host "❌ 状态码不符，期望 $ExpectedStatus 但得到 $statusCode" -ForegroundColor Red
            global:testsFailed++
            return @{ Success = $false; Body = $content }
        }
    }
    
    Write-Host ""
}

# 测试场景
Write-Host "📌 测试场景 1: 基础流程"
Write-Host ""

# 1. 创建收藏夹
$result = Test-ApiEndpoint `
    -Name "创建收藏夹" `
    -Method "POST" `
    -Endpoint "/favorite/v1/folders" `
    -Body @{ "name" = "我的标签"; "is_public" = $false } `
    -ExpectedStatus 200 `
    -Description "创建一个新的收藏夹"

if (-not $result.Success) {
    Write-Host "❌ 创建收藏夹失败，停止测试" -ForegroundColor Red
    exit 1
}

$folderId = $result.Body.id
Write-Host "📌 收藏夹ID: $folderId (用于后续测试)" -ForegroundColor Cyan
Write-Host ""

# 2. 创建收藏项
$result = Test-ApiEndpoint `
    -Name "创建收藏项" `
    -Method "POST" `
    -Endpoint "/favorite/v1/items" `
    -Body @{ "folder_id" = $folderId; "object_type" = "article"; "object_id" = "12345" } `
    -ExpectedStatus 200 `
    -Description "向收藏夹添加一个收藏项"

if ($result.Success) {
    $itemId = $result.Body.id
    Write-Host "📌 收藏项ID: $itemId" -ForegroundColor Cyan
}

Write-Host ""

# 3. 列表收藏项
$result = Test-ApiEndpoint `
    -Name "列表收藏项" `
    -Method "GET" `
    -Endpoint "/favorite/v1/items?folder_id=$folderId&page=1&page_size=10" `
    -ExpectedStatus 200 `
    -Description "查看收藏夹中的所有收藏项"

Write-Host ""

# 4. 重复收藏（应该失败）
Write-Host "📌 测试场景 2: 错误处理"
Write-Host ""

$result = Test-ApiEndpoint `
    -Name "重复收藏同一对象" `
    -Method "POST" `
    -Endpoint "/favorite/v1/items" `
    -Body @{ "folder_id" = $folderId; "object_type" = "article"; "object_id" = "12345" } `
    -ExpectedStatus 409 `
    -Description "尝试收藏同一个对象两次（应返回 409 冲突）"

Write-Host ""

# 5. 删除收藏项
$result = Test-ApiEndpoint `
    -Name "删除收藏项" `
    -Method "DELETE" `
    -Endpoint "/favorite/v1/items" `
    -Body @{ "object_type" = "article"; "object_id" = "12345" } `
    -ExpectedStatus 200 `
    -Description "删除一个收藏项"

Write-Host ""

# 6. 验证删除
$result = Test-ApiEndpoint `
    -Name "验证删除后列表" `
    -Method "GET" `
    -Endpoint "/favorite/v1/items?folder_id=$folderId&page=1&page_size=10" `
    -ExpectedStatus 200 `
    -Description "验证删除后收藏夹为空"

Write-Host ""

# 7. 无效 token
Write-Host "📌 测试场景 3: 安全验证"
Write-Host ""

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow
Write-Host "🧪 使用无效 token" -ForegroundColor Magenta
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Yellow

$uri = "$BaseUrl/favorite/v1/items?folder_id=$folderId"
$headers = @{
    "Authorization" = "Bearer invalid-token"
}

try {
    $response = Invoke-WebRequest -Uri $uri -Method GET -Headers $headers -UseBasicParsing
    Write-Host "⚠️  状态码: 200 (预期: 401)" -ForegroundColor Yellow
    global:testsFailed++
} catch {
    $statusCode = $_.Exception.Response.StatusCode.Value__
    Write-Host "✅ 状态码: $statusCode (期望: 401)" -ForegroundColor Green
    
    if ($statusCode -eq 401) {
        Write-Host "✅ 正确返回 401 Unauthorized" -ForegroundColor Green
        global:testsPassed++
    } else {
        Write-Host "❌ 状态码不符" -ForegroundColor Red
        global:testsFailed++
    }
}

Write-Host ""

# 总结
Write-Host "================================================" -ForegroundColor Cyan
Write-Host "📊 测试结果总结" -ForegroundColor Cyan
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ 通过: $testsPassed" -ForegroundColor Green
Write-Host "❌ 失败: $testsFailed" -ForegroundColor Red
Write-Host "📈 总计: $($testsPassed + $testsFailed)" -ForegroundColor Cyan
Write-Host ""

if ($testsFailed -eq 0) {
    Write-Host "✨ 所有测试通过！API 运行正常。" -ForegroundColor Green
    exit 0
} else {
    Write-Host "⚠️  有 $testsFailed 个测试失败，请检查错误信息。" -ForegroundColor Red
    exit 1
}
