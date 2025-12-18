package main

import (
	"context"
	"fmt"
	"service_test/rpc/trade/tradeclient"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	// 使用 round_robin 负载均衡策略连接两个 trade 服务实例
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: "127.0.0.1:9003,127.0.0.1:9004", // 两个服务实例地址
	}, zrpc.WithGRPCOptions(
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 关键配置
	))

	tradeClient := tradeclient.NewTrade(client.Conn())

	fmt.Println("==========================================")
	fmt.Println("验证 round_robin 轮询是否生效")
	fmt.Println("==========================================")
	fmt.Println("")

	// 发送 20 个请求，观察轮询效果
	instance1Count := 0
	instance2Count := 0

	for i := 1; i <= 20; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := tradeClient.GetInstanceInfo(ctx, &tradeclient.GetInstanceInfoRequest{})
		cancel()

		if err != nil {
			fmt.Printf("请求 %2d: ❌ 失败 - %v\n", i, err)
			continue
		}

		if resp.InstanceId == "trade-instance-1" {
			instance1Count++
			fmt.Printf("请求 %2d: ✅ trade-instance-1 (端口 9003)\n", i)
		} else if resp.InstanceId == "trade-instance-2" {
			instance2Count++
			fmt.Printf("请求 %2d: ✅ trade-instance-2 (端口 9004)\n", i)
		} else {
			fmt.Printf("请求 %2d: ⚠️  未知实例: %s\n", i, resp.InstanceId)
		}

		time.Sleep(100 * time.Millisecond) // 稍微延迟，便于观察
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
			fmt.Println("✅ 轮询生效！两个实例被均匀调用")
			fmt.Println("✅ round_robin 负载均衡策略工作正常")
		} else {
			fmt.Printf("⚠️  轮询不均匀，差异: %d 次\n", diff)
		}
	} else {
		fmt.Println("❌ 轮询未生效！所有请求都路由到同一个实例")
		fmt.Println("   请检查负载均衡配置")
	}
	fmt.Println("==========================================")
}

