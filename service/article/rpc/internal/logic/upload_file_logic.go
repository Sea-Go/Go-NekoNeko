package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"mime" // 引入 mime 包
	"path/filepath"

	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadFileLogic) UploadFile(in *__.UploadFileRequest) (*__.UploadFileResponse, error) {
	ext := filepath.Ext(in.FileName)
	objectName := fmt.Sprintf("%s%s%s", l.svcCtx.Config.MinIO.ImagePath, uuid.New().String(), ext)

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := l.svcCtx.MinioClient.PutObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, objectName,
		bytes.NewReader(in.Content), int64(len(in.Content)),
		minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		return nil, err
	}

	fileUrl := fmt.Sprintf("http://%s/%s/%s", l.svcCtx.Config.MinIO.Endpoint, l.svcCtx.Config.MinIO.BucketName, objectName)

	return &__.UploadFileResponse{
		FileUrl: fileUrl,
	}, nil
}
