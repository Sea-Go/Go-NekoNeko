package logic

import (
	"context"
	"fmt"
	"sea-try-go/service/comment/rpc/internal/svc"
	"sea-try-go/service/comment/rpc/pb"

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
	if in.CommentId == 0 {
		return nil, fmt.Errorf("参数错误:评论ID不能为空")
	}
	var bitOffset uint
	var isSet bool
	switch in.ActionType {
	case pb.ManageType_MANAGE_PIN:
		bitOffset = 2
		isSet = true
	case pb.ManageType_MANAGE_UNPIN:
		bitOffset = 2
		isSet = false
	case pb.ManageType_MANAGE_FEATURE:
		bitOffset = 3
		isSet = true
	case pb.ManageType_MANAGE_UNFEATURE:
		bitOffset = 3
		isSet = false
	}
	err := l.svcCtx.CommentModel.ManageCommentAttribute(l.ctx, in.CommentId, bitOffset, isSet)
	if err != nil {
		return nil, err
	}
	return &pb.ManageCommentResp{
		Success: true,
	}, nil
}
