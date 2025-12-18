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
	// todo: add your logic here and delete this line

	return &product.ProductInfoResponse{}, nil
}
