<template>
  <div class="login-container">
    <!-- macOS 27 Golden Gate Fluid Aura Canvas -->
    <div class="macos-aura-bg">
      <div class="aura-orb aura-blue"></div>
      <div class="aura-orb aura-indigo"></div>
      <div class="aura-orb aura-peach"></div>
    </div>

    <!-- macOS 27 Window Stage -->
    <div class="macos-stage">
      <!-- Main Application Window Card -->
      <div class="macos-window-card">
        <!-- Window Traffic Light & Title Header -->
        <div class="window-chrome">
          <div class="traffic-lights">
            <span class="light red"></span>
            <span class="light yellow"></span>
            <span class="light green"></span>
          </div>
          <div class="window-title-chip">
            <svg viewBox="0 0 32 32" fill="none" class="title-chip-icon" style="width: 15px; height: 15px;">
              <path d="M16 5.5L25 10.5L16 15.5L7 10.5L16 5.5Z" fill="#0A7CFF" />
              <path d="M7 10.5L16 15.5V25.5L7 20.5V10.5Z" fill="#0055D6" />
              <path d="M16 15.5L25 10.5V20.5L16 25.5V15.5Z" fill="#2563EB" />
            </svg>
            <span>EAM-Agent · Enterprise Asset Brain</span>
          </div>
          <div class="window-actions-spacer">
            <span class="version-tag">v2.1 Pro</span>
          </div>
        </div>

        <!-- Window Body Split -->
        <div class="window-body">
          <!-- Left Hero Sidebar -->
          <div class="window-sidebar">
            <div class="app-identity">
              <div class="app-squircle-icon">
                <svg viewBox="0 0 32 32" fill="none" style="width: 36px; height: 36px;">
                  <defs>
                    <linearGradient id="login-top-grad" x1="16" y1="5" x2="16" y2="15" gradientUnits="userSpaceOnUse">
                      <stop stop-color="#0A7CFF" />
                      <stop offset="1" stop-color="#60A5FA" />
                    </linearGradient>
                    <linearGradient id="login-left-grad" x1="8" y1="12" x2="16" y2="25" gradientUnits="userSpaceOnUse">
                      <stop stop-color="#0055D6" />
                      <stop offset="1" stop-color="#1D4ED8" />
                    </linearGradient>
                    <linearGradient id="login-right-grad" x1="16" y1="15" x2="24" y2="25" gradientUnits="userSpaceOnUse">
                      <stop stop-color="#2563EB" />
                      <stop offset="1" stop-color="#1E40AF" />
                    </linearGradient>
                  </defs>
                  <!-- Isometric Asset Brain Cube -->
                  <path d="M16 5.5L25 10.5L16 15.5L7 10.5L16 5.5Z" fill="url(#login-top-grad)" />
                  <path d="M7 10.5L16 15.5V25.5L7 20.5V10.5Z" fill="url(#login-left-grad)" />
                  <path d="M16 15.5L25 10.5V20.5L16 25.5V15.5Z" fill="url(#login-right-grad)" />
                  <!-- Inner Cognitive Core -->
                  <circle cx="16" cy="14" r="2.8" fill="#FFFFFF" />
                  <circle cx="16" cy="14" r="4.6" stroke="#FFFFFF" stroke-opacity="0.8" stroke-width="0.8" stroke-dasharray="1.5 1.5" />
                </svg>
              </div>
              <h1 class="app-name">EAM-Agent</h1>
              <p class="app-tagline">新一代工业设备资产大脑与运维系统</p>
            </div>

            <div class="feature-pills">
              <div class="pill-item">
                <div class="pill-icon">⚡</div>
                <div class="pill-text">
                  <div class="pill-title">全闭环健康预测</div>
                  <div class="pill-desc">基于设备机台运行数据的动态 RUL 与预警</div>
                </div>
              </div>
              <div class="pill-item">
                <div class="pill-icon">🧠</div>
                <div class="pill-text">
                  <div class="pill-title">MCP 协议原生互联</div>
                  <div class="pill-desc">Claude Desktop、Cursor 与 Cline 智能直连</div>
                </div>
              </div>
              <div class="pill-item">
                <div class="pill-icon">🔒</div>
                <div class="pill-text">
                  <div class="pill-title">细粒度 RBAC & 租户</div>
                  <div class="pill-desc">端到端数据隔离与严苛工业审计流</div>
                </div>
              </div>
            </div>

            <div class="sidebar-dock-footer">
              <span class="indicator-dot"></span>
              <span>InstaBox 云端实例在线 · 响应延迟 &lt; 15ms</span>
            </div>
          </div>

          <!-- Right Interactive Pane -->
          <div class="window-content-pane">
            <div class="pane-header">
              <h2>{{ isLoginMode ? '欢迎登录' : '申请账号' }}</h2>
              <p>{{ isLoginMode ? '输入工号凭证访问工厂中枢系统' : '填写您的基本信息，管理员将实时审核' }}</p>
            </div>

            <!-- Login / Apply Form -->
            <el-form
              ref="formRef"
              :model="isLoginMode ? loginForm : applyForm"
              :rules="isLoginMode ? loginRules : applyRules"
              class="macos-glass-form"
            >
              <!-- Username -->
              <el-form-item prop="username">
                <div class="input-wrap">
                  <label class="input-label">账号 / 工号</label>
                  <el-input
                    v-model="(isLoginMode ? loginForm : applyForm).username"
                    placeholder="请输入工号或用户名"
                    size="large"
                    :prefix-icon="UserIcon"
                  />
                </div>
              </el-form-item>

              <!-- Apply Mode Extra Fields -->
              <template v-if="!isLoginMode">
                <el-form-item prop="name">
                  <div class="input-wrap">
                    <label class="input-label">真实姓名</label>
                    <el-input
                      v-model="applyForm.name"
                      placeholder="例如：张工"
                      size="large"
                    />
                  </div>
                </el-form-item>

                <el-form-item prop="role">
                  <div class="input-wrap">
                    <label class="input-label">申请岗位</label>
                    <el-select
                      v-model="applyForm.role"
                      placeholder="请选择岗位职责"
                      size="large"
                      style="width: 100%;"
                    >
                      <el-option label="操作工 (Operator)" value="operator" />
                      <el-option label="维修工 (Maintenance)" value="maintenance" />
                      <el-option label="设备工程师 (Engineer)" value="engineer" />
                      <el-option label="设备主管 (Supervisor)" value="supervisor" />
                    </el-select>
                  </div>
                </el-form-item>

                <el-form-item prop="phone">
                  <div class="input-wrap">
                    <label class="input-label">联系手机</label>
                    <el-input
                      v-model="applyForm.phone"
                      placeholder="选填，用于突发告警短信"
                      size="large"
                    />
                  </div>
                </el-form-item>
              </template>

              <!-- Password -->
              <el-form-item prop="password">
                <div class="input-wrap">
                  <label class="input-label">系统密码</label>
                  <el-input
                    v-model="(isLoginMode ? loginForm : applyForm).password"
                    type="password"
                    :placeholder="isLoginMode ? '请输入登录密码' : '初始密码（至少 6 位）'"
                    size="large"
                    show-password
                    :prefix-icon="LockIcon"
                    @keyup.enter="isLoginMode ? handleLogin() : handleApply()"
                  />
                </div>
              </el-form-item>

              <!-- Submit Button (macOS 27 Pill) -->
              <div class="action-btn-area">
                <button
                  type="button"
                  class="macos-primary-btn"
                  :disabled="loading"
                  @click="isLoginMode ? handleLogin() : handleApply()"
                >
                  <span v-if="loading" class="spinner"></span>
                  <span v-else>{{ isLoginMode ? '登 录 系 统' : '提 交 申 请' }}</span>
                </button>
              </div>
            </el-form>

            <!-- Test Credentials Card (Shoji Inset Box) -->
            <div class="shoji-inset-box">
              <div class="inset-header">
                <span class="inset-badge">快捷体验</span>
                <span class="inset-tip">预设环境测试凭据已就绪</span>
              </div>
              <div class="credential-row">
                <div class="cred-item">
                  <span class="cred-label">管理员</span>
                  <kbd>admin</kbd>
                </div>
                <div class="cred-item">
                  <span class="cred-label">密码</span>
                  <kbd>admin123</kbd>
                </div>
                <button type="button" class="quick-fill-btn" @click="fillAdmin">一键填入</button>
              </div>
            </div>

            <!-- Mode Switch Link -->
            <div class="pane-footer-switch">
              <span class="switch-prompt">{{ isLoginMode ? '首次使用没有账号？' : '已有员工账号？' }}</span>
              <button type="button" class="switch-link-btn" @click="toggleMode">
                {{ isLoginMode ? '在线申请加入' : '直接返回登录' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi, type ApplyAccountRequest } from '@/api/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isLoginMode = ref(true)

// SVG Icons as render functions
const UserIcon = () => h('svg', {
  viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.8',
  style: { width: '16px', height: '16px', color: '#86868b' }
}, [
  h('path', { d: 'M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2' }),
  h('circle', { cx: '12', cy: '7', r: '4' })
])

const LockIcon = () => h('svg', {
  viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '1.8',
  style: { width: '16px', height: '16px', color: '#86868b' }
}, [
  h('rect', { x: '3', y: '11', width: '18', height: '11', rx: '2', ry: '2' }),
  h('path', { d: 'M7 11V7a5 5 0 0 1 10 0v4' })
])

const loginForm = reactive({
  username: 'admin',
  password: 'admin123',
})

const applyForm = reactive({
  username: '',
  password: '',
  name: '',
  role: 'operator',
  phone: '',
})

const loginRules: FormRules = {
  username: [{ required: true, message: '请输入工号或用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const applyRules: FormRules = {
  username: [
    { required: true, message: '请输入工号或用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '长度需在 3 到 50 个字符之间', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码不能少于 6 个字符', trigger: 'blur' },
  ],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择岗位职责', trigger: 'change' }],
}

function fillAdmin() {
  loginForm.username = 'admin'
  loginForm.password = 'admin123'
  // ElMessage.success('已自动填入管理员凭据')
}

function toggleMode() {
  isLoginMode.value = !isLoginMode.value
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

async function handleLogin() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const response = await authStore.login(loginForm)
      // ElMessage.success('登录成功')

      if (response.must_change_password) {
        router.push('/change-password')
        return
      }

      const redirect = (route.query.redirect as string) || '/dashboard'
      router.push(redirect)
    } catch {
      // Handled by request interceptor
    } finally {
      loading.value = false
    }
  })
}

async function handleApply() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    try {
      const data: ApplyAccountRequest = {
        username: applyForm.username,
        password: applyForm.password,
        name: applyForm.name,
        role: applyForm.role,
        phone: applyForm.phone || undefined,
      }
      await authApi.applyAccount(data)
      ElMessage.success('申请提交成功，请联系系统管理员通过审核')
      Object.assign(applyForm, {
        username: '',
        password: '',
        name: '',
        role: 'operator',
        phone: '',
      })
      toggleMode()
    } catch {
      // Handled by request interceptor
    } finally {
      loading.value = false
    }
  })
}

onMounted(() => {
  loginForm.username = 'admin'
  loginForm.password = 'admin123'
})
</script>

<style scoped>
/* Fullscreen stage */
.login-container {
  min-height: 100vh;
  width: 100vw;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background-color: #f2f3f7;
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Text", "SF Pro Display", "PingFang SC", sans-serif;
  color: #1d1d1f;
  -webkit-font-smoothing: antialiased;
}

/* macOS 27 Golden Gate Aura Canvas */
.macos-aura-bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
  background:
    radial-gradient(900px 620px at 15% 10%, rgba(120, 164, 255, 0.45), transparent 65%),
    radial-gradient(800px 600px at 90% 25%, rgba(255, 170, 205, 0.35), transparent 60%),
    radial-gradient(900px 700px at 50% 100%, rgba(255, 206, 150, 0.35), transparent 65%),
    linear-gradient(160deg, #dce6f9 0%, #f4e4ed 50%, #fbf1e6 100%);
  filter: blur(40px);
}

.aura-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.6;
  animation: floatOrb 12s ease-in-out infinite alternate;
}

.aura-blue {
  width: 500px;
  height: 500px;
  background: #3b82f6;
  top: -100px;
  left: -100px;
}

.aura-indigo {
  width: 450px;
  height: 450px;
  background: #8b5cf6;
  bottom: -150px;
  right: -50px;
  animation-delay: -4s;
}

.aura-peach {
  width: 350px;
  height: 350px;
  background: #f43f5e;
  top: 30%;
  right: 25%;
  opacity: 0.25;
  animation-delay: -7s;
}

@keyframes floatOrb {
  0% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(40px, -30px) scale(1.08); }
  100% { transform: translate(-30px, 30px) scale(0.95); }
}

/* Stage & Floating Glass Window */
.macos-stage {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 960px;
  padding: 24px;
}

.macos-window-card {
  background: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(45px) saturate(190%);
  -webkit-backdrop-filter: blur(45px) saturate(190%);
  border-radius: 26px;
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.7),
    0 0 0 1.5px rgba(0, 0, 0, 0.06),
    0 24px 64px -12px rgba(15, 23, 42, 0.18),
    0 12px 28px -8px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

/* Window Chrome (Titlebar) */
.window-chrome {
  height: 48px;
  background: rgba(248, 249, 252, 0.65);
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 18px;
  user-select: none;
}

.traffic-lights {
  display: flex;
  gap: 8px;
  align-items: center;
  width: 70px;
}

.light {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  box-shadow: inset 0 0 0 0.5px rgba(0, 0, 0, 0.15);
}

.light.red { background: #ff5f57; }
.light.yellow { background: #febc2e; }
.light.green { background: #28c840; }

.window-title-chip {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 4px 12px;
  background: rgba(0, 0, 0, 0.035);
  border-radius: 999px;
  font-size: 12px;
  font-weight: 550;
  color: #48484a;
  letter-spacing: -0.01em;
}

.title-chip-icon {
  width: 14px;
  height: 14px;
  color: #0a7cff;
}

.window-actions-spacer {
  width: 70px;
  display: flex;
  justify-content: flex-end;
}

.version-tag {
  font-size: 11px;
  font-weight: 600;
  color: #86868b;
  background: rgba(0, 0, 0, 0.04);
  padding: 2px 8px;
  border-radius: 6px;
}

/* Window Body */
.window-body {
  display: grid;
  grid-template-columns: 380px 1fr;
  min-height: 520px;
}

/* Sidebar */
.window-sidebar {
  background: rgba(246, 247, 251, 0.5);
  border-right: 1px solid rgba(0, 0, 0, 0.06);
  padding: 40px 32px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.app-squircle-icon {
  width: 60px;
  height: 60px;
  border-radius: 16px;
  background: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 0 0 1px rgba(0, 0, 0, 0.06),
    0 8px 20px -4px rgba(10, 124, 255, 0.25);
  margin-bottom: 18px;
}

.app-squircle-icon svg {
  width: 34px;
  height: 34px;
}

.app-name {
  margin: 0 0 6px 0;
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: #1d1d1f;
}

.app-tagline {
  margin: 0;
  font-size: 13.5px;
  color: #636366;
  line-height: 1.45;
}

.feature-pills {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: 28px 0;
}

.pill-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.pill-icon {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.85);
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.pill-text {
  flex: 1;
}

.pill-title {
  font-size: 13px;
  font-weight: 600;
  color: #1d1d1f;
  margin-bottom: 2px;
}

.pill-desc {
  font-size: 12px;
  color: #86868b;
  line-height: 1.35;
}

.sidebar-dock-footer {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11.5px;
  color: #636366;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.indicator-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #34c759;
  box-shadow: 0 0 6px rgba(52, 199, 89, 0.8);
}

/* Content Pane */
.window-content-pane {
  padding: 40px 44px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  background: #ffffff;
}

.pane-header {
  margin-bottom: 26px;
}

.pane-header h2 {
  margin: 0 0 6px 0;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: #1d1d1f;
}

.pane-header p {
  margin: 0;
  font-size: 13.5px;
  color: #86868b;
}

.macos-glass-form {
  margin-bottom: 20px;
}

.input-wrap {
  width: 100%;
}

.input-label {
  display: block;
  font-size: 12px;
  font-weight: 600;
  color: #48484a;
  margin-bottom: 6px;
}

/* Custom styled inputs for macOS 27 */
:deep(.el-input__wrapper) {
  background: #f4f4f7 !important;
  border: 1px solid transparent !important;
  border-radius: 12px !important;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.04) !important;
  height: 44px;
  transition: all 0.2s ease;
}

:deep(.el-input__wrapper:hover) {
  background: #ebebf0 !important;
}

:deep(.el-input__wrapper.is-focus) {
  background: #ffffff !important;
  border-color: #0a7cff !important;
  box-shadow:
    0 0 0 3.5px rgba(10, 124, 255, 0.15),
    inset 0 1px 2px rgba(0, 0, 0, 0.04) !important;
}

:deep(.el-input__inner) {
  font-size: 14px;
  color: #1d1d1f;
  font-weight: 500;
}

.action-btn-area {
  margin-top: 10px;
}

/* macOS 27 Pill Button */
.macos-primary-btn {
  width: 100%;
  height: 46px;
  background: #0a7cff;
  color: #ffffff;
  border: none;
  border-radius: 999px;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: 0.02em;
  cursor: pointer;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.1),
    0 4px 14px -2px rgba(10, 124, 255, 0.4);
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.macos-primary-btn:hover {
  background: #006ce6;
  box-shadow:
    0 2px 4px rgba(0, 0, 0, 0.12),
    0 8px 20px -3px rgba(10, 124, 255, 0.55);
  transform: translateY(-0.5px);
}

.macos-primary-btn:active {
  transform: translateY(0.5px);
}

.macos-primary-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Shoji Inset Box */
.shoji-inset-box {
  background: #f7f7f9;
  border-radius: 14px;
  padding: 12px 16px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.inset-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.inset-badge {
  font-size: 11px;
  font-weight: 700;
  color: #0a7cff;
  background: rgba(10, 124, 255, 0.1);
  padding: 2px 8px;
  border-radius: 999px;
}

.inset-tip {
  font-size: 11.5px;
  color: #86868b;
}

.credential-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.cred-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.cred-label {
  font-size: 12px;
  color: #636366;
}

kbd {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 12px;
  font-weight: 600;
  color: #1d1d1f;
  background: #ffffff;
  border-radius: 6px;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.06),
    0 0 0 1px rgba(0, 0, 0, 0.08);
}

.quick-fill-btn {
  background: transparent;
  border: none;
  color: #0a7cff;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.15s;
}

.quick-fill-btn:hover {
  background: rgba(10, 124, 255, 0.08);
}

/* Switch Link */
.pane-footer-switch {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
}

.switch-prompt {
  color: #86868b;
}

.switch-link-btn {
  background: transparent;
  border: none;
  color: #0a7cff;
  font-weight: 600;
  cursor: pointer;
  padding: 0;
  text-decoration: underline;
  text-underline-offset: 3px;
}

.switch-link-btn:hover {
  color: #006ce6;
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Responsive */
@media (max-width: 820px) {
  .window-body {
    grid-template-columns: 1fr;
  }
  .window-sidebar {
    display: none;
  }
  .window-content-pane {
    padding: 30px 24px;
  }
}

/* ============================================
   macOS Dark Mode Adaptation for Login
   ============================================ */
:global([data-theme="dark"]) .login-container {
  background-color: #0d0e12;
  color: #f5f5f7;
}

:global([data-theme="dark"]) .macos-aura-bg {
  background:
    radial-gradient(900px 620px at 15% 10%, rgba(30, 64, 175, 0.35), transparent 65%),
    radial-gradient(800px 600px at 90% 25%, rgba(136, 19, 55, 0.25), transparent 60%),
    radial-gradient(900px 700px at 50% 100%, rgba(120, 53, 15, 0.25), transparent 65%),
    linear-gradient(160deg, #090d16 0%, #150e18 50%, #16100d 100%);
}

:global([data-theme="dark"]) .macos-window-card {
  background: rgba(28, 28, 34, 0.88);
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.12),
    0 24px 64px -12px rgba(0, 0, 0, 0.65),
    0 12px 28px -8px rgba(0, 0, 0, 0.45);
}

:global([data-theme="dark"]) .window-chrome {
  background: rgba(22, 22, 26, 0.75);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .window-title-chip {
  background: rgba(255, 255, 255, 0.06);
  color: #d1d1d6;
}

:global([data-theme="dark"]) .version-tag {
  color: #a1a1a6;
  background: rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .window-sidebar {
  background: rgba(20, 20, 24, 0.6);
  border-right: 1px solid rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .app-squircle-icon {
  background: #24242a;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.12);
}

:global([data-theme="dark"]) .app-name {
  color: #ffffff;
}

:global([data-theme="dark"]) .app-tagline {
  color: #a1a1a6;
}

:global([data-theme="dark"]) .pill-icon {
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .pill-title {
  color: #ffffff;
}

:global([data-theme="dark"]) .pill-desc {
  color: #8e8e93;
}

:global([data-theme="dark"]) .sidebar-dock-footer {
  color: #8e8e93;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .window-content-pane {
  background: #18181c;
}

:global([data-theme="dark"]) .pane-header h2 {
  color: #ffffff;
}

:global([data-theme="dark"]) .pane-header p {
  color: #8e8e93;
}

:global([data-theme="dark"]) .input-label {
  color: #d1d1d6;
}

:global([data-theme="dark"]) :deep(.el-input__wrapper) {
  background: #24242a !important;
  border-color: rgba(255, 255, 255, 0.08) !important;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

:global([data-theme="dark"]) :deep(.el-input__wrapper:hover) {
  background: #2c2c34 !important;
}

:global([data-theme="dark"]) :deep(.el-input__wrapper.is-focus) {
  background: #1e1e24 !important;
  border-color: #0a7cff !important;
  box-shadow:
    0 0 0 3.5px rgba(10, 124, 255, 0.25),
    inset 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

:global([data-theme="dark"]) :deep(.el-input__inner) {
  color: #ffffff !important;
}

:global([data-theme="dark"]) :deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.3) !important;
}

:global([data-theme="dark"]) .shoji-inset-box {
  background: #202026;
  border-color: rgba(255, 255, 255, 0.08);
}

:global([data-theme="dark"]) .inset-tip {
  color: #8e8e93;
}

:global([data-theme="dark"]) .cred-label {
  color: #a1a1a6;
}

:global([data-theme="dark"]) kbd {
  background: #2c2c34;
  color: #ffffff;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.12);
}

:global([data-theme="dark"]) .switch-prompt {
  color: #8e8e93;
}

</style>

<style>
[data-theme="dark"] .login-container,
.dark .login-container {
  background-color: #0c0d12 !important;
  color: #f5f5f7 !important;
}

[data-theme="dark"] .macos-aura-bg,
.dark .macos-aura-bg {
  background:
    radial-gradient(900px 620px at 15% 10%, rgba(30, 64, 175, 0.4), transparent 65%),
    radial-gradient(800px 600px at 90% 25%, rgba(136, 19, 55, 0.3), transparent 60%),
    radial-gradient(900px 700px at 50% 100%, rgba(120, 53, 15, 0.3), transparent 65%),
    linear-gradient(160deg, #090d16 0%, #150e18 50%, #16100d 100%) !important;
}

[data-theme="dark"] .macos-window-card,
.dark .macos-window-card {
  background: rgba(28, 28, 34, 0.9) !important;
  box-shadow:
    0 0 0 1px rgba(255, 255, 255, 0.12),
    0 24px 64px -12px rgba(0, 0, 0, 0.75) !important;
}

[data-theme="dark"] .window-chrome,
.dark .window-chrome {
  background: rgba(22, 22, 26, 0.8) !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .window-title-chip,
.dark .window-title-chip {
  background: rgba(255, 255, 255, 0.08) !important;
  color: #d1d1d6 !important;
}

[data-theme="dark"] .version-tag,
.dark .version-tag {
  color: #a1a1a6 !important;
  background: rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .window-sidebar,
.dark .window-sidebar {
  background: rgba(20, 20, 24, 0.65) !important;
  border-right: 1px solid rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .app-squircle-icon,
.dark .app-squircle-icon {
  background: #24242a !important;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.12) !important;
}

[data-theme="dark"] .app-name,
.dark .app-name {
  color: #ffffff !important;
}

[data-theme="dark"] .app-tagline,
.dark .app-tagline {
  color: #a1a1a6 !important;
}

[data-theme="dark"] .pill-icon,
.dark .pill-icon {
  background: rgba(255, 255, 255, 0.08) !important;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .pill-title,
.dark .pill-title {
  color: #ffffff !important;
}

[data-theme="dark"] .pill-desc,
.dark .pill-desc {
  color: #8e8e93 !important;
}

[data-theme="dark"] .sidebar-dock-footer,
.dark .sidebar-dock-footer {
  color: #8e8e93 !important;
  border-top: 1px solid rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .window-content-pane,
.dark .window-content-pane {
  background: #18181c !important;
}

[data-theme="dark"] .pane-header h2,
.dark .pane-header h2 {
  color: #ffffff !important;
}

[data-theme="dark"] .pane-header p,
.dark .pane-header p {
  color: #8e8e93 !important;
}

[data-theme="dark"] .input-label,
.dark .input-label {
  color: #d1d1d6 !important;
}

[data-theme="dark"] .el-form-item,
.dark .el-form-item {
  background: transparent !important;
}

[data-theme="dark"] .login-container .el-input__wrapper,
.dark .login-container .el-input__wrapper {
  background: #24242a !important;
  border-color: rgba(255, 255, 255, 0.08) !important;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

[data-theme="dark"] .login-container .el-input__wrapper:hover,
.dark .login-container .el-input__wrapper:hover {
  background: #2c2c34 !important;
}

[data-theme="dark"] .login-container .el-input__wrapper.is-focus,
.dark .login-container .el-input__wrapper.is-focus {
  background: #1e1e24 !important;
  border-color: #0a7cff !important;
  box-shadow:
    0 0 0 3.5px rgba(10, 124, 255, 0.25),
    inset 0 1px 2px rgba(0, 0, 0, 0.3) !important;
}

[data-theme="dark"] .login-container .el-input__inner,
.dark .login-container .el-input__inner {
  color: #ffffff !important;
}

[data-theme="dark"] .login-container .el-input__inner::placeholder,
.dark .login-container .el-input__inner::placeholder {
  color: rgba(255, 255, 255, 0.3) !important;
}

[data-theme="dark"] .shoji-inset-box,
.dark .shoji-inset-box {
  background: #202026 !important;
  border-color: rgba(255, 255, 255, 0.08) !important;
}

[data-theme="dark"] .inset-tip,
.dark .inset-tip {
  color: #8e8e93 !important;
}

[data-theme="dark"] .cred-label,
.dark .cred-label {
  color: #a1a1a6 !important;
}

[data-theme="dark"] .login-container kbd,
.dark .login-container kbd {
  background: #2c2c34 !important;
  color: #ffffff !important;
  box-shadow:
    0 1px 2px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.12) !important;
}

[data-theme="dark"] .switch-prompt,
.dark .switch-prompt {
  color: #8e8e93 !important;
}
</style>
