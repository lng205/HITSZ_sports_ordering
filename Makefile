# 项目配置
APP_NAME := sports-order
DB_NAME := sports-order.db
SOURCE_DIR := source
DB_DIR := database

# 编译 Go 程序
build: config
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

# 生成配置文件
config:
	@echo "# 用户配置" > config.yaml
	@echo "user:" >> config.yaml
	@echo "  student_id: \"20231234567\"    # 学号" >> config.yaml
	@echo "  name: \"张三\"                  # 姓名" >> config.yaml
	@echo "  phone: \"13800138000\"         # 手机号" >> config.yaml
	@echo "  image_url: \"\"                # 头像URL（使用HTTPS抓包获取，时效30天+）" >> config.yaml
	@echo "  token: \"hEBOountLgwjBUl4FW9Vv2GGOpmoIQR1FLzRT2TuFROh9gW36DLe2VY5L8Jzp0m7-oVsbQ\"                    # 认证令牌（使用HTTPS抓包获取，时效48小时）" >> config.yaml
	@echo "" >> config.yaml
	@echo "# 数据库配置" >> config.yaml
	@echo "database:" >> config.yaml
	@echo "  path: \"sports-order.db\"" >> config.yaml
