# Go-Zero 微服务示例

这是一个使用 go-zero 框架构建的微服务示例项目，包含以下服务：

## 服务架构

1. **web** - HTTP API 网关服务，提供对外 HTTP 接口
2. **user** - 用户信息服务（RPC）
3. **product** - 产品信息服务（RPC）
4. **trade** - 交易服务（RPC），依赖 user 和 product 服务

## 项目结构

```
service_test/
├── api/web/          # HTTP API 服务
├── rpc/user/         # 用户 RPC 服务
├── rpc/product/      # 产品 RPC 服务
└── rpc/trade/        # 交易 RPC 服务
```

## 快速开始

### 1. 安装 go-zero 工具链

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 使用 goctl 生成代码

#### 方式一：使用脚本（推荐）

```bash
./gen.sh
```

#### 方式二：手动生成

```bash
# 生成 user 服务
cd rpc/user
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 生成 product 服务
cd ../product
goctl rpc protoc product.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 生成 trade 服务
cd ../trade
goctl rpc protoc trade.proto --go_out=. --go-grpc_out=. --zrpc_out=.

# 生成 web API 服务
cd ../../api/web
goctl api go -api web.api -dir . -style gozero
```

### 4. 启动 etcd（服务发现）

```bash
# 使用 Docker 启动 etcd
docker run -d --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  quay.io/coreos/etcd:v3.5.0 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
```

### 5. 启动服务（按顺序）

```bash
# 终端1: 启动 user 服务
cd rpc/user
go run user.go -f etc/user.yaml

# 终端2: 启动 product 服务
cd rpc/product
go run product.go -f etc/product.yaml

# 终端3: 启动 trade 服务
cd rpc/trade
go run trade.go -f etc/trade.yaml

# 终端4: 启动 web 服务
cd api/web
go run web.go -f etc/web.yaml
```

### 6. 测试接口

```bash
# 获取用户信息
curl http://localhost:8888/api/user/1

# 获取产品信息
curl http://localhost:8888/api/product/1

# 创建交易
curl -X POST http://localhost:8888/api/trade/create \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"product_id":1,"quantity":2}'
```

## 服务端口

- web: 8888
- user: 9001
- product: 9002
- trade: 9003

## 负载均衡测试

### 测试 round_robin 负载均衡策略

项目包含负载均衡测试功能，可以验证多个服务实例间的轮询是否生效。

#### 1. 启动两个 trade 服务实例

```bash
./start_trade_instances.sh
```

这将启动两个 trade 服务实例：
- 实例1: 端口 9003, 实例ID: trade-instance-1
- 实例2: 端口 9004, 实例ID: trade-instance-2

#### 2. 运行负载均衡测试客户端

```bash
cd test_loadbalance
go mod tidy
go run main.go
```

测试客户端会发送 10 个请求，如果看到两个不同的实例ID交替出现，说明轮询负载均衡生效。

#### 3. 测试代码说明

测试客户端使用以下配置连接两个服务实例：

```go
client := zrpc.MustNewClient(zrpc.RpcClientConf{
    Target: "127.0.0.1:9003,127.0.0.1:9004", // 两个服务实例地址
}, zrpc.WithGRPCOptions(
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
))
```

## 说明

- 本项目使用内存数据模拟，无需外部数据库
- 服务间通过 gRPC 进行通信
- 使用 etcd 进行服务发现
- web 服务作为 API 网关，统一对外提供 HTTP 接口
- 使用 go-zero 工具链自动生成大部分代码，只需编写业务逻辑
- 支持负载均衡测试，可验证 round_robin 策略


# Jenkins test Sat Dec 27 00:28:34 CST 2025
# Jenkins test Sat Dec 27 00:30:57 CST 2025
