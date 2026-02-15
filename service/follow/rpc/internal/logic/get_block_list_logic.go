package logic

import (
	"context"

	"sea-try-go/service/follow/rpc/internal/svc"
	"sea-try-go/service/follow/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBlockListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockListLogic {
	return &GetBlockListLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetBlockListLogic) GetBlockList(in *pb.ListReq) (*pb.UserListResp, error) {
	ids, err := l.svcCtx.FollowModel.GetBlockList(l.ctx, in.UserId, in.Offset, in.Limit)
	if err != nil {
		l.Logger.Errorf("GetBlockList db err: %v", err)
		return &pb.UserListResp{Code: 500, Msg: "DB Error"}, err
	}
	return &pb.UserListResp{Code: 0, Msg: "success", UserIds: ids}, nil
}
