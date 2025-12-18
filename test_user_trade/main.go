package main

import (
	"context"
	"fmt"
	"service_test/rpc/user/userclient"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	// 连接 user 服务
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:9001",
	})

	userClient := userclient.NewUser(client)

	fmt.Println("==========================================")
	fmt.Println("验证 user -> trade 负载均衡轮询")
	fmt.Println("==========================================")
	fmt.Println("通过 user 服务调用 trade 服务")
	fmt.Println("trade 服务启动了两个实例 (9003, 9004)")
	fmt.Println("")

	// 统计信息
	instance1Count := 0
	instance2Count := 0

	// 发送 20 个请求
	for i := 1; i <= 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := userClient.CallTrade(ctx, &userclient.CallTradeRequest{
			UserId:    1,
			ProductId: 1,
			Quantity:  1,
		})
		cancel()

		if err != nil {
			fmt.Printf("请求 %2d: ❌ 失败 - %v\n", i, err)
			continue
		}

		if resp.TradeInstanceId == "trade-instance-1" {
			instance1Count++
			fmt.Printf("请求 %2d: ✅ trade-instance-1 (端口 9003)\n", i)
		} else if resp.TradeInstanceId == "trade-instance-2" {
			instance2Count++
			fmt.Printf("请求 %2d: ✅ trade-instance-2 (端口 9004)\n", i)
		} else {
			fmt.Printf("请求 %2d: ⚠️  未知实例: %s\n", i, resp.TradeInstanceId)
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n==========================================")
	fmt.Println("轮询统计结果")
	fmt.Println("==========================================")
	fmt.Printf("trade-instance-1: %d 次\n", instance1Count)
	fmt.Printf("trade-instance-2: %d 次\n", instance2Count)
	fmt.Println("")

	// 判断轮询是否生效
	if instance1Count > 0 && instance2Count > 0 {
		diff := instance1Count - instance2Count
		if diff < 0 {
			diff = -diff
		}
		if diff <= 2 {
			fmt.Println("✅ 轮询生效！两个 trade 实例被均匀调用")
			fmt.Println("✅ round_robin 负载均衡策略工作正常")
			fmt.Println("✅ user -> trade 服务间调用负载均衡验证成功")
		} else {
			fmt.Printf("⚠️  轮询不均匀，差异: %d 次\n", diff)
		}
	} else {
		fmt.Println("❌ 轮询未生效！所有请求都路由到同一个 trade 实例")
		fmt.Println("   请检查负载均衡配置")
	}
	fmt.Println("==========================================")
}

