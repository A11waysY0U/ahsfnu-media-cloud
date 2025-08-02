# AHSFNU Media Cloud API 文档

## 概述

AHSFNU Media Cloud 是一个媒体文件管理系统，提供文件上传、管理、标签分类、工作流协作等功能。

### 基础信息

- **基础URL**: `http://localhost:8080`
- **API版本**: v1
- **认证方式**: JWT Token
- **数据格式**: JSON

### 通用响应格式

#### 成功响应
```json
{
  "data": "响应数据",
  "message": "成功消息"
}
```

#### 错误响应
```json
{
  "error": "错误信息"
}
```

#### 分页响应
```json
{
  "data": "数据列表",
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

---

## 1. 认证相关 API

### 1.1 获取验证码

**接口地址**: `GET /api/v1/auth/captcha`

**请求参数**: 无

**响应格式**:
```json
{
  "captcha_id": "验证码ID",
  "captcha_b64": "base64编码的验证码图片",
  "auth_token": "认证token"
}
```

### 1.2 验证验证码

**接口地址**: `POST /api/v1/auth/verify-captcha`

**请求格式**:
```json
{
  "captcha_id": "验证码ID",
  "captcha_code": "用户输入的验证码"
}
```

**响应格式**:
```json
{
  "auth_token": "新的认证token",
  "message": "验证码验证成功"
}
```

### 1.3 用户登录

**接口地址**: `POST /api/v1/auth/login`

**请求格式**:
```json
{
  "username": "用户名或邮箱",
  "password": "密码",
  "auth_token": "认证token"
}
```

**响应格式**:
```json
{
  "token": "JWT访问令牌",
  "user": {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "role": "用户角色",
    "created_at": "创建时间"
  }
}
```

### 1.4 用户注册

**接口地址**: `POST /api/v1/auth/register`

**请求格式**:
```json
{
  "username": "用户名",
  "email": "邮箱地址",
  "password": "密码",
  "invite_code": "邀请码",
  "auth_token": "认证token"
}
```

**响应格式**:
```json
{
  "token": "JWT访问令牌",
  "user": {
    "id": 1,
    "username": "用户名",
    "email": "邮箱地址",
    "role": "user",
    "created_at": "创建时间"
  }
}
```

---

## 2. 素材管理 API

### 2.1 上传素材

**接口地址**: `POST /api/v1/materials`

**认证要求**: 需要JWT Token

**请求格式**: `multipart/form-data`

**请求参数**:
- `file`: 文件数据（必需）
- `workflow_id`: 工作流ID（可选）
- `is_public`: 是否公开（可选，默认false）

**响应格式**:
```json
{
  "id": 1,
  "filename": "存储的文件名",
  "original_filename": "原始文件名",
  "file_path": "文件路径",
  "file_size": 1024000,
  "file_type": "image",
  "mime_type": "image/jpeg",
  "width": 1920,
  "height": 1080,
  "duration": null,
  "uploaded_by": 1,
  "workflow_id": null,
  "upload_time": "2024-01-01T00:00:00Z",
  "is_starred": false,
  "is_public": false,
  "thumbnail_path": "/uploads/thumbnails/example.jpg",
  "uploader": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 2.2 更新素材

**接口地址**: `PUT /api/v1/materials/:id`

**认证要求**: 需要JWT Token（只能更新自己的素材或管理员）

**URL参数**:
- `id`: 素材ID

**请求格式**:
```json
{
  "workflow_id": 1,
  "is_public": true
}
```

### 2.3 获取素材详情

**接口地址**: `GET /api/v1/materials/:id`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 素材ID

### 2.4 删除素材

**接口地址**: `DELETE /api/v1/materials/:id`

**认证要求**: 需要JWT Token（只能删除自己的素材或管理员）

**URL参数**:
- `id`: 素材ID

### 2.5 搜索素材

**接口地址**: `GET /api/v1/materials`

**认证要求**: 需要JWT Token

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `keyword`: 搜索关键词（可选）
- `file_type`: 文件类型过滤（可选）
- `workflow_id`: 工作流ID过滤（可选）
- `tags`: 标签ID列表，逗号分隔（可选）
- `public`: 是否只显示公开素材（可选，true/false）

### 2.6 收藏/取消收藏素材

**接口地址**: `POST /api/v1/materials/:id/star`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 素材ID

---

## 3. 用户管理 API

### 3.1 获取个人资料

**接口地址**: `GET /api/v1/profile`

**认证要求**: 需要JWT Token

**响应格式**:
```json
{
  "id": 1,
  "username": "用户名",
  "email": "邮箱地址",
  "role": "用户角色",
  "created_at": "创建时间"
}
```

### 3.2 更新个人资料

**接口地址**: `PUT /api/v1/profile`

**认证要求**: 需要JWT Token

**请求格式**:
```json
{
  "email": "新邮箱地址"
}
```

### 3.3 修改密码

**接口地址**: `PUT /api/v1/profile/password`

**认证要求**: 需要JWT Token

**请求格式**:
```json
{
  "old_password": "旧密码",
  "new_password": "新密码"
}
```

### 3.4 获取用户列表

**接口地址**: `GET /api/v1/users`

**认证要求**: 需要JWT Token（管理员权限）

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `keyword`: 搜索关键词（可选）

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "username": "用户名",
      "email": "邮箱地址",
      "role": "用户角色",
      "created_at": "创建时间"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 3.5 更新用户角色

**接口地址**: `PUT /api/v1/users/:id/role`

**认证要求**: 需要JWT Token（管理员权限）

**URL参数**:
- `id`: 用户ID

**请求格式**:
```json
{
  "role": "admin"
}
```

### 3.6 删除用户

**接口地址**: `DELETE /api/v1/users/:id`

**认证要求**: 需要JWT Token（管理员权限）

**URL参数**:
- `id`: 用户ID

---

## 4. 邀请码管理 API

### 4.1 生成邀请码

**接口地址**: `POST /api/v1/invite_codes`

**认证要求**: 需要JWT Token（管理员权限）

**请求格式**:
```json
{
  "max_uses": 10,
  "expires_at": "2024-12-31T23:59:59Z",
  "description": "邀请码描述"
}
```

**响应格式**:
```json
{
  "id": 1,
  "code": "INVITE123456",
  "max_uses": 10,
  "current_uses": 0,
  "expires_at": "2024-12-31T23:59:59Z",
  "description": "邀请码描述",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "is_active": true
}
```

### 4.2 获取邀请码列表

**接口地址**: `GET /api/v1/invite_codes`

**认证要求**: 需要JWT Token（管理员权限）

**查询参数**:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `is_active`: 是否激活（true/false）

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "code": "INVITE123456",
      "max_uses": 10,
      "current_uses": 3,
      "expires_at": "2024-12-31T23:59:59Z",
      "description": "邀请码描述",
      "created_by": {
        "id": 1,
        "username": "admin"
      },
      "created_at": "2024-01-01T00:00:00Z",
      "is_active": true
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

---

## 5. 标签管理 API

### 5.1 获取标签列表

**接口地址**: `GET /api/v1/tags`

**认证要求**: 需要JWT Token

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "标签名称",
      "color": "#409EFF",
      "created_by": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "creator": {
        "id": 1,
        "username": "testuser"
      },
      "material_tags": [
        {
          "id": 1,
          "material": {
            "id": 1,
            "filename": "example.jpg"
          }
        }
      ]
    }
  ]
}
```

### 5.2 创建标签

**接口地址**: `POST /api/v1/tags`

**认证要求**: 需要JWT Token

**请求格式**:
```json
{
  "name": "标签名称",
  "color": "#409EFF"
}
```

**请求参数说明**:
- `name`: 标签名称（必需）
- `color`: 标签颜色（可选，默认#409EFF）

**响应格式**:
```json
{
  "id": 1,
  "name": "标签名称",
  "color": "#409EFF",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 5.3 更新标签

**接口地址**: `PUT /api/v1/tags/:id`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 标签ID

**请求格式**:
```json
{
  "name": "更新后的标签名称",
  "color": "#FF6B6B"
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "更新后的标签名称",
  "color": "#FF6B6B",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 5.4 删除标签

**接口地址**: `DELETE /api/v1/tags/:id`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 标签ID

**响应格式**:
```json
{
  "message": "标签删除成功"
}
```

### 5.5 为素材添加标签

**接口地址**: `POST /api/v1/tags/:id/materials/:materialId`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 标签ID
- `materialId`: 素材ID

**响应格式**:
```json
{
  "message": "标签添加成功"
}
```

### 5.6 从素材移除标签

**接口地址**: `DELETE /api/v1/tags/:id/materials/:materialId`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 标签ID
- `materialId`: 素材ID

**响应格式**:
```json
{
  "message": "标签移除成功"
}
```

---

## 认证说明

### JWT Token 使用

在登录或注册成功后，客户端会收到一个JWT token。在后续的API请求中，需要在请求头中包含此token：

```
Authorization: Bearer <your_jwt_token>
```

### 错误码说明

- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 认证失败或token无效
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误
