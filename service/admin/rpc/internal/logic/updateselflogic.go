package logic

import (
	"context"
	"errors"

	"sea-try-go/service/admin/rpc/internal/model"
	"sea-try-go/service/admin/rpc/internal/svc"
	"sea-try-go/service/admin/rpc/pb"
	"sea-try-go/service/common/cryptx"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSelfLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSelfLogic {
	return &UpdateSelfLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSelfLogic) UpdateSelf(in *pb.UpdateSelfReq) (*pb.UpdateSelfResp, error) {

	updates := make(map[string]interface{})
	if len(in.Username) > 0 {
		updates["username"] = in.Username
	}

	if len(in.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(in.Password)
		if e != nil {
			return nil, e
		}
		updates["password"] = newPassword
	}

	if len(in.Email) > 0 {
		updates["email"] = in.Email
	}
	if in.ExtraInfo != nil {
		updates["extra_info"] = in.ExtraInfo
	}

	if len(updates) > 0 {
		err := l.svcCtx.DB.Model(&model.Admin{}).Where("id = ?", in.Id).Updates(updates).Error
		if err != nil {
			return nil, errors.New("更新失败:" + err.Error())
		}
	}
	var newAdmin model.Admin
	err := l.svcCtx.DB.Model(&model.Admin{}).Where("id = ?", in.Id).First(&newAdmin).Error
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSelfResp{
		Success: true,
		Admin: &pb.AdminInfo{
			Id:        newAdmin.Id,
			Username:  newAdmin.Username,
			Email:     newAdmin.Email,
			ExtraInfo: newAdmin.ExtraInfo,
		},
	}, nil
}
