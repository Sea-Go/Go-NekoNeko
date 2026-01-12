// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"errors"

	"sea-try-go/service/admin/api/internal/model"
	"sea-try-go/service/admin/api/internal/svc"
	"sea-try-go/service/admin/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBanuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanuserLogic {
	return &BanuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BanuserLogic) Banuser(req *types.BanUserReq) (resp *types.BanUserResp, err error) {
	id := req.Id
	result := l.svcCtx.DB.Model(&model.User{}).Where("id = ?", id).Update("status", 1)
	if result.Error != nil {
		return nil, errors.New("封禁失败" + result.Error.Error())
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}
	return &types.BanUserResp{
		Success: true,
	}, nil
}
