<template>
  <div class="workflow-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.back()" :icon="ArrowLeft" text>
          返回
        </el-button>
        <h1>工作流详情</h1>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="editMode = !editMode" :icon="Edit">
          {{ editMode ? '取消编辑' : '编辑' }}
        </el-button>
        <el-button type="danger" @click="deleteWorkflow" :icon="Delete">
          删除
        </el-button>
      </div>
    </div>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <div v-else-if="workflow" class="workflow-content">
      <el-row :gutter="24">
        <!-- 左侧：工作流信息 -->
        <el-col :xs="24" :lg="16">
          <el-card class="info-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>基本信息</h3>
                <div class="workflow-status">
                  <el-tag :type="workflow.is_active ? 'success' : 'info'">
                    {{ workflow.is_active ? '启用' : '禁用' }}
                  </el-tag>
                </div>
              </div>
            </template>

            <el-form
              ref="formRef"
              :model="form"
              :rules="formRules"
              label-width="120px"
              :disabled="!editMode"
            >
              <el-form-item label="工作流名称" prop="name">
                <el-input v-model="form.name" />
              </el-form-item>
              
              <el-form-item label="工作流类型" prop="type">
                <el-select v-model="form.type" placeholder="请选择工作流类型">
                  <el-option label="图片处理" value="image_processing" />
                  <el-option label="视频处理" value="video_processing" />
                  <el-option label="文件转换" value="file_conversion" />
                  <el-option label="批量操作" value="batch_operation" />
                  <el-option label="自定义" value="custom" />
                </el-select>
              </el-form-item>
              
              <el-form-item label="工作流颜色" prop="color">
                <el-color-picker
                  v-model="form.color"
                  show-alpha
                  :predefine="predefineColors"
                />
              </el-form-item>
              
              <el-form-item label="描述" prop="description">
                <el-input
                  v-model="form.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入工作流描述"
                />
              </el-form-item>
              
              <el-form-item label="状态" prop="is_active">
                <el-switch
                  v-model="form.is_active"
                  active-text="启用"
                  inactive-text="禁用"
                />
              </el-form-item>
              
              <el-form-item label="配置参数" prop="config">
                <el-input
                  v-model="form.config"
                  type="textarea"
                  :rows="6"
                  placeholder="请输入工作流配置参数（JSON格式）"
                />
              </el-form-item>
              
              <el-form-item v-if="editMode">
                <el-button type="primary" @click="saveWorkflow" :loading="saving">
                  保存更改
                </el-button>
                <el-button @click="cancelEdit">取消</el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 成员管理 -->
          <el-card class="members-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>成员管理</h3>
                <el-button
                  v-if="editMode"
                  type="primary"
                  size="small"
                  @click="showAddMemberDialog = true"
                  :icon="Plus"
                >
                  添加成员
                </el-button>
              </div>
            </template>

            <el-table :data="workflow.members || []" stripe>
              <el-table-column label="用户" min-width="200">
                <template #default="{ row }">
                  <div class="user-info">
                    <el-avatar :size="32" :src="getUserAvatar(row.user)">
                      {{ row.user?.username?.charAt(0).toUpperCase() }}
                    </el-avatar>
                    <div class="user-details">
                      <div class="username">{{ row.user?.username }}</div>
                      <div class="email">{{ row.user?.email }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="角色" width="120">
                <template #default="{ row }">
                  <el-tag :type="getRoleTagType(row.role)">
                    {{ getRoleName(row.role) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120" v-if="editMode">
                <template #default="{ row }">
                  <el-button
                    type="danger"
                    size="small"
                    @click="removeMember(row)"
                  >
                    移除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <el-empty
              v-if="!workflow.members?.length"
              description="暂无成员"
              :image-size="80"
            />
          </el-card>
        </el-col>

        <!-- 右侧：统计信息 -->
        <el-col :xs="24" :lg="8">
          <el-card class="stats-card" shadow="never">
            <template #header>
              <h3>统计信息</h3>
            </template>

            <div class="stats-list">
              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#409EFF">
                    <Picture />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ workflow.material_count || 0 }}</div>
                  <div class="stat-label">关联素材</div>
                </div>
              </div>

              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#67C23A">
                    <UserIcon />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ workflow.members?.length || 0 }}</div>
                  <div class="stat-label">成员数量</div>
                </div>
              </div>

              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#E6A23C">
                    <Clock />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ formatDate(workflow.created_at) }}</div>
                  <div class="stat-label">创建时间</div>
                </div>
              </div>
            </div>
          </el-card>

          <!-- 工作流图标 -->
          <el-card class="icon-card" shadow="never">
            <template #header>
              <h3>工作流图标</h3>
            </template>

            <div class="workflow-icon-display">
              <div class="icon-container" :style="{ backgroundColor: workflow.color + '20' }">
                <el-icon :size="48" :color="workflow.color">
                  <component :is="getWorkflowIcon(workflow.type)" />
                </el-icon>
              </div>
              <div class="icon-info">
                <div class="icon-name">{{ getWorkflowTypeName(workflow.type) }}</div>
                <div class="icon-color">{{ workflow.color }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 添加成员对话框 -->
    <el-dialog
      v-model="showAddMemberDialog"
      title="添加成员"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="memberFormRef"
        :model="memberForm"
        :rules="memberRules"
        label-width="100px"
      >
        <el-form-item label="选择用户" prop="user_id">
          <el-select
            v-model="memberForm.user_id"
            placeholder="请选择用户"
            filterable
          >
            <el-option
              v-for="user in availableUsers"
              :key="user.id"
              :label="user.username"
              :value="user.id"
            >
              <div class="user-option">
                <el-avatar :size="24" :src="getUserAvatar(user)">
                  {{ user.username.charAt(0).toUpperCase() }}
                </el-avatar>
                <span>{{ user.username }}</span>
                <span class="user-email">({{ user.email }})</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="memberForm.role" placeholder="请选择角色">
            <el-option label="管理员" value="admin" />
            <el-option label="编辑者" value="editor" />
            <el-option label="查看者" value="viewer" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddMemberDialog = false">取消</el-button>
          <el-button type="primary" @click="addMember" :loading="addingMember">
            添加
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  ArrowLeft, 
  Edit, 
  Delete, 
  Plus, 
  Picture, 
  User as UserIcon, 
  Clock,
  VideoPlay,
  Refresh,
  Operation,
  Setting
} from '@element-plus/icons-vue'
import type { WorkflowGroup, WorkflowMember, User } from '@/types'
import { workflowAPI, userAPI } from '@/api'

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const addingMember = ref(false)
const editMode = ref(false)
const showAddMemberDialog = ref(false)
const workflow = ref<WorkflowGroup | null>(null)
const allUsers = ref<User[]>([])

// 表单数据
const formRef = ref()
const memberFormRef = ref()
const form = reactive({
  name: '',
  type: '',
  color: '#409EFF',
  description: '',
  is_active: true,
  config: '',
})

const memberForm = reactive({
  user_id: null as number | null,
  role: 'viewer',
})

// 表单验证规则
const formRules = {
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

const memberRules = {
  user_id: [
    { required: true, message: '请选择用户', trigger: 'change' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
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
const workflowId = computed(() => Number(route.params.id))

const availableUsers = computed(() => {
  if (!workflow.value?.members) return allUsers.value
  const memberIds = workflow.value.members.map(m => m.user_id)
  return allUsers.value.filter(user => !memberIds.includes(user.id))
})

// 方法
const loadWorkflow = async () => {
  loading.value = true
  try {
    const response = await workflowAPI.getDetail(workflowId.value)
    workflow.value = response.data
    updateForm()
  } catch (error) {
    ElMessage.error('加载工作流详情失败')
    router.push('/workflows')
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await userAPI.getList()
    allUsers.value = response.data
  } catch (error) {
    ElMessage.error('加载用户列表失败')
  }
}

const updateForm = () => {
  if (workflow.value) {
    form.name = workflow.value.name
    form.type = workflow.value.type
    form.color = workflow.value.color
    form.description = workflow.value.description || ''
    form.is_active = workflow.value.is_active
    form.config = workflow.value.config || ''
  }
}

const saveWorkflow = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true
    
    // 验证JSON格式
    if (form.config) {
      try {
        JSON.parse(form.config)
      } catch (e) {
        ElMessage.error('配置参数必须是有效的JSON格式')
        return
      }
    }
    
    await workflowAPI.update(workflowId.value, form)
    ElMessage.success('保存成功')
    editMode.value = false
    loadWorkflow()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error('保存失败')
    }
  } finally {
    saving.value = false
  }
}

const cancelEdit = () => {
  editMode.value = false
  updateForm()
  formRef.value?.clearValidate()
}

const deleteWorkflow = async () => {
  if (!workflow.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除工作流 "${workflow.value.name}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await workflowAPI.delete(workflowId.value)
    ElMessage.success('删除成功')
    router.push('/workflows')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const addMember = async () => {
  try {
    await memberFormRef.value?.validate()
    addingMember.value = true
    
    await workflowAPI.addMember(workflowId.value, memberForm)
    ElMessage.success('添加成员成功')
    showAddMemberDialog.value = false
    memberForm.user_id = null
    memberForm.role = 'viewer'
    loadWorkflow()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error('添加成员失败')
    }
  } finally {
    addingMember.value = false
  }
}

const removeMember = async (member: WorkflowMember) => {
  try {
    await ElMessageBox.confirm(
      `确定要移除成员 "${member.user?.username}" 吗？`,
      '确认移除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await workflowAPI.removeMember(workflowId.value, member.user_id)
    ElMessage.success('移除成员成功')
    loadWorkflow()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('移除成员失败')
    }
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

const getRoleName = (role: string) => {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    editor: '编辑者',
    viewer: '查看者',
  }
  return roleMap[role] || role
}

const getRoleTagType = (role: string) => {
  const tagMap: Record<string, string> = {
    admin: 'danger',
    editor: 'warning',
    viewer: 'info',
  }
  return tagMap[role] || ''
}

const getUserAvatar = (user?: User) => {
  // 这里可以返回用户的头像URL，暂时返回空字符串使用默认头像
  return ''
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 监听编辑模式变化
watch(editMode, (newVal) => {
  if (!newVal) {
    updateForm()
    formRef.value?.clearValidate()
  }
})

// 生命周期
onMounted(() => {
  loadWorkflow()
  loadUsers()
})
</script>

<style scoped>
.workflow-detail-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-left h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #1a1a1a;
}

.header-right {
  display: flex;
  gap: 12px;
}

.loading-container {
  padding: 40px;
}

.workflow-content {
  margin-bottom: 24px;
}

.info-card,
.members-card,
.stats-card,
.icon-card {
  margin-bottom: 24px;
  border: none;
  border-radius: 12px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.workflow-status {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-details {
  flex: 1;
}

.username {
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 2px;
}

.email {
  font-size: 12px;
  color: #666;
}

.stats-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: rgba(64, 158, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

.workflow-icon-display {
  text-align: center;
  padding: 20px;
}

.icon-container {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.icon-info {
  text-align: center;
}

.icon-name {
  font-size: 16px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.icon-color {
  font-size: 12px;
  color: #666;
  font-family: 'Courier New', monospace;
}

.user-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.user-email {
  color: #999;
  font-size: 12px;
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
    flex: 1;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .stats-list {
    gap: 12px;
  }
  
  .stat-item {
    padding: 10px;
  }
  
  .stat-value {
    font-size: 16px;
  }
}
</style>
