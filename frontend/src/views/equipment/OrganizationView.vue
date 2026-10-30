<template>
  <div class="organization-view">
    <el-card class="org-main-card" shadow="never">
      <template #header>
        <div class="org-header-flex">
          <div>
            <h3 class="org-title">组织架构与制造层级</h3>
            <span class="org-subtitle">统一管理集团基地 (HQ)、制造工厂 (Factory) 与作业车间 (Workshop)</span>
          </div>
          <div>
            <el-button v-if="activeTab === 'bases'" type="primary" size="small" @click="showBaseDialog = true">
              <el-icon><Plus /></el-icon> 新增集团基地
            </el-button>
            <el-button v-else-if="activeTab === 'factories'" type="primary" size="small" @click="showFactoryDialog = true">
              <el-icon><Plus /></el-icon> 新增制造工厂
            </el-button>
            <el-button v-else type="primary" size="small" @click="showWorkshopDialog = true">
              <el-icon><Plus /></el-icon> 新增作业车间
            </el-button>
          </div>
        </div>
      </template>

      <el-tabs v-model="activeTab" class="org-tabs">
        <!-- 基地 Tab -->
        <el-tab-pane :label="'集团基地 (' + bases.length + ')'" name="bases">
          <div class="tab-pane-content">
            <el-table :data="bases" border stripe style="width: 100%">
              <el-table-column prop="code" label="基地编号" width="160" />
              <el-table-column prop="name" label="基地名称" min-width="200" />
              <el-table-column label="操作" width="160" align="center" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" link type="primary" @click="editBase(row)">编辑</el-button>
                  <el-button size="small" link type="danger" @click="deleteBase(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- 工厂 Tab -->
        <el-tab-pane :label="'制造工厂 (' + factories.length + ')'" name="factories">
          <div class="tab-pane-content">
            <el-table :data="factories" border stripe style="width: 100%">
              <el-table-column prop="code" label="工厂编号" width="160" />
              <el-table-column prop="name" label="工厂名称" min-width="200" />
              <el-table-column prop="base_name" label="归属集团基地" min-width="180">
                <template #default="{ row }">
                  <el-tag size="small" type="info">{{ row.base_name || '未关联' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="160" align="center" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" link type="primary" @click="editFactory(row)">编辑</el-button>
                  <el-button size="small" link type="danger" @click="deleteFactory(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <!-- 车间 Tab -->
        <el-tab-pane :label="'作业车间 (' + workshops.length + ')'" name="workshops">
          <div class="tab-pane-content">
            <el-table :data="workshops" border stripe style="width: 100%">
              <el-table-column prop="code" label="车间编号" width="160" />
              <el-table-column prop="name" label="车间名称" min-width="200" />
              <el-table-column prop="factory_name" label="归属工厂" min-width="180">
                <template #default="{ row }">
                  <el-tag size="small" type="success">{{ row.factory_name || '未关联' }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="160" align="center" fixed="right">
                <template #default="{ row }">
                  <el-button size="small" link type="primary" @click="editWorkshop(row)">编辑</el-button>
                  <el-button size="small" link type="danger" @click="deleteWorkshop(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Base Dialog -->
    <el-dialog v-model="showBaseDialog" :title="editingBase ? '编辑基地' : '新增基地'" width="500px">
      <el-form ref="baseFormRef" :model="baseForm" :rules="baseRules" label-width="100px">
        <el-form-item label="基地编号" prop="code">
          <el-input v-model="baseForm.code" placeholder="如: BASE01" />
        </el-form-item>
        <el-form-item label="基地名称" prop="name">
          <el-input v-model="baseForm.name" placeholder="如: 华东基地" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBaseDialog = false">取消</el-button>
        <el-button type="primary" @click="saveBase">确定</el-button>
      </template>
    </el-dialog>

    <!-- Factory Dialog -->
    <el-dialog v-model="showFactoryDialog" :title="editingFactory ? '编辑工厂' : '新增工厂'" width="500px">
      <el-form ref="factoryFormRef" :model="factoryForm" :rules="factoryRules" label-width="100px">
        <el-form-item label="所属基地" prop="base_id">
          <el-select v-model="factoryForm.base_id" placeholder="请选择基地" style="width: 100%">
            <el-option v-for="base in bases" :key="base.id" :label="base.name" :value="base.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="工厂编号" prop="code">
          <el-input v-model="factoryForm.code" placeholder="如: FAC01" />
        </el-form-item>
        <el-form-item label="工厂名称" prop="name">
          <el-input v-model="factoryForm.name" placeholder="如: 苏州工厂" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFactoryDialog = false">取消</el-button>
        <el-button type="primary" @click="saveFactory">确定</el-button>
      </template>
    </el-dialog>

    <!-- Workshop Dialog -->
    <el-dialog v-model="showWorkshopDialog" :title="editingWorkshop ? '编辑车间' : '新增车间'" width="500px">
      <el-form ref="workshopFormRef" :model="workshopForm" :rules="workshopRules" label-width="100px">
        <el-form-item label="所属工厂" prop="factory_id">
          <el-select v-model="workshopForm.factory_id" placeholder="请选择工厂" style="width: 100%">
            <el-option v-for="factory in factories" :key="factory.id" :label="factory.name" :value="factory.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="车间编号" prop="code">
          <el-input v-model="workshopForm.code" placeholder="如: WS01" />
        </el-form-item>
        <el-form-item label="车间名称" prop="name">
          <el-input v-model="workshopForm.name" placeholder="如: 机加车间" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWorkshopDialog = false">取消</el-button>
        <el-button type="primary" @click="saveWorkshop">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { orgApi, type Base, type Factory, type Workshop } from '@/api/equipment'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'

const activeTab = ref('bases')
const bases = ref<Base[]>([])
const factories = ref<Factory[]>([])
const workshops = ref<Workshop[]>([])

// Base
const showBaseDialog = ref(false)
const editingBase = ref<Base | null>(null)
const baseFormRef = ref<FormInstance>()
const baseForm = reactive({ code: '', name: '' })
const baseRules: FormRules = {
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

// Factory
const showFactoryDialog = ref(false)
const editingFactory = ref<Factory | null>(null)
const factoryFormRef = ref<FormInstance>()
const factoryForm = reactive({ base_id: undefined as number | undefined, code: '', name: '' })
const factoryRules: FormRules = {
  base_id: [{ required: true, message: '请选择基地', trigger: 'change' }],
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

// Workshop
const showWorkshopDialog = ref(false)
const editingWorkshop = ref<Workshop | null>(null)
const workshopFormRef = ref<FormInstance>()
const workshopForm = reactive({ factory_id: undefined as number | undefined, code: '', name: '' })
const workshopRules: FormRules = {
  factory_id: [{ required: true, message: '请选择工厂', trigger: 'change' }],
  code: [{ required: true, message: '请输入编号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
}

async function loadBases() {
  try {
    const response = await orgApi.getBases()
    bases.value = response.data
  } catch (error) {
    // Error handled by interceptor
  }
}

async function loadFactories() {
  try {
    const response = await orgApi.getFactories()
    factories.value = response.data
  } catch (error) {
    // Error handled by interceptor
  }
}

async function loadWorkshops() {
  try {
    const response = await orgApi.getWorkshops()
    workshops.value = response.data
  } catch (error) {
    // Error handled by interceptor
  }
}

// Base CRUD
function editBase(row: Base) {
  editingBase.value = row
  baseForm.code = row.code
  baseForm.name = row.name
  showBaseDialog.value = true
}

async function saveBase() {
  if (!baseFormRef.value) return

  await baseFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingBase.value) {
        await orgApi.updateBase(editingBase.value.id, baseForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createBase(baseForm)
        ElMessage.success('创建成功')
      }

      showBaseDialog.value = false
      resetBaseForm()
      loadBases()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetBaseForm() {
  editingBase.value = null
  baseForm.code = ''
  baseForm.name = ''
  baseFormRef.value?.clearValidate()
}

async function deleteBase(row: Base) {
  try {
    await ElMessageBox.confirm(`确定要删除基地 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteBase(row.id)
    ElMessage.success('删除成功')
    loadBases()
    loadFactories()
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

// Factory CRUD
function editFactory(row: Factory) {
  editingFactory.value = row
  factoryForm.base_id = row.base_id
  factoryForm.code = row.code
  factoryForm.name = row.name
  showFactoryDialog.value = true
}

async function saveFactory() {
  if (!factoryFormRef.value) return

  await factoryFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingFactory.value) {
        await orgApi.updateFactory(editingFactory.value.id, factoryForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createFactory(factoryForm)
        ElMessage.success('创建成功')
      }

      showFactoryDialog.value = false
      resetFactoryForm()
      loadFactories()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetFactoryForm() {
  editingFactory.value = null
  factoryForm.base_id = undefined
  factoryForm.code = ''
  factoryForm.name = ''
  factoryFormRef.value?.clearValidate()
}

async function deleteFactory(row: Factory) {
  try {
    await ElMessageBox.confirm(`确定要删除工厂 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteFactory(row.id)
    ElMessage.success('删除成功')
    loadFactories()
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

// Workshop CRUD
function editWorkshop(row: Workshop) {
  editingWorkshop.value = row
  workshopForm.factory_id = row.factory_id
  workshopForm.code = row.code
  workshopForm.name = row.name
  showWorkshopDialog.value = true
}

async function saveWorkshop() {
  if (!workshopFormRef.value) return

  await workshopFormRef.value.validate(async (valid) => {
    if (!valid) return

    try {
      if (editingWorkshop.value) {
        await orgApi.updateWorkshop(editingWorkshop.value.id, workshopForm)
        ElMessage.success('更新成功')
      } else {
        await orgApi.createWorkshop(workshopForm)
        ElMessage.success('创建成功')
      }

      showWorkshopDialog.value = false
      resetWorkshopForm()
      loadWorkshops()
    } catch (error) {
      // Error handled by interceptor
    }
  })
}

function resetWorkshopForm() {
  editingWorkshop.value = null
  workshopForm.factory_id = undefined
  workshopForm.code = ''
  workshopForm.name = ''
  workshopFormRef.value?.clearValidate()
}

async function deleteWorkshop(row: Workshop) {
  try {
    await ElMessageBox.confirm(`确定要删除车间 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    await orgApi.deleteWorkshop(row.id)
    ElMessage.success('删除成功')
    loadWorkshops()
  } catch {
    // User cancelled
  }
}

onMounted(() => {
  loadBases()
  loadFactories()
  loadWorkshops()
})
</script>

<style scoped>
.organization {
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>

<style scoped>
.organization-view {
  max-width: 1360px;
  margin: 0 auto;
}

.org-main-card {
  border-radius: 18px !important;
  background: var(--color-bg-card) !important;
  border: 1px solid var(--color-border) !important;
  min-height: 600px;
}

.org-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.org-title {
  margin: 0;
  font-size: 16px;
  font-weight: 650;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.org-subtitle {
  font-size: 12.5px;
  color: var(--color-text-tertiary);
}

.tab-pane-content {
  padding-top: 10px;
}
</style>
