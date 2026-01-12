// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"errors"

	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/user/api/internal/model"
	"sea-try-go/service/user/api/internal/svc"
	"sea-try-go/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

//register代码存在一定问题:
//用First和Create可能会在高并发情境下产生错误,出现两个人注册了同一个的情况

func (l *RegisterLogic) Register(req *types.CreateUserReq) (resp *types.CreateUserResp, err error) {
	user := model.User{}

	isExist := l.svcCtx.DB.Where("username = ?", req.Username).First(&user).Error
	if isExist == nil {
		return nil, errors.New("用户名已存在")
	}

	truePassword, e := cryptx.PasswordEncrypt(req.Password)
	if e != nil {
		return nil, e
	}

	user = model.User{
		Username:  req.Username,
		Password:  truePassword,
		Email:     req.Email,
		Status:    0,
		ExtraInfo: req.Extrainfo,
	}
	err = l.svcCtx.DB.Create(&user).Error

	if err != nil {
		return nil, err
	}

	return &types.CreateUserResp{
		Id: user.Id,
	}, nil
}
