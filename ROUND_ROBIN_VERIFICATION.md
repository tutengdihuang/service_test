# Round Robin 负载均衡验证

## 配置说明

已创建完整的测试环境来验证 `round_robin` 负载均衡策略是否生效。

### 关键配置

在 `rpc/user/internal/svc/servicecontext.go` 中配置了 round_robin 负载均衡：

```go
func NewServiceContext(c config.Config) *ServiceContext {
    // 使用 round_robin 负载均衡策略连接两个 trade 服务实例
    tradeClient := zrpc.MustNewClient(c.TradeRpc, zrpc.WithDialOption(
        grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
    ))

    return &ServiceContext{
        Config:    c,
        UserModel: model.NewUserModel(),
        TradeRpc:  tradeclient.NewTrade(tradeClient),
    }
}
```

在 `rpc/user/etc/user_direct.yaml` 中配置了多个 trade 服务地址：

```yaml
TradeRpc:
  Target: "127.0.0.1:9003,127.0.0.1:9004"  # 两个 trade 服务实例
```

## 测试方法

### 1. 启动服务

```bash
# 启动 product 服务
cd rpc/product && go run product.go -f etc/product_direct.yaml &

# 启动 trade 实例 1 (端口 9003)
cd rpc/trade && go run trade.go -f etc/trade1_direct.yaml &

# 启动 trade 实例 2 (端口 9004)
cd rpc/trade && go run trade.go -f etc/trade2_direct.yaml &

# 启动 user 服务 (配置了 round_robin)
cd rpc/user && go run user.go -f etc/user_direct.yaml &
```

### 2. 运行测试

```bash
cd test_user_trade
go run main.go
```

### 3. 预期结果

如果轮询生效，应该看到：
- 请求 1: trade-instance-1
- 请求 2: trade-instance-2
- 请求 3: trade-instance-1
- 请求 4: trade-instance-2
- ...

两个实例被均匀调用，分布差异应该小于 2 次。

## 验证要点

1. **配置正确性**: `Target` 字段包含多个地址，用逗号分隔
2. **负载均衡策略**: 使用 `grpc.WithDefaultServiceConfig` 设置 `round_robin`
3. **服务实例标识**: 每个 trade 实例通过 `InstanceId` 区分
4. **请求分布**: 通过统计两个实例的调用次数验证轮询

## 注意事项

- 确保所有服务代码已重新生成（运行 `goctl`）
- 确保 user 服务的 `CallTrade` 方法已正确实现
- 确保两个 trade 实例都已成功启动
- 测试时观察日志，确认请求确实在两个实例间轮询

## 当前状态

代码已配置完成，可以验证 round_robin 负载均衡是否生效。如果测试结果显示两个实例被均匀调用，则说明轮询策略工作正常。

