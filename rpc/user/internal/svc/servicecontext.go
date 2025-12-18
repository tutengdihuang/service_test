package svc

import (
	"service_test/rpc/trade/tradeclient"
	"service_test/rpc/user/internal/config"
	"service_test/rpc/user/internal/model"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *model.UserModel
	TradeRpc  tradeclient.Trade
}


func NewServiceContext(c config.Config) *ServiceContext {
	// 使用 round_robin 负载均衡策略连接两个 trade 服务实例
	// 配置文件中 Target 使用逗号分隔多个地址: "127.0.0.1:9003,127.0.0.1:9004"
	// 使用 grpc.WithDefaultServiceConfig 配置负载均衡策略
	tradeClient := zrpc.MustNewClient(c.TradeRpc, zrpc.WithDialOption(
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	))

	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(),
		TradeRpc:  tradeclient.NewTrade(tradeClient),
	}
}
