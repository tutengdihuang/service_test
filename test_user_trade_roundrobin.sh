#!/bin/bash

# 测试 user -> trade 负载均衡轮询

set -e

echo "=========================================="
echo "测试 user -> trade 负载均衡轮询"
echo "=========================================="

cleanup() {
    echo ""
    echo "清理进程..."
    kill $USER_PID $PRODUCT_PID $TRADE1_PID $TRADE2_PID 2>/dev/null || true
    wait 2>/dev/null || true
}

trap cleanup EXIT INT TERM

# 启动 product 服务
echo "启动 product 服务 (端口 9002)..."
cd rpc/product
go run product.go -f etc/product_direct.yaml > /tmp/product.log 2>&1 &
PRODUCT_PID=$!
cd ../..
sleep 2

# 启动 trade 实例 1
echo "启动 trade 实例 1 (端口 9003)..."
cd rpc/trade
go run trade.go -f etc/trade1_direct.yaml > /tmp/trade1.log 2>&1 &
TRADE1_PID=$!
cd ../..
sleep 2

# 启动 trade 实例 2
echo "启动 trade 实例 2 (端口 9004)..."
cd rpc/trade
go run trade.go -f etc/trade2_direct.yaml > /tmp/trade2.log 2>&1 &
TRADE2_PID=$!
cd ../..
sleep 2

# 启动 user 服务（配置了 round_robin 负载均衡连接两个 trade 实例）
echo "启动 user 服务 (端口 9001)..."
cd rpc/user
go run user.go -f etc/user_direct.yaml > /tmp/user.log 2>&1 &
USER_PID=$!
cd ../..
sleep 3

echo ""
echo "所有服务已启动！"
echo "Product PID: $PRODUCT_PID (端口 9002)"
echo "Trade1 PID: $TRADE1_PID (端口 9003)"
echo "Trade2 PID: $TRADE2_PID (端口 9004)"
echo "User PID: $USER_PID (端口 9001)"
echo ""

# 运行测试
echo "=========================================="
echo "开始测试..."
echo "=========================================="
echo ""

cd test_user_trade
go mod tidy > /dev/null 2>&1
go run main.go

