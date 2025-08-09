<template>
  <div class="tags-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>标签管理</h1>
        <p class="subtitle">管理和组织您的素材标签</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog = true" :icon="Plus">
          创建标签
        </el-button>
      </div>
    </div>

    <!-- 搜索区域 -->
    <div class="search-section">
      <el-card shadow="never">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索标签名称..."
          clearable
          :prefix-icon="Search"
          @input="handleSearch"
          style="max-width: 400px;"
        />
      </el-card>
    </div>

    <!-- 标签列表 -->
    <div class="tags-content">
      <el-row :gutter="20">
        <el-col
          v-for="tag in filteredTags"
          :key="tag.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
          :xl="4"
        >
          <el-card class="tag-card" shadow="hover">
            <div class="tag-header">
              <div class="tag-color" :style="{ backgroundColor: tag.color }"></div>
              <div class="tag-actions">
                <el-button
                  type="primary"
                  size="small"
                  circle
                  @click="editTag(tag)"
                >
                  <el-icon><Edit /></el-icon>
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  circle
                  @click="deleteTag(tag)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
              </div>
            </div>
            <div class="tag-info">
              <h3 class="tag-name">{{ tag.name }}</h3>
              <p class="tag-description" v-if="tag.description">
                {{ tag.description }}
              </p>
              <div class="tag-stats">
                <span class="material-count">
                  {{ tag.material_count || 0 }} 个素材
                </span>
                <span class="created-time">
                  {{ formatDate(tag.created_at) }}
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 空状态 -->
      <el-empty
        v-if="filteredTags.length === 0 && !loading"
        description="暂无标签"
        :image-size="200"
      >
        <el-button type="primary" @click="showCreateDialog = true">
          创建第一个标签
        </el-button>
      </el-empty>
    </div>

    <!-- 创建/编辑标签对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="isEditing ? '编辑标签' : '创建标签'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="tagFormRef"
        :model="tagForm"
        :rules="tagRules"
        label-width="80px"
      >
        <el-form-item label="标签名称" prop="name">
          <el-input
            v-model="tagForm.name"
            placeholder="请输入标签名称"
            maxlength="20"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="标签颜色" prop="color">
          <el-color-picker
            v-model="tagForm.color"
            show-alpha
            :predefine="predefineColors"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="tagForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入标签描述（可选）"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="submitTag" :loading="submitting">
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
import { Plus, Search, Edit, Delete } from '@element-plus/icons-vue'
import type { Tag } from '@/types'
import { tagAPI } from '@/api'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const tags = ref<Tag[]>([])
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const isEditing = ref(false)
const editingTagId = ref<number | null>(null)

// 表单数据
const tagFormRef = ref()
const tagForm = reactive({
  name: '',
  color: '#409EFF',
  description: '',
})

// 表单验证规则
const tagRules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' },
    { min: 1, max: 20, message: '标签名称长度在 1 到 20 个字符', trigger: 'blur' },
  ],
  color: [
    { required: true, message: '请选择标签颜色', trigger: 'change' },
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
const filteredTags = computed(() => {
  if (!searchKeyword.value) {
    return tags.value
  }
  return tags.value.filter(tag =>
    tag.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
    (tag.description && tag.description.toLowerCase().includes(searchKeyword.value.toLowerCase()))
  )
})

// 方法
const loadTags = async () => {
  loading.value = true
  try {
    const response = await tagAPI.getList()
    tags.value = Array.isArray(response.data) ? response.data : []
  } catch (error) {
    ElMessage.error('加载标签失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // 搜索功能通过计算属性实现
}

const resetForm = () => {
  tagForm.name = ''
  tagForm.color = '#409EFF'
  tagForm.description = ''
  isEditing.value = false
  editingTagId.value = null
  tagFormRef.value?.clearValidate()
}

const editTag = (tag: Tag) => {
  isEditing.value = true
  editingTagId.value = tag.id
  tagForm.name = tag.name
  tagForm.color = tag.color
  tagForm.description = tag.description || ''
  showCreateDialog.value = true
}

const deleteTag = async (tag: Tag) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除标签 "${tag.name}" 吗？删除后该标签将从所有素材中移除。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await tagAPI.delete(tag.id)
    ElMessage.success('删除成功')
    loadTags()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const submitTag = async () => {
  try {
    await tagFormRef.value?.validate()
    submitting.value = true
    
    if (isEditing.value && editingTagId.value) {
      await tagAPI.update(editingTagId.value, tagForm)
      ElMessage.success('更新成功')
    } else {
      await tagAPI.create(tagForm)
      ElMessage.success('创建成功')
    }
    
    showCreateDialog.value = false
    resetForm()
    loadTags()
  } catch (error) {
    if (error !== false) { // 表单验证失败时 error 为 false
      ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 监听对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 生命周期
onMounted(() => {
  loadTags()
})
</script>

<style scoped>
.tags-page {
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

.tags-content {
  margin-bottom: 24px;
}

.tag-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
  border: none;
  border-radius: 12px;
  overflow: hidden;
}

.tag-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.tag-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 16px 0 16px;
}

.tag-color {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tag-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.tag-card:hover .tag-actions {
  opacity: 1;
}

.tag-info {
  padding: 16px;
}

.tag-name {
  margin: 0 0 8px 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tag-description {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #666;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.tag-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #999;
}

.material-count {
  background: #f5f5f5;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 500;
}

.created-time {
  color: #bbb;
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
}
</style>
