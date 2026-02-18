// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"sea-try-go/service/common/logger"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"
	"sea-try-go/service/user/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetuserLogic {
	return &GetuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetuserLogic) Getuser(req *types.GetUserReq) (resp *types.GetUserResp, code int) {

	userId, ok := l.ctx.Value("userId").(json.Number)
	if !ok {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorTokenRuntime, fmt.Errorf("ctx userId is not json.Number"))
		return nil, errmsg.ErrorTokenRuntime
	}
	uid, err := userId.Int64()

	if err != nil {
		logger.LogBusinessErr(l.ctx, errmsg.ErrorTokenRuntime, fmt.Errorf("parse userId to int64 failed: %v", err))
		return nil, errmsg.ErrorTokenRuntime
	}

	rpcReq := &pb.GetUserReq{
		Uid: uid,
	}

	rpcResp, err := l.svcCtx.UserRpc.GetUser(l.ctx, rpcReq)

	if err != nil || !rpcResp.Found {
		logger.LogBusinessErr(l.ctx, errmsg.Error, err)
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			return nil, errmsg.ErrorUserNotExist
		case codes.Internal:
			return nil, errmsg.ErrorServerCommon
		default:
			return nil, errmsg.CodeServerBusy
		}
	}

	return &types.GetUserResp{
		User: types.UserInfo{
			Uid:       rpcResp.User.Uid,
			Username:  rpcResp.User.Username,
			Email:     rpcResp.User.Email,
			Extrainfo: rpcResp.User.ExtraInfo,
		},
		Found: true,
	}, errmsg.Success
}
