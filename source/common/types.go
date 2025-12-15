package common

// ============================================================================
// 对外 API 响应结构（我们从外部接口收到的 JSON）
// ============================================================================

// ProfileResponse 对应 profile 接口响应。
type ProfileResponse struct {
	Code int `json:"code"`
	Data struct {
		Version int `json:"version"`
	} `json:"data"`
}

// CatalogResponse 对应 catalog 接口响应。
type CatalogResponse struct {
	Code int `json:"code"`
	Data struct {
		Catalogs []Catalog `json:"catalogs"`
	} `json:"data"`
}

// Catalog 表示一个表单目录节点。
type Catalog struct {
	Cid          string        `json:"cid"`
	Type         string        `json:"type"`
	FormCatalogs []FormCatalog `json:"formCatalogs"`
}

// FormCatalog 表示某个目录节点下的可选项/日期节点。
type FormCatalog struct {
	Cid           string         `json:"cid"`
	Content       string         `json:"content"`
	Role          string         `json:"role"`
	ChildCatalogs []ChildCatalog `json:"childCatalogs,omitempty"`
}

// ChildCatalog 表示时段节点及其元数据。
type ChildCatalog struct {
	Cid     string `json:"cid"`
	Content struct {
		StartTime int `json:"startTime"`
		EndTime   int `json:"endTime"`
	} `json:"content"`
}

// BookingResponse 对应提交预约后的响应。
type BookingResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// ============================================================================
// 对外 API 请求结构（我们提交给外部接口的 JSON）
// ============================================================================

// BookingRequest 是提交预约的请求体。
type BookingRequest struct {
	Catalogs      []RequestField `json:"catalogs"`
	ShowQuestions []string       `json:"showQuestions"`
	FormVersion   int            `json:"formVersion"`
}

// RequestField 表示表单里的一个字段。
type RequestField struct {
	Type  string `json:"type"`
	Cid   string `json:"cid"`
	Value any    `json:"value"`
}

// ReservationValue 描述"预约字段"的具体值（日期、时段、场地等）。
type ReservationValue struct {
	DateID     string `json:"dateId"`
	TimeID     string `json:"timeId"`
	OptionID   string `json:"optionId"`
	Count      int    `json:"count"`
	DateStr    string `json:"dateStr"`
	TimeStr    string `json:"timeStr"`
	OptionName string `json:"optionName"`
}

// ============================================================================
// 内部业务结构（用于更方便地做业务计算）
// ============================================================================

// CatalogData 保存解析后的配置：版本号、场地选项，以及日期到时段的映射。
type CatalogData struct {
	FormVersion int                 // 表单版本号
	Options     []string            // 场地选项 ID 列表
	DateMap     map[string]DateInfo // 日期 -> 时段映射
}

// DateInfo 表示某一天的时段映射。
type DateInfo struct {
	DateID  string         // 外部 API 的日期标识
	TimeMap map[int]string // 小时 -> 时段 ID
}

// BookingSlot 表示一个预约时段的参数。
type BookingSlot struct {
	Date  string // 预约日期，如 "2025-12-15"
	Hour  int    // 开始小时，如 14 表示 14:00-15:00
	Venue int    // 场地号，从 1 开始
}
