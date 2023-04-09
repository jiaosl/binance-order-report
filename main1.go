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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	client := futures.NewClient(apiKey, secretKey)

	// 设置查询条件
	symbol := "BTCUSDT"                                         // 币种
	limit := 100                                                // 查询记录数
	startTime := time.Now().Add(-30*24*time.Hour).Unix() * 1000 // 查询起始时间（30天前）
	endTime := time.Now().Unix() * 1000                         // 查询结束时间（当前时间）

	// 获取永续合约已完成订单
	orders, err := client.NewListOrdersService().
		Symbol(symbol).
		Limit(limit).
		StartTime(startTime).
		EndTime(endTime).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// 打印订单信息
	for _, order := range orders {
		fmt.Println(order)
	}
}
