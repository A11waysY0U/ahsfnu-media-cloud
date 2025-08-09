<template>
  <div class="material-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.back()" :icon="ArrowLeft" text class="back-button">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <div class="title-section">
          <h1 class="page-title">{{ material?.original_filename || '素材详情' }}</h1>
          <div class="title-meta">
            <el-tag :type="material?.file_type === 'image' ? 'success' : 'warning'" size="small">
              {{ material?.file_type === 'image' ? '图片' : '视频' }}
            </el-tag>
            <span class="file-size" v-if="material">{{ formatFileSize(material.file_size) }}</span>
            <span class="material-id" v-if="material">#{{ material.id }}</span>
          </div>
        </div>
      </div>
      <div class="header-right">
        <el-button-group>
          <el-button 
            type="primary" 
            @click="editMode = !editMode" 
            :icon="Edit"
            :class="{ 'is-active': editMode }"
          >
            <el-icon><Edit /></el-icon>
            {{ editMode ? '取消编辑' : '编辑' }}
          </el-button>
          <el-button 
            type="danger" 
            @click="deleteMaterial" 
            :icon="Delete"
            class="delete-button"
          >
            <el-icon><Delete /></el-icon>
            删除
          </el-button>
        </el-button-group>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <!-- 素材内容 -->
    <div v-else-if="material" class="material-content">
      <el-row :gutter="32">
        <!-- 左侧：素材预览 -->
        <el-col :xs="24" :lg="16">
          <el-card class="preview-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <h3 class="card-title">
                  <el-icon><Picture v-if="material.file_type === 'image'" /><VideoPlay v-else /></el-icon>
                  素材预览
                </h3>
                <div class="card-actions">
                  <el-button 
                    type="primary" 
                    size="small" 
                    @click="downloadOriginal(material)"
                    :icon="Download"
                  >
                    下载原文件
                  </el-button>
                </div>
              </div>
            </template>
            
            <div class="preview-container">
              <!-- 图片预览 -->
              <div v-if="material.file_type === 'image'" class="image-preview">
                <div class="image-wrapper">
                  <img
                    :src="getMaterialUrl(material)"
                    :alt="material.original_filename"
                    class="preview-image"
                    @click="showImageViewer = true"
                    loading="lazy"
                  />
                  <div class="image-overlay">
                    <el-button 
                      type="primary" 
                      size="large" 
                      circle
                      @click="showImageViewer = true"
                    >
                      <el-icon><ZoomIn /></el-icon>
                    </el-button>
                  </div>
                </div>
                <div class="image-info">
                  <div class="info-item" v-if="material.width && material.height">
                    <el-icon><Picture /></el-icon>
                    <span>{{ material.width }} × {{ material.height }}</span>
                  </div>
                  <div class="info-item">
                    <el-icon><Document /></el-icon>
                    <span>{{ formatFileSize(material.file_size) }}</span>
                  </div>
                  <div class="info-item">
                    <el-icon><Calendar /></el-icon>
                    <span>{{ formatDate(material.upload_time) }}</span>
                  </div>
                </div>
              </div>
              
              <!-- 视频预览 -->
              <div v-else class="video-preview">
                <div class="video-wrapper">
                  <video
                    :src="getMaterialUrl(material)"
                    controls
                    class="preview-video"
                    :poster="getMaterialUrl(material)"
                    preload="metadata"
                  />
                </div>
                <div class="video-info">
                  <div class="info-item" v-if="material.width && material.height">
                    <el-icon><VideoPlay /></el-icon>
                    <span>{{ material.width }} × {{ material.height }}</span>
                  </div>
                  <div class="info-item" v-if="material.duration">
                    <el-icon><Clock /></el-icon>
                    <span>{{ formatDuration(material.duration) }}</span>
                  </div>
                  <div class="info-item">
                    <el-icon><Document /></el-icon>
                    <span>{{ formatFileSize(material.file_size) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧：素材信息 -->
        <el-col :xs="24" :lg="8">
          <!-- 基本信息卡片 -->
          <el-card class="info-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <h3 class="card-title">
                  <el-icon><InfoFilled /></el-icon>
                  基本信息
                </h3>
              </div>
            </template>

            <el-form
              ref="formRef"
              :model="form"
              :rules="formRules"
              label-width="100px"
              :disabled="!editMode"
              class="info-form"
            >

              <el-form-item label="文件名" prop="original_filename">
                <el-input v-model="form.original_filename" />
              </el-form-item>
              
              <el-form-item label="存储文件名">
                <span class="readonly-text">{{ material.filename }}</span>
              </el-form-item>
              
              <el-form-item label="文件类型">
                <el-tag :type="material.file_type === 'image' ? 'success' : 'warning'" size="large">
                  <el-icon><Picture v-if="material.file_type === 'image'" /><VideoPlay v-else /></el-icon>
                  {{ material.file_type === 'image' ? '图片' : '视频' }}
                </el-tag>
              </el-form-item>
              
              <el-form-item label="MIME类型">
                <span class="readonly-text">{{ material.mime_type }}</span>
              </el-form-item>
              
              <el-form-item label="文件大小">
                <span class="readonly-text">{{ formatFileSize(material.file_size) }}</span>
              </el-form-item>

              <el-form-item label="分辨率" v-if="material.width && material.height">
                <span class="readonly-text">{{ material.width }} × {{ material.height }}</span>
              </el-form-item>

              <el-form-item label="时长" v-if="material.duration">
                <span class="readonly-text">{{ formatDuration(material.duration) }}</span>
              </el-form-item>
              
              <el-form-item label="上传时间">
                <span class="readonly-text">{{ formatDate(material.upload_time) }}</span>
              </el-form-item>
              
              <el-form-item label="上传者">
                <div class="uploader-info">
                  <el-avatar :size="24">
                    {{ material.uploader?.username?.charAt(0).toUpperCase() }}
                  </el-avatar>
                  <span class="uploader-name">{{ material.uploader?.username || '未知' }}</span>
                  <el-tag v-if="material.uploader?.role === 'admin'" type="danger" size="small">管理员</el-tag>
                </div>
              </el-form-item>
              
              <el-form-item label="星标">
                <el-switch v-model="form.is_starred" />
              </el-form-item>
              
              <el-form-item label="公开">
                <el-switch v-model="form.is_public" />
              </el-form-item>
              
              <el-form-item label="工作流" prop="workflow_id">
                <el-select v-model="form.workflow_id" placeholder="选择工作流" clearable>
                  <el-option
                    v-for="workflow in workflows"
                    :key="workflow.id"
                    :label="workflow.name"
                    :value="workflow.id"
                  >
                    <div class="workflow-option">
                      <div class="workflow-color" :style="{ backgroundColor: workflow.color }"></div>
                      <span>{{ workflow.name }}</span>
                    </div>
                  </el-option>
                </el-select>
              </el-form-item>
              
              <el-form-item label="标签" prop="tag_ids">
                <el-select
                  v-model="form.tag_ids"
                  placeholder="选择标签"
                  multiple
                  clearable
                  filterable
                >
                  <el-option
                    v-for="tag in tags"
                    :key="tag.id"
                    :label="tag.name"
                    :value="tag.id"
                  >
                    <div class="tag-option">
                      <div class="tag-color" :style="{ backgroundColor: tag.color }"></div>
                      <span>{{ tag.name }}</span>
                    </div>
                  </el-option>
                </el-select>
              </el-form-item>
              
              <el-form-item v-if="editMode">
                <div class="form-actions">
                  <el-button type="primary" @click="saveMaterial" :loading="saving">
                    <el-icon><Check /></el-icon>
                    保存更改
                  </el-button>
                  <el-button @click="cancelEdit">
                    <el-icon><Close /></el-icon>
                    取消
                  </el-button>
                </div>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 标签展示卡片 -->
          <el-card class="tags-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <h3 class="card-title">
                  <el-icon><Collection /></el-icon>
                  当前标签
                </h3>
              </div>
            </template>
            
            <div class="tags-display">
              <el-tag
                v-for="materialTag in material.material_tags"
                :key="materialTag.tag?.id || materialTag.id"
                :color="materialTag.tag?.color || '#409EFF'"
                effect="dark"
                class="material-tag"
                size="large"
              >
                <el-icon><Collection /></el-icon>
                {{ materialTag.tag?.name || '未知标签' }}
              </el-tag>
              <el-empty
                v-if="!material.material_tags?.length"
                description="暂无标签"
                :image-size="60"
              />
            </div>
          </el-card>

          <!-- 工作流信息卡片 -->
          <el-card v-if="material.workflow" class="workflow-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <h3 class="card-title">
                  <el-icon><Connection /></el-icon>
                  关联工作流
                </h3>
              </div>
            </template>
            
            <div class="workflow-info">
              <div class="workflow-item">
                <div class="workflow-header">
                  <div class="workflow-color" :style="{ backgroundColor: material.workflow.color }"></div>
                  <span class="workflow-name">{{ material.workflow.name }}</span>
                  <el-tag :type="material.workflow.is_active ? 'success' : 'info'" size="small">
                    {{ material.workflow.is_active ? '活跃' : '非活跃' }}
                  </el-tag>
                </div>
                <div class="workflow-description" v-if="material.workflow.description">
                  {{ material.workflow.description }}
                </div>
                <div class="workflow-meta">
                  <span class="workflow-type">{{ material.workflow.type }}</span>
                  <span class="workflow-created">{{ formatDate(material.workflow.created_at) }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 图片查看器 -->
    <el-image-viewer
      v-if="showImageViewer && material"
      :url-list="[getMaterialUrl(material)]"
      :initial-index="0"
      @close="showImageViewer = false"
    />
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
  Download, 
  Picture, 
  VideoPlay, 
  InfoFilled, 
  Collection,
  ZoomIn,
  Document,
  Calendar,
  Clock,
  Check,
  Close,
  Connection
} from '@element-plus/icons-vue'
import type { Material, Tag, WorkflowGroup } from '@/types'
import { materialAPI, tagAPI, workflowAPI } from '@/api'

const route = useRoute()
const router = useRouter()

// 响应式数据
const loading = ref(false)
const saving = ref(false)
const editMode = ref(false)
const showImageViewer = ref(false)
const material = ref<Material | null>(null)
const tags = ref<Tag[]>([])
const workflows = ref<WorkflowGroup[]>([])

// 表单数据
const formRef = ref()
const form = reactive({
  original_filename: '',
  is_starred: false,
  is_public: false,
  workflow_id: null as number | null,
  tag_ids: [] as number[],
})

// 表单验证规则
const formRules = {
  original_filename: [
    { required: true, message: '请输入文件名', trigger: 'blur' },
    { min: 1, max: 255, message: '文件名长度在 1 到 255 个字符', trigger: 'blur' },
  ],
}

// 计算属性
const materialId = computed(() => Number(route.params.id))

// 方法
const loadMaterial = async () => {
  loading.value = true
  try {
    const response = await materialAPI.getDetail(materialId.value)
    // API返回的数据结构是 { data: Material }
    material.value = response.data.data
    updateForm()
  } catch (error) {
    ElMessage.error('加载素材详情失败')
    router.push('/materials')
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

const updateForm = () => {
  if (material.value) {
    form.original_filename = material.value.original_filename
    form.is_starred = material.value.is_starred
    form.is_public = material.value.is_public
    form.workflow_id = material.value.workflow_id || null
    form.tag_ids = material.value.material_tags?.map(mt => mt.tag?.id).filter((id): id is number => id !== undefined && id !== null) || []
  }
}

const saveMaterial = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true
    
    await materialAPI.update(materialId.value, form)
    ElMessage.success('保存成功')
    editMode.value = false
    loadMaterial()
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

const deleteMaterial = async () => {
  if (!material.value) return
  
  try {
    await ElMessageBox.confirm(
      `确定要删除素材 "${material.value.original_filename}" 吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    await materialAPI.delete(materialId.value)
    ElMessage.success('删除成功')
    router.push('/materials')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

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
  
  // 如果是相对路径，添加 /uploads 前缀
  if (material.thumbnail_path.startsWith('/')) {
    return material.thumbnail_path
  }
  
  return `/uploads/${material.thumbnail_path}`
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDuration = (seconds: number) => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  const secs = seconds % 60
  
  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  } else {
    return `${minutes}:${secs.toString().padStart(2, '0')}`
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

const downloadOriginal = (material: Material) => {
  const url = getMaterialUrl(material)
  const link = document.createElement('a')
  link.href = url
  link.download = material.original_filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('开始下载原文件')
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
  loadMaterial()
  loadTags()
  loadWorkflows()
})
</script>

<style scoped>
.material-detail-page {
  padding: 24px;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding: 24px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.back-button:hover {
  background: #f0f2f5;
  transform: translateX(-4px);
}

.title-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.page-title {
  margin: 0;
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1.2;
}

.title-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-size {
  color: #666;
  font-size: 14px;
}

.material-id {
  color: #666;
  font-size: 14px;
  background: #f0f2f5;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 500;
}

.header-right {
  display: flex;
  gap: 12px;
}

.header-right .el-button-group {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.header-right .el-button {
  border-radius: 12px;
  padding: 12px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.header-right .el-button.is-active {
  background: #409EFF;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
}

.delete-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(245, 108, 108, 0.4);
}

.loading-container {
  padding: 40px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.material-content {
  margin-bottom: 32px;
}

.preview-card,
.info-card,
.tags-card,
.workflow-card {
  margin-bottom: 24px;
  border: none;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.preview-card:hover,
.info-card:hover,
.tags-card:hover,
.workflow-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-actions {
  display: flex;
  gap: 8px;
}

.preview-container {
  text-align: center;
}

.image-preview,
.video-preview {
  position: relative;
}

.image-wrapper,
.video-wrapper {
  position: relative;
  border-radius: 12px;
  overflow: hidden;
  background: #f8f9fa;
}

.preview-image {
  max-width: 100%;
  max-height: 600px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: block;
}

.preview-image:hover {
  transform: scale(1.02);
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
  border-radius: 12px;
}

.image-wrapper:hover .image-overlay {
  opacity: 1;
}

.preview-video {
  max-width: 100%;
  max-height: 600px;
  border-radius: 12px;
  display: block;
}

.image-info,
.video-info {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 20px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 12px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #666;
  font-size: 14px;
}

.info-item .el-icon {
  color: #409EFF;
}

.readonly-text {
  color: #666;
  line-height: 32px;
  font-weight: 500;
}

.file-path {
  word-break: break-all;
  font-size: 14px;
  color: #333;
  background: #f8f9fa;
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.uploader-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.uploader-name {
  font-weight: 500;
  color: #1a1a1a;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.material-tag {
  margin: 0;
  padding: 8px 16px;
  border-radius: 20px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s ease;
}

.material-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.tag-option,
.workflow-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-color,
.workflow-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.workflow-info {
  padding: 16px;
}

.workflow-item {
  border: 1px solid #ebeef5;
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 16px;
  background-color: #f9fafc;
  transition: all 0.3s ease;
}

.workflow-item:hover {
  border-color: #409EFF;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.1);
}

.workflow-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.workflow-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  flex-grow: 1;
}

.workflow-description {
  font-size: 14px;
  color: #666;
  line-height: 1.5;
  margin-bottom: 12px;
}

.workflow-meta {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #909399;
}

.workflow-type {
  font-weight: 500;
}

.workflow-created {
  font-weight: 400;
}

.info-form {
  margin-top: 16px;
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.form-actions .el-button {
  flex: 1;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 500;
}

@media (max-width: 768px) {
  .material-detail-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 20px;
  }
  
  .header-right {
    width: 100%;
  }
  
  .header-right .el-button-group {
    width: 100%;
  }
  
  .header-right .el-button {
    flex: 1;
  }
  
  .page-title {
    font-size: 24px;
  }
  
  .preview-image,
  .preview-video {
    max-height: 400px;
  }
  
  .image-info,
  .video-info {
    flex-direction: column;
    gap: 12px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .workflow-meta {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
