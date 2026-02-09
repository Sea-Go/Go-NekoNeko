package logic

import (
	"context"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	// 获取现有用户
	user, err := l.svcCtx.UserModel.GetByUid(in.Uid)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	if user == nil {
		return &pb.UpdateUserResp{User: nil}, nil
	}

	// 更新字段
	if in.Username != "" {
		user.Username = in.Username
	}
	if in.Password != "" {
		hashedPassword, err := cryptx.PasswordEncrypt(in.Password)
		if err != nil {
			l.Errorf("密码加密失败: %v", err)
			return nil, err
		}
		user.Password = hashedPassword
	}
	if in.Email != "" {
		user.Email = in.Email
	}
	if in.ExtraInfo != nil {
		user.ExtraInfo = in.ExtraInfo
	}

	// 保存更改
	err = l.svcCtx.UserModel.Update(user)
	if err != nil {
		l.Errorf("更新用户失败: %v", err)
		return nil, err
	}

	return &pb.UpdateUserResp{
		User: &pb.UserInfo{
			Uid:       user.Uid,
			Score:     uint32(user.Score),
			Username:  user.Username,
			Email:     user.Email,
			ExtraInfo: user.ExtraInfo,
		},
	}, nil
}
