// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"sea-try-go/service/user/admin/api/internal/svc"
	"sea-try-go/service/user/admin/api/internal/types"
	"sea-try-go/service/user/admin/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetuserlistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetuserlistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetuserlistLogic {
	return &GetuserlistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetuserlistLogic) Getuserlist(req *types.GetUserListReq) (resp *types.GetUserListResp, err error) {
	rpcReq := &pb.GetUserListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Keyword:  req.Keyword,
	}

	rpcResp, err := l.svcCtx.AdminRpc.GetUserList(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	var list []types.UserInfo
	if rpcResp.List != nil {
		for _, v := range rpcResp.List {
			list = append(list, types.UserInfo{
				Uid:       v.Uid,
				Username:  v.Username,
				Email:     v.Email,
				Status:    int64(v.Status),
				Extrainfo: v.ExtraInfo,
			})
		}
	}
	return &types.GetUserListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}
