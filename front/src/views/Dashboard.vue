<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
              <el-icon><Picture /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.totalMaterials }}</div>
              <div class="stats-label">总素材数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
              <el-icon><Star /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.starredMaterials }}</div>
              <div class="stats-label">收藏素材</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
              <el-icon><Collection /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.totalTags }}</div>
              <div class="stats-label">标签数量</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :sm="12" :md="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-content">
            <div class="stats-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
              <el-icon><Connection /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-number">{{ stats.totalWorkflows }}</div>
              <div class="stats-label">工作流</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近上传的素材 -->
    <el-row :gutter="20" class="content-row">
      <el-col :xs="24" :lg="16">
        <el-card class="content-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近上传的素材</span>
              <el-button type="primary" text @click="$router.push('/materials')">
                查看全部
              </el-button>
            </div>
          </template>
          
          <div v-if="recentMaterials.length === 0" class="empty-state">
            <el-empty description="暂无素材" />
          </div>
          
          <div v-else class="materials-grid">
            <div
              v-for="material in recentMaterials"
              :key="material.id"
              class="material-item"
              @click="viewMaterial(material.id)"
            >
              <div class="material-preview">
                <img
                  v-if="material.file_type === 'image'"
                  :src="getThumbnailUrl(material)"
                  :alt="material.original_filename"
                  class="material-image"
                />
                <div v-else class="material-video">
                  <el-icon><VideoPlay /></el-icon>
                </div>
                <div v-if="material.is_starred" class="material-star">
                  <el-icon><Star /></el-icon>
                </div>
              </div>
              <div class="material-info">
                <div class="material-name">{{ material.original_filename }}</div>
                <div class="material-meta">
                  <span>{{ formatFileSize(material.file_size) }}</span>
                  <span>{{ formatDate(material.upload_time) }}</span>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="8">
        <el-card class="content-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span>快速操作</span>
            </div>
          </template>
          
          <div class="quick-actions">
            <el-button
              type="primary"
              size="large"
              class="action-button"
              @click="showUploadDialog = true"
            >
              <el-icon><Upload /></el-icon>
              上传素材
            </el-button>
            
            <el-button
              type="success"
              size="large"
              class="action-button"
              @click="$router.push('/materials')"
            >
              <el-icon><Picture /></el-icon>
              管理素材
            </el-button>
            
            <el-button
              type="warning"
              size="large"
              class="action-button"
              @click="$router.push('/tags')"
            >
              <el-icon><Collection /></el-icon>
              管理标签
            </el-button>
            
            <el-button
              type="info"
              size="large"
              class="action-button"
              @click="$router.push('/workflows')"
            >
              <el-icon><Connection /></el-icon>
              管理工作流
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 上传对话框 -->
    <el-dialog
      v-model="showUploadDialog"
      title="上传素材"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :on-change="handleFileChange"
        :file-list="fileList"
        drag
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 jpg/png/gif/mp4/mov/avi 格式文件，单个文件不超过50MB
          </div>
        </template>
      </el-upload>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showUploadDialog = false">取消</el-button>
          <el-button type="primary" @click="handleUpload" :loading="uploading">
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { UploadFile, UploadFiles } from 'element-plus'
import { materialAPI, tagAPI, workflowAPI } from '@/api'
import type { Material } from '@/types'

const router = useRouter()

// 统计数据
const stats = reactive({
  totalMaterials: 0,
  starredMaterials: 0,
  totalTags: 0,
  totalWorkflows: 0
})

// 最近素材
const recentMaterials = ref<Material[]>([])

// 上传相关
const showUploadDialog = ref(false)
const uploading = ref(false)
const fileList = ref<UploadFile[]>([])
const uploadRef = ref()

// 获取统计数据
const fetchStats = async () => {
  try {
    // 获取总素材数
    const materialsResponse = await materialAPI.getList({ page: 1, page_size: 1 })
    stats.totalMaterials = materialsResponse.data.pagination?.total || 0
    
    // 获取收藏素材数 - 通过获取所有素材并过滤
    const allMaterialsResponse = await materialAPI.getList({ page: 1, page_size: 1000 })
    const allMaterials = allMaterialsResponse.data.data || []
    stats.starredMaterials = allMaterials.filter(material => material.is_starred).length
    
    // 获取标签数量
    const tagsResponse = await tagAPI.getList()
    stats.totalTags = Array.isArray(tagsResponse.data.data) ? tagsResponse.data.data.length : 0
    
    // 获取工作流数量
    const workflowsResponse = await workflowAPI.getList({ page: 1, page_size: 1 })
    stats.totalWorkflows = workflowsResponse.data.pagination?.total || 0
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 获取最近素材
const fetchRecentMaterials = async () => {
  try {
    const response = await materialAPI.getList({ page: 1, page_size: 8 })
    recentMaterials.value = response.data.data || []
  } catch (error) {
    console.error('获取最近素材失败:', error)
  }
}

// 查看素材详情
const viewMaterial = (id: number) => {
  router.push(`/materials/${id}`)
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 处理文件选择
const handleFileChange = (file: UploadFile, files: UploadFiles) => {
  fileList.value = files
}

// 处理上传
const handleUpload = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('请选择要上传的文件')
    return
  }
  
  uploading.value = true
  
  try {
    for (const file of fileList.value) {
      if (file.raw) {
        await materialAPI.upload(file.raw)
      }
    }
    
    ElMessage.success('上传成功')
    showUploadDialog.value = false
    fileList.value = []
    
    // 刷新数据
    await fetchStats()
    await fetchRecentMaterials()
  } catch (error) {
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 获取缩略图URL
const getThumbnailUrl = (material: Material) => {
  if (!material.thumbnail_path) {
    // 如果没有缩略图，返回原图
    return material.file_path.startsWith('http') ? material.file_path : `/uploads/${material.file_path}`
  }
  
  if (material.thumbnail_path.startsWith('http')) {
    return material.thumbnail_path
  }
  
  if (material.thumbnail_path.startsWith('/')) {
    return material.thumbnail_path
  }
  
  return `/uploads/${material.thumbnail_path}`
}

onMounted(() => {
  fetchStats()
  fetchRecentMaterials()
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 24px;
}

.stats-card {
  border: none;
  border-radius: 12px;
  overflow: hidden;
}

.stats-content {
  display: flex;
  align-items: center;
  padding: 8px 0;
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: white;
  font-size: 24px;
}

.stats-info {
  flex: 1;
}

.stats-number {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  line-height: 1;
  margin-bottom: 4px;
}

.stats-label {
  font-size: 14px;
  color: #666;
}

.content-row {
  margin-bottom: 24px;
}

.content-card {
  border: none;
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.empty-state {
  padding: 40px 0;
}

.materials-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.material-item {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
  border: 1px solid #f0f0f0;
}

.material-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.material-preview {
  position: relative;
  height: 120px;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.material-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.material-video {
  color: #999;
  font-size: 32px;
}

.material-star {
  position: absolute;
  top: 8px;
  right: 8px;
  color: #ffd700;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.material-info {
  padding: 12px;
}

.material-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.material-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action-button {
  width: 100%;
  height: 48px;
  justify-content: flex-start;
  padding-left: 20px;
}

.action-button .el-icon {
  margin-right: 8px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .materials-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
  
  .material-preview {
    height: 100px;
  }
  
  .stats-content {
    padding: 4px 0;
  }
  
  .stats-icon {
    width: 50px;
    height: 50px;
    font-size: 20px;
  }
  
  .stats-number {
    font-size: 24px;
  }
}
</style>
