// 用户相关类型
export interface User {
  id: number
  username: string
  email: string
  role: 'admin' | 'user'
  created_at: string
  updated_at: string
  inviter_id?: number
}

// 素材相关类型
export interface Material {
  id: number
  filename: string
  original_filename: string
  file_path: string
  file_size: number
  file_type: 'image' | 'video'
  mime_type: string
  width?: number
  height?: number
  duration?: number
  uploaded_by: number
  workflow_id?: number
  upload_time: string
  is_starred: boolean
  is_public: boolean
  thumbnail_path?: string
  uploader?: User
  workflow?: WorkflowGroup
  material_tags?: MaterialTag[]
}

// 标签相关类型
export interface Tag {
  id: number
  name: string
  color: string
  description?: string
  material_count?: number
  created_by: number
  created_at: string
  creator?: User
  material_tags?: MaterialTag[]
}

export interface MaterialTag {
  id: number
  material_id: number
  tag_id: number
  created_by: number
  created_at: string
  material?: Material
  tag?: Tag
  creator?: User
}

// 工作流相关类型
export interface WorkflowGroup {
  id: number
  name: string
  description?: string
  type: string
  color: string
  is_active: boolean
  config?: string
  material_count?: number
  created_by: number
  created_at: string
  creator?: User
  members?: WorkflowMember[]
}

export interface WorkflowMember {
  id: number
  workflow_id: number
  user_id: number
  role: string
  user?: User
}

// 邀请码相关类型
export interface InviteCode {
  id: number
  code: string
  status: number
  created_by: number
  used_by?: number
  created_at: string
  used_at?: string
  creator?: User
  user?: User
}

// 认证相关类型
export interface LoginRequest {
  username: string
  password: string
  auth_token: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  invite_code: string
  auth_token: string
}

export interface AuthResponse {
  token: string
  user: User
}

// 分页相关类型
export interface Pagination {
  page: number
  page_size: number
  total: number
}

export interface ApiResponse<T> {
  data: T
  pagination?: Pagination
}

export interface ErrorResponse {
  error: string
}

// 文件上传相关类型
export interface UploadProgress {
  loaded: number
  total: number
  percentage: number
}
