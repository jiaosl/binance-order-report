package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	// 从环境变量中读取 apiKey 和 secretKey
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	// 创建客户端
	client := futures.NewClient(apiKey, secretKey)

	// 定义起止日期
	startTimeStr := "2023-03-01"
	endTimeStr := "2023-04-09"
	startTime, _ := time.Parse("2006-01-02", startTimeStr)
	endTime, _ := time.Parse("2006-01-02", endTimeStr)

	// 计算时间间隔
	diff := endTime.Sub(startTime)
	days := int(diff.Hours() / 24)

	// 获取永续合约已完成订单
	var allOrders []*futures.Order
	var fromID int64 = 0
	for i := 0; i <= days/7; i++ {
		endTimeUnix := endTime.Unix() * 1000
		startTimeUnix := endTime.AddDate(0, 0, -7).Unix() * 1000
		if startTimeUnix < startTime.Unix()*1000 {
			startTimeUnix = startTime.Unix() * 1000
		}

		// 分页获取已完成订单
		for {
			listOrdersResponse, err := client.NewListOrdersService().
				Symbol("BTCUSDT").
				Limit(1000).
				StartTime(startTimeUnix).
				EndTime(endTimeUnix).
				OrderID(fromID).
				Do(context.Background())
			if err != nil {
				log.Fatalf("Error fetching orders: %s\n", err)
			}

			// 将已完成订单添加到数组中
			for _, order := range listOrdersResponse {
				if order.Status == futures.OrderStatusTypeFilled {
					allOrders = append(allOrders, order)
				}
			}

			// 如果订单已经全部获取完毕，则退出循环
			if len(listOrdersResponse) < 1000 {
				break
			}

			// 如果还有未获取的订单，则更新 fromID
			fromID = listOrdersResponse[len(listOrdersResponse)-1].OrderID + 1
		}

		// 更新 endTime
		endTime = endTime.AddDate(0, 0, -7)
	}

	// 输出已完成订单
	for _, order := range allOrders {
		fmt.Println(order)
	}
	fmt.Println("Total Orders:", len(allOrders))
}
