# 集成测试脚本 - 完整测试收藏系统 API
# 使用方法: .\integration_test.ps1 -BaseUrl "http://localhost:8888" -JwtSecret "favorite-secret-key"

param(
    [string]$BaseUrl = "http://localhost:8888",
    [string]$JwtSecret = "favorite-secret-key"
)

$ErrorActionPreference = "Stop"

# 颜色定义
$Colors = @{
    Success = 'Green'
    Error   = 'Red'
    Warning = 'Yellow'
    Info    = 'Cyan'
    Test    = 'Magenta'
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ️  $Message" -ForegroundColor $Colors.Info
}

function Write-Test {
    param([string]$Message)
    Write-Host "🧪 $Message" -ForegroundColor $Colors.Test
}

function Write-Success {
    param([string]$Message)
    Write-Host "✅ $Message" -ForegroundColor $Colors.Success
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "❌ $Message" -ForegroundColor $Colors.Error
}

function Write-Warning-Custom {
    param([string]$Message)
    Write-Host "⚠️  $Message" -ForegroundColor $Colors.Warning
}

# 生成 JWT Token
function Generate-JwtToken {
    param(
        [int64]$UserId = 1,
        [int]$ExpiryHours = 24
    )
    
    Write-Info "生成 JWT Token..."
    
    # 使用 PowerShell 的 System.Security.Cryptography 生成 HS256 签名
    $header = @{
        "alg" = "HS256"
        "typ" = "JWT"
    } | ConvertTo-Json -Compress | ConvertTo-Base64Url
    
    $now = Get-Date -AsUTC
    $payload = @{
        "user_id" = $UserId
        "exp"     = [int]($now.AddHours($ExpiryHours) | Get-Date -UFormat %s)
        "iat"     = [int]($now | Get-Date -UFormat %s)
    } | ConvertTo-Json -Compress | ConvertTo-Base64Url
    
    $signatureInput = "$header.$payload"
    
    # HMAC-SHA256 签名
    $hmac = New-Object System.Security.Cryptography.HMACSHA256
    $hmac.Key = [System.Text.Encoding]::UTF8.GetBytes($JwtSecret)
    $signature = $hmac.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($signatureInput)) | ConvertTo-Base64Url
    
    $token = "$signatureInput.$signature"
    Write-Success "JWT Token 已生成"
    
    return $token
}

function ConvertTo-Base64Url {
    param(
        [Parameter(ValueFromPipeline = $true)]
        [string]$Text
    )
    
    if ($Text -is [string]) {
        $bytes = [System.Text.Encoding]::UTF8.GetBytes($Text)
    } else {
        $bytes = $Text
    }
    
    [Convert]::ToBase64String($bytes) -replace '\+', '-' -replace '/', '_' -replace '=+$', ''
}

# API 请求函数
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Endpoint,
        [object]$Body,
        [string]$Token,
        [int]$ExpectedStatus = 200
    )
    
    $uri = "$BaseUrl$Endpoint"
    $headers = @{
        "Content-Type"  = "application/json"
        "Authorization" = "Bearer $Token"
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
        
        if ($response.StatusCode -eq $ExpectedStatus) {
            return @{
                Success = $true
                Status  = $response.StatusCode
                Body    = $response.Content | ConvertFrom-Json
            }
        } else {
            return @{
                Success = $false
                Status  = $response.StatusCode
                Body    = $response.Content
            }
        }
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.Value__
        $content = $_.Exception.Response.Content.ReadAsStream() | ForEach-Object { [System.IO.StreamReader]::new($_).ReadToEnd() }
        
        if ($statusCode -eq $ExpectedStatus) {
            return @{
                Success = $true
                Status  = $statusCode
                Body    = $content | ConvertFrom-Json
            }
        } else {
            return @{
                Success = $false
                Status  = $statusCode
                Body    = $content
            }
        }
    }
}

# 测试场景
function Test-Api {
    Write-Host ""
    Write-Host "================================================" -ForegroundColor $Colors.Test
    Write-Host "🚀 启动收藏系统集成测试" -ForegroundColor $Colors.Test
    Write-Host "================================================" -ForegroundColor $Colors.Test
    Write-Host ""
    
    # 生成 token
    $token = Generate-JwtToken
    Write-Host ""
    
    # 测试场景 1: 创建收藏夹
    Write-Test "测试 1: 创建收藏夹"
    $response = Invoke-ApiRequest -Method POST -Endpoint "/favorite/v1/folders" `
        -Body @{
            "name"      = "我的标签"
            "is_public" = $false
        } `
        -Token $token -ExpectedStatus 200
    
    if ($response.Success) {
        Write-Success "收藏夹创建成功"
        $folderId = $response.Body.id
        Write-Info "收藏夹ID: $folderId"
    } else {
        Write-Error-Custom "收藏夹创建失败"
        Write-Error-Custom "状态码: $($response.Status)"
        return
    }
    Write-Host ""
    
    # 测试场景 2: 创建收藏项
    Write-Test "测试 2: 创建收藏项（有效的收藏）"
    $response = Invoke-ApiRequest -Method POST -Endpoint "/favorite/v1/items" `
        -Body @{
            "folder_id"   = $folderId
            "object_type" = "article"
            "object_id"   = "12345"
        } `
        -Token $token -ExpectedStatus 200
    
    if ($response.Success) {
        Write-Success "收藏项创建成功"
        $itemId = $response.Body.id
        Write-Info "收藏项ID: $itemId"
    } else {
        Write-Error-Custom "收藏项创建失败"
        Write-Error-Custom "状态码: $($response.Status)"
        return
    }
    Write-Host ""
    
    # 测试场景 3: 重复收藏（应该返回 409）
    Write-Test "测试 3: 重复收藏同一对象（应返回 409）"
    $response = Invoke-ApiRequest -Method POST -Endpoint "/favorite/v1/items" `
        -Body @{
            "folder_id"   = $folderId
            "object_type" = "article"
            "object_id"   = "12345"
        } `
        -Token $token -ExpectedStatus 409
    
    if ($response.Status -eq 409) {
        Write-Success "正确返回 409 Conflict"
        Write-Info "错误信息: $($response.Body.message)"
    } else {
        Write-Warning-Custom "预期 409，实际 $($response.Status)"
    }
    Write-Host ""
    
    # 测试场景 4: 列表收藏项
    Write-Test "测试 4: 列表收藏项"
    $response = Invoke-ApiRequest -Method GET -Endpoint "/favorite/v1/items?folder_id=$folderId&page=1&page_size=10" `
        -Token $token -ExpectedStatus 200
    
    if ($response.Success) {
        Write-Success "收藏项列表获取成功"
        Write-Info "总数: $($response.Body.total)"
        Write-Info "项目数: $($response.Body.items.Count)"
    } else {
        Write-Error-Custom "获取收藏项列表失败"
        Write-Error-Custom "状态码: $($response.Status)"
    }
    Write-Host ""
    
    # 测试场景 5: 删除收藏项
    Write-Test "测试 5: 删除收藏项"
    $response = Invoke-ApiRequest -Method DELETE -Endpoint "/favorite/v1/items" `
        -Body @{
            "object_type" = "article"
            "object_id"   = "12345"
        } `
        -Token $token -ExpectedStatus 200
    
    if ($response.Success) {
        Write-Success "收藏项删除成功"
    } else {
        Write-Error-Custom "收藏项删除失败"
        Write-Error-Custom "状态码: $($response.Status)"
    }
    Write-Host ""
    
    # 测试场景 6: 验证删除（应该返回空列表）
    Write-Test "测试 6: 验证删除后列表为空"
    $response = Invoke-ApiRequest -Method GET -Endpoint "/favorite/v1/items?folder_id=$folderId&page=1&page_size=10" `
        -Token $token -ExpectedStatus 200
    
    if ($response.Success -and $response.Body.total -eq 0) {
        Write-Success "验证成功，列表为空"
    } else {
        Write-Warning-Custom "预期空列表，实际包含 $($response.Body.total) 项"
    }
    Write-Host ""
    
    # 测试场景 7: 无效 token（应该返回 401）
    Write-Test "测试 7: 使用无效 token（应返回 401）"
    $response = Invoke-ApiRequest -Method GET -Endpoint "/favorite/v1/items?folder_id=$folderId" `
        -Token "invalid-token" -ExpectedStatus 401
    
    if ($response.Status -eq 401) {
        Write-Success "正确返回 401 Unauthorized"
    } else {
        Write-Warning-Custom "预期 401，实际 $($response.Status)"
    }
    Write-Host ""
    
    # 测试场景 8: 缺失 token（应该返回 401）
    Write-Test "测试 8: 缺失 Authorization header（应返回 401）"
    
    try {
        $uri = "$BaseUrl/favorite/v1/items?folder_id=$folderId"
        $response = Invoke-WebRequest -Uri $uri -Method GET -UseBasicParsing
        Write-Warning-Custom "预期 401，但请求成功"
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.Value__
        if ($statusCode -eq 401) {
            Write-Success "正确返回 401 Unauthorized"
        } else {
            Write-Warning-Custom "预期 401，实际 $statusCode"
        }
    }
    
    Write-Host ""
    Write-Host "================================================" -ForegroundColor $Colors.Success
    Write-Host "✨ 测试完成！" -ForegroundColor $Colors.Success
    Write-Host "================================================" -ForegroundColor $Colors.Success
}

# Main
try {
    Test-Api
} catch {
    Write-Error-Custom "测试中发生错误: $_"
    exit 1
}
