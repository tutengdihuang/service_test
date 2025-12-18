#!/bin/bash

# Go-Zero 微服务启动脚本

echo "启动 Go-Zero 微服务..."

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

# 启动 user 服务
echo "启动 user 服务..."
cd rpc/user
go run user.go -f etc/user.yaml &
USER_PID=$!
cd ../..

sleep 1

# 启动 product 服务
echo "启动 product 服务..."
cd rpc/product
go run product.go -f etc/product.yaml &
PRODUCT_PID=$!
cd ../..

sleep 1

# 启动 trade 服务
echo "启动 trade 服务..."
cd rpc/trade
go run trade.go -f etc/trade.yaml &
TRADE_PID=$!
cd ../..

sleep 1

# 启动 web 服务
echo "启动 web 服务..."
cd api/web
go run web.go -f etc/web.yaml &
WEB_PID=$!
cd ../..

echo "所有服务已启动！"
echo "User PID: $USER_PID"
echo "Product PID: $PRODUCT_PID"
echo "Trade PID: $TRADE_PID"
echo "Web PID: $WEB_PID"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待中断信号
trap "kill $USER_PID $PRODUCT_PID $TRADE_PID $WEB_PID; exit" INT TERM
wait

