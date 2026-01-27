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

type UpdateSelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSelfLogic {
	return &UpdateSelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSelfLogic) UpdateSelf(in *pb.UpdateSelfReq) (*pb.UpdateSelfResp, error) {

	toUpdate := &model.Admin{}
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
	err := l.svcCtx.AdminModel.UpdateOneAdminByUid(l.ctx, in.Uid, toUpdate)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	var newAdmin *model.Admin
	newAdmin, err = l.svcCtx.AdminModel.FindOneAdminByUid(l.ctx, in.Uid)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSelfResp{
		Success: true,
		Admin: &pb.AdminInfo{
			Uid:       newAdmin.Uid,
			Username:  newAdmin.Username,
			Email:     newAdmin.Email,
			ExtraInfo: newAdmin.ExtraInfo,
		},
	}, nil
}
