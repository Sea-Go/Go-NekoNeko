package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/metric"

	"sea-try-go/service/common/cryptx"
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

var metricUserCreated = metric.NewCounterVec(&metric.CounterVecOpts{
	Namespace: "user_rpc",
	Subsystem: "logic",
	Name:      "user_created_total",
	Help:      "Total number of users created",
	Labels:    []string{"result"},
})

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	var user model.User
	err := l.svcCtx.DB.Where("username = ?", in.Username).First(&user).Error
	if err == nil {
		return nil, errors.New("用户名已存在")
	}
	truePassword, er := cryptx.PasswordEncrypt(in.Password)
	if er != nil {
		return nil, er
	}
	newUser := model.User{
		Username:  in.Username,
		Password:  truePassword,
		Email:     in.Email,
		Score:     0,
		Status:    0,
		ExtraInfo: in.ExtraInfo,
	}
	err = l.svcCtx.DB.Create(&newUser).Error
	if err != nil {
		return nil, err
	}
	metricUserCreated.Inc("success")
	return &pb.CreateUserResp{
		Id: newUser.Id,
	}, nil
}
