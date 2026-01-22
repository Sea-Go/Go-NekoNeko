package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"
	"sea-try-go/service/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanUserLogic {
	return &BanUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BanUserLogic) BanUser(in *pb.BanUserReq) (*pb.BanUserResp, error) {

	err := l.svcCtx.AdminModel.UpdateUserStatusByUid(l.ctx, in.Uid, 1)
	if err != nil {
		if err == model.ErrorNotFound {
			return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserNotExist))
		}
		return nil, errors.New("封禁失败" + err.Error())
	}
	return &pb.BanUserResp{
		Success: true,
	}, nil
}
