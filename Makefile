.PHONY: help build test init init-db add-order config

# 项目配置
APP_NAME := sports-order
DB_NAME := sports-order.db
SOURCE_DIR := source
DB_DIR := database

# 默认执行 help
.DEFAULT_GOAL := help

help: ## 显示帮助信息
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

init: config init-db ## 初始化项目 (生成配置 + 初始化数据库)

build: ## 编译 Go 程序
	cd $(SOURCE_DIR) && go build -o ../$(APP_NAME) .

test: ## 运行测试 (访问真实 API)
	cd $(SOURCE_DIR) && go test -v

init-db: ## 初始化数据库
	sqlite3 $(DB_NAME) < $(DB_DIR)/init.sql

add-order: ## 交互式添加订单
	@./$(DB_DIR)/add_order.sh

config: ## 生成配置文件 (安全模式：仅不存在时生成)
	@if [ -f config.yaml ]; then \
		echo "⚠️  config.yaml 已存在，跳过生成以防止覆盖您的配置。"; \
	else \
		echo "📝 生成默认配置文件 config.yaml..."; \
		echo "# 用户配置" > config.yaml; \
		echo "user:" >> config.yaml; \
		echo "  student_id: \"20231234567\"    # 学号" >> config.yaml; \
		echo "  name: \"张三\"                  # 姓名" >> config.yaml; \
		echo "  phone: \"13800138000\"         # 手机号" >> config.yaml; \
		echo "  image_url: \"\"                # 头像URL（使用HTTPS抓包获取，时效30天+）" >> config.yaml; \
		echo "  token: \"hEBOountLgwjBUl4FW9Vv2GGOpmoIQR1FLzRT2TuFROh9gW36DLe2VY5L8Jzp0m7-oVsbQ\"                    # 认证令牌（使用HTTPS抓包获取，时效48小时）" >> config.yaml; \
		echo "" >> config.yaml; \
		echo "# 数据库配置" >> config.yaml; \
		echo "database:" >> config.yaml; \
		echo "  path: \"sports-order.db\"" >> config.yaml; \
	fi
