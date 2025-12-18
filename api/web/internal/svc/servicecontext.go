package svc

import (
	"service_test/api/web/internal/config"
	"service_test/rpc/user/userclient"
	"service_test/rpc/product/productclient"
	"service_test/rpc/trade/tradeclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   userclient.User
	ProductRpc productclient.Product
	TradeRpc  tradeclient.Trade
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserRpc:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc).Conn()),
		ProductRpc: productclient.NewProduct(zrpc.MustNewClient(c.ProductRpc).Conn()),
		TradeRpc:  tradeclient.NewTrade(zrpc.MustNewClient(c.TradeRpc).Conn()),
	}
}

