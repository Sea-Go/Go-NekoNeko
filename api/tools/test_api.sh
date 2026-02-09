#!/bin/bash
# 测试脚本：收藏夹 API 测试

# 设置变量
BASE_URL="http://localhost:8888"
JWT_TOKEN="Bearer your-jwt-token-here"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== 收藏夹 API 测试脚本 ===${NC}\n"

# 1. 创建收藏项
echo -e "${GREEN}1. 创建收藏项${NC}"
curl -X POST \
  "$BASE_URL/favorite/v1/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: $JWT_TOKEN" \
  -d '{
    "folder_id": 1,
    "object_type": "article",
    "object_id": "12345"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

# 2. 列表收藏项
echo -e "${GREEN}2. 列表收藏项${NC}"
curl -X GET \
  "$BASE_URL/favorite/v1/items?folder_id=1&page=1&page_size=10" \
  -H "Authorization: $JWT_TOKEN" \
  -w "\nHTTP Status: %{http_code}\n\n"

# 3. 删除收藏项
echo -e "${GREEN}3. 删除收藏项${NC}"
curl -X DELETE \
  "$BASE_URL/favorite/v1/items" \
  -H "Content-Type: application/json" \
  -H "Authorization: $JWT_TOKEN" \
  -d '{
    "object_type": "article",
    "object_id": "12345"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

# 4. 测试无效的 JWT token（应该返回 401）
echo -e "${GREEN}4. 测试无效 JWT token（应返回 401）${NC}"
curl -X GET \
  "$BASE_URL/favorite/v1/items?folder_id=1" \
  -H "Authorization: Bearer invalid-token" \
  -w "\nHTTP Status: %{http_code}\n\n"

echo -e "${YELLOW}=== 测试完成 ===${NC}"
