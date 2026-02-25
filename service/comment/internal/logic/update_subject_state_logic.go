package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &pb.UpdateSubjectStateResp{}, nil
}
