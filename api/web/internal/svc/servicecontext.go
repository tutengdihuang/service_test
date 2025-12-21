package svc

import (
	"service_test/api/web/internal/config"
	"service_test/rpc/product/productclient"
	"service_test/rpc/trade/tradeclient"
	"service_test/rpc/user/userclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    userclient.User
	ProductRpc productclient.Product
	TradeRpc   tradeclient.Trade
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc)),
		TradeRpc:   tradeclient.NewTrade(zrpc.MustNewClient(c.TradeRpc)),
	}
}
