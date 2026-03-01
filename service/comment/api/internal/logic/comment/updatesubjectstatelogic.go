// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package comment

import (
	"context"

	"sea-try-go/service/comment/api/internal/svc"
	"sea-try-go/service/comment/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSubjectStateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSubjectStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSubjectStateLogic {
	return &UpdateSubjectStateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSubjectStateLogic) UpdateSubjectState(req *types.UpdateSubjectStateReq) (resp *types.UpdateSubjectStateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
