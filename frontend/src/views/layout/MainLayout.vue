<template>
  <div class="main-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <!-- Logo 区域 -->
      <div class="sidebar-header">
        <div class="logo-container">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
              <!-- 外层工业六角框架 -->
              <path d="M12 2l9 5v10l-9 5-9-5V7l9-5z" />
              <!-- 核心齿轮 -->
              <circle cx="12" cy="12" r="5" />
              <!-- 智能运维脉冲线 (代表实时监测与诊断) -->
              <path d="M9.5 12h1l1-2 1.5 4 1-2h1" stroke-width="1.2" />
            </svg>
          </div>
          <transition name="logo-text">
            <span v-show="!sidebarCollapsed" class="logo-text">EAM-Agent</span>
          </transition>
        </div>
        <!-- 折叠按钮 - 固定在右上角 -->
        <button class="sidebar-collapse-btn" @click="toggleSidebar" title="折叠/展开侧边栏">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 14px; height: 14px;">
            <polyline points="15 18 9 12 15 6"/>
          </svg>
        </button>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">基础功能</div>
          <router-link to="/dashboard" class="nav-item" :class="{ active: isActive('/dashboard') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="7" height="7" rx="1"/>
                <rect x="14" y="3" width="7" height="7" rx="1"/>
                <rect x="3" y="14" width="7" height="7" rx="1"/>
                <rect x="14" y="14" width="7" height="7" rx="1"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">首页</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>

          <router-link to="/equipment" class="nav-item" :class="{ active: isActive('/equipment') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                <polyline points="3.27 6.96 12 12.01 20.73 6.96"/>
                <line x1="12" y1="22.08" x2="12" y2="12"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">设备台账</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>

          <div v-if="canManageOrg" class="nav-group" :class="{ active: isGroupActive('/organization') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('organization')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M3 21h18"/>
                  <rect x="5" y="3" width="4" height="14" rx="1"/>
                  <rect x="15" y="8" width="4" height="9" rx="1"/>
                  <path d="M10 21v-4"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">组织架构</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.organization }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.organization" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/organization" class="nav-subitem" :class="{ active: isActive('/organization') && !isActive('/organization/users') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">组织管理</span>
                </router-link>
                <router-link to="/organization/users" class="nav-subitem" :class="{ active: isActive('/organization/users') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">人员管理</span>
                </router-link>
              </div>
            </transition>
          </div>

          <router-link v-if="!canManageOrg" to="/organization" class="nav-item" :class="{ active: isActive('/organization') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 21h18"/>
                <rect x="5" y="3" width="4" height="14" rx="1"/>
                <rect x="15" y="8" width="4" height="9" rx="1"/>
                <path d="M10 21v-4"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">组织架构</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>

        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">设备维护</div>

          <!-- 点检管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/inspection') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('inspection')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                  <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">点检管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.inspection }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.inspection" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/inspection/tasks" class="nav-subitem" :class="{ active: isActive('/inspection/tasks') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">点检任务</span>
                </router-link>
                <router-link v-if="canManageOrg" to="/inspection/templates" class="nav-subitem" :class="{ active: isActive('/inspection/templates') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">模板配置</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- 维修管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/repair') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('repair')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">维修管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.repair }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.repair" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link to="/repair/orders" class="nav-subitem" :class="{ active: isActive('/repair/orders') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">维修工单</span>
                </router-link>
              </div>
            </transition>
          </div>

          <!-- 保养管理 -->
          <div class="nav-group" :class="{ active: isGroupActive('/maintenance') }">
            <button class="nav-item nav-group-toggle" @click="toggleGroup('maintenance')">
              <div class="nav-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
              </div>
              <transition name="nav-text">
                <span v-show="!sidebarCollapsed" class="nav-text">保养管理</span>
              </transition>
              <transition name="arrow">
                <svg v-show="!sidebarCollapsed" class="nav-arrow" :class="{ expanded: expandedGroups.maintenance }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"/>
                </svg>
              </transition>
            </button>
            <transition name="submenu">
              <div v-show="sidebarCollapsed || expandedGroups.maintenance" class="nav-submenu" :class="{ 'always-show': sidebarCollapsed }">
                <router-link v-if="canManageOrg" to="/maintenance/plans" class="nav-subitem" :class="{ active: isActive('/maintenance/plans') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">计划配置</span>
                </router-link>
                <router-link to="/maintenance/tasks" class="nav-subitem" :class="{ active: isActive('/maintenance/tasks') }">
                  <span class="submenu-dot"></span>
                  <span class="submenu-text">任务管理</span>
                </router-link>
              </div>
            </transition>
          </div>

          <router-link to="/spareparts" class="nav-item" :class="{ active: isActive('/spareparts') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">备件管理</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>

        <div class="nav-section">
          <div v-if="!sidebarCollapsed" class="nav-section-title">数据分析</div>
          <router-link to="/agent/assistant" class="nav-item" :class="{ active: isActive('/agent/assistant') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">管理助手</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
          <router-link to="/agent/integration" class="nav-item" :class="{ active: isActive('/agent/integration') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101" />
                <path d="M10.172 13.828a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.102 1.101" />
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">Agent 集成</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
          <router-link to="/analytics" class="nav-item" :class="{ active: isActive('/analytics') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <line x1="18" y1="20" x2="18" y2="10"/>
                <line x1="12" y1="20" x2="12" y2="4"/>
                <line x1="6" y1="20" x2="6" y2="14"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">统计分析</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
          <router-link to="/knowledge" class="nav-item" :class="{ active: isActive('/knowledge') }">
            <div class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/>
                <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/>
              </svg>
            </div>
            <transition name="nav-text">
              <span v-show="!sidebarCollapsed" class="nav-text">知识库</span>
            </transition>
            <div v-if="!sidebarCollapsed" class="nav-indicator"></div>
          </router-link>
        </div>
      </nav>

      <!-- 侧边栏底部 -->
      <div class="sidebar-footer">
        <!-- Footer content if needed -->
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 顶部栏 -->
      <header class="header">
        <div class="header-left">
          <button class="mobile-menu-btn" @click="toggleSidebar">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <line x1="3" y1="12" x2="21" y2="12"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
              <line x1="3" y1="18" x2="21" y2="18"/>
            </svg>
          </button>
          <div class="breadcrumb">
            <span class="breadcrumb-current">{{ currentTitle }}</span>
          </div>
        </div>

        <div class="header-right">
          <!-- 主题切换 -->
          <button class="theme-toggle" @click="themeStore.toggleTheme()" :title="themeStore.theme === 'dark' ? '切换到浅色模式' : '切换到深色模式'">
            <svg v-if="themeStore.theme === 'dark'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="5"/>
              <line x1="12" y1="1" x2="12" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="23"/>
              <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
            </svg>
          </button>

          <!-- 通知 -->
          <button class="header-icon-btn" @click="handleNotificationClick" title="待处理申请">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
              <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
            </svg>
            <span v-if="pendingCount > 0" class="notification-badge">{{ pendingCount }}</span>
          </button>
          <!-- 用户信息 -->
          <div class="user-section">
            <div class="user-avatar">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
            </div>
            <div class="user-info">
              <span class="user-name">{{ authStore.userName }}</span>
              <span v-if="roleText !== authStore.userName" class="user-role">{{ roleText }}</span>
            </div>
          </div>

          <!-- 退出按钮 -->
          <button class="logout-btn" @click="handleLogout" title="退出登录">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
              <polyline points="16 17 21 12 16 7"/>
              <line x1="21" y1="12" x2="9" y2="12"/>
            </svg>
          </button>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- 移动端遮罩 -->
    <div
      v-if="!sidebarCollapsed"
      class="sidebar-overlay"
      @click="toggleSidebar"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import { ElMessageBox } from 'element-plus'
import { userApi } from '@/api/user'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const themeStore = useThemeStore()

const sidebarCollapsed = ref(false)
const pendingCount = ref(0)
const expandedGroups = ref({
  organization: true,
  inspection: true,
  repair: false,
  maintenance: false,
})

const canManageOrg = computed(() => authStore.hasRole('admin', 'engineer'))
const isAdmin = computed(() => authStore.hasRole('admin'))

const roleTextMap: Record<string, string> = {
  admin: '系统管理员',
  supervisor: '设备主管',
  engineer: '设备工程师',
  maintenance: '维修工',
  operator: '操作工',
}

const roleText = computed(() => roleTextMap[authStore.userRole] || '未知')

const currentTitle = computed(() => route.meta?.title as string || '首页')

function isActive(path: string): boolean {
  return route.path === path || route.path.startsWith(path + '/')
}

function isGroupActive(path: string): boolean {
  return route.path.startsWith(path + '/')
}

function toggleGroup(group: keyof typeof expandedGroups) {
  if (sidebarCollapsed.value) return
  expandedGroups.value[group] = !expandedGroups.value[group]
}

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function handleNotificationClick() {
  if (isAdmin.value) {
    router.push('/organization/users')
  }
}

async function fetchPendingCount() {
  if (isAdmin.value) {
    try {
      const response = await userApi.getPendingApplications()
      pendingCount.value = response.data.length
    } catch (error) {
      console.error('Failed to fetch pending applications count:', error)
    }
  }
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
    authStore.logout()
    router.push('/login')
  } catch {
    // User cancelled
  }
}

// 响应式处理
function handleResize() {
  if (window.innerWidth < 1024) {
    sidebarCollapsed.value = true
  }
}

onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
  fetchPendingCount()
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
</script>

<style>
/* Global overrides */
html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  overflow-x: hidden;
}

* {
  box-sizing: border-box;
}

/* Custom scrollbar for main content areas only */
.main-content::-webkit-scrollbar,
.el-aside::-webkit-scrollbar {
  width: 6px;
}

.main-content::-webkit-scrollbar-thumb,
.el-aside::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.15);
  border-radius: 3px;
}

.main-content::-webkit-scrollbar-thumb:hover,
.el-aside::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.25);
}

/* Visible focus styles for accessibility */
*:focus-visible {
  outline: 2px solid var(--color-primary, #409eff);
  outline-offset: 2px;
}
</style>

<style scoped>
/* 主布局 */
.main-layout {
  display: flex;
  min-height: 100vh;
  height: 100vh;
  background: var(--color-bg-primary);
  border: none !important;
  margin: 0;
  padding: 0;
  width: 100%;
  overflow: hidden;
}

/* 侧边栏 */
/* macOS 27 Golden Gate Refined Liquid Glass Sidebar */
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 250px;
  height: 100vh;
  background: rgba(246, 247, 251, 0.82);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border-right: 1px solid rgba(0, 0, 0, 0.07);
  display: flex;
  flex-direction: column;
  z-index: 100;
  transition: transform var(--transition-base), width var(--transition-base);
  box-sizing: border-box;
  overflow: hidden;
}

.sidebar.collapsed {
  width: 72px;
}

.sidebar-header {
  padding: 18px 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.sidebar-collapse-btn {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.04);
  border: none;
  border-radius: 7px;
  color: #86868b;
  cursor: pointer;
  transition: all 0.15s;
  padding: 0;
}

.sidebar-collapse-btn:hover {
  background: rgba(0, 0, 0, 0.08);
  color: #1d1d1f;
}

.sidebar-collapse-btn svg {
  width: 14px;
  height: 14px;
  transition: transform 0.2s;
}

.sidebar.collapsed .sidebar-collapse-btn svg {
  transform: rotate(180deg);
}

.logo-icon {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  border-radius: 9px;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.06), 0 2px 8px rgba(10, 124, 255, 0.2);
  color: #0a7cff;
  flex-shrink: 0;
}

.logo-icon svg {
  width: 20px;
  height: 20px;
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #1d1d1f;
}

/* 侧边栏导航 */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 14px 10px;
  scrollbar-width: none !important;
}

.sidebar-nav::-webkit-scrollbar {
  display: none !important;
}

.nav-section {
  margin-bottom: 20px;
}

.nav-section:last-child {
  margin-bottom: 0;
}

.nav-section-title {
  font-size: 11px;
  font-weight: 650;
  color: #86868b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0 10px 6px;
}

/* macOS 27 Floating Highlight Item */
.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 10px;
  margin: 1.5px 0;
  border-radius: 8px;
  font-weight: 500;
  color: #48484a;
  text-decoration: none;
  transition: all 0.15s cubic-bezier(0.16, 1, 0.3, 1);
  position: relative;
  cursor: pointer;
  background: transparent;
  border: none;
  width: 100%;
  text-align: left;
}

.nav-item:hover {
  background: rgba(0, 0, 0, 0.045);
  color: #1d1d1f;
}

.nav-item.active {
  background: #0a7cff;
  color: #ffffff !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08), 0 4px 12px -2px rgba(10, 124, 255, 0.35);
}

.nav-item.active .nav-icon {
  color: #ffffff !important;
}

.nav-indicator {
  display: none;
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #636366;
  transition: color 0.15s;
}

.nav-item:hover .nav-icon {
  color: #1d1d1f;
}

.nav-icon svg {
  width: 17px;
  height: 17px;
}

.nav-text {
  flex: 1;
  font-size: 13px;
  font-weight: 550;
  letter-spacing: -0.01em;
}

.nav-arrow {
  width: 14px;
  height: 14px;
  transition: transform 0.2s;
  color: #86868b;
}

.nav-item.active .nav-arrow {
  color: rgba(255, 255, 255, 0.85);
}

.nav-arrow.expanded {
  transform: rotate(180deg);
}

/* 子菜单 */
.nav-group {
  margin: 1.5px 0;
  position: relative;
}

.nav-submenu {
  padding-left: 28px;
  overflow: hidden;
}

.sidebar.collapsed .nav-submenu.always-show {
  position: absolute;
  left: 72px;
  top: 0;
  min-width: 180px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(25px);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  padding: 6px;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  z-index: 50;
  opacity: 0;
  visibility: hidden;
  pointer-events: none;
  transition: opacity 0.2s, visibility 0.2s;
}

.sidebar.collapsed .nav-group:hover .nav-submenu.always-show {
  opacity: 1;
  visibility: visible;
  pointer-events: auto;
}

.nav-subitem {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  margin: 1.5px 0;
  border-radius: 6px;
  color: #636366;
  text-decoration: none;
  font-size: 12.5px;
  font-weight: 500;
  transition: all 0.15s;
}

.nav-subitem:hover {
  color: #1d1d1f;
  background: rgba(0, 0, 0, 0.04);
}

.nav-subitem.active {
  color: #0a7cff;
  font-weight: 600;
  background: rgba(10, 124, 255, 0.08);
}

.submenu-dot {
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.6;
}

.nav-subitem.active .submenu-dot {
  opacity: 1;
  box-shadow: 0 0 6px currentColor;
}

/* 主内容区 */
.main-content {
  flex: 1;
  margin-left: 250px;
  display: flex;
  flex-direction: column;
  height: 100vh;
  transition: margin-left var(--transition-base);
  box-sizing: border-box;
  background: #fbfbfd;
  overflow: hidden;
}

.sidebar.collapsed + .main-content {
  margin-left: 72px;
}

/* 顶部栏 - macOS 27 Header Bar */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  height: 52px;
  background: rgba(251, 251, 253, 0.82);
  backdrop-filter: blur(25px) saturate(180%);
  -webkit-backdrop-filter: blur(25px) saturate(180%);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.breadcrumb-current {
  color: #1d1d1f;
  font-size: 15px;
  font-weight: 650;
  letter-spacing: -0.02em;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon-btn, .theme-toggle, .logout-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.035);
  border: 1px solid rgba(0, 0, 0, 0.04);
  border-radius: 8px;
  color: #636366;
  cursor: pointer;
  transition: all 0.15s;
  padding: 0;
}

.header-icon-btn:hover, .theme-toggle:hover {
  background: rgba(0, 0, 0, 0.07);
  color: #1d1d1f;
}

.logout-btn:hover {
  background: rgba(255, 59, 48, 0.1);
  color: #ff3b30;
  border-color: rgba(255, 59, 48, 0.2);
}

.user-section {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 10px 4px 4px;
  background: rgba(0, 0, 0, 0.035);
  border: 1px solid rgba(0, 0, 0, 0.04);
  border-radius: 999px;
  cursor: pointer;
  transition: all 0.15s;
}

.user-section:hover {
  background: rgba(0, 0, 0, 0.07);
}

.user-avatar {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #0a7cff;
  border-radius: 50%;
  color: #ffffff;
}

.user-avatar svg {
  width: 14px;
  height: 14px;
}

.user-name {
  font-size: 12.5px;
  font-weight: 600;
  color: #1d1d1f;
}

.user-role {
  font-size: 11px;
  color: #86868b;
  margin-left: 2px;
}

/* 内容区 */
.content {
  flex: 1;
  padding: 24px 32px;
  overflow-y: auto;
  overflow-x: hidden;
  scrollbar-width: thin;
}
</style>
