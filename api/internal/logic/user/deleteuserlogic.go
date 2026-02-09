// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"sea-try-go/api/internal/svc"
	"sea-try-go/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteuserLogic {
	return &DeleteuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteuserLogic) Deleteuser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	// todo: add your logic here and delete this line

	return
}
