// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package article

import (
	"context"
	"io"
	"mime/multipart"
	"sea-try-go/service/article/api/internal/svc"
	"sea-try-go/service/article/api/internal/types"
	"sea-try-go/service/article/rpc/articleservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImageLogic) UploadImage(file multipart.File, header *multipart.FileHeader) (resp *types.UploadImageResp, err error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	rpcResp, err := l.svcCtx.ArticleRpc.UploadFile(l.ctx, &articleservice.UploadFileRequest{
		Content:  content,
		FileName: header.Filename,
	})
	if err != nil {
		return nil, err
	}

	return &types.UploadImageResp{
		ImageUrl: rpcResp.FileUrl,
	}, nil
}
