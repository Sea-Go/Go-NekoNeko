package logic

import (
	"context"

	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageCommentLogic {
	return &ManageCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ManageCommentLogic) ManageComment(in *pb.ManageCommentReq) (*pb.ManageCommentResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ManageCommentResp{}, nil
}
