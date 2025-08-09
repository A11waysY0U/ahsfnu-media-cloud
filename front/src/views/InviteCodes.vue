<template>
  <div class="invite-codes-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>邀请码管理</h1>
        <p class="subtitle">管理系统邀请码的生成和使用</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="generateInviteCode" :icon="Plus" :loading="generating">
          生成邀请码
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-section">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card" shadow="never">
            <div class="stat-content">
              <div class="stat-icon total">
                <el-icon :size="24"><Ticket /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-label">总邀请码</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card" shadow="never">
            <div class="stat-content">
              <div class="stat-icon unused">
                <el-icon :size="24"><Document /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.unused }}</div>
                <div class="stat-label">未使用</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card" shadow="never">
            <div class="stat-content">
              <div class="stat-icon used">
                <el-icon :size="24"><UserIcon /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.used }}</div>
                <div class="stat-label">已使用</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card class="stat-card" shadow="never">
            <div class="stat-content">
              <div class="stat-icon expired">
                <el-icon :size="24"><Clock /></el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-value">{{ stats.expired }}</div>
                <div class="stat-label">已过期</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="filter-section">
      <el-card shadow="never">
        <div class="filter-row">
          <div class="filter-item">
            <el-input
              v-model="filters.keyword"
              placeholder="搜索邀请码..."
              clearable
              :prefix-icon="Search"
              @input="handleSearch"
            />
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.status"
              placeholder="状态筛选"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="未使用" :value="0" />
              <el-option label="已使用" :value="1" />
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

    <!-- 邀请码列表 -->
    <div class="invite-codes-content">
      <el-table
        :data="inviteCodes"
        v-loading="loading"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="邀请码" width="200">
          <template #default="{ row }">
            <div class="code-display">
              <span class="code-text">{{ row.code }}</span>
              <el-button
                type="text"
                size="small"
                @click="copyCode(row.code)"
                :icon="CopyDocument"
              >
                复制
              </el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'info'">
              {{ row.status === 0 ? '未使用' : '已使用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建者" width="120">
          <template #default="{ row }">
            <span v-if="row.creator">{{ row.creator.username }}</span>
            <span v-else>{{ getCreatorName(row.created_by) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用者" width="120">
          <template #default="{ row }">
            <span v-if="row.user">{{ row.user.username }}</span>
            <span v-else-if="row.used_by">{{ getUsedByUserName(row.used_by) }}</span>
            <span v-else class="no-user">未使用</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              size="small"
              @click="deleteInviteCode(row)"
              :disabled="row.status === 1"
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

    <!-- 批量生成邀请码对话框 -->
    <el-dialog
      v-model="showGenerateDialog"
      title="生成邀请码"
      width="400px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="generateFormRef"
        :model="generateForm"
        :rules="generateRules"
        label-width="100px"
      >
        <el-form-item label="生成数量" prop="count">
          <el-input-number
            v-model="generateForm.count"
            :min="1"
            :max="100"
            placeholder="请输入生成数量"
          />
        </el-form-item>
        <el-form-item label="备注" prop="note">
          <el-input
            v-model="generateForm.note"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息（可选）"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showGenerateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitGenerate" :loading="generating">
            生成
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, CopyDocument, Ticket, Document, User as UserIcon, Clock } from '@element-plus/icons-vue'
import type { InviteCode, User } from '@/types'
import { inviteCodeAPI, userAPI } from '@/api'

// 响应式数据
const loading = ref(false)
const generating = ref(false)
const inviteCodes = ref<InviteCode[]>([])
const users = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const showGenerateDialog = ref(false)

// 筛选条件
const filters = reactive({
  keyword: '',
  status: null as number | null,
  dateRange: [] as string[],
})

// 统计信息
const stats = ref({
  total: 0,
  unused: 0,
  used: 0,
  expired: 0,
})

// 表单数据
const generateFormRef = ref()
const generateForm = reactive({
  count: 1 as number,
  note: '' as string,
})

// 表单验证规则
const generateRules = {
  count: [
    { required: true, message: '请输入生成数量', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '生成数量在 1 到 100 之间', trigger: 'blur' },
  ],
}

// 方法
const loadInviteCodes = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filters.keyword || undefined,
      status: filters.status,
      start_date: filters.dateRange[0] || undefined,
      end_date: filters.dateRange[1] || undefined,
    }
    
    const response = await inviteCodeAPI.getList(params)
    // 后端现在返回分页格式
    if (response.data && 'data' in response.data) {
      inviteCodes.value = response.data.data || []
      total.value = response.data.pagination?.total || 0
    } else {
      // 兼容旧格式
      inviteCodes.value = Array.isArray(response.data) ? response.data : []
      total.value = inviteCodes.value.length
    }
  } catch (error) {
    console.error('加载邀请码列表失败:', error)
    ElMessage.error('加载邀请码列表失败')
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await userAPI.getList({ page: 1, page_size: 1000 })
    // 后端现在返回分页格式
    if (response.data && 'data' in response.data) {
      users.value = response.data.data || []
    } else {
      // 兼容旧格式
      users.value = Array.isArray(response.data) ? response.data : []
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
  }
}

const loadStats = async () => {
  try {
    const response = await inviteCodeAPI.getStats()
    if (response.data && response.data.data) {
      stats.value = response.data.data
    } else {
      console.error('统计数据结构不正确:', response)
    }
  } catch (error) {
    console.error('加载统计信息失败:', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadInviteCodes()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadInviteCodes()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadInviteCodes()
}

const generateInviteCode = () => {
  generateForm.count = 1
  generateForm.note = ''
  showGenerateDialog.value = true
}

const submitGenerate = async () => {
  try {
    await generateFormRef.value?.validate()
    generating.value = true
    
    // 确保 count 是数字类型
    const count = Number(generateForm.count)
    if (isNaN(count) || count < 1 || count > 100) {
      ElMessage.error('生成数量必须在 1 到 100 之间')
      return
    }
    
    await inviteCodeAPI.generate(count)
    ElMessage.success(`成功生成 ${count} 个邀请码`)
    
    showGenerateDialog.value = false
    loadInviteCodes()
    loadStats()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error('生成邀请码失败')
    }
  } finally {
    generating.value = false
  }
}

const deleteInviteCode = async (inviteCode: InviteCode) => {
  if (inviteCode.status === 1) {
    ElMessage.warning('已使用的邀请码不能删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除邀请码 "${inviteCode.code}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await inviteCodeAPI.delete(inviteCode.id)
    ElMessage.success('删除成功')
    loadInviteCodes()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const copyCode = async (code: string) => {
  try {
    await navigator.clipboard.writeText(code)
    ElMessage.success('邀请码已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const getCreatorName = (creatorId: number) => {
  // 首先尝试从邀请码对象中获取creator信息
  const inviteCode = inviteCodes.value.find(code => code.created_by === creatorId)
  if (inviteCode?.creator) {
    return inviteCode.creator.username
  }
  
  // 如果邀请码对象中没有creator信息，则从users数组中查找
  const creator = users.value.find(user => user.id === creatorId)
  return creator ? creator.username : `用户${creatorId}`
}

const getUsedByUserName = (userId: number) => {
  // 首先尝试从邀请码对象中获取user信息
  const inviteCode = inviteCodes.value.find(code => code.used_by === userId)
  if (inviteCode?.user) {
    return inviteCode.user.username
  }
  
  // 如果邀请码对象中没有user信息，则从users数组中查找
  const user = users.value.find(user => user.id === userId)
  return user ? user.username : `用户${userId}`
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadInviteCodes()
  loadUsers()
  loadStats()
})
</script>

<style scoped>
.invite-codes-page {
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

.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  border-radius: 12px;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.unused {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.used {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.expired {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
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

.invite-codes-content {
  margin-bottom: 24px;
}

.code-display {
  display: flex;
  align-items: center;
  gap: 8px;
}

.code-text {
  font-family: 'Courier New', monospace;
  font-weight: 500;
  color: #409EFF;
  background: #f0f9ff;
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #e0f2fe;
}

.no-user {
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
  
  .stat-content {
    gap: 12px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>
