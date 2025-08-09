<template>
  <div class="profile-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1>个人资料</h1>
      <p class="subtitle">管理您的账户信息和设置</p>
    </div>

    <div class="profile-content">
      <el-row :gutter="24">
        <!-- 左侧：个人信息 -->
        <el-col :xs="24" :lg="16">
          <el-card class="profile-card" shadow="never">
            <template #header>
              <div class="card-header">
                <h3>基本信息</h3>
                <el-button
                  type="primary"
                  size="small"
                  @click="toggleEdit"
                  :icon="isEditing ? Close : Edit"
                >
                  {{ isEditing ? '取消' : '编辑' }}
                </el-button>
              </div>
            </template>

            <el-form
              ref="profileFormRef"
              :model="profileForm"
              :rules="profileRules"
              label-width="100px"
              :disabled="!isEditing"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="profileForm.username" />
              </el-form-item>
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="profileForm.email" type="email" />
              </el-form-item>
              <el-form-item label="角色">
                <el-tag :type="user?.role === 'admin' ? 'danger' : 'primary'">
                  {{ user?.role === 'admin' ? '管理员' : '普通用户' }}
                </el-tag>
              </el-form-item>
              <el-form-item label="注册时间">
                <span class="readonly-text">{{ formatDate(user?.created_at) }}</span>
              </el-form-item>
              <el-form-item v-if="isEditing">
                <el-button type="primary" @click="saveProfile" :loading="saving">
                  保存更改
                </el-button>
                <el-button @click="cancelEdit">取消</el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 修改密码 -->
          <el-card class="password-card" shadow="never">
            <template #header>
              <h3>修改密码</h3>
            </template>

            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              label-width="120px"
            >
              <el-form-item label="当前密码" prop="currentPassword">
                <el-input
                  v-model="passwordForm.currentPassword"
                  type="password"
                  show-password
                  placeholder="请输入当前密码"
                />
              </el-form-item>
              <el-form-item label="新密码" prop="newPassword">
                <el-input
                  v-model="passwordForm.newPassword"
                  type="password"
                  show-password
                  placeholder="请输入新密码"
                />
              </el-form-item>
              <el-form-item label="确认新密码" prop="confirmPassword">
                <el-input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  show-password
                  placeholder="请再次输入新密码"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="changePassword" :loading="changingPassword">
                  修改密码
                </el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-col>

        <!-- 右侧：统计信息 -->
        <el-col :xs="24" :lg="8">
          <el-card class="stats-card" shadow="never">
            <template #header>
              <h3>账户统计</h3>
            </template>

            <div class="stats-list">
              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#409EFF">
                    <Picture />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ stats.totalMaterials }}</div>
                  <div class="stat-label">上传素材</div>
                </div>
              </div>

              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#67C23A">
                    <Star />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ stats.starredMaterials }}</div>
                  <div class="stat-label">星标素材</div>
                </div>
              </div>

              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#E6A23C">
                    <Collection />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ stats.totalTags }}</div>
                  <div class="stat-label">创建标签</div>
                </div>
              </div>

              <div class="stat-item">
                <div class="stat-icon">
                  <el-icon :size="24" color="#F56C6C">
                    <Operation />
                  </el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">{{ stats.totalWorkflows }}</div>
                  <div class="stat-label">工作流</div>
                </div>
              </div>
            </div>
          </el-card>

          <!-- 最近活动 -->
          <el-card class="activity-card" shadow="never">
            <template #header>
              <h3>最近活动</h3>
            </template>

            <div class="activity-list">
              <div
                v-for="activity in recentActivities"
                :key="activity.id"
                class="activity-item"
              >
                <div class="activity-icon">
                  <el-icon :size="16" :color="activity.color">
                    <component :is="activity.icon" />
                  </el-icon>
                </div>
                <div class="activity-content">
                  <div class="activity-text">{{ activity.text }}</div>
                  <div class="activity-time">{{ formatTime(activity.time) }}</div>
                </div>
              </div>

              <el-empty
                v-if="recentActivities.length === 0"
                description="暂无活动记录"
                :image-size="80"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Edit, Close, Picture, Star, Collection, Operation, Upload, Delete, EditPen } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { materialAPI, tagAPI, workflowAPI } from '@/api'

const authStore = useAuthStore()

// 响应式数据
const isEditing = ref(false)
const saving = ref(false)
const changingPassword = ref(false)
const stats = ref({
  totalMaterials: 0,
  starredMaterials: 0,
  totalTags: 0,
  totalWorkflows: 0,
})

// 表单引用
const profileFormRef = ref()
const passwordFormRef = ref()

// 表单数据
const profileForm = reactive({
  username: '',
  email: '',
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
})

// 表单验证规则
const profileRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' },
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
}

// 计算属性
const user = computed(() => authStore.user)

// 模拟最近活动数据
const recentActivities = ref([
  {
    id: 1,
    text: '上传了图片素材 "banner.jpg"',
    time: new Date(Date.now() - 1000 * 60 * 30), // 30分钟前
    icon: Upload,
    color: '#409EFF',
  },
  {
    id: 2,
    text: '创建了标签 "设计素材"',
    time: new Date(Date.now() - 1000 * 60 * 60 * 2), // 2小时前
    icon: EditPen,
    color: '#67C23A',
  },
  {
    id: 3,
    text: '删除了视频素材 "intro.mp4"',
    time: new Date(Date.now() - 1000 * 60 * 60 * 24), // 1天前
    icon: Delete,
    color: '#F56C6C',
  },
])

// 方法
const loadStats = async () => {
  try {
    // 加载素材统计
    const materialsResponse = await materialAPI.getList({ page: 1, page_size: 1 })
    stats.value.totalMaterials = materialsResponse.data.pagination?.total || 0

    // 加载星标素材统计
    const starredResponse = await materialAPI.getList({ page: 1, page_size: 1, is_starred: true })
    stats.value.starredMaterials = starredResponse.data.pagination?.total || 0

    // 加载标签统计
    const tagsResponse = await tagAPI.getList()
    stats.value.totalTags = tagsResponse.data.data.length

    // 加载工作流统计
    const workflowsResponse = await workflowAPI.getList()
    stats.value.totalWorkflows = workflowsResponse.data.data.length
  } catch (error) {
    console.error('加载统计信息失败:', error)
  }
}

const toggleEdit = () => {
  if (isEditing.value) {
    cancelEdit()
  } else {
    startEdit()
  }
}

const startEdit = () => {
  if (user.value) {
    profileForm.username = user.value.username
    profileForm.email = user.value.email
  }
  isEditing.value = true
}

const cancelEdit = () => {
  isEditing.value = false
  profileFormRef.value?.clearValidate()
}

const saveProfile = async () => {
  try {
    await profileFormRef.value?.validate()
    saving.value = true

    await authStore.updateProfile(profileForm)
    ElMessage.success('个人信息更新成功')
    isEditing.value = false
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error('更新失败')
    }
  } finally {
    saving.value = false
  }
}

const changePassword = async () => {
  try {
    await passwordFormRef.value?.validate()
    changingPassword.value = true

    await authStore.changePassword(
      passwordForm.currentPassword,
      passwordForm.newPassword
    )
    
    ElMessage.success('密码修改成功')
    
    // 清空表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    passwordFormRef.value?.clearValidate()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error('密码修改失败')
    }
  } finally {
    changingPassword.value = false
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('zh-CN')
}

const formatTime = (date: Date) => {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    return `${days}天前`
  }
}

// 生命周期
onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.profile-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
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

.profile-content {
  margin-bottom: 24px;
}

.profile-card,
.password-card,
.stats-card,
.activity-card {
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

.readonly-text {
  color: #666;
  line-height: 32px;
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
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: rgba(64, 158, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-text {
  font-size: 14px;
  color: #1a1a1a;
  line-height: 1.4;
}

.activity-time {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

@media (max-width: 768px) {
  .profile-content {
    margin-bottom: 16px;
  }
  
  .stats-list {
    gap: 12px;
  }
  
  .stat-item {
    padding: 10px;
  }
  
  .stat-value {
    font-size: 18px;
  }
}
</style>
