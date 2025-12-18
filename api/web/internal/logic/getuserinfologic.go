package logic

import (
	"context"
	"service_test/api/web/internal/svc"
	"service_test/api/web/internal/types"
	"service_test/rpc/user/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoRequest) (resp *types.UserInfoResponse, err error) {
	userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		UserId:    userResp.UserId,
		Username:  userResp.Username,
		Email:     userResp.Email,
		CreatedAt: userResp.CreatedAt,
	}, nil
}

