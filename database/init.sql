-- 数据库库表配置

-- 启用外键约束
PRAGMA foreign_keys = ON;

-- 订单表: 存储场地预约订单（用户信息从 config.yaml 读取）
CREATE TABLE IF NOT EXISTS `orders` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,   -- 订单ID（主键，自增）
    `date` TEXT NOT NULL,                      -- 预约日期
    `hour` INTEGER NOT NULL,                   -- 预约时段（小时，如15表示15:00-16:00）
    `venue` INTEGER NOT NULL DEFAULT 4,        -- 场地编号
    `status` TEXT NOT NULL DEFAULT 'PENDING',  -- 订单状态: PENDING-待处理, SUCCESS-成功, FAILED-失败
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP   -- 更新时间
);

-- 日志表: 存储应用日志和预约记录
CREATE TABLE IF NOT EXISTS `logs` (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,   -- 日志ID（主键，自增）
    `level` TEXT NOT NULL,                     -- 日志级别: INFO, WARN, ERROR等
    `message` TEXT NOT NULL,                   -- 日志消息内容
    `order_id` INTEGER,                        -- 关联订单ID（可为空）
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    FOREIGN KEY (`order_id`) REFERENCES `orders`(`id`)         -- 关联订单表
);
