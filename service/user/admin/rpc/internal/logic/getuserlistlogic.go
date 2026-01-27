package logic

import (
	"context"
	"errors"
	"sea-try-go/service/user/admin/rpc/internal/svc"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *pb.GetUserListReq) (*pb.GetUserListResp, error) {

	users, total, err := l.svcCtx.AdminModel.FindUserListByKeyword(l.ctx, in.Page, in.PageSize, in.Keyword)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbSelect))
	}

	list := make([]*pb.UserInfo, 0)

	for _, user := range users {
		list = append(list, &pb.UserInfo{
			Uid:       user.Uid,
			Username:  user.Username,
			Email:     user.Email,
			Status:    uint64(user.Status),
			ExtraInfo: user.ExtraInfo,
		})
	}
	return &pb.GetUserListResp{
		List:  list,
		Total: total,
	}, nil
}
