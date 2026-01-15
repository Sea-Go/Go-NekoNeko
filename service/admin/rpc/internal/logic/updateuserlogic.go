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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	user := model.User{}
	err := l.svcCtx.DB.Where("id = ?", in.Id).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	updates := make(map[string]interface{})
	if len(in.Username) > 0 {
		updates["username"] = in.Username
	}
	if len(in.Password) > 0 {
		newPassword, e := cryptx.PasswordEncrypt(in.Password)
		if e != nil {
			return nil, errors.New("密码加密失败")
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
		err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.Id).Updates(updates).Error
		if err != nil {
			return nil, errors.New("更新失败" + err.Error())
		}
	}
	var newUser model.User
	err = l.svcCtx.DB.Model(&model.User{}).Where("id = ?", in.Id).First(&newUser).Error
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResp{
		User: &pb.UserInfo{
			Id:        newUser.Id,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Status:    uint64(newUser.Status),
			ExtraInfo: newUser.ExtraInfo,
		},
	}, nil
}
