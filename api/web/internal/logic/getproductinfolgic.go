package logic

import (
	"context"
	"service_test/api/web/internal/svc"
	"service_test/api/web/internal/types"
	"service_test/rpc/product/productclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductInfoLogic {
	return &GetProductInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductInfoLogic) GetProductInfo(req *types.GetProductInfoRequest) (resp *types.ProductInfoResponse, err error) {
	productResp, err := l.svcCtx.ProductRpc.GetProductInfo(l.ctx, &productclient.ProductInfoRequest{
		ProductId: req.ProductId,
	})
	if err != nil {
		return nil, err
	}

	return &types.ProductInfoResponse{
		ProductId:   productResp.ProductId,
		Name:        productResp.Name,
		Description: productResp.Description,
		Price:       float64(productResp.Price),
		Stock:       productResp.Stock,
	}, nil
}
