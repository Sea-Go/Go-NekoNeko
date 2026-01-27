package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/internal/model"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"

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
	_, err := l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err == model.ErrorNotFound {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
	}
	if err != nil {
		return nil, err
	}
	toUpdate := &model.User{}
	if len(in.Username) > 0 {
		toUpdate.Username = in.Username
	}
	if len(in.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(in.Password)
		if e != nil {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
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
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	var newUser *model.User
	newUser, err = l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
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
