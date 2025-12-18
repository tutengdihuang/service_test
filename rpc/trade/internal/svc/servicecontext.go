package svc

import (
	"service_test/rpc/trade/internal/config"
	"service_test/rpc/trade/internal/model"
	"service_test/rpc/user/userclient"
	"service_test/rpc/product/productclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    userclient.User
	ProductRpc productclient.Product
	TradeModel *model.TradeModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		TradeModel: model.NewTradeModel(),
	}
}

