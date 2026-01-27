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

type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(in *pb.ResetUserPasswordReq) (*pb.ResetUserPasswordResp, error) {

	_, err := l.svcCtx.AdminModel.FindOneUserByUid(l.ctx, in.Uid)
	if err != nil {
		if err == model.ErrorNotFound {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
		}
		//可能需要写日志来记录错误
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	var password string
	password, err = cryptx.PasswordEncrypt(l.svcCtx.Config.System.DefaultPassword)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorServerCommon))
	}
	err = l.svcCtx.AdminModel.UpdateUserPasswordByUid(l.ctx, in.Uid, password)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbUpdate))
	}
	return &pb.ResetUserPasswordResp{
		Success: true,
	}, nil
}
