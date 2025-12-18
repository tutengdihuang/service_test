package svc

import (
	"service_test/rpc/product/internal/config"
	"service_test/rpc/product/internal/model"
)

type ServiceContext struct {
	Config      config.Config
	ProductModel *model.ProductModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ProductModel: model.NewProductModel(),
	}
}

