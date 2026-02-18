package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/admin/rpc/internal/model"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	_, err := l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserNotExist, err)
			return nil, status.Error(codes.NotFound, "用户不存在")
		}
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}
	toUpdate := &model.User{}
	if len(in.Username) > 0 {
		existUser, err := l.svcCtx.AdminModel.FindOneUserByUsername(l.ctx, in.Username)
		fmt.Println("123456789")
		fmt.Println(err)
		if err == nil && in.Uid != existUser.Uid {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorUserExist, err)
			return nil, status.Error(codes.AlreadyExists, "用户名已存在")
		}
		if err != nil && err != model.ErrorNotFound {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
			return nil, status.Error(codes.Internal, "DB查询失败")
		}
		toUpdate.Username = in.Username
	}
	if len(in.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(in.Password)
		if e != nil {
			logger.LogBusinessErr(l.ctx, errmsg.ErrorServerCommon, e)
			return nil, status.Error(codes.Internal, "密码加密失败")
		}
		toUpdate.Password = newPassword
	}
	if len(in.Email) > 0 {
		toUpdate.Email = in.Email
	}
	if in.ExtraInfo != nil {
		toUpdate.ExtraInfo = in.ExtraInfo
	}
	err = l.svcCtx.AdminModel.UpdateOneUserByUid(l.ctx, in.Uid, toUpdate)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbUpdate, err)
		return nil, status.Error(codes.Internal, "DB更新失败")
	}
	var newUser *model.User
	newUser, err = l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorDbSelect, err)
		return nil, status.Error(codes.Internal, "DB查询失败")
	}
	logger.LogInfo(l.ctx, fmt.Sprintf("update success,uid : %d", in.Uid))
	return &pb.UpdateUserResp{
		User: &pb.UserInfo{
			Uid:       newUser.Uid,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Status:    uint64(newUser.Status),
			ExtraInfo: newUser.ExtraInfo,
		},
	}, nil
}
