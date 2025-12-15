# 项目配置
APP_NAME := sports-order
DB_NAME := sports-order.db
SOURCE_DIR := source
DB_DIR := database

# 编译 Go 程序
build:
	cd $(SOURCE_DIR) && go build -o ../$(APP_NAME) .

# 运行测试（访问真实 API）
test:
	cd $(SOURCE_DIR) && go test -v

# 初始化数据库
init-db:
	sqlite3 $(DB_NAME) < $(DB_DIR)/init.sql

# 交互式添加订单
add-order:
	@./$(DB_DIR)/add_order.sh
