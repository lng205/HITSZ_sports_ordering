package main

import (
	"fmt"
	"os"

	"sports_order/common"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ============================================================================
// 配置加载
// ============================================================================

// LoadConfig 从默认路径 config.yaml 加载配置
func LoadConfig() (*common.Config, error) {
	return LoadConfigFrom("config.yaml")
}

// LoadConfigFrom 从指定路径加载配置
func LoadConfigFrom(path string) (*common.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config common.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// ============================================================================
// 数据库仓储
// ============================================================================

// Repository 封装数据库操作。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 Repository。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// FindOrdersByDate 查询指定日期的待处理订单。
func (r *Repository) FindOrdersByDate(date string) ([]*common.Order, error) {
	var orders []*common.Order
	return orders, r.db.Where("date = ? AND status = ?", date, common.OrderStatusPending).Find(&orders).Error
}

// UpdateOrderStatus 更新订单状态。
func (r *Repository) UpdateOrderStatus(id uint, status common.OrderStatus) error {
	updates := map[string]any{
		"status": string(status),
	}
	return r.db.Model(&common.Order{}).Where("id = ?", id).Updates(updates).Error
}

// CreateLog 写入一条日志记录。
func (r *Repository) CreateLog(level common.LogLevel, message string, orderID *int) error {
	entry := &common.Log{
		Level:   string(level),
		Message: message,
		OrderID: orderID,
	}
	return r.db.Create(entry).Error
}

// CreateLogf 写入格式化日志记录。
func (r *Repository) CreateLogf(level common.LogLevel, orderID *int, format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return r.CreateLog(level, message, orderID)
}

// ============================================================================
// 数据库初始化
// ============================================================================

// InitDB 使用 GORM 打开 SQLite 数据库并返回连接。
func InitDB(config *common.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}

// CloseDB 关闭数据库连接。
func CloseDB(db *gorm.DB) {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
