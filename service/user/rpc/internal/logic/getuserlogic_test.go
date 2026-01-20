package logic

import (
	"testing"

	pb "sea-try-go/service/user/rpc/pb"
)

func TestGetUser_Success(t *testing.T) {
	db := setupTestDB()
	cleanupTestUsers(db)

	// 创建测试用户
	testUser := createTestUser(db, "getuser", "password123", "getuser@example.com")

	svcCtx := setupTestServiceContext(db)
	ctx := newTestContext()

	logic := NewGetUserLogic(ctx, svcCtx)

	req := &pb.GetUserReq{
		Id: testUser.Id,
	}

	resp, err := logic.GetUser(req)

	if err != nil {
		t.Fatalf("获取用户请求失败: %v", err)
	}

	if !resp.Found {
		t.Error("应该找到用户")
	}

	if resp.User == nil {
		t.Fatal("User 不应为 nil")
	}

	if resp.User.Username != "getuser" {
		t.Errorf("用户名不匹配: 期望 %s, 实际 %s", "getuser", resp.User.Username)
	}

	if resp.User.Email != "getuser@example.com" {
		t.Errorf("邮箱不匹配: 期望 %s, 实际 %s", "getuser@example.com", resp.User.Email)
	}

	t.Logf("✅ 获取用户成功，用户名: %s", resp.User.Username)
}

func TestGetUser_NotFound(t *testing.T) {
	db := setupTestDB()
	cleanupTestUsers(db)

	db = db.Table("test_users")
	svcCtx := setupTestServiceContext(db)
	ctx := newTestContext()

	logic := NewGetUserLogic(ctx, svcCtx)

	req := &pb.GetUserReq{
		Id: 99999, // 不存在的 ID
	}

	resp, err := logic.GetUser(req)

	// GetUser 在用户不存在时会返回 error
	if err == nil {
		if resp.Found {
			t.Error("用户不存在时 Found 应为 false")
		}
	}

	if resp != nil && resp.Found {
		t.Error("用户不存在时 Found 应为 false")
	}

	t.Log("✅ 用户不存在被正确识别")
}
