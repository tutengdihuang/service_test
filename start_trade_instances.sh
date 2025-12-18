#!/bin/bash

# 启动两个 trade 服务实例用于测试负载均衡

echo "启动两个 trade 服务实例..."

# 检查 etcd 是否运行
if ! docker ps | grep -q etcd; then
    echo "启动 etcd 服务..."
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
    sleep 2
fi

# 检查 user 和 product 服务是否运行
echo "请确保 user 和 product 服务已启动..."

# 启动 trade 实例 1 (端口 9003)
echo "启动 trade 实例 1 (端口 9003)..."
cd rpc/trade
go run trade.go -f etc/trade1.yaml &
TRADE1_PID=$!
cd ../..

sleep 1

# 启动 trade 实例 2 (端口 9004)
echo "启动 trade 实例 2 (端口 9004)..."
cd rpc/trade
go run trade.go -f etc/trade2.yaml &
TRADE2_PID=$!
cd ../..

sleep 2

echo "两个 trade 服务实例已启动！"
echo "Trade 实例 1 PID: $TRADE1_PID (端口 9003)"
echo "Trade 实例 2 PID: $TRADE2_PID (端口 9004)"
echo ""
echo "现在可以运行测试客户端验证负载均衡:"
echo "cd test_loadbalance && go run main.go"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待中断信号
trap "kill $TRADE1_PID $TRADE2_PID 2>/dev/null; exit" INT TERM
wait

