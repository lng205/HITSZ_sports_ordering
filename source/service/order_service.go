package service

import (
	"fmt"
	"sync"

	"sports_order/common"
)

// maxConcurrentOrders 用于限制同时处理订单的并发数。
const maxConcurrentOrders = 10

// OrderProcessor 负责从数据库取单、并发执行预约、并落订单状态与日志。
type OrderProcessor struct {
	repo           common.Repository
	bookingService *BookingService
}

// NewOrderProcessor 创建订单处理服务。
func NewOrderProcessor(
	apiClient common.APIClient,
	repo common.Repository,
	user *common.User,
) *OrderProcessor {
	booking := NewBookingService(apiClient, repo, user)
	return &OrderProcessor{
		repo:           repo,
		bookingService: booking,
	}
}

// ProcessOrdersForDate 并发处理某一天所有待预约订单。
func (s *OrderProcessor) ProcessOrdersForDate(targetDate string) error {
	// 查询指定日期下待处理订单
	orders, err := s.repo.FindOrdersByDate(targetDate)
	if err != nil {
		return fmt.Errorf("查询订单失败: %v", err)
	}

	if len(orders) == 0 {
		s.repo.CreateLogf(common.LogLevelInfo, nil, "无订单: %s", targetDate)
		return nil
	}

	s.repo.CreateLogf(common.LogLevelInfo, nil, "找到 %d 个订单，日期: %s", len(orders), targetDate)
	// 获取表单配置（版本号、日期/时段/场地映射等）。
	catalogData, err := s.bookingService.GetCatalogData()
	if err != nil {
		s.repo.CreateLogf(common.LogLevelError, nil, "获取预约元数据失败: %v", err)
		return err
	}

	// 使用信号量限制并发
	semaphore := make(chan struct{}, maxConcurrentOrders)
	var wg sync.WaitGroup

	// 并发处理每一条订单
	for _, order := range orders {
		wg.Add(1)
		go func(order *common.Order, data *common.CatalogData) {
			defer wg.Done()
			// 获取并发令牌
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // 释放并发令牌

			s.processSingleOrder(order, data)
		}(order, catalogData)
	}

	// 等待全部订单处理完成
	wg.Wait()

	return nil
}

// processSingleOrder 尝试处理单条订单：失败落 FAILED，成功落 SUCCESS。
func (s *OrderProcessor) processSingleOrder(order *common.Order, data *common.CatalogData) {
	orderID := int(order.ID)
	s.repo.CreateLogf(common.LogLevelInfo, &orderID, "开始预约订单 %d: %s %d:00-%d:00 场地 %d",
		order.ID, order.Date, order.Hour, order.Hour+1, order.Venue)

	// 执行预约
	err := s.bookingService.BookTimeSlot(data, common.BookingSlot{Date: order.Date, Hour: order.Hour, Venue: order.Venue})
	if err != nil {
		s.repo.UpdateOrderStatus(order.ID, common.OrderStatusFailed)
		s.repo.CreateLogf(common.LogLevelError, &orderID, "订单 %d 失败: %v", order.ID, err)
	} else {
		// 预约成功
		s.repo.UpdateOrderStatus(order.ID, common.OrderStatusSuccess)
		s.repo.CreateLogf(common.LogLevelInfo, &orderID, "订单 %d 预约成功", order.ID)
	}
}
