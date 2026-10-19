package tool

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ems/backend/internal/model"
	"github.com/ems/backend/pkg/config"
	"github.com/ems/backend/pkg/database"
)

type SQLAnalystTool struct {
	SchemaDescription string
}

func NewSQLAnalystTool() *SQLAnalystTool {
	return &SQLAnalystTool{
		SchemaDescription: `
可供查询的业务表结构：
1. equipments: 设备表 (id, code, name, model, spec, status, purchase_price, purchase_date, workshop_id, type_id)
2. repair_orders: 报修单 (id, equipment_id, fault_description, status, reporter_id, assignee_id, priority, created_at, completed_at)
3. maintenance_tasks: 保养任务 (id, equipment_id, plan_id, status, scheduled_date, completed_at, actual_hours)
4. inspection_tasks: 点检任务 (id, equipment_id, status, scheduled_date, completed_at)
5. spare_parts: 备件表 (id, code, name, safety_stock, factory_id)
6. workshops: 车间表 (id, factory_id, name)
7. factories: 工厂表 (id, name)
8. equipment_types: 设备类型 (id, name, category)

安全与隔离规则：
- 所有查询必须是只读的单个 SELECT 语句。
- 严禁查询系统用户、密钥与认证相关表 (users, user_api_keys 等)。
- 非系统管理员查询请务必关联并包含 factory_id 或 workshop_id。
- 请使用标准 PostgreSQL 语法。
`,
	}
}

// 单词边界正则匹配破坏性关键字
var forbiddenRegex = regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|TRUNCATE|ALTER|CREATE|GRANT|REVOKE|EXEC|EXECUTE|COPY)\b`)

// 禁止访问的敏感表和系统元数据视图
var sensitiveTablesRegex = regexp.MustCompile(`(?i)\b(users|user_api_keys|agent_sessions|agent_messages|agent_artifacts|pg_shadow|pg_authid|pg_user|pg_proc|information_schema|pg_catalog)\b`)

func (t *SQLAnalystTool) ExecuteQuery(query string, user model.User) (interface{}, error) {
	trimmed := strings.TrimSpace(query)

	// 1. 禁止多语句执行 (Prevent stacked queries)
	if strings.Contains(trimmed, ";") {
		// 允许末尾有单个分号，但禁止中间有多条命令
		semiIdx := strings.Index(trimmed, ";")
		if semiIdx != len(trimmed)-1 {
			return nil, fmt.Errorf("安全拒绝：禁止执行多语句查询")
		}
		trimmed = strings.TrimSuffix(trimmed, ";")
	}

	// 2. 严格 SELECT 开头
	upperQuery := strings.ToUpper(trimmed)
	if !strings.HasPrefix(upperQuery, "SELECT") {
		return nil, fmt.Errorf("无效查询：必须以 SELECT 开头")
	}

	// 3. 关键字安全拦截 (使用词边界正则，避免误杀字段名中的 update 等)
	if forbiddenRegex.MatchString(trimmed) {
		return nil, fmt.Errorf("安全拒绝：仅允许执行只读分析查询 (SELECT)")
	}

	// 4. 敏感表白名单与黑名单防护
	if sensitiveTablesRegex.MatchString(trimmed) {
		return nil, fmt.Errorf("安全拒绝：禁止查询系统核心认证或凭据表")
	}

	// 5. 存储模式判断
	if config.Cfg.Storage.Mode == "memory" {
		return t.executeInMemory(trimmed, user)
	}

	return t.executeInDB(trimmed, user)
}

func (t *SQLAnalystTool) executeInDB(query string, user model.User) (interface{}, error) {
	db := database.GetDB()

	rows := []map[string]interface{}{}

	// Limit results to 100 to prevent memory exhaustion
	safeQuery := fmt.Sprintf("SELECT * FROM (%s) AS subquery LIMIT 100", query)

	err := db.Raw(safeQuery).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("SQL 执行失败: %v", err)
	}

	// 4. 多租户工厂隔离过滤 (Default Deny 策略)
	if user.Role != "admin" && user.FactoryID != nil {
		filteredRows := []map[string]interface{}{}
		for _, row := range rows {
			if t.isRowAuthorized(row, user) {
				filteredRows = append(filteredRows, row)
			}
		}
		return filteredRows, nil
	}

	return rows, nil
}

func (t *SQLAnalystTool) executeInMemory(query string, user model.User) (interface{}, error) {
	return nil, fmt.Errorf("数据分析工具目前仅在数据库模式下可用，请切换至生产数据库环境")
}

func (t *SQLAnalystTool) isRowAuthorized(row map[string]interface{}, user model.User) bool {
	if user.Role == "admin" || user.FactoryID == nil {
		return true
	}

	fid := *user.FactoryID

	// 检查直接包含 factory_id 的情况
	if val, ok := row["factory_id"]; ok {
		return t.toUint(val) == fid
	}

	// 检查包含 workshop_id 的情况
	if val, ok := row["workshop_id"]; ok {
		wsID := t.toUint(val)
		if wsID > 0 {
			var workshop model.Workshop
			if err := database.GetDB().Select("factory_id").First(&workshop, wsID).Error; err == nil {
				return workshop.FactoryID == fid
			}
		}
	}

	// 检查包含 equipment_id 的情况
	if val, ok := row["equipment_id"]; ok {
		eqID := t.toUint(val)
		if eqID > 0 {
			var eq model.Equipment
			if err := database.GetDB().Preload("Workshop").First(&eq, eqID).Error; err == nil && eq.Workshop != nil {
				return eq.Workshop.FactoryID == fid
			}
		}
	}

	// 默认安全原则：对于非管理员用户，没有明确所属工厂标记的行，一律拒绝返回（Default Deny）
	return false
}

func (t *SQLAnalystTool) toUint(val interface{}) uint {
	switch v := val.(type) {
	case int:
		return uint(v)
	case int64:
		return uint(v)
	case uint:
		return v
	case float64:
		return uint(v)
	}
	return 0
}
