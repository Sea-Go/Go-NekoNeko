package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kqtypes "sea-try-go/service/comment/common/types"
	"sea-try-go/service/comment/internal/svc"
	"sea-try-go/service/comment/pb"
	"sea-try-go/service/common/snowflake"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCommentLogic) CreateComment(in *pb.CreateCommentReq) (*pb.CreateCommentResp, error) {
	if in.TargetId == "" || in.TargetType == "" {
		return nil, fmt.Errorf("参数错误:类型与ID不能为空")
	}
	if in.Content == "" {
		return nil, fmt.Errorf("参数错误:内容不能为空")
	}
	//查询owner还有bug
	ownerId, err := l.svcCtx.CommentModel.GetOwnerId(l.ctx, in.TargetType, in.TargetId)
	if err != nil {
		return nil, fmt.Errorf("查询Owner失败")
	}
	commentId, err := snowflake.GetID()
	if err != nil {
		return nil, fmt.Errorf("雪花算法生成ID出错")
	}
	now := time.Now()
	msg := kqtypes.CommentKafkaMsg{
		CommentId:  commentId,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		UserId:     in.UserId,
		RootId:     in.RootId,
		ParentId:   in.ParentId,
		Content:    in.Content,
		OwnerId:    int64(ownerId),
		Meta:       in.Meta,
		Attribute:  0,
		CreateTime: now.Unix(),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Kafka消息序列化失败:%v", err)
	}
	err = l.svcCtx.KqPusherClient.Push(l.ctx, string(msgBytes))
	if err != nil {
		return nil, fmt.Errorf("Kafka推送消息失败:%v", err)
	}

	return &pb.CreateCommentResp{
		Id: commentId,
		//将时间格式化为"2006-01-02 15:04:05"
		CreatedAt: now.Format(time.DateTime),
		//Gemini说返回0,"客户端乐观更新"
		SubjectCount: 0,
	}, nil
}
