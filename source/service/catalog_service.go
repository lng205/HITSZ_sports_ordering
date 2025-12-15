package service

import (
	"encoding/json"
	"fmt"

	"sports_order/common"
)

// GetCatalogData 拉取并解析表单元数据（版本号、场地选项、日期/时段映射）。
func (s *BookingService) GetCatalogData() (*common.CatalogData, error) {
	// profile：获取表单版本等信息
	versionResp, err := s.apiClient.Get(common.FormURL + common.ProfileEndpoint)
	if err != nil {
		return nil, fmt.Errorf("请求表单配置失败: %v", err)
	}

	var profileResp common.ProfileResponse
	if err := json.Unmarshal(versionResp, &profileResp); err != nil {
		return nil, fmt.Errorf("反序列化表单配置失败: %v", err)
	}

	// catalog：获取可选场地与可预约日期/时段配置
	catalogResp, err := s.apiClient.Get(common.FormURL + common.CatalogEndpoint)
	if err != nil {
		return nil, fmt.Errorf("请求场地目录失败: %v", err)
	}

	var catalogResponse common.CatalogResponse
	if err := json.Unmarshal(catalogResp, &catalogResponse); err != nil {
		return nil, fmt.Errorf("反序列化场地目录失败: %v", err)
	}

	// 组装业务侧更好用的数据结构
	data := &common.CatalogData{
		FormVersion: profileResp.Data.Version,
		DateMap:     make(map[string]common.DateInfo),
	}

	// 场地预约配置固定在 catalogs[4]
	const venueCatalogIndex = 4
	if len(catalogResponse.Data.Catalogs) <= venueCatalogIndex {
		return nil, fmt.Errorf("表单长度错误，长度: %d", len(catalogResponse.Data.Catalogs))
	}
	venueCatalog := catalogResponse.Data.Catalogs[venueCatalogIndex]
	if venueCatalog.Type != common.TypeReservation {
		return nil, fmt.Errorf("表单场地项类型错误: %s", venueCatalog.Type)
	}

	// 抽取场地选项与日期/时段映射
	for _, formCatalog := range venueCatalog.FormCatalogs {
		switch common.CatalogRole(formCatalog.Role) {
		case common.RoleOption:
			// 场地选项（例如 1 号场、2 号场...）
			data.Options = append(data.Options, formCatalog.Cid)
		case common.RoleReservationDate:
			// 日期与时段映射
			data.DateMap[formCatalog.Content] = parseDateInfo(formCatalog)
		}
	}

	return data, nil
}

// parseDateInfo 从 FormCatalog 解析出日期信息及其时段映射。
func parseDateInfo(fc common.FormCatalog) common.DateInfo {
	dateInfo := common.DateInfo{
		DateID:  fc.Cid,
		TimeMap: make(map[int]string),
	}
	for _, child := range fc.ChildCatalogs {
		hour := child.Content.StartTime / 100
		dateInfo.TimeMap[hour] = child.Cid
	}
	return dateInfo
}
