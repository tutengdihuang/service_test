package logic

import (
	"context"
	"service_test/rpc/product/internal/svc"
	"service_test/rpc/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReduceStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReduceStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReduceStockLogic {
	return &ReduceStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReduceStockLogic) ReduceStock(in *product.ReduceStockRequest) (*product.ReduceStockResponse, error) {
	success := l.svcCtx.ProductModel.ReduceStock(in.ProductId, in.Quantity)
	if success {
		return &product.ReduceStockResponse{
			Success: true,
			Message: "stock reduced successfully",
		}, nil
	}
	return &product.ReduceStockResponse{
		Success: false,
		Message: "insufficient stock or product not found",
	}, nil
}

