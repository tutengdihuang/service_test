package logic

import (
	"context"
	"service_test/rpc/trade/internal/svc"
	"service_test/rpc/trade/trade"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInstanceInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInstanceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInstanceInfoLogic {
	return &GetInstanceInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInstanceInfoLogic) GetInstanceInfo(in *trade.GetInstanceInfoRequest) (*trade.GetInstanceInfoResponse, error) {
	return &trade.GetInstanceInfoResponse{
		InstanceId: l.svcCtx.Config.InstanceId,
		ListenOn:   l.svcCtx.Config.ListenOn,
	}, nil
}

