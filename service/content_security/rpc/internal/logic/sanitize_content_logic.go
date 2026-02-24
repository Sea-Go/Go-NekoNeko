package logic

import (
	"context"

	"sea-try-go/service/content_security/rpc/internal/svc"
	"sea-try-go/service/content_security/rpc/pb/sea-try-go/service/content-security/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SanitizeContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSanitizeContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SanitizeContentLogic {
	return &SanitizeContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SanitizeContentLogic) SanitizeContent(in *pb.SanitizeContentRequest) (*pb.SanitizeContentResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.SanitizeContentResponse{}, nil
}
