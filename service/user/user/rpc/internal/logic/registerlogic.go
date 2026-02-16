package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/rpc/internal/model"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	_, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.Username)
	if err == nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorUserExist, fmt.Errorf("username has existed"))
		return nil, status.Error(codes.AlreadyExists, "用户名已存在")
	}
	if err != model.ErrorNotFound {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}
	truePassword, err := cryptx.PasswordEncrypt(in.Password)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorServerCommon, err)
		return nil, status.Error(codes.Internal, "密码生成失败")
	}

	var uid int64
	uid, err = snowflake.GetID()
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorServerCommon, err)
		return nil, status.Error(codes.Internal, "uid生成失败")
	}
	newUser := model.User{
		Uid:       uid,
		Username:  in.Username,
		Password:  truePassword,
		Email:     in.Email,
		Score:     0,
		Status:    0,
		ExtraInfo: in.ExtraInfo,
	}
	err = l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return nil, status.Error(codes.Internal, "DB更新失败")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("register success,uid : %d", uid))
	return &pb.CreateUserResp{
		Uid: newUser.Uid,
	}, nil
}
