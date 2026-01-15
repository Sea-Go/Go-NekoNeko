package logic

import (
	"context"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"

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
	admin := model.Admin{}
	err := l.svcCtx.DB.Where("id = ?", uint64(in.Id)).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &pb.GetSelfResp{
		Admin: &pb.AdminInfo{
			Id:        admin.Id,
			Username:  admin.Username,
			Email:     admin.Email,
			ExtraInfo: admin.ExtraInfo,
		},
	}, nil
}
