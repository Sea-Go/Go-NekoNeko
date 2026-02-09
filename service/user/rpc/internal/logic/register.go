package logic

import (
	"context"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	// 检查用户名是否存在
	existing, err := l.svcCtx.UserModel.GetByUsername(in.Username)
	if err != nil {
		l.Errorf("检查用户名失败: %v", err)
		return nil, err
	}

	if existing != nil {
		l.Infof("用户名已存在: %s", in.Username)
		return &pb.CreateUserResp{Uid: 0}, nil
	}

	// 生成UID
	uid, err := snowflake.GetID()
	if err != nil {
		l.Errorf("生成UID失败: %v", err)
		return nil, err
	}

	// 密码加密
	hashedPassword, err := cryptx.PasswordEncrypt(in.Password)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Uid:       uid,
		Username:  in.Username,
		Password:  hashedPassword,
		Email:     in.Email,
		Status:    0,
		Score:     0,
		ExtraInfo: in.ExtraInfo,
	}

	err = l.svcCtx.UserModel.Create(user)
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, err
	}

	return &pb.CreateUserResp{Uid: uid}, nil
}
