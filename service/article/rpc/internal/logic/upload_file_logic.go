package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime"
	"path/filepath"

	"sea-try-go/service/article/common/errmsg"
	"sea-try-go/service/article/rpc/internal/svc"
	"sea-try-go/service/article/rpc/pb"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/common/snowflake"

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
	id, err := snowflake.GetID()
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.Error, fmt.Errorf("generate snowflake id failed: %w", err)) // 雪花ID生成失败，暂时用通用错误
		return nil, err
	}

	ext := filepath.Ext(in.FileName)
	objectName := fmt.Sprintf("%s%d%s", l.svcCtx.Config.MinIO.ImagePath, id, ext)

	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = l.svcCtx.MinioClient.PutObject(l.ctx, l.svcCtx.Config.MinIO.BucketName, objectName,
		bytes.NewReader(in.Content), int64(len(in.Content)),
		minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorMinioUpload, fmt.Errorf("minio put object failed: %w", err))
		return nil, err
	}

	fileUrl := fmt.Sprintf("http://%s/%s/%s", l.svcCtx.Config.MinIO.Endpoint, l.svcCtx.Config.MinIO.BucketName, objectName)

	return &__.UploadFileResponse{
		FileUrl: fileUrl,
	}, nil
}
