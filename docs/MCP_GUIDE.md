# EAM-Agent MCP (Model Context Protocol) 接入指南

EAM-Agent 原生支持 Anthropic 开源的 **MCP (Model Context Protocol 2024-11-05)** 协议。
任何支持 MCP 的宿主客户端（如 **Claude Desktop**, **Cursor**, **Hermes**, **Cline** 等）都可以通过 **本地 stdio 桥接** 或 **原生云端 SSE (Server-Sent Events)** 两种模式直连 EAM-Agent，让大模型直接感知并操作工业设备管理系统。

---

## 一、支持的 MCP Tools (核心工具集)

通过 MCP，外部 Agent 能够调用经过多租户与 RBAC 安全隔离的工业级工具（共计 17 个）：

| 工具分类 | 工具名称 | 功能描述 | 核心入参 |
| :--- | :--- | :--- | :--- |
| **设备台账** | `search_equipment` | 根据名称、设备编码或规格型号检索设备档案 | `keyword` (string) |
| | `get_equipment_profile` | 获取设备基础档案、出厂规格与位置 | `equipment_id` (int) |
| | `get_equipment_financials` | 查询设备原值、折旧净值与停机小时损失 | `equipment_id` (int) |
| **健康与预测** | `get_equipment_health` | 获取设备实时健康评分、劣化状态与 RUL 寿命预测 | `equipment_id` (int) |
| | `predict_remaining_life` | 专用 RUL 剩余使用寿命预测算法 | `equipment_id` (int) |
| | `detect_symptoms` | 设备亚健康与隐患征兆识别 | `equipment_id` (int) |
| | `get_retirement_recommendation` | 报废与以旧换新决策评估分析 | `equipment_id` (int) |
| **工单全流程** | `report_repair` | 设备发生异常时自动创建报修工单 | `equipment_id`, `fault_description`, `priority` |
| | `update_repair_progress` | 记录排故维修进度、排查动作与过程日志 | `order_id`, `progress_note` |
| | `close_repair_order` | 故障消缺结案，记录原因与解决方案并关闭工单 | `order_id`, `solution`, `fault_reason` |
| | `get_failure_stats` | 设备历史故障频率、MTTR、MTBF 统计 | `equipment_id` (int) |
| | `get_repair_costs` | 工单累计维修成本（人工+备件）统计 | `equipment_id` (int) |
| **备件与维保** | `get_spare_part_inventory` | 检查指定备件在指定厂区的库存保有量与警戒状态 | `spare_part_id`, `factory_id` |
| | `get_maintenance_compliance`| 点检保养依从性与计划执行率评估 | `equipment_id` (int) |
| **知识与手册** | `search_manual_knowledge` | 检索设备技术说明书切片、排故手册与图纸知识 | `query` (string) |
| **深度分析** | `sql_data_analyst` | 安全只读 SQL 引擎，用于跨表宏观多维数据洞察 | `sql_query` (string) |

---

## 二、获取访问令牌 (API 密钥 / PAT)

EAM-Agent 提供了完整的个人访问令牌（Personal Access Token）管理能力：

1. **在网页端一键生成**：
   * 登录 EAM-Agent 网页端，进入 **「Agent 外部集成」→「API 密钥」** 页面；
   * 点击 **「创建新密钥」**，设置密钥名称、权限范围（如 `read:equipment`, `write:repair`, `read:knowledge`）及有效期（支持 30天/90天/1年/永不过期）；
   * 系统将生成以 `ems_...` 开头的长效密钥（安全弹窗仅显示一次，支持一键复制）。
2. **多租户权限穿透**：
   * 密钥完全绑定生成者的企业角色与工厂归属；
   * 车间工人生成的密钥只能检索和报修所属车间的设备，绝无法跨厂越权或篡改配置。

---

## 三、客户端接入配置（支持两种模式）

### 模式一：本地 stdio 桥接模式（推荐用于 Claude Desktop / 本地命令行）

编译出的 `eam-mcp` 仅 **8.3MB** 单文件，无需在宿主机器安装任何额外语言运行时。

编辑 Claude Desktop 配置文件：
* **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
* **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "eam-agent": {
      "command": "/path/to/EAM-Agent/backend/bin/eam-mcp",
      "env": {
        "EAM_BASE_URL": "https://<your-eam-domain>",
        "EAM_TOKEN": "ems_xxxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
```

---

### 模式二：原生云端 SSE 直连模式（推荐用于 Cursor / Web / 无需本地二进制）

EAM-Agent 后端已原生内置 MCP SSE 端点：

* **SSE 接入地址**：`https://<your-eam-domain>/api/v1/mcp/sse`
* **认证方式**：在请求头传入 `Authorization: Bearer ems_xxxxxxxxxxxxxxxxxxxxxx`

在 Cursor 等客户端的 MCP 设置中配置该 URL 与 Token 即可立即唤醒。

---

## 四、自然语言交互示例

在支持 MCP 的对话窗口中，直接输入自然语言即可触发工具调用：

* **查设备台账**：*“帮我检索一下涂布车间名称包含‘曝光机’的设备，查看运行状态。”*
* **看健康与寿命**：*“调出设备 ID 为 1002 的健康诊断报告和剩余寿命预测。”*
* **查手册与对策**：*“报错代码 E-301 代表什么故障？调出说明书排故对策。”*
* **自然语言一键报修**：*“1002 号曝光机导轨出现剧烈异响并卡顿，立刻提交紧急报修工单。”*
* **流转更新维修记录**：*“在 9005 号工单记录排故过程：已清洗导轨并补充润滑油，正在试机。”*
* **结案关闭工单**：*“9005 号工单排故完毕，原因为导轨积尘严重，目前试机正常，关闭工单。”*
