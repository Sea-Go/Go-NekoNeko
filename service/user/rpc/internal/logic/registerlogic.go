package logic

import (
	"context"
	"errors"
	"sea-try-go/service/common/cryptx"
	"sea-try-go/service/common/errmsg"
	"sea-try-go/service/common/snowflake"
	"sea-try-go/service/user/rpc/internal/metrics"
	"sea-try-go/service/user/rpc/internal/model"
	"sea-try-go/service/user/rpc/internal/svc"
	pb "sea-try-go/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.CreateUserReq) (*pb.CreateUserResp, error) {

	_, err := l.svcCtx.UserModel.FindOneByUserName(l.ctx, in.Username)
	if err == nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorUserExist))
	}
	if err != model.ErrorNotFound {
		return nil, err
	}
	truePassword, er := cryptx.PasswordEncrypt(in.Password)
	if er != nil {
		return nil, er
	}

	var uid int64
	uid, err = snowflake.GetID()
	if err != nil {
		l.Logger.Errorf("生成UID失败:%v", err)
		return nil, errors.New("ID生成失败")
	}
	newUser := model.User{
		Uid:       uid,
		Username:  in.Username,
		Password:  truePassword,
		Email:     in.Email,
		Score:     0,
		Status:    0,
		ExtraInfo: in.ExtraInfo,
	}
	err = l.svcCtx.UserModel.Insert(l.ctx, &newUser)
	if err != nil {
		return nil, errors.New(errmsg.GetErrMsg(errmsg.ErrorDbUpdate))
	}
	metrics.RegisterSuccessTotal.Inc()
	return &pb.CreateUserResp{
		Uid: newUser.Uid,
	}, nil
}
