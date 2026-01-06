#!/bin/bash
# 交互式添加订单脚本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 数据库路径（相对于项目根目录）
DB_PATH="sports-order.db"

# 获取脚本所在目录的父目录（项目根目录）
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DB_FULL_PATH="$PROJECT_ROOT/$DB_PATH"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 检查数据库是否存在
check_database() {
    if [ ! -f "$DB_FULL_PATH" ]; then
        print_error "数据库文件不存在: $DB_FULL_PATH"
        print_info "请先运行 'make init-db' 初始化数据库"
        exit 1
    fi
}

# 验证日期格式 (YYYY-MM-DD)
validate_date() {
    local date=$1
    if [[ ! $date =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}$ ]]; then
        return 1
    fi
    # 验证日期是否有效
    date -d "$date" > /dev/null 2>&1
    return $?
}

# 验证小时 (7-22)
validate_hour() {
    local hour=$1
    if [[ ! $hour =~ ^[0-9]+$ ]]; then
        return 1
    fi
    if [ "$hour" -lt 7 ] || [ "$hour" -gt 22 ]; then
        return 1
    fi
    return 0
}

# 验证场地编号 (1-10)
validate_venue() {
    local venue=$1
    if [[ ! $venue =~ ^[0-9]+$ ]]; then
        return 1
    fi
    if [ "$venue" -lt 1 ] || [ "$venue" -gt 10 ]; then
        return 1
    fi
    return 0
}

# 显示现有订单
show_orders() {
    print_info "当前待处理订单:"
    echo "----------------------------------------"
    sqlite3 -header -column "$DB_FULL_PATH" \
        "SELECT id, date, hour || ':00-' || (hour+1) || ':00' AS time_slot, venue, status, created_at 
         FROM orders 
         WHERE status = 'PENDING' 
         ORDER BY date, hour;"
    echo "----------------------------------------"
}

# 主菜单
show_menu() {
    echo ""
    echo -e "${BLUE}========== 订单管理 ==========${NC}"
    echo "1. 添加新订单"
    echo "2. 查看待处理订单"
    echo "3. 查看所有订单"
    echo "4. 删除订单"
    echo "5. 批量添加订单"
    echo "0. 退出"
    echo -e "${BLUE}==============================${NC}"
    echo -n "请选择操作 [0-5]: "
}

# 添加单个订单
add_single_order() {
    echo ""
    print_info "=== 添加新订单 ==="
    
    # 输入日期
    while true; do
        echo -n "请输入预约日期 (YYYY-MM-DD，回车默认后天): "
        read -r input_date
        if [ -z "$input_date" ]; then
            input_date=$(date -d "+2 days" +%Y-%m-%d)
            print_info "使用默认日期: $input_date"
        fi
        if validate_date "$input_date"; then
            break
        else
            print_error "日期格式无效，请使用 YYYY-MM-DD 格式"
        fi
    done
    
    # 输入小时
    while true; do
        echo -n "请输入预约时段 (7-22，如 15 表示 15:00-16:00): "
        read -r input_hour
        if validate_hour "$input_hour"; then
            break
        else
            print_error "时段无效，请输入 7-22 之间的数字"
        fi
    done
    
    # 输入场地
    while true; do
        echo -n "请输入场地编号 (1-10，回车默认 4): "
        read -r input_venue
        if [ -z "$input_venue" ]; then
            input_venue=4
            print_info "使用默认场地: $input_venue"
        fi
        if validate_venue "$input_venue"; then
            break
        else
            print_error "场地编号无效，请输入 1-10 之间的数字"
        fi
    done
    
    # 确认信息
    echo ""
    print_info "订单信息确认:"
    echo "  日期: $input_date"
    echo "  时段: ${input_hour}:00-$((input_hour+1)):00"
    echo "  场地: $input_venue 号场"
    echo -n "确认添加? (Y/n): "
    read -r confirm
    
    if [ "$confirm" = "n" ] || [ "$confirm" = "N" ]; then
        print_warning "已取消添加"
        return
    fi
    
    # 插入数据库
    sqlite3 "$DB_FULL_PATH" \
        "INSERT INTO orders (date, hour, venue, status) VALUES ('$input_date', $input_hour, $input_venue, 'PENDING');"
    
    if [ $? -eq 0 ]; then
        print_success "订单添加成功!"
    else
        print_error "订单添加失败"
    fi
}

# 查看所有订单
show_all_orders() {
    echo ""
    print_info "所有订单:"
    echo "----------------------------------------"
    sqlite3 -header -column "$DB_FULL_PATH" \
        "SELECT id, date, hour || ':00-' || (hour+1) || ':00' AS time_slot, venue, status, created_at 
         FROM orders 
         ORDER BY date DESC, hour;"
    echo "----------------------------------------"
}

# 删除订单
delete_order() {
    echo ""
    show_orders
    echo -n "请输入要删除的订单ID (输入 0 取消): "
    read -r order_id
    
    if [ "$order_id" = "0" ]; then
        print_warning "已取消删除"
        return
    fi
    
    # 检查订单是否存在
    exists=$(sqlite3 "$DB_FULL_PATH" "SELECT COUNT(*) FROM orders WHERE id = $order_id;")
    if [ "$exists" = "0" ]; then
        print_error "订单不存在"
        return
    fi
    
    echo -n "确认删除订单 #$order_id? (y/N): "
    read -r confirm
    
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        sqlite3 "$DB_FULL_PATH" "DELETE FROM orders WHERE id = $order_id;"
        print_success "订单 #$order_id 已删除"
    else
        print_warning "已取消删除"
    fi
}

# 批量添加订单
batch_add_orders() {
    echo ""
    print_info "=== 批量添加订单 ==="
    
    # 输入日期
    while true; do
        echo -n "请输入预约日期 (YYYY-MM-DD，回车默认后天): "
        read -r input_date
        if [ -z "$input_date" ]; then
            input_date=$(date -d "+2 days" +%Y-%m-%d)
            print_info "使用默认日期: $input_date"
        fi
        if validate_date "$input_date"; then
            break
        else
            print_error "日期格式无效，请使用 YYYY-MM-DD 格式"
        fi
    done
    
    # 输入时段范围
    while true; do
        echo -n "请输入起始时段 (7-22): "
        read -r start_hour
        if validate_hour "$start_hour"; then
            break
        else
            print_error "时段无效，请输入 7-22 之间的数字"
        fi
    done
    
    while true; do
        echo -n "请输入结束时段 (${start_hour}-22): "
        read -r end_hour
        if validate_hour "$end_hour" && [ "$end_hour" -ge "$start_hour" ]; then
            break
        else
            print_error "时段无效，请输入 ${start_hour}-22 之间的数字"
        fi
    done
    
    # 输入场地
    while true; do
        echo -n "请输入场地编号 (1-10，回车默认 4): "
        read -r input_venue
        if [ -z "$input_venue" ]; then
            input_venue=4
            print_info "使用默认场地: $input_venue"
        fi
        if validate_venue "$input_venue"; then
            break
        else
            print_error "场地编号无效，请输入 1-10 之间的数字"
        fi
    done
    
    # 计算订单数量
    order_count=$((end_hour - start_hour + 1))
    
    # 确认信息
    echo ""
    print_info "批量订单信息确认:"
    echo "  日期: $input_date"
    echo "  时段: ${start_hour}:00 - $((end_hour+1)):00 (共 $order_count 个时段)"
    echo "  场地: $input_venue 号场"
    echo -n "确认添加 $order_count 个订单? (Y/n): "
    read -r confirm
    
    if [ "$confirm" = "n" ] || [ "$confirm" = "N" ]; then
        print_warning "已取消添加"
        return
    fi
    
    # 批量插入
    success_count=0
    for hour in $(seq "$start_hour" "$end_hour"); do
        sqlite3 "$DB_FULL_PATH" \
            "INSERT INTO orders (date, hour, venue, status) VALUES ('$input_date', $hour, $input_venue, 'PENDING');"
        if [ $? -eq 0 ]; then
            success_count=$((success_count + 1))
        fi
    done
    
    print_success "成功添加 $success_count 个订单!"
}

# 主程序
main() {
    check_database
    
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════╗"
    echo "║      体育场馆预约订单管理工具         ║"
    echo "╚═══════════════════════════════════════╝"
    echo -e "${NC}"
    
    while true; do
        show_menu
        read -r choice
        
        case $choice in
            1)
                add_single_order
                ;;
            2)
                show_orders
                ;;
            3)
                show_all_orders
                ;;
            4)
                delete_order
                ;;
            5)
                batch_add_orders
                ;;
            0)
                print_info "再见!"
                exit 0
                ;;
            *)
                print_error "无效选择，请重试"
                ;;
        esac
    done
}

# 运行主程序
main
