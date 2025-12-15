package service

import (
	"encoding/json"
	"fmt"

	"sports_order/common"
)

// BookingService 封装"预约下单"相关的业务逻辑。
type BookingService struct {
	apiClient common.APIClient
	repo      common.Repository
	user      *common.User // 从配置文件加载的用户信息
}

// NewBookingService 创建预约服务。
func NewBookingService(apiClient common.APIClient, repo common.Repository, user *common.User) *BookingService {
	return &BookingService{apiClient: apiClient, repo: repo, user: user}
}

// BookTimeSlot 针对某一天某一小时提交一次预约请求。
func (s *BookingService) BookTimeSlot(data *common.CatalogData, slot common.BookingSlot) error {
	dateInfo, exists := data.DateMap[slot.Date]
	if !exists {
		return fmt.Errorf("日期 %s 不可预约", slot.Date)
	}

	if _, exists := dateInfo.TimeMap[slot.Hour]; !exists {
		return fmt.Errorf("时段 %d:00 在 %s 不可预约", slot.Hour, slot.Date)
	}

	if slot.Venue < 1 || slot.Venue > len(data.Options) {
		return fmt.Errorf("无效的场地号 %d", slot.Venue)
	}

	request := buildBookingRequest(s.user, data, slot)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	resp, err := s.apiClient.Post(common.FormDataURL, jsonData, s.user.Token)
	if err != nil {
		return fmt.Errorf("提交预约请求失败: %v", err)
	}

	var bookingResp common.BookingResponse
	if err := json.Unmarshal(resp, &bookingResp); err != nil {
		return fmt.Errorf("解析预约响应失败: %v", err)
	}

	if bookingResp.Code != common.ResponseCodeSuccess {
		return fmt.Errorf("预约失败: %s", bookingResp.Message)
	}

	return nil
}

// buildBookingRequest 构造对外 API 需要的预约请求体。
func buildBookingRequest(user *common.User, catalog *common.CatalogData, slot common.BookingSlot) common.BookingRequest {
	dateInfo := catalog.DateMap[slot.Date]
	return common.BookingRequest{
		Catalogs: []common.RequestField{
			{Type: common.TypeWord, Cid: common.CIDName, Value: user.Name},
			{Type: common.TypeTelephone, Cid: common.CIDPhone, Value: user.Phone},
			{Type: common.TypeWord, Cid: common.CIDStudentID, Value: user.StudentID},
			{Type: common.TypeImage, Cid: common.CIDImage, Value: []string{user.ImageURL}},
			{Type: common.TypeReservation, Cid: common.CIDReservation, Value: []common.ReservationValue{
				{
					DateID:     dateInfo.DateID,
					TimeID:     dateInfo.TimeMap[slot.Hour],
					OptionID:   catalog.Options[slot.Venue-1],
					Count:      common.DefaultVenueCount,
					DateStr:    slot.Date,
					TimeStr:    fmt.Sprintf("%02d:00-%02d:00", slot.Hour, slot.Hour+1),
					OptionName: fmt.Sprintf("%d号", slot.Venue),
				},
			}},
		},
		ShowQuestions: common.ShowQuestions,
		FormVersion:   catalog.FormVersion,
	}
}
