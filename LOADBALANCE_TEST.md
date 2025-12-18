# 负载均衡轮询测试指南

本指南说明如何验证 go-zero 的 round_robin 负载均衡策略是否生效。

## 前置条件

1. 确保已安装 go-zero 工具链
2. 确保已生成所有服务代码（运行 `./gen.sh`）
3. 确保 user 和 product 服务已启动

## 测试步骤

### 1. 重新生成 trade 服务代码（如果已修改 proto）

由于我们在 `trade.proto` 中添加了 `GetInstanceInfo` 方法，需要重新生成代码：

```bash
cd rpc/trade
goctl rpc protoc trade.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd ../..
```

### 2. 启动两个 trade 服务实例

```bash
./start_trade_instances.sh
```

这将启动：
- **trade 实例 1**: 端口 9003, 实例ID: `trade-instance-1`
- **trade 实例 2**: 端口 9004, 实例ID: `trade-instance-2`

### 3. 运行负载均衡测试客户端

```bash
cd test_loadbalance
go mod tidy
go run main.go
```

### 4. 观察输出

如果负载均衡正常工作，你应该看到类似以下的输出：

```
开始测试负载均衡轮询...
发送 10 个请求，观察是否在两个实例间轮询

请求  1: 实例ID=trade-instance-1, 监听地址=0.0.0.0:9003
请求  2: 实例ID=trade-instance-2, 监听地址=0.0.0.0:9004
请求  3: 实例ID=trade-instance-1, 监听地址=0.0.0.0:9003
请求  4: 实例ID=trade-instance-2, 监听地址=0.0.0.0:9004
请求  5: 实例ID=trade-instance-1, 监听地址=0.0.0.0:9003
请求  6: 实例ID=trade-instance-2, 监听地址=0.0.0.0:9004
请求  7: 实例ID=trade-instance-1, 监听地址=0.0.0.0:9003
请求  8: 实例ID=trade-instance-2, 监听地址=0.0.0.0:9004
请求  9: 实例ID=trade-instance-1, 监听地址=0.0.0.0:9003
请求 10: 实例ID=trade-instance-2, 监听地址=0.0.0.0:9004

测试完成！
如果看到两个不同的实例ID交替出现，说明轮询负载均衡生效。
```

## 关键配置说明

测试客户端使用了以下关键配置：

```go
client := zrpc.MustNewClient(zrpc.RpcClientConf{
    Target: "127.0.0.1:9003,127.0.0.1:9004", // 多个服务地址用逗号分隔
}, zrpc.WithGRPCOptions(
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 轮询策略
))
```

### 配置要点：

1. **Target**: 多个服务地址用逗号分隔
2. **loadBalancingPolicy**: 设置为 `"round_robin"` 启用轮询策略
3. **服务实例标识**: 每个 trade 实例通过配置文件中的 `InstanceId` 字段区分

## 故障排查

### 问题1: 所有请求都路由到同一个实例

**可能原因**:
- gRPC 客户端没有正确配置负载均衡策略
- 服务地址配置错误

**解决方法**:
- 检查 `grpc.WithDefaultServiceConfig` 配置是否正确
- 确认两个服务实例都已成功启动

### 问题2: 连接失败

**可能原因**:
- 服务实例未启动
- 端口被占用
- 网络连接问题

**解决方法**:
- 检查服务是否正常运行: `ps aux | grep trade`
- 检查端口是否被占用: `lsof -i :9003` 和 `lsof -i :9004`
- 确认防火墙设置

### 问题3: 找不到 GetInstanceInfo 方法

**可能原因**:
- 未重新生成代码
- proto 文件未更新

**解决方法**:
- 重新运行 `goctl rpc protoc` 生成代码
- 检查 `trade.proto` 文件是否包含 `GetInstanceInfo` 方法定义

## 扩展测试

你可以修改测试客户端来测试其他场景：

1. **增加请求数量**: 修改循环次数，观察轮询是否持续
2. **测试并发请求**: 使用 goroutine 发送并发请求
3. **测试服务下线**: 停止一个实例，观察客户端如何处理

