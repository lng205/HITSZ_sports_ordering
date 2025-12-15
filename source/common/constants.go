package common

// API 地址与路径常量
const (
	BaseURL         = "https://form.qun100.com"
	FormID          = "1627049420674297856"
	FormURL         = BaseURL + "/v1/form/" + FormID + "/"
	FormDataURL     = BaseURL + "/v1/" + FormID + "/form_data"
	ProfileEndpoint = "profile"
	CatalogEndpoint = "catalog"
)

// 表单字段 CID（Content ID）
const (
	CIDName        = "1627049422343630849"
	CIDPhone       = "1627049422343630851"
	CIDStudentID   = "1628873244736188417"
	CIDImage       = "1627056232015765504"
	CIDReservation = "1627049422343630855"
)

// 表单字段类型
const (
	TypeWord        = "WORD"
	TypeTelephone   = "TELEPHONE"
	TypeImage       = "IMAGE"
	TypeReservation = "RESERVATION"
)

// ShowQuestions 指定提交预约时需要展示/提交的字段 CID 列表。
var ShowQuestions = []string{CIDName, CIDPhone, CIDStudentID, CIDImage, CIDReservation}

// HTTPHeaders 是调用外部表单 API 时需要携带的一组固定请求头。
var HTTPHeaders = map[string]string{
	"client-form-id": FormID,
	"client-app-id":  "wxfc4ef6d539d03373",
	"Content-Type":   "application/json",
}

// 对外 API 返回码
const ResponseCodeSuccess = 0

// 默认值
const (
	DefaultVenueCount = 1
	DefaultTimeoutSec = 30
	DaysAhead         = 2 // 处理"今天 + N 天"的订单
)

// CatalogRole 表示 catalog 节点的角色类型。
type CatalogRole string

const (
	RoleOption          CatalogRole = "OPTION"
	RoleReservationDate CatalogRole = "RESERVATION_DATE"
)

// LogLevel 表示日志级别。
type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
)

// OrderStatus 表示订单状态。
type OrderStatus string

const (
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusSuccess OrderStatus = "SUCCESS"
	OrderStatusFailed  OrderStatus = "FAILED"
)
