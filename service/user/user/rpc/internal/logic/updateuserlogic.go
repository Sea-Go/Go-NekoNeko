package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/rpc/internal/model"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

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

	err := l.svcCtx.UserModel.UpdateUserById(l.ctx, in.Uid, toUpdate)

	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}

	newUser, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, in.Uid)

	if err != nil {
		if err == model.ErrorNotFound {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
		}
		return nil, err
	}

	return &pb.UpdateUserResp{
		User: &pb.UserInfo{
			Uid:       newUser.Uid,
			Username:  newUser.Username,
			Email:     newUser.Email,
			ExtraInfo: newUser.ExtraInfo,
		},
	}, nil
}
