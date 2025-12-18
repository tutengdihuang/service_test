package logic

import (
	"context"
	"service_test/rpc/trade/internal/svc"
	"service_test/rpc/trade/trade"
	"service_test/rpc/user/userclient"
	"service_test/rpc/product/productclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateTradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTradeLogic {
	return &CreateTradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTradeLogic) CreateTrade(in *trade.CreateTradeRequest) (*trade.CreateTradeResponse, error) {
	// 调用 user 服务获取用户信息
	userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: in.UserId,
	})
	if err != nil || userResp.UserId == 0 {
		return &trade.CreateTradeResponse{
			Status: "failed: user not found",
		}, nil
	}

	// 调用 product 服务获取产品信息
	productResp, err := l.svcCtx.ProductRpc.GetProductInfo(l.ctx, &productclient.ProductInfoRequest{
		ProductId: in.ProductId,
	})
	if err != nil || productResp.ProductId == 0 {
		return &trade.CreateTradeResponse{
			Status: "failed: product not found",
		}, nil
	}

	// 检查库存
	if productResp.Stock < in.Quantity {
		return &trade.CreateTradeResponse{
			Status: "failed: insufficient stock",
		}, nil
	}

	// 计算总金额 (productResp.Price 是 float32)
	totalAmount := float32(productResp.Price * float32(in.Quantity))

	// 通过 RPC 调用减少库存
	reduceResp, err := l.svcCtx.ProductRpc.ReduceStock(l.ctx, &productclient.ReduceStockRequest{
		ProductId: in.ProductId,
		Quantity:  in.Quantity,
	})
	if err != nil || !reduceResp.Success {
		return &trade.CreateTradeResponse{
			Status: "failed: stock update failed",
		}, nil
	}

	// 创建交易记录
	t := l.svcCtx.TradeModel.CreateTrade(in.UserId, in.ProductId, in.Quantity, float64(totalAmount))

	return &trade.CreateTradeResponse{
		TradeId:     t.TradeId,
		UserId:      t.UserId,
		ProductId:   t.ProductId,
		Quantity:    t.Quantity,
		TotalAmount: totalAmount,
		Status:      t.Status,
		InstanceId:  l.svcCtx.Config.InstanceId, // 返回实例ID
	}, nil
}

