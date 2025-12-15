package common

// ============================================================================
// 接口定义（便于替换实现与单元测试）
// ============================================================================

// APIClient 抽象对外 HTTP 调用，便于 mock。
type APIClient interface {
	Get(url string) ([]byte, error)
	Post(url string, data []byte, auth string) ([]byte, error)
}

// Repository 抽象数据库操作。
type Repository interface {
	// 订单相关
	FindOrdersByDate(date string) ([]*Order, error)
	UpdateOrderStatus(id uint, status OrderStatus) error
	// 日志相关
	CreateLog(level LogLevel, message string, orderID *int) error
	CreateLogf(level LogLevel, orderID *int, format string, args ...any) error
}
