# EAM-Agent MCP (Model Context Protocol) 接入指南

EAM-Agent 原生支持 Anthropic 开源的 **MCP (Model Context Protocol)** 协议。任何支持 MCP 的宿主客户端（如 **Claude Desktop**, **Cursor**, **Hermes**, **Cline** 等）都可以通过 stdio 管道直连 EAM-Agent，让大模型直接感知并操作工业设备管理系统。

---

## 一、支持的 MCP Tools (核心工具集)

通过 MCP，外部 Agent 能够调用以下 4 个经过多租户与 RBAC 安全隔离的工业级工具：

| 工具名称 | 功能描述 | 核心入参 |
| :--- | :--- | :--- |
| `search_equipment` | 根据名称、设备编码或规格型号检索设备档案 | `keyword` (string) |
| `get_equipment_health` | 获取设备实时健康评分、劣化状态与剩余使用寿命 (RUL) 预测 | `equipment_id` (int) |
| `get_spare_part_inventory` | 检查指定备件在指定厂区的库存保有量与警戒状态 | `spare_part_id`, `factory_id` |
| `report_repair` | 设备发生异常时自动创建报修工单并通知车间维修人员 | `equipment_id`, `fault_description`, `priority` |

---

## 二、获取访问令牌 (Token)

EAM-Agent MCP 严格执行身份穿透与工厂数据隔离。

1. **登录获取 Token**：
   调用 EAM 登录接口或在网页端登录后获取 JWT Token：
   ```bash
   curl -X POST https://<your-eam-domain>/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"your_username","password":"your_password"}'
   ```
   返回结果中的 `token` 即为访问凭据。
2. 用户角色（工人、维修工、主管、管理员）将直接决定 Agent 能检索哪些设备与执行哪些流转操作，避免越权。

---

## 三、客户端接入配置

### 1. Claude Desktop 接入

编辑 Claude Desktop 配置文件：
* **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
* **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

添加 `eam-agent` 服务配置：

```json
{
  "mcpServers": {
    "eam-agent": {
      "command": "/绝对路径/EAM-Agent/backend/bin/eam-mcp",
      "env": {
        "EAM_BASE_URL": "https://<your-eam-domain>",
        "EAM_TOKEN": "your_bearer_token_here"
      }
    }
  }
}
```

重启 Claude Desktop，界面右下角会亮起锤子 🔨 图标，显示 4 个可用的 EAM 设备工具。

---

### 2. Cursor / Cline 接入

在 Cursor 的 `Settings -> Features -> MCP` 中添加：
* **Type**: `command`
* **Command**: `/path/to/eam-mcp`
* **Environment Variables**:
  * `EAM_BASE_URL`: `https://<your-eam-domain>`
  * `EAM_TOKEN`: `your_bearer_token_here`

---

## 四、自然语言交互示例

在支持 MCP 的对话窗口中，直接输入自然语言即可触发工具调用：

* **查设备**：*“帮我检索一下涂布车间名称包含‘曝光机’的设备，查看运行状态。”*
* **看健康与寿命**：*“调出设备 ID 为 1002 的健康诊断报告和剩余寿命预测。”*
* **查备件**：*“看一下备件库里 SKF 轴承（ID: 201）在 1 号工厂还有多少库存？”*
* **自然语言一键报修**：*“1002 号曝光机导轨出现剧烈异响并卡顿，立刻提交高优先级的报修工单。”*
