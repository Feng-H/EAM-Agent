<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section animate-fade-in">
      <div class="welcome-content">
        <h1 class="welcome-title">
          <span class="title-greeting">欢迎回来，</span>
          <span class="title-name">{{ authStore.userName }}</span>
        </h1>
        <p class="welcome-subtitle">EAM-Agent · 工业资产大脑与智能运维中枢</p>
      </div>
      <div class="welcome-time">
        <div class="time-display">{{ currentTime }}</div>
        <div class="date-display">{{ currentDate }}</div>
      </div>
    </div>

    <!-- 设备统计卡片 -->
    <div class="stats-grid animate-fade-in animate-delay-1">
      <div class="stat-card" data-type="total">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
              <line x1="12" y1="22.08" x2="12" y2="12"/>
            </svg>
          </div>
          <div class="stat-trend">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.total }}</div>
        <div class="stat-label">设备总数</div>
        <div class="stat-bar">
          <div class="bar-fill" style="width: 100%"></div>
        </div>
      </div>

      <div class="stat-card" data-type="running">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
          </div>
          <div class="stat-trend up">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/>
              <polyline points="17 6 23 6 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.running }}</div>
        <div class="stat-label">运行中</div>
        <div class="stat-bar">
          <div class="bar-fill success" :style="{ width: runningPercent + '%' }"></div>
        </div>
      </div>

      <div class="stat-card" data-type="maintenance">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
            </svg>
          </div>
          <div class="stat-trend">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="6" x2="12" y2="12"/>
              <line x1="12" y1="12" x2="16" y2="14"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.maintenance }}</div>
        <div class="stat-label">维修中</div>
        <div class="stat-bar">
          <div class="bar-fill warning" :style="{ width: maintenancePercent + '%' }"></div>
        </div>
      </div>

      <div class="stat-card" data-type="stopped">
        <div class="stat-header">
          <div class="stat-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="6" y="4" width="4" height="16"/>
              <rect x="14" y="4" width="4" height="16"/>
            </svg>
          </div>
          <div class="stat-trend down">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 18 13.5 8.5 8.5 13.5 1 6"/>
              <polyline points="17 18 23 18 23 12"/>
            </svg>
          </div>
        </div>
        <div class="stat-value">{{ statistics.stopped }}</div>
        <div class="stat-label">已停机</div>
        <div class="stat-bar">
          <div class="bar-fill danger" :style="{ width: stoppedPercent + '%' }"></div>
        </div>
      </div>
    </div>

    <!-- 点检统计面板 -->
    <div class="inspection-panel animate-fade-in animate-delay-2">
      <div class="panel-header">
        <h3 class="panel-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          点检统计
        </h3>
        <router-link to="/inspection/tasks" class="panel-link">
          查看全部
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </router-link>
      </div>

      <div class="inspection-stats">
        <div class="istat-item">
          <div class="istat-value">{{ inspectionStats.total_tasks }}</div>
          <div class="istat-label">总任务数</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value pending">{{ inspectionStats.pending_tasks }}</div>
          <div class="istat-label">待执行</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value progress">{{ inspectionStats.in_progress_tasks }}</div>
          <div class="istat-label">进行中</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value completed">{{ inspectionStats.completed_tasks }}</div>
          <div class="istat-label">已完成</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item">
          <div class="istat-value today">{{ inspectionStats.today_completed }}</div>
          <div class="istat-label">今日完成</div>
        </div>
        <div class="istat-divider"></div>
        <div class="istat-item highlight">
          <div class="istat-value">{{ inspectionStats.completion_rate.toFixed(1) }}%</div>
          <div class="istat-label">完成率</div>
        </div>
      </div>

      <!-- 进度条 -->
      <div class="progress-track">
        <div
          class="progress-fill"
          :style="{ width: inspectionStats.completion_rate + '%' }"
        ></div>
      </div>
    </div>

    <!-- 下方区域 -->
    <div class="dashboard-grid animate-fade-in animate-delay-3">
      <!-- 快速操作 -->
      <div class="dashboard-card quick-actions">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          快速操作
        </div>
        <div class="action-list">
          <router-link to="/equipment" class="action-item">
            <div class="action-icon primary">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">设备台账</div>
              <div class="action-desc">查看所有设备信息</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/inspection/execute" class="action-item">
            <div class="action-icon success">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">开始点检</div>
              <div class="action-desc">执行设备点检任务</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/inspection/tasks" class="action-item">
            <div class="action-icon warning">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="8" y1="6" x2="21" y2="6"/>
                <line x1="8" y1="12" x2="21" y2="12"/>
                <line x1="8" y1="18" x2="21" y2="18"/>
                <line x1="3" y1="6" x2="3.01" y2="6"/>
                <line x1="3" y1="12" x2="3.01" y2="12"/>
                <line x1="3" y1="18" x2="3.01" y2="18"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">点检任务</div>
              <div class="action-desc">管理点检任务列表</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>

          <router-link to="/analytics" class="action-item">
            <div class="action-icon info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="18" y1="20" x2="18" y2="10"/>
                <line x1="12" y1="20" x2="12" y2="4"/>
                <line x1="6" y1="20" x2="6" y2="14"/>
              </svg>
            </div>
            <div class="action-content">
              <div class="action-title">数据统计</div>
              <div class="action-desc">查看数据分析报告</div>
            </div>
            <svg class="action-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 18 15 12 9 6"/>
            </svg>
          </router-link>
        </div>
      </div>

      <!-- 我的今日任务 -->
      <div class="dashboard-card my-tasks">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <polyline points="12 6 12 12 16 14"/>
          </svg>
          我的今日任务
        </div>
        <div v-if="myStats.loading" class="tasks-loading">
          <div class="loading-spinner"></div>
        </div>
        <div v-else class="tasks-summary">
          <div class="task-metric">
            <div class="metric-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.pending_count }}</div>
              <div class="metric-label">待执行</div>
            </div>
          </div>
          <div class="task-metric">
            <div class="metric-icon active">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.in_progress_count }}</div>
              <div class="metric-label">进行中</div>
            </div>
          </div>
          <div class="task-metric">
            <div class="metric-icon success">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                <polyline points="22 4 12 14.01 9 11.01"/>
              </svg>
            </div>
            <div class="metric-content">
              <div class="metric-value">{{ myStats.today_tasks }}</div>
              <div class="metric-label">今日完成</div>
            </div>
          </div>
        </div>
        <router-link to="/inspection/execute" class="tasks-action">
          <span>开始执行</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"/>
          </svg>
        </router-link>
      </div>

      <!-- 待办事项 -->
      <div class="dashboard-card todo-list">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <polyline points="10 9 9 9 8 9"/>
          </svg>
          待办事项
        </div>
        <div v-if="pendingTasks.length === 0" class="todo-empty">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <p>暂无待办任务</p>
        </div>
        <div v-else class="todo-items">
          <router-link
            v-for="task in pendingTasks.slice(0, 4)"
            :key="task.id"
            :to="`/inspection/execute/${task.id}`"
            class="todo-item"
          >
            <div class="todo-status"></div>
            <div class="todo-content">
              <div class="todo-title">{{ task.equipment_name || task.equipment_code }}</div>
              <div class="todo-time">{{ task.scheduled_date }}</div>
            </div>
            <div class="todo-tag">待点检</div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { equipmentApi, type EquipmentStatistics } from '@/api/equipment'
import {
  inspectionTaskApi,
  type InspectionStatistics,
  type MyTasksStatistics,
  type InspectionTask
} from '@/api/inspection'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const statistics = ref<EquipmentStatistics>({
  total: 0,
  running: 0,
  stopped: 0,
  maintenance: 0,
  scrapped: 0,
})

const inspectionStats = ref<InspectionStatistics>({
  total_tasks: 0,
  pending_tasks: 0,
  in_progress_tasks: 0,
  completed_tasks: 0,
  overdue_tasks: 0,
  today_completed: 0,
  completion_rate: 0,
})

const myStats = ref<{
  pending_count: number
  in_progress_count: number
  today_tasks: number
  loading: boolean
}>({
  pending_count: 0,
  in_progress_count: 0,
  today_tasks: 0,
  loading: true
})

const pendingTasks = ref<InspectionTask[]>([])
const currentTime = ref('')
const currentDate = ref('')

// 计算百分比
const runningPercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.running / statistics.value.total) * 100
})

const maintenancePercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.maintenance / statistics.value.total) * 100
})

const stoppedPercent = computed(() => {
  if (statistics.value.total === 0) return 0
  return (statistics.value.stopped / statistics.value.total) * 100
})

// 更新时间
function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  currentDate.value = now.toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' })
}

async function loadStatistics() {
  try {
    const response = await equipmentApi.getStatistics()
    statistics.value = response.data
  } catch (error) {
    ElMessage.error('获取设备统计数据失败')
  }
}

async function loadInspectionStats() {
  try {
    const response = await inspectionTaskApi.getStatistics()
    inspectionStats.value = response.data
  } catch (error) {
    console.error('获取点检统计失败', error)
  }
}

async function loadMyStats() {
  myStats.value.loading = true
  try {
    const response = await inspectionTaskApi.getMyStats()
    myStats.value = {
      ...response.data,
      loading: false
    }
  } catch (error) {
    console.error('获取我的任务统计失败', error)
    myStats.value.loading = false
  }
}

async function loadPendingTasks() {
  try {
    const response = await inspectionTaskApi.getMyTasks()
    if (Array.isArray(response.data)) {
      pendingTasks.value = response.data.filter(t => t.status === 'pending')
    } else if (response.data && (response.data as any).items) {
      pendingTasks.value = (response.data as any).items.filter((t: any) => t.status === 'pending')
    }
  } catch (error) {
    console.error('获取待办任务失败', error)
  }
}

let timeInterval: ReturnType<typeof setInterval>

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)
  loadStatistics()
  loadInspectionStats()
  loadMyStats()
  loadPendingTasks()
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
})
</script>

<style scoped>
.dashboard {
  max-width: 1360px;
  margin: 0 auto;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "SF Pro Display", "PingFang SC", sans-serif;
  color: #1d1d1f;
}

/* macOS 27 Golden Gate Fluid Welcome Banner */
.welcome-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding: 24px 28px;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(35px) saturate(180%);
  border-radius: 20px;
  box-shadow: 
    0 0 0 1px rgba(0, 0, 0, 0.05),
    0 4px 20px -2px rgba(10, 124, 255, 0.06);
  position: relative;
  overflow: hidden;
}

.welcome-title {
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.02em;
  margin: 0 0 4px 0;
}

.title-greeting {
  color: #86868b;
  font-weight: 500;
}

.title-name {
  color: #1d1d1f;
  font-weight: 700;
}

.welcome-subtitle {
  font-size: 13px;
  color: #86868b;
  margin: 0;
  letter-spacing: 0.01em;
}

.welcome-time {
  text-align: right;
}

.time-display {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 28px;
  font-weight: 700;
  color: #0a7cff;
  line-height: 1;
}

.date-display {
  font-size: 12px;
  font-weight: 500;
  color: #86868b;
  margin-top: 4px;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  position: relative;
  padding: 20px;
  background: #ffffff;
  border-radius: 18px;
  box-shadow: 
    0 0 0 1px rgba(0, 0, 0, 0.05),
    0 2px 8px rgba(0, 0, 0, 0.03);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 0 0 1px rgba(10, 124, 255, 0.2),
    0 8px 24px -4px rgba(0, 0, 0, 0.08);
}

.stat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.stat-icon {
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f4f4f7;
  border-radius: 10px;
  color: #0a7cff;
}

.stat-icon svg {
  width: 20px;
  height: 20px;
}

.stat-value {
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 30px;
  font-weight: 700;
  color: #1d1d1f;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.stat-label {
  font-size: 12.5px;
  font-weight: 550;
  color: #86868b;
  margin-top: 4px;
  margin-bottom: 12px;
}

.stat-bar {
  height: 4px;
  background: #f0f0f3;
  border-radius: 999px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 999px;
}

.stat-card[data-type="total"] .bar-fill { background: #0a7cff; }
.stat-card[data-type="running"] .bar-fill { background: #34c759; }
.stat-card[data-type="maintenance"] .bar-fill { background: #ff9500; }
.stat-card[data-type="stopped"] .bar-fill { background: #ff3b30; }

/* Dashboard Grid */
.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 24px;
}

.panel-card {
  background: #ffffff;
  border-radius: 20px;
  padding: 22px;
  box-shadow: 
    0 0 0 1px rgba(0, 0, 0, 0.05),
    0 2px 8px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18px;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 650;
  color: #1d1d1f;
  letter-spacing: -0.01em;
}

.panel-icon {
  width: 18px;
  height: 18px;
  color: #0a7cff;
}

.todo-list {
  grid-column: span 2;
}

.task-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-radius: 12px;
  background: #f8f9fa;
  margin-bottom: 8px;
  text-decoration: none;
  transition: all 0.15s;
}

.task-item:hover {
  background: #f0f4ff;
  transform: translateX(2px);
}

.task-name {
  font-size: 13.5px;
  font-weight: 550;
  color: #1d1d1f;
}

.todo-tag {
  padding: 3px 10px;
  background: rgba(10, 124, 255, 0.1);
  color: #0a7cff;
  font-size: 11.5px;
  font-weight: 600;
  border-radius: 999px;
}

/* 响应式 */
@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  .todo-list {
    grid-column: span 1;
  }
}

@media (max-width: 600px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
  .welcome-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
}
</style>
