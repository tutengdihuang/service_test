package logic

import (
	"context"
	"service_test/api/web/internal/svc"
	"service_test/api/web/internal/types"
	"service_test/rpc/trade/tradeclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTradeLogic {
	return &CreateTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTradeLogic) CreateTrade(req *types.CreateTradeRequest) (resp *types.CreateTradeResponse, err error) {
	tradeResp, err := l.svcCtx.TradeRpc.CreateTrade(l.ctx, &tradeclient.CreateTradeRequest{
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateTradeResponse{
		TradeId:     tradeResp.TradeId,
		UserId:      tradeResp.UserId,
		ProductId:   tradeResp.ProductId,
		Quantity:    tradeResp.Quantity,
		TotalAmount: float64(tradeResp.TotalAmount),
		Status:      tradeResp.Status,
	}, nil
}
