package main

import (
	"log"
	"time"

	"sports_order/common"
	"sports_order/service"
)

// main 使用依赖注入组装各层依赖，并处理目标日期的待预约订单
func main() {
	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	db, err := InitDB(config)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer CloseDB(db)

	// 初始化仓储层
	repo := NewRepository(db)

	// 记录启动日志
	repo.CreateLog(common.LogLevelInfo, "应用启动", nil)

	// 初始化服务层
	apiClient := NewHTTPClient()
	orderProcessor := service.NewOrderProcessor(apiClient, repo, &config.User)

	// 计算目标日期
	targetDate := time.Now().AddDate(0, 0, common.DaysAhead).Format("2006-01-02")

	// 处理目标日期的订单
	if err := orderProcessor.ProcessOrdersForDate(targetDate); err != nil {
		repo.CreateLogf(common.LogLevelError, nil, "处理订单失败: %v", err)
		log.Fatalf("处理订单失败: %v", err)
	}

	// 记录完成日志
	repo.CreateLogf(common.LogLevelInfo, nil, "订单处理完成，目标日期: %s", targetDate)
}
