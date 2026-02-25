package logic

import (
	"context"
	"fmt"

	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubjectStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSubjectStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubjectStateLogic {
	return &UpdateSubjectStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSubjectStateLogic) UpdateSubjectState(in *pb.UpdateSubjectStateReq) (*pb.UpdateSubjectStateResp, error) {
	if in.TargetType == "" || in.TargetId == "" {
		return nil, fmt.Errorf("Type和Id不能为空")
	}
	err := l.svcCtx.CommentModel.UpdateSubjectState(l.ctx, in.TargetType, in.TargetId, int32(in.State))
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSubjectStateResp{
		Success: true,
	}, nil
}
