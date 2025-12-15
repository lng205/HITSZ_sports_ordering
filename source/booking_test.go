package main

import (
	"testing"

	"sports_order/common"
	"sports_order/service"
)

// sortStrings 简单的字符串排序
func sortStrings(s []string) {
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// sortInts 简单的整数排序
func sortInts(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// TestReadCatalogAndBook 从真实 API 读取 Catalog 并尝试预定一个可用选项
func TestReadCatalogAndBook(t *testing.T) {
	// 从 config.yaml 读取配置（测试从 source/ 目录运行）
	config, err := LoadConfigFrom("../config.yaml")
	if err != nil {
		t.Fatalf("读取配置失败: %v", err)
	}

	t.Logf("用户: %s (学号: %s)", config.User.Name, config.User.StudentID)
	if config.User.Token == "" {
		t.Log("警告: Token 为空，预定请求将会失败")
	}

	httpClient := NewHTTPClient()
	bookingService := service.NewBookingService(httpClient, nil, &config.User)

	// Step 1: 读取 Catalog
	t.Log("\n=== Step 1: 从真实 API 读取 Catalog ===")
	catalogData, err := bookingService.GetCatalogData()
	if err != nil {
		t.Fatalf("读取 Catalog 失败: %v", err)
	}

	t.Logf("表单版本: %d", catalogData.FormVersion)
	t.Logf("可用场地: %d 个", len(catalogData.Options))
	t.Logf("可预约日期: %d 天", len(catalogData.DateMap))

	// 打印所有可用日期和时段
	t.Log("\n所有可预约选项:")
	dates := make([]string, 0, len(catalogData.DateMap))
	for date := range catalogData.DateMap {
		dates = append(dates, date)
	}
	sortStrings(dates)

	for _, date := range dates {
		dateInfo := catalogData.DateMap[date]
		hours := make([]int, 0, len(dateInfo.TimeMap))
		for hour := range dateInfo.TimeMap {
			hours = append(hours, hour)
		}
		sortInts(hours)
		t.Logf("  %s: %d 个时段可选", date, len(hours))
		for _, hour := range hours {
			t.Logf("    - %02d:00-%02d:00", hour, hour+1)
		}
	}

	// Step 2: 选择一个可用时段尝试预定
	t.Log("\n=== Step 2: 尝试预定一个可用时段 ===")

	if len(dates) == 0 {
		t.Fatal("没有可预约的日期")
	}

	// 选择最后一天的最晚时段，避免真的预约成功
	selectedDate := dates[len(dates)-1]
	dateInfo := catalogData.DateMap[selectedDate]
	hours := make([]int, 0, len(dateInfo.TimeMap))
	for hour := range dateInfo.TimeMap {
		hours = append(hours, hour)
	}
	sortInts(hours)
	selectedHour := hours[len(hours)-1]

	slot := common.BookingSlot{
		Date:  selectedDate,
		Hour:  selectedHour,
		Venue: 1,
	}

	t.Logf("尝试预定:")
	t.Logf("  日期: %s", slot.Date)
	t.Logf("  时段: %02d:00-%02d:00", slot.Hour, slot.Hour+1)
	t.Logf("  场地: %d号", slot.Venue)

	// 尝试预定
	err = bookingService.BookTimeSlot(catalogData, slot)
	if err != nil {
		t.Logf("预定失败: %v", err)
	} else {
		t.Log("预定成功！")
	}
}
