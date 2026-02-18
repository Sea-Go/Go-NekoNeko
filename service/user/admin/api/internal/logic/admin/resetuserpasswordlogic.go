// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/admin/api/internal/svc"
	"sea-try-go/service/user/admin/api/internal/types"
	"sea-try-go/service/user/admin/rpc/pb"
	"sea-try-go/service/user/common/errmsg"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResetuserpasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetuserpasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetuserpasswordLogic {
	return &ResetuserpasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetuserpasswordLogic) Resetuserpassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordResp, code int) {
	rpcReq := &pb.ResetUserPasswordReq{
		Uid: req.Uid,
	}
	rpcResp, err := l.svcCtx.AdminRpc.ResetUserPassword(l.ctx, rpcReq)
	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.Error, err)
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, errmsg.ErrorUserNotExist
		case codes.Internal:
			return nil, errmsg.ErrorDbSelect
		default:
			return nil, errmsg.CodeServerBusy
		}
	}
	return &types.ResetUserPasswordResp{
		Success: rpcResp.Success,
	}, errmsg.Success
}
