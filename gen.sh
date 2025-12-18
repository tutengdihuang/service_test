#!/bin/bash

# Go-Zero 代码生成脚本

echo "开始生成 Go-Zero 代码..."

# 检查 goctl 是否安装
if ! command -v goctl &> /dev/null; then
    echo "错误: goctl 未安装，请先安装:"
    echo "go install github.com/zeromicro/go-zero/tools/goctl@latest"
    exit 1
fi

# 生成 user RPC 服务
echo "生成 user RPC 服务..."
cd rpc/user
goctl rpc protoc user.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd ../..

# 生成 product RPC 服务
echo "生成 product RPC 服务..."
cd rpc/product
goctl rpc protoc product.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd ../..

# 生成 trade RPC 服务
echo "生成 trade RPC 服务..."
cd rpc/trade
goctl rpc protoc trade.proto --go_out=. --go-grpc_out=. --zrpc_out=.
cd ../..

# 生成 web API 服务
echo "生成 web API 服务..."
cd api/web
goctl api go -api web.api -dir . -style gozero
cd ../..

echo "代码生成完成！"
echo ""
echo "接下来可以运行:"
echo "1. go mod tidy  # 安装依赖"
echo "2. ./start.sh   # 启动所有服务"

