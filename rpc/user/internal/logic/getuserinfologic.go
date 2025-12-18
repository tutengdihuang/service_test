package logic

import (
	"context"
	"service_test/rpc/user/internal/svc"
	"service_test/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	u, ok := l.svcCtx.UserModel.GetUser(in.UserId)
	if !ok {
		return &user.UserInfoResponse{
			UserId:    0,
			Username:  "",
			Email:     "",
			CreatedAt: 0,
		}, nil
	}

	return &user.UserInfoResponse{
		UserId:    u.UserId,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

