<template>
  <div class="material-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <el-button @click="$router.back()" :icon="ArrowLeft" text>
          返回
        </el-button>
        <h1>素材详情</h1>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="editMode = !editMode" :icon="Edit">
          {{ editMode ? '取消编辑' : '编辑' }}
        </el-button>
        <el-button type="danger" @click="deleteMaterial" :icon="Delete">
          删除
        </el-button>
      </div>
    </div>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="10" animated />
    </div>

    <div v-else-if="material" class="material-content">
      <el-row :gutter="24">
        <!-- 左侧：素材预览 -->
        <el-col :xs="24" :lg="16">
          <el-card class="preview-card" shadow="never">
            <template #header>
              <h3>素材预览</h3>
            </template>
            
            <div class="preview-container">
              <!-- 图片预览 -->
              <div v-if="material.file_type === 'image'" class="image-preview">
                <img
                  :src="getMaterialUrl(material)"
                  :alt="material.original_filename"
                  class="preview-image"
                  @click="showImageViewer = true"
                />
                <div class="image-info">
                  <span>{{ material.width }} × {{ material.height }}</span>
                  <span>{{ formatFileSize(material.file_size) }}</span>
                </div>
              </div>
              
              <!-- 视频预览 -->
              <div v-else class="video-preview">
                <video
                  :src="getMaterialUrl(material)"
                  controls
                  class="preview-video"
                  :poster="material.thumbnail_path"
                />
                <div class="video-info">
                  <span>{{ material.width }} × {{ material.height }}</span>
                  <span v-if="material.duration">{{ formatDuration(material.duration) }}</span>
                  <span>{{ formatFileSize(material.file_size) }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧：素材信息 -->
        <el-col :xs="24" :lg="8">
          <el-card class="info-card" shadow="never">
            <template #header>
              <h3>素材信息</h3>
            </template>

            <el-form
              ref="formRef"
              :model="form"
              :rules="formRules"
              label-width="100px"
              :disabled="!editMode"
            >
              <el-form-item label="文件名" prop="original_filename">
                <el-input v-model="form.original_filename" />
              </el-form-item>
              
              <el-form-item label="文件类型">
                <el-tag :type="material.file_type === 'image' ? 'success' : 'warning'">
                  {{ material.file_type === 'image' ? '图片' : '视频' }}
                </el-tag>
              </el-form-item>
              
              <el-form-item label="MIME类型">
                <span class="readonly-text">{{ material.mime_type }}</span>
              </el-form-item>
              
              <el-form-item label="文件大小">
                <span class="readonly-text">{{ formatFileSize(material.file_size) }}</span>
              </el-form-item>
              
              <el-form-item label="上传时间">
                <span class="readonly-text">{{ formatDate(material.upload_time) }}</span>
              </el-form-item>
              
              <el-form-item label="上传者">
                <span class="readonly-text">{{ material.uploader?.username || '未知' }}</span>
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
                  />
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
                <el-button type="primary" @click="saveMaterial" :loading="saving">
                  保存更改
                </el-button>
                <el-button @click="cancelEdit">取消</el-button>
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 标签展示 -->
          <el-card class="tags-card" shadow="never">
            <template #header>
              <h3>当前标签</h3>
            </template>
            
            <div class="tags-display">
              <el-tag
                v-for="materialTag in material.material_tags"
                :key="materialTag.tag.id"
                :color="materialTag.tag.color"
                effect="dark"
                class="material-tag"
              >
                {{ materialTag.tag.name }}
              </el-tag>
              <el-empty
                v-if="!material.material_tags?.length"
                description="暂无标签"
                :image-size="60"
              />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 图片查看器 -->
    <el-image-viewer
      v-if="showImageViewer"
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
import { ArrowLeft, Edit, Delete } from '@element-plus/icons-vue'
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
    material.value = response.data
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
    tags.value = response.data
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
    form.workflow_id = material.value.workflow_id
    form.tag_ids = material.value.material_tags?.map(mt => mt.tag.id) || []
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
  return material.file_path.startsWith('http') 
    ? material.file_path 
    : `/uploads${material.file_path}`
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

.material-content {
  margin-bottom: 24px;
}

.preview-card,
.info-card,
.tags-card {
  margin-bottom: 24px;
  border: none;
  border-radius: 12px;
}

.preview-card :deep(.el-card__header) h3,
.info-card :deep(.el-card__header) h3,
.tags-card :deep(.el-card__header) h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.preview-container {
  text-align: center;
}

.image-preview,
.video-preview {
  position: relative;
}

.preview-image {
  max-width: 100%;
  max-height: 600px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.preview-image:hover {
  transform: scale(1.02);
}

.preview-video {
  max-width: 100%;
  max-height: 600px;
  border-radius: 8px;
}

.image-info,
.video-info {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-top: 12px;
  font-size: 14px;
  color: #666;
}

.readonly-text {
  color: #666;
  line-height: 32px;
}

.tags-display {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.material-tag {
  margin: 0;
}

.tag-option {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tag-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
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
  
  .preview-image,
  .preview-video {
    max-height: 400px;
  }
}
</style>
