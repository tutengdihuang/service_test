package logic

import (
	"context"
	"service_test/rpc/product/internal/svc"
	"service_test/rpc/product/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductInfoLogic {
	return &GetProductInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductInfoLogic) GetProductInfo(in *product.ProductInfoRequest) (*product.ProductInfoResponse, error) {
	p, ok := l.svcCtx.ProductModel.GetProduct(in.ProductId)
	if !ok {
		return &product.ProductInfoResponse{
			ProductId:   0,
			Name:        "",
			Description: "",
			Price:       0,
			Stock:       0,
		}, nil
	}

	return &product.ProductInfoResponse{
		ProductId:   p.ProductId,
		Name:        p.Name,
		Description: p.Description,
		Price:       float32(p.Price),
		Stock:       p.Stock,
	}, nil
}
