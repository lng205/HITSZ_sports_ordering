package common

import (
	"time"
)

// ============================================================================
// 数据库模型（GORM）
// ============================================================================

type Log struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Level   string `json:"level" gorm:"not null"`
	Message string `json:"message" gorm:"not null"`

	OrderID *int `json:"order_id"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
}

type Order struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Date   string `json:"date" gorm:"not null"`
	Hour   int    `json:"hour" gorm:"not null"`
	Venue  int    `json:"venue" gorm:"not null"`
	Status string `json:"status" gorm:"not null"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}

// ============================================================================
// 配置模型
// ============================================================================

// User 用户信息（从配置文件读取）
type User struct {
	StudentID string `yaml:"student_id"`
	Name      string `yaml:"name"`
	Phone     string `yaml:"phone"`
	ImageURL  string `yaml:"image_url"`
	Token     string `yaml:"token"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// Config 应用配置
type Config struct {
	User     User           `yaml:"user"`
	Database DatabaseConfig `yaml:"database"`
}
