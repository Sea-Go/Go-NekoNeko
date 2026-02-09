package logic

import (
	"context"

	"sea-try-go/service/user/rpc/internal/svc"
	"sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.GetUserResp, error) {
	user, err := l.svcCtx.UserModel.GetByUid(in.Uid)
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.GetUserResp{Found: false}, err
	}

	if user == nil {
		return &pb.GetUserResp{Found: false}, nil
	}

	return &pb.GetUserResp{
		User: &pb.UserInfo{
			Uid:       user.Uid,
			Score:     uint32(user.Score),
			Username:  user.Username,
			Email:     user.Email,
			ExtraInfo: user.ExtraInfo,
		},
		Found: true,
	}, nil
}
