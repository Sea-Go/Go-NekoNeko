package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"

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

	result := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.Id).Update("status", 1)
	if result.Error != nil {
		return nil, errors.New("封禁失败" + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}
	return &pb.BanUserResp{
		Success: true,
	}, nil
}
