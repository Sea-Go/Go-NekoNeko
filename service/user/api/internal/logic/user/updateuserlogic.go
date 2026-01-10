// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"encoding/json"
	"errors"

	"sea-try-go/service/user/api/internal/model"
	"sea-try-go/service/user/api/internal/svc"
	"sea-try-go/service/user/api/internal/types"
	"sea-try-go/service/user/common/cryptx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateuserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateuserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateuserLogic {
	return &UpdateuserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateuserLogic) Updateuser(req *types.UpdateUserReq) (resp *types.UpdateUserResp, err error) {
	userId := l.ctx.Value("userId").(json.Number)
	id, er := userId.Int64()
	if er != nil {
		return nil, er
	}
	updates := make(map[string]interface{})
	if len(req.Username) > 0 {
		updates["username"] = req.Username
	}

	if len(req.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(req.Password)
		if e != nil {
			return nil, e
		}
		updates["password"] = newPassword
	}

	if len(req.Email) > 0 {
		updates["email"] = req.Email
	}

	if req.Extrainfo != nil {
		updates["extra_info"] = req.Extrainfo
	}

	if len(updates) > 0 {
		err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
		if err != nil {
			return nil, errors.New("更新失败:" + err.Error())
		}
	}

	var newUser model.User
	err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", id).First(&newUser).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &types.UpdateUserResp{
		User: types.UserInfo{
			Id:        newUser.Id,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Extrainfo: newUser.ExtraInfo,
		},
	}, nil

}
