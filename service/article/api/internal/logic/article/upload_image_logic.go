// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"net/http"
	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	r      *http.Request
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext, r *http.Request) *UploadImageLogic {
	return &UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		r:      r,
	}
}

func (l *UploadImageLogic) UploadImage(req *types.UploadImageReq) (resp *types.UploadImageResp, err error) {
	return
}
