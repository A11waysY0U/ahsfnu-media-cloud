<template>
  <div class="workflows-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>工作流管理</h1>
        <p class="subtitle">管理和组织您的工作流程</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog = true" :icon="Plus">
          创建工作流
        </el-button>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <el-card shadow="never">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索工作流名称..."
          clearable
          :prefix-icon="Search"
          @input="handleSearch"
          style="max-width: 400px;"
        />
      </el-card>
    </div>

    <!-- 工作流列表 -->
    <div class="workflows-content">
      <el-row :gutter="20">
        <el-col
          v-for="workflow in filteredWorkflows"
          :key="workflow.id"
          :xs="24"
          :sm="12"
          :lg="8"
          :xl="6"
        >
          <el-card class="workflow-card" shadow="hover">
            <div class="workflow-header">
              <div class="workflow-icon">
                <el-icon :size="24" :color="workflow.color">
                  <component :is="getWorkflowIcon(workflow.type)" />
                </el-icon>
              </div>
              <div class="workflow-actions">
                <el-button
                  type="primary"
                  size="small"
                  circle
                  @click="editWorkflow(workflow)"
                >
                  <el-icon><Edit /></el-icon>
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  circle
                  @click="deleteWorkflow(workflow)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
            <div class="workflow-info">
              <h3 class="workflow-name">{{ workflow.name }}</h3>
              <p class="workflow-description" v-if="workflow.description">
                {{ workflow.description }}
              </p>
              <div class="workflow-meta">
                <el-tag :type="getWorkflowTypeTag(workflow.type)" size="small">
                  {{ getWorkflowTypeName(workflow.type) }}
                </el-tag>
                <span class="material-count">
                  {{ workflow.material_count || 0 }} 个素材
                </span>
              </div>
              <div class="workflow-stats">
                <span class="created-time">
                  创建于 {{ formatDate(workflow.created_at) }}
                </span>
                <span class="status" :class="workflow.is_active ? 'active' : 'inactive'">
                  {{ workflow.is_active ? '启用' : '禁用' }}
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 空状态 -->
      <el-empty
        v-if="filteredWorkflows.length === 0 && !loading"
        description="暂无工作流"
        :image-size="200"
      >
        <el-button type="primary" @click="showCreateDialog = true">
          创建第一个工作流
        </el-button>
      </el-empty>
    </div>

    <!-- 创建/编辑工作流对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="isEditing ? '编辑工作流' : '创建工作流'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="workflowFormRef"
        :model="workflowForm"
        :rules="workflowRules"
        label-width="100px"
      >
        <el-form-item label="工作流名称" prop="name">
          <el-input
            v-model="workflowForm.name"
            placeholder="请输入工作流名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="工作流类型" prop="type">
          <el-select v-model="workflowForm.type" placeholder="请选择工作流类型">
            <el-option label="图片处理" value="image_processing" />
            <el-option label="视频处理" value="video_processing" />
            <el-option label="文件转换" value="file_conversion" />
            <el-option label="批量操作" value="batch_operation" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="工作流颜色" prop="color">
          <el-color-picker
            v-model="workflowForm.color"
            show-alpha
            :predefine="predefineColors"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="workflowForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入工作流描述（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="状态" prop="is_active">
          <el-switch
            v-model="workflowForm.is_active"
            active-text="启用"
            inactive-text="禁用"
          />
        </el-form-item>
        <el-form-item label="配置参数" prop="config">
          <el-input
            v-model="workflowForm.config"
            type="textarea"
            :rows="4"
            placeholder="请输入工作流配置参数（JSON格式）"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitWorkflow" :loading="submitting">
            {{ isEditing ? '保存' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  Edit, 
  Delete,
  Picture,
  VideoPlay,
  Refresh,
  Operation,
  Setting
} from '@element-plus/icons-vue'
import type { WorkflowGroup } from '@/types'
import { workflowAPI } from '@/api'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const workflows = ref<WorkflowGroup[]>([])
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const isEditing = ref(false)
const editingWorkflowId = ref<number | null>(null)

// 表单数据
const workflowFormRef = ref()
const workflowForm = reactive({
  name: '',
  type: '',
  color: '#409EFF',
  description: '',
  is_active: true,
  config: '',
})

// 表单验证规则
const workflowRules = {
  name: [
    { required: true, message: '请输入工作流名称', trigger: 'blur' },
    { min: 1, max: 50, message: '工作流名称长度在 1 到 50 个字符', trigger: 'blur' },
  ],
  type: [
    { required: true, message: '请选择工作流类型', trigger: 'change' },
  ],
  color: [
    { required: true, message: '请选择工作流颜色', trigger: 'change' },
  ],
}

// 预定义颜色
const predefineColors = [
  '#409EFF',
  '#67C23A',
  '#E6A23C',
  '#F56C6C',
  '#909399',
  '#9C27B0',
  '#FF9800',
  '#795548',
  '#607D8B',
  '#3F51B5',
]

// 计算属性
const filteredWorkflows = computed(() => {
  if (!workflows.value || workflows.value.length === 0) {
    return []
  }
  
  if (!searchKeyword.value) {
    return workflows.value
  }
  
  return workflows.value.filter(workflow =>
    workflow.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    (workflow.description && workflow.description.toLowerCase().includes(searchKeyword.value.toLowerCase()))
  )
})

// 方法
const loadWorkflows = async () => {
  loading.value = true
  try {
    const response = await workflowAPI.getList()
    workflows.value = response.data.data || []
  } catch (error) {
    ElMessage.error('加载工作流失败')
    workflows.value = []
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // 搜索功能通过计算属性实现
}

const resetForm = () => {
  workflowForm.name = ''
  workflowForm.type = ''
  workflowForm.color = '#409EFF'
  workflowForm.description = ''
  workflowForm.is_active = true
  workflowForm.config = ''
  isEditing.value = false
  editingWorkflowId.value = null
  workflowFormRef.value?.clearValidate()
}

const editWorkflow = (workflow: WorkflowGroup) => {
  isEditing.value = true
  editingWorkflowId.value = workflow.id
  workflowForm.name = workflow.name
  workflowForm.type = workflow.type
  workflowForm.color = workflow.color
  workflowForm.description = workflow.description || ''
  workflowForm.is_active = workflow.is_active
  workflowForm.config = workflow.config || ''
  showCreateDialog.value = true
}

const deleteWorkflow = async (workflow: WorkflowGroup) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除工作流 "${workflow.name}" 吗？删除后该工作流将从所有素材中移除。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await workflowAPI.delete(workflow.id)
    ElMessage.success('删除成功')
    loadWorkflows()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const submitWorkflow = async () => {
  try {
    await workflowFormRef.value?.validate()
    submitting.value = true
    
    // 验证JSON格式
    if (workflowForm.config) {
      try {
        JSON.parse(workflowForm.config)
      } catch (e) {
        ElMessage.error('配置参数必须是有效的JSON格式')
        return
      }
    }
    
    if (isEditing.value && editingWorkflowId.value) {
      await workflowAPI.update(editingWorkflowId.value, workflowForm)
      ElMessage.success('更新成功')
    } else {
      await workflowAPI.create(workflowForm)
      ElMessage.success('创建成功')
    }
    
    showCreateDialog.value = false
    resetForm()
    loadWorkflows()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

const getWorkflowIcon = (type: string) => {
  const iconMap: Record<string, any> = {
    image_processing: Picture,
    video_processing: VideoPlay,
    file_conversion: Refresh,
    batch_operation: Operation,
    custom: Setting,
  }
  return iconMap[type] || Setting
}

const getWorkflowTypeName = (type: string) => {
  const typeMap: Record<string, string> = {
    image_processing: '图片处理',
    video_processing: '视频处理',
    file_conversion: '文件转换',
    batch_operation: '批量操作',
    custom: '自定义',
  }
  return typeMap[type] || '未知'
}

const getWorkflowTypeTag = (type: string) => {
  const tagMap: Record<string, string> = {
    image_processing: 'success',
    video_processing: 'warning',
    file_conversion: 'info',
    batch_operation: 'danger',
    custom: '',
  }
  return tagMap[type] || ''
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadWorkflows()
})
</script>

<style scoped>
.workflows-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left h1 {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
}

.subtitle {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.search-section {
  margin-bottom: 24px;
}

.workflows-content {
  margin-bottom: 24px;
}

.workflow-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
  border: none;
  border-radius: 12px;
  overflow: hidden;
}

.workflow-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.workflow-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 16px 0 16px;
}

.workflow-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(64, 158, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.workflow-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.workflow-card:hover .workflow-actions {
  opacity: 1;
}

.workflow-info {
  padding: 16px;
}

.workflow-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.workflow-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #666;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.workflow-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.material-count {
  font-size: 12px;
  color: #999;
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 12px;
}

.workflow-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
}

.created-time {
  color: #bbb;
}

.status {
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.status.active {
  background: #f0f9ff;
  color: #409EFF;
}

.status.inactive {
  background: #fef2f2;
  color: #f56c6c;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .header-right {
    width: 100%;
  }
  
  .header-right .el-button {
    width: 100%;
  }
  
  .workflow-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
}
</style>
