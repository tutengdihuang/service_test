package main

import (
	"context"
	"fmt"
	"service_test/rpc/trade/tradeclient"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

func main() {
	// 使用 round_robin 负载均衡策略连接两个 trade 服务实例
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Endpoints: []string{"127.0.0.1:9003", "127.0.0.1:9004"}, // 两个服务实例地址
		NonBlock:  true,
		Timeout:   300000,
	}, zrpc.WithDialOption(
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 关键配置
	))

	tradeClient := tradeclient.NewTrade(client)

	fmt.Println("==========================================")
	fmt.Println("负载均衡轮询测试")
	fmt.Println("==========================================")
	fmt.Println("配置: round_robin 负载均衡策略")
	fmt.Println("目标: 127.0.0.1:9003,127.0.0.1:9004")
	fmt.Println("")

	const totalRequests = 100
	fmt.Printf("发送 %d 个请求，验证轮询分布...\n\n", totalRequests)

	// 统计信息
	var (
		instance1Count int
		instance2Count int
		successCount   int
		failCount      int
		mu             sync.Mutex
		wg             sync.WaitGroup
	)

	// 并发发送请求
	startTime := time.Now()
	for i := 1; i <= totalRequests; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			resp, err := tradeClient.GetInstanceInfo(ctx, &tradeclient.GetInstanceInfoRequest{})
			cancel()

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				failCount++
				if reqNum <= 10 {
					fmt.Printf("请求 %2d: ❌ 失败 - %v\n", reqNum, err)
				}
				return
			}

			successCount++
			if resp.InstanceId == "trade-instance-1" {
				instance1Count++
			} else if resp.InstanceId == "trade-instance-2" {
				instance2Count++
			}

			// 只打印前10个请求的详细信息
			if reqNum <= 10 {
				fmt.Printf("请求 %2d: ✅ 实例ID=%s, 监听地址=%s\n", reqNum, resp.InstanceId, resp.ListenOn)
			}
		}(i)

		// 控制并发数，避免过快
		if i%10 == 0 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	wg.Wait()
	duration := time.Since(startTime)

	// 打印统计结果
	fmt.Println("\n==========================================")
	fmt.Println("测试结果统计")
	fmt.Println("==========================================")
	fmt.Printf("总请求数:     %d\n", totalRequests)
	fmt.Printf("成功请求数:   %d\n", successCount)
	fmt.Printf("失败请求数:   %d\n", failCount)
	fmt.Printf("成功率:       %.2f%%\n", float64(successCount)/float64(totalRequests)*100)
	fmt.Printf("总耗时:       %v\n", duration)
	fmt.Println("")
	fmt.Println("实例分布:")
	fmt.Printf("  trade-instance-1: %d 次 (%.2f%%)\n", instance1Count, float64(instance1Count)/float64(successCount)*100)
	fmt.Printf("  trade-instance-2: %d 次 (%.2f%%)\n", instance2Count, float64(instance2Count)/float64(successCount)*100)
	fmt.Println("")

	// 验证轮询是否生效
	diff := instance1Count - instance2Count
	if diff < 0 {
		diff = -diff
	}
	balanceRatio := float64(diff) / float64(successCount) * 100

	fmt.Println("==========================================")
	fmt.Println("负载均衡验证")
	fmt.Println("==========================================")
	if successCount == totalRequests {
		fmt.Println("✅ 所有请求成功")
	} else {
		fmt.Printf("⚠️  有 %d 个请求失败\n", failCount)
	}

	if balanceRatio < 10 {
		fmt.Printf("✅ 负载均衡良好 (分布差异: %.2f%%)\n", balanceRatio)
		fmt.Println("✅ round_robin 轮询策略生效！")
	} else {
		fmt.Printf("⚠️  负载分布不均匀 (分布差异: %.2f%%)\n", balanceRatio)
		fmt.Println("⚠️  可能需要检查负载均衡配置")
	}
	fmt.Println("==========================================")
}
