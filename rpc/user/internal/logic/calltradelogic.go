package logic

import (
	"context"
	"service_test/rpc/user/internal/svc"
	"service_test/rpc/user/user"
	"service_test/rpc/trade/tradeclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CallTradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCallTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallTradeLogic {
	return &CallTradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CallTradeLogic) CallTrade(in *user.CallTradeRequest) (*user.CallTradeResponse, error) {
	// 调用 trade 服务获取实例信息（用于验证负载均衡）
	tradeResp, err := l.svcCtx.TradeRpc.GetInstanceInfo(l.ctx, &tradeclient.GetInstanceInfoRequest{})
	if err != nil {
		return &user.CallTradeResponse{
			TradeInstanceId: "",
			Status:          "failed: " + err.Error(),
		}, nil
	}

	return &user.CallTradeResponse{
		TradeInstanceId: tradeResp.InstanceId,
		Status:          "success",
	}, nil
}

