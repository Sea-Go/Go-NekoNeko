// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"sea-try-go/service/user/common/errmsg"
	"sea-try-go/service/user/user/api/internal/svc"
	"sea-try-go/service/user/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, code int) {

	//不是很懂logout具体逻辑怎么写,先这样子吧
	return &types.LogoutResp{
		Success: true,
	}, errmsg.Success
}
