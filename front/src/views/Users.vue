<template>
  <div class="users-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>用户管理</h1>
        <p class="subtitle">管理系统中的所有用户账户</p>
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
        <el-table-column label="邀请人" width="120">
          <template #default="{ row }">
            <span v-if="row.inviter_id">{{ getInviterName(row.inviter_id) }}</span>
            <span v-else class="no-inviter">无</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import type { User } from '@/types'
import { userAPI } from '@/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 筛选条件
const filters = reactive({
  keyword: '',
  role: '',
  dateRange: [] as string[],
})

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
    // 后端现在返回分页格式
    if (response.data && 'data' in response.data) {
      users.value = response.data.data || []
      total.value = response.data.pagination?.total || 0
    } else {
      // 兼容旧格式
      users.value = Array.isArray(response.data) ? response.data : []
      total.value = users.value.length
    }
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
    
    await userAPI.updateRole(user.id, newRole)
    ElMessage.success(`${actionText}成功`)
    loadUsers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
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
