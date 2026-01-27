package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/rpc/internal/model"
	"sea-try-go/service/user/user/rpc/internal/svc"
	"sea-try-go/service/user/user/rpc/pb"

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

	user, err := l.svcCtx.UserModel.FindOneByUid(l.ctx, in.Uid)
	if err == model.ErrorNotFound {
		return &pb.GetUserResp{
			Found: false,
		}, nil
	}
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbSelect))
	}

	return &pb.GetUserResp{
		User: &pb.UserInfo{
			Uid:       user.Uid,
			Username:  user.Username,
			Email:     user.Email,
			ExtraInfo: user.ExtraInfo,
		},
		Found: true,
	}, nil
}
