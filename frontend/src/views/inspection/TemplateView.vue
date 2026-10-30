<template>
  <div class="template-view">
    <el-row :gutter="24">
      <!-- 模板列表 (左栏 38% 舒适宽度) -->
      <el-col :span="10">
        <el-card class="template-list-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <div>
                <h3 class="panel-main-title">点检模板库</h3>
                <span class="panel-sub-title">共 {{ templates.length }} 个标准模板</span>
              </div>
              <el-button type="primary" size="small" @click="showCreateDialog = true">
                <el-icon><Plus /></el-icon> 新建模板
              </el-button>
            </div>
          </template>

          <div class="template-items-stream" v-loading="loading">
            <div
              v-for="tpl in templates"
              :key="tpl.id"
              class="tpl-card-item"
              :class="{ active: currentTemplate && currentTemplate.id === tpl.id }"
              @click="selectTemplate(tpl)"
            >
              <div class="tpl-item-main">
                <div class="tpl-name-row">
                  <span class="tpl-title">{{ tpl.name }}</span>
                  <el-tag size="small" type="primary" effect="plain" round>{{ tpl.equipment_type_name || '通用' }}</el-tag>
                </div>
                <div class="tpl-meta-row">
                  <span class="tpl-count">
                    <el-icon style="margin-right: 3px"><Document /></el-icon>
                    {{ tpl.item_count || 0 }} 个点检项
                  </span>
                  <el-button
                    type="danger"
                    size="small"
                    link
                    class="tpl-del-btn"
                    @click.stop="deleteTemplate(tpl)"
                  >
                    删除
                  </el-button>
                </div>
              </div>
            </div>
            <div v-if="templates.length === 0 && !loading" class="empty-list-hint">
              暂无点检模板，点击右上角新建
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 模板详情 (右栏 62% 沉浸大卡片) -->
      <el-col :span="14">
        <el-card class="template-detail-card" shadow="never">
          <template #header>
            <div class="card-header-flex" v-if="currentTemplate">
              <div>
                <h3 class="panel-main-title">{{ currentTemplate.name }}</h3>
                <span class="panel-sub-title">适用类别：{{ currentTemplate.equipment_type_name || '所有设备' }}</span>
              </div>
              <el-button type="primary" size="small" @click="showAddItemDialog = true">
                <el-icon><Plus /></el-icon> 添加点检项
              </el-button>
            </div>
            <div class="card-header-flex" v-else>
              <h3 class="panel-main-title">模板配置详情</h3>
            </div>
          </template>

          <div v-if="currentTemplate">
            <el-table :data="templateItems" v-loading="itemsLoading" border stripe style="width: 100%">
              <el-table-column prop="sequence_order" label="顺序" width="65" align="center" />
              <el-table-column prop="name" label="检查项目" min-width="120" show-overflow-tooltip />
              <el-table-column prop="method" label="检查方法" min-width="120" show-overflow-tooltip />
              <el-table-column prop="criteria" label="判定基准" min-width="130" show-overflow-tooltip />
              <el-table-column label="操作" width="120" align="center" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" size="small" link @click="editItem(row)">编辑</el-button>
                  <el-button type="danger" size="small" link @click="deleteItem(row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div v-else class="empty-detail-state">
            <el-empty description="请从左侧列表选择一个模板查看或配置点检项目" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 创建模板对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      title="新建点检模板"
      width="500px"
    >
      <el-form :model="templateForm" :rules="templateRules" ref="templateFormRef" label-width="100px">
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="设备类型" prop="equipment_type_id">
          <el-select
            v-model="templateForm.equipment_type_id"
            placeholder="请选择设备类型"
            style="width: 100%"
          >
            <el-option
              v-for="type in equipmentTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="createTemplate" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 添加/编辑项目对话框 -->
    <el-dialog
      v-model="showAddItemDialog"
      :title="editingItem ? '编辑点检项目' : '添加点检项目'"
      width="600px"
    >
      <el-form :model="itemForm" :rules="itemRules" ref="itemFormRef" label-width="100px">
        <el-form-item label="项目名称" prop="name">
          <el-input v-model="itemForm.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="检查序号" prop="sequence_order">
          <el-input-number v-model="itemForm.sequence_order" :min="1" />
        </el-form-item>
        <el-form-item label="检查方法" prop="method">
          <el-input
            v-model="itemForm.method"
            type="textarea"
            :rows="2"
            placeholder="请输入检查方法，如：目视检查、听诊等"
          />
        </el-form-item>
        <el-form-item label="判定标准" prop="criteria">
          <el-input
            v-model="itemForm.criteria"
            type="textarea"
            :rows="2"
            placeholder="请输入判定标准，如：无异常、无泄漏等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddItemDialog = false">取消</el-button>
        <el-button type="primary" @click="saveItem" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  inspectionTemplateApi,
  inspectionItemApi,
  type InspectionTemplate,
  type InspectionTemplateDetail,
  type InspectionItem
} from '@/api/inspection'
import { equipmentApi, type EquipmentType } from '@/api/equipment'

const loading = ref(false)
const itemsLoading = ref(false)
const submitting = ref(false)
const templates = ref<InspectionTemplate[]>([])
const templateItems = ref<InspectionItem[]>([])
const currentTemplate = ref<InspectionTemplateDetail | null>(null)
const equipmentTypes = ref<EquipmentType[]>([])

const showCreateDialog = ref(false)
const showAddItemDialog = ref(false)
const editingItem = ref<InspectionItem | null>(null)

const templateFormRef = ref<FormInstance>()
const itemFormRef = ref<FormInstance>()

const templateForm = reactive({
  name: '',
  equipment_type_id: undefined as number | undefined
})

const itemForm = reactive({
  id: undefined as number | undefined,
  name: '',
  method: '',
  criteria: '',
  sequence_order: 1
})

const templateRules: FormRules = {
  name: [{ required: true, message: '请输入模板名称', trigger: 'blur' }],
  equipment_type_id: [{ required: true, message: '请选择设备类型', trigger: 'change' }]
}

const itemRules: FormRules = {
  name: [{ required: true, message: '请输入项目名称', trigger: 'blur' }],
  sequence_order: [{ required: true, message: '请输入检查序号', trigger: 'blur' }]
}

// 加载模板列表
const loadTemplates = async () => {
  loading.value = true
  try {
    const response = await inspectionTemplateApi.getTemplates()
    templates.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载模板列表失败')
  } finally {
    loading.value = false
  }
}

// 加载设备类型
const loadEquipmentTypes = async () => {
  try {
    const response = await equipmentApi.getTypes()
    equipmentTypes.value = response.data
  } catch (error: any) {
    ElMessage.error(error.message || '加载设备类型失败')
  }
}

// 选择模板
const selectTemplate = async (row: InspectionTemplate) => {
  itemsLoading.value = true
  try {
    const response = await inspectionTemplateApi.getTemplate(row.id)
    currentTemplate.value = response.data
    templateItems.value = response.data.items || []
  } catch (error: any) {
    ElMessage.error(error.message || '加载模板详情失败')
  } finally {
    itemsLoading.value = false
  }
}

// 创建模板
const createTemplate = async () => {
  if (!templateFormRef.value) return
  await templateFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await inspectionTemplateApi.createTemplate({
        name: templateForm.name,
        equipment_type_id: templateForm.equipment_type_id!
      })
      ElMessage.success('模板创建成功')
      showCreateDialog.value = false
      templateForm.name = ''
      templateForm.equipment_type_id = undefined
      loadTemplates()
    } catch (error: any) {
      ElMessage.error(error.message || '创建模板失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除模板
const deleteTemplate = async (row: InspectionTemplate) => {
  try {
    await ElMessageBox.confirm('确定要删除该模板吗？此操作不可恢复。', '提示', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    await inspectionTemplateApi.deleteTemplate(row.id)
    ElMessage.success('模板删除成功')
    if (currentTemplate.value?.id === row.id) {
      currentTemplate.value = null
      templateItems.value = []
    }
    loadTemplates()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

// 编辑项目
const editItem = (row: InspectionItem) => {
  editingItem.value = row
  itemForm.id = row.id
  itemForm.name = row.name
  itemForm.method = row.method || ''
  itemForm.criteria = row.criteria || ''
  itemForm.sequence_order = row.sequence_order
  showAddItemDialog.value = true
}

// 保存项目
const saveItem = async () => {
  if (!itemFormRef.value) return
  if (!currentTemplate.value) return

  await itemFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      if (editingItem.value) {
        await inspectionItemApi.updateItem(editingItem.value.id, {
          name: itemForm.name,
          method: itemForm.method,
          criteria: itemForm.criteria,
          sequence_order: itemForm.sequence_order
        })
        ElMessage.success('更新成功')
      } else {
        await inspectionItemApi.createItem({
          template_id: currentTemplate.value.id,
          name: itemForm.name,
          method: itemForm.method,
          criteria: itemForm.criteria,
          sequence_order: itemForm.sequence_order
        })
        ElMessage.success('添加成功')
      }
      showAddItemDialog.value = false
      resetItemForm()
      // 重新加载模板详情
      selectTemplate(currentTemplate.value)
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '操作失败')
    } finally {
      submitting.value = false
    }
  })
}

// 删除项目
const deleteItem = async (row: InspectionItem) => {
  try {
    await ElMessageBox.confirm('确定要删除该项目吗？', '提示', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    await inspectionItemApi.deleteItem(row.id)
    ElMessage.success('项目删除成功')
    if (currentTemplate.value) {
      selectTemplate(currentTemplate.value)
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const resetItemForm = () => {
  editingItem.value = null
  itemForm.id = undefined
  itemForm.name = ''
  itemForm.method = ''
  itemForm.criteria = ''
  itemForm.sequence_order = templateItems.value.length + 1
}

onMounted(() => {
  loadTemplates()
  loadEquipmentTypes()
})
</script>

<style scoped>
.template-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.template-list-card,
.template-detail-card {
  height: calc(100vh - 200px);
}

.template-detail-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

.template-list-card :deep(.el-card__body) {
  height: calc(100% - 60px);
  overflow-y: auto;
}

:deep(.el-table__row) {
  cursor: pointer;
}
</style>

<style scoped>
.template-view {
  max-width: 1360px;
  margin: 0 auto;
}

.template-list-card,
.template-detail-card {
  border-radius: 18px !important;
  background: var(--color-bg-card) !important;
  border: 1px solid var(--color-border) !important;
  min-height: 600px;
}

.card-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-main-title {
  margin: 0;
  font-size: 16px;
  font-weight: 650;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.panel-sub-title {
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.template-items-stream {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tpl-card-item {
  padding: 14px 16px;
  border-radius: 12px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border);
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.tpl-card-item:hover {
  background: var(--color-bg-tertiary);
  transform: translateY(-1px);
}

.tpl-card-item.active {
  background: rgba(10, 124, 255, 0.08);
  border-color: #0a7cff;
  box-shadow: 0 0 0 1px #0a7cff;
}

.tpl-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.tpl-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
}

.tpl-meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: var(--color-text-tertiary);
}

.tpl-count {
  display: inline-flex;
  align-items: center;
}

.empty-detail-state {
  padding: 80px 0;
}

.empty-list-hint {
  text-align: center;
  padding: 40px 0;
  font-size: 13px;
  color: var(--color-text-tertiary);
}
</style>
