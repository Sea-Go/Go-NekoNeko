package logic

import (
	"context"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user, err := l.svcCtx.UserModel.GetByUsername(in.Username)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return nil, err
	}

	if user == nil {
		l.Infof("用户不存在: %s", in.Username)
		return &pb.LoginResp{Uid: 0, Status: 0}, nil
	}

	// 验证密码
	if !cryptx.PasswordVerify(in.Password, user.Password) {
		l.Infof("密码错误: %s", in.Username)
		return &pb.LoginResp{Uid: 0, Status: 0}, nil
	}

	return &pb.LoginResp{
		Uid:    user.Uid,
		Status: uint64(user.Status),
	}, nil
}
