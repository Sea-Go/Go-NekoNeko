package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelfLogic {
	return &GetSelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSelfLogic) GetSelf(in *pb.GetSelfReq) (*pb.GetSelfResp, error) {
	admin, err := l.svcCtx.AdminModel.FindOneAdminByUid(l.ctx, in.Uid)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbSelect))
	}
	return &pb.GetSelfResp{
		Admin: &pb.AdminInfo{
			Uid:       admin.Uid,
			Username:  admin.Username,
			Email:     admin.Email,
			ExtraInfo: admin.ExtraInfo,
		},
	}, nil
}
