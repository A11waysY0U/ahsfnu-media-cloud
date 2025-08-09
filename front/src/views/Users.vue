<template>
  <div class="users-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>用户管理</h1>
        <p class="subtitle">管理系统中的所有用户账户</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog = true" :icon="Plus">
          创建用户
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="filter-section">
      <el-card shadow="never">
        <div class="filter-row">
          <div class="filter-item">
            <el-input
              v-model="filters.keyword"
              placeholder="搜索用户名或邮箱..."
              clearable
              :prefix-icon="Search"
              @input="handleSearch"
            />
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.role"
              placeholder="角色筛选"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="管理员" value="admin" />
              <el-option label="普通用户" value="user" />
            </el-select>
          </div>
          <div class="filter-item">
            <el-date-picker
              v-model="filters.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              @change="handleSearch"
            />
          </div>
        </div>
      </el-card>
    </div>

    <!-- 用户列表 -->
    <div class="users-content">
      <el-table
        :data="users"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户信息" min-width="200">
          <template #default="{ row }">
            <div class="user-info">
              <div class="user-avatar">
                <el-avatar :size="40" :src="getUserAvatar(row)">
                  {{ row.username.charAt(0).toUpperCase() }}
                </el-avatar>
              </div>
              <div class="user-details">
                <div class="username">{{ row.username }}</div>
                <div class="email">{{ row.email }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="角色" width="120">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'primary'">
              {{ row.role === 'admin' ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="最后更新" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>
        <el-table-column label="邀请人" width="120">
          <template #default="{ row }">
            <span v-if="row.inviter_id">{{ getInviterName(row.inviter_id) }}</span>
            <span v-else class="no-inviter">无</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click="editUser(row)"
            >
              编辑
            </el-button>
            <el-button
              :type="row.role === 'admin' ? 'warning' : 'success'"
              size="small"
              @click="toggleRole(row)"
            >
              {{ row.role === 'admin' ? '取消管理员' : '设为管理员' }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="deleteUser(row)"
              :disabled="row.id === currentUser?.id"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 创建/编辑用户对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="isEditing ? '编辑用户' : '创建用户'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="userForm.username"
            placeholder="请输入用户名"
            :disabled="isEditing"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="userForm.email"
            type="email"
            placeholder="请输入邮箱地址"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEditing">
          <el-input
            v-model="userForm.password"
            type="password"
            show-password
            placeholder="请输入密码"
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role" placeholder="请选择角色">
            <el-option label="普通用户" value="user" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        <el-form-item label="邀请码" prop="inviteCode" v-if="!isEditing">
          <el-input
            v-model="userForm.inviteCode"
            placeholder="请输入邀请码"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitUser" :loading="submitting">
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
import { Plus, Search } from '@element-plus/icons-vue'
import type { User } from '@/types'
import { userAPI } from '@/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const showCreateDialog = ref(false)
const isEditing = ref(false)
const editingUserId = ref<number | null>(null)

// 筛选条件
const filters = reactive({
  keyword: '',
  role: '',
  dateRange: [] as string[],
})

// 表单数据
const userFormRef = ref()
const userForm = reactive({
  username: '',
  email: '',
  password: '',
  role: 'user',
  inviteCode: '',
})

// 表单验证规则
const userRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于 6 个字符', trigger: 'blur' },
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' },
  ],
  inviteCode: [
    { required: true, message: '请输入邀请码', trigger: 'blur' },
  ],
}

// 计算属性
const currentUser = computed(() => authStore.user)

// 方法
const loadUsers = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filters.keyword || undefined,
      role: filters.role || undefined,
      start_date: filters.dateRange[0] || undefined,
      end_date: filters.dateRange[1] || undefined,
    }
    
    const response = await userAPI.getList(params)
    users.value = response.data.data
    total.value = response.data.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadUsers()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadUsers()
}

const resetForm = () => {
  userForm.username = ''
  userForm.email = ''
  userForm.password = ''
  userForm.role = 'user'
  userForm.inviteCode = ''
  isEditing.value = false
  editingUserId.value = null
  userFormRef.value?.clearValidate()
}

const editUser = (user: User) => {
  isEditing.value = true
  editingUserId.value = user.id
  userForm.username = user.username
  userForm.email = user.email
  userForm.role = user.role
  showCreateDialog.value = true
}

const deleteUser = async (user: User) => {
  if (user.id === currentUser.value?.id) {
    ElMessage.warning('不能删除自己的账户')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await userAPI.delete(user.id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const toggleRole = async (user: User) => {
  if (user.id === currentUser.value?.id) {
    ElMessage.warning('不能修改自己的角色')
    return
  }

  const newRole = user.role === 'admin' ? 'user' : 'admin'
  const actionText = newRole === 'admin' ? '设为管理员' : '取消管理员'

  try {
    await ElMessageBox.confirm(
      `确定要${actionText} "${user.username}" 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await userAPI.update(user.id, { role: newRole })
    ElMessage.success(`${actionText}成功`)
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const submitUser = async () => {
  try {
    await userFormRef.value?.validate()
    submitting.value = true
    
    if (isEditing.value && editingUserId.value) {
      await userAPI.update(editingUserId.value, {
        email: userForm.email,
        role: userForm.role,
      })
      ElMessage.success('更新成功')
    } else {
      await userAPI.create(userForm)
      ElMessage.success('创建成功')
    }
    
    showCreateDialog.value = false
    resetForm()
    loadUsers()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

const getUserAvatar = (user: User) => {
  // 这里可以返回用户的头像URL，暂时返回空字符串使用默认头像
  return ''
}

const getInviterName = (inviterId: number) => {
  const inviter = users.value.find(user => user.id === inviterId)
  return inviter ? inviter.username : `用户${inviterId}`
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.users-page {
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

.filter-section {
  margin-bottom: 24px;
}

.filter-row {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-item {
  flex: 1;
  min-width: 200px;
}

.users-content {
  margin-bottom: 24px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  flex-shrink: 0;
}

.user-details {
  flex: 1;
}

.username {
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.email {
  font-size: 12px;
  color: #666;
}

.no-inviter {
  color: #999;
  font-style: italic;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
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
  
  .filter-row {
    flex-direction: column;
  }
  
  .filter-item {
    min-width: auto;
  }
}
</style>
