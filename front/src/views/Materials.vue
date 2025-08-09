<template>
  <div class="materials-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1>素材管理</h1>
        <p class="subtitle">管理和组织您的媒体素材</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showUploadDialog = true" :icon="Plus">
          上传素材
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
              placeholder="搜索素材名称..."
              clearable
              :prefix-icon="Search"
              @input="handleSearch"
            />
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.fileType"
              placeholder="文件类型"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="图片" value="image" />
              <el-option label="视频" value="video" />
            </el-select>
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.isStarred"
              placeholder="星标状态"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option label="已星标" :value="true" />
              <el-option label="未星标" :value="false" />
            </el-select>
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.workflowId"
              placeholder="工作流"
              clearable
              @change="handleSearch"
            >
              <el-option label="全部" value="" />
              <el-option
                v-for="workflow in workflows"
                :key="workflow.id"
                :label="workflow.name"
                :value="workflow.id"
              />
            </el-select>
          </div>
          <div class="filter-item">
            <el-select
              v-model="filters.tagIds"
              placeholder="标签"
              clearable
              multiple
              collapse-tags
              @change="handleSearch"
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 素材列表 -->
    <div class="materials-grid">
      <el-row :gutter="20">
        <el-col
          v-for="material in materials"
          :key="material.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
          :xl="4"
        >
          <el-card
            class="material-card"
            shadow="hover"
            @click="viewMaterial(material)"
          >
            <div class="material-preview">
              <img
                v-if="material.file_type === 'image'"
                :src="getThumbnailUrl(material)"
                :alt="material.original_filename"
                class="preview-image"
              />
              <div
                v-else
                class="video-preview"
                :style="{ backgroundImage: `url(${material.thumbnail_path || '/video-placeholder.jpg'})` }"
              >
                <el-icon class="play-icon"><VideoPlay /></el-icon>
              </div>
              <div class="material-overlay">
                <div class="overlay-actions">
                  <el-button
                    type="primary"
                    size="small"
                    circle
                    @click.stop="viewMaterial(material)"
                  >
                    <el-icon><View /></el-icon>
                  </el-button>
                  <el-button
                    v-if="material.file_type === 'image'"
                    type="success"
                    size="small"
                    circle
                    @click.stop="downloadOriginal(material)"
                    title="下载原图"
                  >
                    <el-icon><Download /></el-icon>
                  </el-button>
                  <el-button
                    :type="material.is_starred ? 'warning' : 'default'"
                    size="small"
                    circle
                    @click.stop="toggleStar(material)"
                  >
                    <el-icon><Star /></el-icon>
                  </el-button>
                  <el-button
                    type="danger"
                    size="small"
                    circle
                    @click.stop="deleteMaterial(material)"
                  >
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
            </div>
            <div class="material-info">
              <h4 class="material-title" :title="material.original_filename">
                {{ material.original_filename }}
              </h4>
              <div class="material-meta">
                <span class="file-size">{{ formatFileSize(material.file_size) }}</span>
                <span class="upload-time">{{ formatDate(material.upload_time) }}</span>
                <span class="uploader" v-if="material.uploader">
                  <el-icon><User /></el-icon>
                  {{ material.uploader.username }}
                </span>
              </div>
              <div class="material-tags" v-if="material.material_tags?.length">
                <el-tag
                  v-for="tag in material.material_tags.slice(0, 2)"
                  :key="tag.tag?.id || tag.id"
                  size="small"
                  class="tag-item"
                >
                  {{ tag.tag?.name || '未知标签' }}
                </el-tag>
                <el-tag
                  v-if="material.material_tags.length > 2"
                  size="small"
                  type="info"
                >
                  +{{ material.material_tags.length - 2 }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 空状态 -->
      <el-empty
        v-if="materials.length === 0 && !loading"
        description="暂无素材"
        :image-size="200"
      >
        <el-button type="primary" @click="showUploadDialog = true">
          上传第一个素材
        </el-button>
      </el-empty>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传素材"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadRef"
        class="upload-area"
        drag
        multiple
        :action="uploadAction"
        :headers="uploadHeaders"
        :data="uploadData"
        :before-upload="beforeUpload"
        :on-progress="onUploadProgress"
        :on-success="onUploadSuccess"
        :on-error="onUploadError"
        :file-list="uploadFileList"
        accept="image/*,video/*"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持图片和视频格式，单个文件不超过100MB
          </div>
        </template>
      </el-upload>

      <div class="upload-options">
        <el-form :model="uploadForm" label-width="80px">
          <el-form-item label="工作流">
            <el-select v-model="uploadForm.workflowId" placeholder="选择工作流" clearable>
              <el-option
                v-for="workflow in workflows"
                :key="workflow.id"
                :label="workflow.name"
                :value="workflow.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="标签">
            <el-select
              v-model="uploadForm.tagIds"
              placeholder="选择标签"
              multiple
              clearable
            >
              <el-option
                v-for="tag in tags"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button type="primary" @click="submitUpload">开始上传</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 素材详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="素材详情"
      width="800px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedMaterial" class="material-detail">
        <div class="detail-preview">
          <img
            v-if="selectedMaterial.file_type === 'image'"
            :src="getMaterialUrl(selectedMaterial)"
            :alt="selectedMaterial.original_filename"
            class="detail-image"
          />
          <video
            v-else
            :src="getMaterialUrl(selectedMaterial)"
            controls
            class="detail-video"
          />
        </div>
        <div class="detail-info">
          <el-form :model="editForm" label-width="100px">
            <el-form-item label="文件名">
              <el-input v-model="editForm.original_filename" />
            </el-form-item>
            <el-form-item label="星标">
              <el-switch v-model="editForm.is_starred" />
            </el-form-item>
            <el-form-item label="公开">
              <el-switch v-model="editForm.is_public" />
            </el-form-item>
            <el-form-item label="工作流">
              <el-select v-model="editForm.workflow_id" placeholder="选择工作流" clearable>
                <el-option
                  v-for="workflow in workflows"
                  :key="workflow.id"
                  :label="workflow.name"
                  :value="workflow.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="标签">
              <el-select
                v-model="editForm.tag_ids"
                placeholder="选择标签"
                multiple
                clearable
              >
                <el-option
                  v-for="tag in tags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.id"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showDetailDialog = false">取消</el-button>
          <el-button type="primary" @click="saveMaterial">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  View, 
  Star, 
  Delete, 
  Edit, 
  VideoPlay, 
  UploadFilled,
  User,
  Download,
} from '@element-plus/icons-vue'
import type { Material, Tag, WorkflowGroup } from '@/types'
import { materialAPI, tagAPI, workflowAPI } from '@/api'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const loading = ref(false)
const materials = ref<Material[]>([])
const tags = ref<Tag[]>([])
const workflows = ref<WorkflowGroup[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(24)

// 筛选条件
const filters = reactive({
  keyword: '',
  fileType: '',
  isStarred: null as boolean | null,
  workflowId: null as number | null,
  tagIds: [] as number[],
})

// 对话框状态
const showUploadDialog = ref(false)
const showDetailDialog = ref(false)
const selectedMaterial = ref<Material | null>(null)

// 上传相关
const uploadRef = ref()
const uploadFileList = ref([])
const uploadForm = reactive({
  workflowId: null as number | null,
  tagIds: [] as number[],
})

// 编辑表单
const editForm = reactive({
  original_filename: '',
  is_starred: false,
  is_public: false,
  workflow_id: null as number | null,
  tag_ids: [] as number[],
})

// 计算属性
const uploadAction = computed(() => '/api/v1/materials')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${authStore.token}`,
}))
const uploadData = computed(() => {
  const data: Record<string, any> = {}
  
  if (uploadForm.workflowId) {
    data.workflow_id = uploadForm.workflowId
  }
  
  if (uploadForm.tagIds && uploadForm.tagIds.length > 0) {
    data.tag_ids = uploadForm.tagIds
  }
  
  return data
})

// 方法
const loadMaterials = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filters.keyword || undefined,
      file_type: filters.fileType || undefined,
      workflow_id: filters.workflowId || undefined,
      tags: filters.tagIds.length > 0 ? filters.tagIds.join(',') : undefined,
    }
    
    const response = await materialAPI.getList(params)
    materials.value = response.data.data || []
    total.value = response.data.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载素材失败')
  } finally {
    loading.value = false
  }
}

const loadTags = async () => {
  try {
    const response = await tagAPI.getList()
    tags.value = response.data.data || []
  } catch (error) {
    ElMessage.error('加载标签失败')
  }
}

const loadWorkflows = async () => {
  try {
    const response = await workflowAPI.getList()
    workflows.value = response.data.data
  } catch (error) {
    ElMessage.error('加载工作流失败')
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadMaterials()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadMaterials()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadMaterials()
}

const viewMaterial = (material: Material) => {
  selectedMaterial.value = material
  editForm.original_filename = material.original_filename
  editForm.is_starred = material.is_starred
  editForm.is_public = material.is_public
  editForm.workflow_id = material.workflow_id || null
  editForm.tag_ids = material.material_tags?.map(mt => mt.tag?.id).filter((id): id is number => id !== undefined && id !== null) || []
  showDetailDialog.value = true
}

const toggleStar = async (material: Material) => {
  try {
    await materialAPI.update(material.id, {
      is_starred: !material.is_starred,
    })
    material.is_starred = !material.is_starred
    ElMessage.success(material.is_starred ? '已添加到星标' : '已取消星标')
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const deleteMaterial = async (material: Material) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除素材 "${material.original_filename}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await materialAPI.delete(material.id)
    ElMessage.success('删除成功')
    loadMaterials()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const saveMaterial = async () => {
  if (!selectedMaterial.value) return
  
  try {
    await materialAPI.update(selectedMaterial.value.id, editForm)
    ElMessage.success('保存成功')
    showDetailDialog.value = false
    loadMaterials()
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 上传相关方法
const beforeUpload = (file: File) => {
  const isValidType = file.type.startsWith('image/') || file.type.startsWith('video/')
  const isLt100M = file.size / 1024 / 1024 < 100

  if (!isValidType) {
    ElMessage.error('只能上传图片或视频文件!')
    return false
  }
  if (!isLt100M) {
    ElMessage.error('文件大小不能超过 100MB!')
    return false
  }
  return true
}

const onUploadProgress = (event: any, file: any) => {
  // 可以在这里显示上传进度
}

const onUploadSuccess = (response: any, file: any) => {
  ElMessage.success(`${file.name} 上传成功`)
  loadMaterials()
}

const onUploadError = (error: any, file: any) => {
  console.error('Upload error:', error)
  let errorMessage = `${file.name} 上传失败`
  
  if (error.response) {
    const errorData = error.response.data
    if (errorData && errorData.error) {
      errorMessage = `${file.name} 上传失败: ${errorData.error}`
    } else if (error.response.status === 413) {
      errorMessage = `${file.name} 上传失败: 文件大小超过限制`
    } else if (error.response.status === 401) {
      errorMessage = `${file.name} 上传失败: 请重新登录`
    } else if (error.response.status === 403) {
      errorMessage = `${file.name} 上传失败: 权限不足`
    }
  } else if (error.message) {
    errorMessage = `${file.name} 上传失败: ${error.message}`
  }
  
  ElMessage.error(errorMessage)
}

const submitUpload = () => {
  uploadRef.value?.submit()
  showUploadDialog.value = false
}

// 工具方法
const getMaterialUrl = (material: Material) => {
  if (material.file_path.startsWith('http')) {
    return material.file_path
  }
  
  // 如果是相对路径，添加 /uploads 前缀
  if (material.file_path.startsWith('/')) {
    return material.file_path
  }
  
  return `/uploads/${material.file_path}`
}

const getThumbnailUrl = (material: Material) => {
  if (!material.thumbnail_path) {
    // 如果没有缩略图，返回原图
    return getMaterialUrl(material)
  }
  
  if (material.thumbnail_path.startsWith('http')) {
    return material.thumbnail_path
  }
  if (material.thumbnail_path.startsWith('/')) {
    return material.thumbnail_path
  }
  return `/uploads/${material.thumbnail_path}`
}

const downloadOriginal = (material: Material) => {
  const url = getMaterialUrl(material)
  const link = document.createElement('a')
  link.href = url
  link.download = material.original_filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('开始下载原图')
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('zh-CN')
}

// 生命周期
onMounted(() => {
  loadMaterials()
  loadTags()
  loadWorkflows()
})
</script>

<style scoped>
.materials-page {
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

.materials-grid {
  margin-bottom: 24px;
}

.material-card {
  margin-bottom: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.material-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
}

.material-preview {
  position: relative;
  height: 200px;
  overflow: hidden;
  border-radius: 8px;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.video-preview {
  width: 100%;
  height: 100%;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.play-icon {
  font-size: 48px;
  color: #fff;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.material-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.material-card:hover .material-overlay {
  opacity: 1;
}

.overlay-actions {
  display: flex;
  gap: 8px;
}

.material-info {
  padding: 12px 0;
}

.material-title {
  margin: 0 0 8px 0;
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.material-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 8px;
  font-size: 12px;
  color: #666;
}

.material-meta .file-size,
.material-meta .upload-time,
.material-meta .uploader {
  display: flex;
  align-items: center;
  gap: 4px;
}

.material-meta .uploader {
  color: #409eff;
  font-weight: 500;
}

.material-meta .uploader .el-icon {
  font-size: 12px;
}

.material-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.tag-item {
  margin: 0;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 32px;
}

.upload-area {
  width: 100%;
}

.upload-options {
  margin-top: 20px;
}

.material-detail {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.detail-preview {
  text-align: center;
}

.detail-image,
.detail-video {
  max-width: 100%;
  max-height: 400px;
  border-radius: 8px;
}

.detail-info {
  padding: 0 16px;
}

@media (max-width: 768px) {
  .material-detail {
    grid-template-columns: 1fr;
  }
  
  .filter-row {
    flex-direction: column;
  }
  
  .filter-item {
    min-width: auto;
  }
}
</style>
