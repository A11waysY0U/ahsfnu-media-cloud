import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  Material,
  Tag,
  WorkflowGroup,
  InviteCode,
  User,
  ApiResponse,
  ErrorResponse
} from '@/types'

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      const errorMessage = (data as ErrorResponse)?.error || '请求失败'
      
      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('权限不足')
          break
        case 404:
          ElMessage.error('资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(errorMessage)
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

// 认证相关API
export const authAPI = {
  // 获取验证码
  getCaptcha: () => api.get('/auth/captcha'),
  
  // 验证验证码
  verifyCaptcha: (captchaId: string, captchaCode: string) =>
    api.post('/auth/verify-captcha', { captcha_id: captchaId, captcha_code: captchaCode }),
  
  // 用户登录
  login: (data: LoginRequest) => api.post<AuthResponse>('/auth/login', data),
  
  // 用户注册
  register: (data: RegisterRequest) => api.post<AuthResponse>('/auth/register', data),
  
  // 获取个人资料
  getProfile: () => api.get<User>('/profile'),
  
  // 更新个人资料
  updateProfile: (data: Partial<User>) => api.put<User>('/profile', data),
  
  // 修改密码
  changePassword: (currentPassword: string, newPassword: string) =>
    api.put('/profile/password', { current_password: currentPassword, new_password: newPassword }),
}

// 素材相关API
export const materialAPI = {
  // 上传素材
  upload: (file: File, workflowId?: number, onProgress?: (progress: number) => void) => {
    const formData = new FormData()
    formData.append('file', file)
    if (workflowId) {
      formData.append('workflow_id', workflowId.toString())
    }
    
    // 获取认证token
    const token = localStorage.getItem('token')
    const headers: Record<string, string> = {}
    if (token) {
      headers.Authorization = `Bearer ${token}`
    }
    
    return api.post<Material>('/materials', formData, {
      headers,
      onUploadProgress: (progressEvent) => {
        if (progressEvent.total && onProgress) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      },
    })
  },
  
  // 获取素材列表
  getList: (params?: {
    page?: number
    page_size?: number
    workflow_id?: number
    file_type?: string
    keyword?: string
    tags?: string
  }) => api.get<ApiResponse<Material[]>>('/materials', { params }),
  
  // 获取素材详情
  getDetail: (id: number) => api.get<Material>(`/materials/${id}`),
  
  // 更新素材
  update: (id: number, data: {
    original_filename?: string
    is_starred?: boolean
    is_public?: boolean
    workflow_id?: number | null
    tag_ids?: number[]
  }) => api.put<Material>(`/materials/${id}`, data),
  
  // 删除素材
  delete: (id: number) => api.delete(`/materials/${id}`),
}

// 标签相关API
export const tagAPI = {
  // 获取标签列表
  getList: () => api.get<Tag[]>('/tags'),
  
  // 创建标签
  create: (data: { name: string; color?: string }) => api.post<Tag>('/tags', data),
  
  // 更新标签
  update: (id: number, data: { name?: string; color?: string }) => api.put<Tag>(`/tags/${id}`, data),
  
  // 删除标签
  delete: (id: number) => api.delete(`/tags/${id}`),
}

// 工作流相关API
export const workflowAPI = {
  // 获取工作流列表
  getList: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    api.get<ApiResponse<WorkflowGroup[]>>('/workflows', { params }),
  
  // 创建工作流
  create: (data: { name: string; description?: string; members?: number[] }) =>
    api.post<WorkflowGroup>('/workflows', data),
  
  // 获取工作流详情
  getDetail: (id: number) => api.get<WorkflowGroup>(`/workflows/${id}`),
  
  // 更新工作流
  update: (id: number, data: { name?: string; description?: string; members?: number[] }) =>
    api.put<WorkflowGroup>(`/workflows/${id}`, data),
  
  // 删除工作流
  delete: (id: number) => api.delete(`/workflows/${id}`),
  
  // 添加工作流成员
  addMember: (id: number, data: { user_id: number; role?: string }) =>
    api.post(`/workflows/${id}/members`, data),
  
  // 移除工作流成员
  removeMember: (id: number, userId: number) =>
    api.delete(`/workflows/${id}/members/${userId}`),
}

// 用户管理API（管理员）
export const userAPI = {
  // 获取用户列表
  getList: (params?: { page?: number; page_size?: number; keyword?: string; role?: string }) => 
    api.get<ApiResponse<User[]>>('/users', { params }),
  
  // 更新用户角色
  updateRole: (id: number, role: 'admin' | 'user') =>
    api.put(`/users/${id}/role`, { role }),
  
  // 删除用户
  delete: (id: number) => api.delete(`/users/${id}`),
}

// 邀请码API（管理员）
export const inviteCodeAPI = {
  // 生成邀请码
  generate: (count: number) => api.post<InviteCode[]>('/invite_codes', { count }),
  
  // 获取邀请码列表
  getList: (params?: { page?: number; page_size?: number }) =>
    api.get<ApiResponse<InviteCode[]>>('/invite_codes', { params }),
  
  // 获取邀请码统计信息
  getStats: () => api.get<{
    data: {
      total: number
      unused: number
      used: number
      expired: number
    }
  }>('/invite_codes/stats'),
  
  // 删除邀请码
  delete: (id: number) => api.delete(`/invite_codes/${id}`),
}

export default api
