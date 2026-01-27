package logic

import (
	"context"
	"errors"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/user/admin/rpc/internal/model"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/cryptx"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminLogic {
	return &CreateAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAdminLogic) CreateAdmin(in *pb.CreateAdminReq) (*pb.CreateAdminResp, error) {

	_, err := l.svcCtx.AdminModel.FindOneAdminByUsername(l.ctx, in.Username)
	if err == nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserExist))
	}

	password, err := cryptx.PasswordEncrypt(in.Password)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	uid, err := snowflake.GetID()
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	admin := &model.Admin{
		Uid:       uid,
		Username:  in.Username,
		Password:  password,
		Email:     in.Email,
		ExtraInfo: in.ExtraInfo,
	}
	err = l.svcCtx.AdminModel.InsertOneAdmin(l.ctx, admin)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbUpdate))
	}
	return &pb.CreateAdminResp{
		Uid: admin.Uid,
	}, nil
}
