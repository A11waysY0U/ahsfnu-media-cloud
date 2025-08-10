# AHSFNU 媒体云平台 API 文档

## 概述

AHSFNU 媒体云平台提供了一套完整的 RESTful API，用于管理用户、素材、标签和工作流。本文档详细说明了所有可用的 API 接口。

**基础URL**: `http://localhost:8080/api/v1`

**认证方式**: JWT Token (Bearer Token)

---

## 目录

1. [认证相关 API](#认证相关-api)
2. [素材管理 API](#素材管理-api)
3. [标签管理 API](#标签管理-api)
4. [工作流管理 API](#工作流管理-api)
5. [用户管理 API](#用户管理-api)
6. [邀请码管理 API](#邀请码管理-api)
7. [通用响应格式](#通用响应格式)

---

## 认证相关 API

### 1. 获取验证码

**接口**: `GET /auth/captcha`

**描述**: 获取图形验证码和认证token

**请求参数**: 无

**响应格式**:
```json
{
  "captcha_id": "string",
  "captcha_b64": "string (base64编码的图片)",
  "auth_token": "string"
}
```

**示例响应**:
```json
{
  "captcha_id": "abc123",
  "captcha_b64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 2. 验证验证码

**接口**: `POST /auth/verify-captcha`

**描述**: 验证用户输入的验证码

**请求格式**:
```json
{
  "captcha_id": "string",
  "captcha_code": "string"
}
```

**响应格式**:
```json
{
  "auth_token": "string",
  "message": "验证码验证成功"
}
```

### 3. 用户登录

**接口**: `POST /auth/login`

**描述**: 用户登录获取JWT token

**请求格式**:
```json
{
  "username": "string (用户名或邮箱)",
  "password": "string",
  "auth_token": "string (从验证码接口获取)"
}
```

**响应格式**:
```json
{
  "token": "string (JWT token)",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "string (admin|user)",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 4. 用户注册

**接口**: `POST /auth/register`

**描述**: 新用户注册

**请求格式**:
```json
{
  "username": "string (3-50字符)",
  "email": "string (有效邮箱格式)",
  "password": "string (最少6位)",
  "invite_code": "string",
  "auth_token": "string (从验证码接口获取)"
}
```

**响应格式**:
```json
{
  "token": "string (JWT token)",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 素材管理 API

### 1. 上传素材

**接口**: `POST /materials`

**描述**: 上传新的素材文件

**认证**: 需要JWT token

**请求格式**: `multipart/form-data`

**请求参数**:
- `file`: 文件 (必需)
- `workflow_id`: 工作流ID (可选)

**响应格式**:
```json
{
  "id": 1,
  "filename": "string",
  "original_filename": "string",
  "file_path": "string",
  "file_size": 1024,
  "file_type": "string",
  "mime_type": "string",
  "width": 1920,
  "height": 1080,
  "duration": 120,
  "uploaded_by": 1,
  "upload_time": "2024-01-01T00:00:00Z",
  "is_starred": false,
  "is_public": false,
  "workflow_id": null,
  "thumbnail_path": "string",
  "uploader": {
    "id": 1,
    "username": "string"
  }
}
```

### 2. 更新素材

**接口**: `PUT /materials/{id}`

**描述**: 更新素材信息（包括星标状态、标签等）

**认证**: 需要JWT token

**请求格式**:
```json
{
  "original_filename": "string (可选)",
  "is_starred": true/false (可选，切换星标状态)",
  "is_public": true/false (可选)",
  "workflow_id": 1 (可选，null表示移除工作流)",
  "tag_ids": [1, 2, 3] (可选，标签ID数组，用于更新素材的标签)"
}
```

**响应格式**:
```json
{
  "data": {
    "id": 1,
    "filename": "string",
    "original_filename": "string",
    "file_path": "string",
    "file_size": 1024,
    "file_type": "string",
    "mime_type": "string",
    "width": 1920,
    "height": 1080,
    "duration": 120,
    "uploaded_by": 1,
    "upload_time": "2024-01-01T00:00:00Z",
    "is_starred": false,
    "is_public": false,
    "workflow_id": null,
    "thumbnail_path": "string",
    "uploader": {
      "id": 1,
      "username": "string"
    },
    "material_tags": [
      {
        "tag": {
          "id": 1,
          "name": "string",
          "color": "string"
        }
      }
    ]
  }
}
```

### 3. 获取素材详情

**接口**: `GET /materials/{id}`

**描述**: 获取单个素材的详细信息

**认证**: 需要JWT token

**响应格式**:
```json
{
  "id": 1,
  "filename": "string",
  "original_filename": "string",
  "file_path": "string",
  "file_size": 1024,
  "file_type": "string",
  "mime_type": "string",
  "width": 1920,
  "height": 1080,
  "duration": 120,
  "uploaded_by": 1,
  "upload_time": "2024-01-01T00:00:00Z",
  "is_starred": false,
  "is_public": false,
  "workflow_id": null,
  "thumbnail_path": "string",
  "uploader": {
    "id": 1,
    "username": "string"
  },
  "material_tags": [
    {
      "tag": {
        "id": 1,
        "name": "string",
        "color": "string"
      }
    }
  ]
}
```

### 4. 删除素材

**接口**: `DELETE /materials/{id}`

**描述**: 删除素材文件

**认证**: 需要JWT token (只能删除自己的素材或管理员)

**响应格式**:
```json
{
  "message": "素材删除成功"
}
```

### 5. 搜索素材

**接口**: `GET /materials`

**描述**: 搜索和筛选素材列表

**认证**: 需要JWT token

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20)
- `workflow_id`: 工作流ID (可选)
- `file_type`: 文件类型 (可选)
- `keyword`: 关键词搜索 (可选)
- `tags`: 标签ID列表，逗号分隔 (可选)

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "filename": "string",
      "original_filename": "string",
      "file_path": "string",
      "file_size": 1024,
      "file_type": "string",
      "mime_type": "string",
      "width": 1920,
      "height": 1080,
      "duration": 120,
      "uploaded_by": 1,
      "upload_time": "2024-01-01T00:00:00Z",
      "is_starred": false,
      "is_public": false,
      "workflow_id": null,
      "thumbnail_path": "string",
      "uploader": {
        "id": 1,
        "username": "string"
      }
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

## 标签管理 API

### 1. 获取标签列表

**接口**: `GET /tags`

**描述**: 获取所有标签

**认证**: 需要JWT token

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "string",
      "color": "string",
      "created_by": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "creator": {
        "id": 1,
        "username": "string"
      },
      "material_tags": [
        {
          "material": {
            "id": 1,
            "original_filename": "string"
          }
        }
      ]
    }
  ]
}
```

### 2. 创建标签

**接口**: `POST /tags`

**描述**: 创建新标签

**认证**: 需要JWT token

**请求格式**:
```json
{
  "name": "string (必需)",
  "color": "string (可选，默认#409EFF)"
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "string",
  "color": "string",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "string"
  }
}
```

### 3. 更新标签

**接口**: `PUT /tags/{id}`

**描述**: 更新标签信息

**认证**: 需要JWT token (只能修改自己创建的标签或管理员)

**请求格式**:
```json
{
  "name": "string (可选)",
  "color": "string (可选)"
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "string",
  "color": "string",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "string"
  }
}
```

### 4. 删除标签

**接口**: `DELETE /tags/{id}`

**描述**: 删除标签

**认证**: 需要JWT token (只能删除自己创建的标签或管理员)

**响应格式**:
```json
{
  "message": "标签删除成功"
}
```

---

## 工作流管理 API

### 1. 获取工作流列表

**接口**: `GET /workflows`

**描述**: 获取工作流列表

**认证**: 需要JWT token

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20)
- `keyword`: 关键词搜索 (可选)

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "name": "string",
      "description": "string",
      "type": "string",
      "color": "string",
      "is_active": true,
      "config": "string",
      "status": "string",
      "created_by": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "creator": {
        "id": 1,
        "username": "string"
      },
      "members": [
        {
          "user_id": 1,
          "role": "string",
          "user": {
            "id": 1,
            "username": "string"
          }
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 2. 创建工作流

**接口**: `POST /workflows`

**描述**: 创建新的工作流

**认证**: 需要JWT token

**请求格式**:
```json
{
  "name": "string (必需)",
  "description": "string (可选)",
  "type": "string (可选，默认custom)",
  "color": "string (可选，默认#409EFF)",
  "is_active": true/false (可选，默认true)",
  "config": "string (可选)",
  "members": [1, 2, 3] (可选，用户ID数组)
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "string",
  "description": "string",
  "type": "custom",
  "color": "#409EFF",
  "is_active": true,
  "config": "string",
  "status": "active",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "string"
  },
  "members": [
    {
      "user_id": 1,
      "role": "member",
      "user": {
        "id": 1,
        "username": "string"
      }
    }
  ]
}
```

### 3. 获取工作流详情

**接口**: `GET /workflows/{id}`

**描述**: 获取单个工作流的详细信息

**认证**: 需要JWT token

**响应格式**:
```json
{
  "id": 1,
  "name": "string",
  "description": "string",
  "type": "string",
  "color": "string",
  "is_active": true,
  "config": "string",
  "status": "string",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "string"
  },
  "members": [
    {
      "user_id": 1,
      "role": "string",
      "user": {
        "id": 1,
        "username": "string"
      }
    }
  ]
}
```

### 4. 更新工作流

**接口**: `PUT /workflows/{id}`

**描述**: 更新工作流信息

**认证**: 需要JWT token (只能修改自己创建的工作流或管理员)

**请求格式**:
```json
{
  "name": "string (可选)",
  "description": "string (可选)",
  "type": "string (可选)",
  "color": "string (可选)",
  "is_active": true/false (可选)",
  "config": "string (可选)",
  "members": [1, 2, 3] (可选，用户ID数组)
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "string",
  "description": "string",
  "type": "string",
  "color": "string",
  "is_active": true,
  "config": "string",
  "status": "string",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "string"
  },
  "members": [
    {
      "user_id": 1,
      "role": "string",
      "user": {
        "id": 1,
        "username": "string"
      }
    }
  ]
}
```

### 5. 删除工作流

**接口**: `DELETE /workflows/{id}`

**描述**: 删除工作流

**认证**: 需要JWT token (只能删除自己创建的工作流或管理员)

**响应格式**:
```json
{
  "message": "删除成功，素材已解除关联"
}
```

### 6. 添加工作流成员

**接口**: `POST /workflows/{id}/members`

**描述**: 为工作流添加成员

**认证**: 需要JWT token (只能为自己创建的工作流添加成员或管理员)

**请求格式**:
```json
{
  "user_id": 1 (必需),
  "role": "string (可选，默认member)"
}
```

**响应格式**:
```json
{
  "workflow_id": 1,
  "user_id": 1,
  "role": "string"
}
```

### 7. 移除工作流成员

**接口**: `DELETE /workflows/{id}/members/{userId}`

**描述**: 从工作流移除成员

**认证**: 需要JWT token (只能为自己创建的工作流移除成员或管理员)

**响应格式**:
```json
{
  "message": "成员已移除"
}
```

---

## 用户管理 API

### 1. 获取个人资料

**接口**: `GET /profile`

**描述**: 获取当前用户的个人资料

**认证**: 需要JWT token

**响应格式**:
```json
{
  "id": 1,
  "username": "string",
  "email": "string",
  "role": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 2. 更新个人资料

**接口**: `PUT /profile`

**描述**: 更新当前用户的个人资料

**认证**: 需要JWT token

**请求格式**:
```json
{
  "username": "string (可选)",
  "email": "string (可选)",
  "avatar": "string (可选)"
}
```

**响应格式**:
```json
{
  "id": 1,
  "username": "string",
  "email": "string",
  "role": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### 3. 修改密码

**接口**: `PUT /profile/password`

**描述**: 修改当前用户的密码

**认证**: 需要JWT token

**请求格式**:
```json
{
  "current_password": "string (必需)",
  "new_password": "string (必需)"
}
```

**响应格式**:
```json
{
  "message": "密码修改成功"
}
```

### 4. 获取用户列表 (管理员)

**接口**: `GET /users`

**描述**: 获取所有用户列表 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**响应格式**:
```json
[
  {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "string",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "inviter": {
      "id": 1,
      "username": "string"
    }
  }
]
```

### 5. 更新用户角色 (管理员)

**接口**: `PUT /users/{id}/role`

**描述**: 修改用户角色 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**请求格式**:
```json
{
  "role": "admin|user (必需)"
}
```

**响应格式**:
```json
{
  "message": "权限修改成功",
  "user": {
    "id": 1,
    "username": "string",
    "email": "string",
    "role": "string"
  }
}
```

### 6. 删除用户 (管理员)

**接口**: `DELETE /users/{id}`

**描述**: 删除用户 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**响应格式**:
```json
{
  "message": "删除成功"
}
```

---

## 邀请码管理 API

### 1. 生成邀请码 (管理员)

**接口**: `POST /invite_codes`

**描述**: 生成新的邀请码 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**请求格式**:
```json
{
  "count": 10 (必需，1-100之间)
}
```

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "code": "string",
      "status": 0,
      "created_by": 1,
      "used_by": null,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 2. 获取邀请码列表 (管理员)

**接口**: `GET /invite_codes`

**描述**: 获取邀请码列表 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 20)

**响应格式**:
```json
{
  "data": [
    {
      "id": 1,
      "code": "string",
      "status": 0,
      "created_by": 1,
      "used_by": null,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 3. 获取邀请码统计信息 (管理员)

**接口**: `GET /invite_codes/stats`

**描述**: 获取邀请码统计信息 (仅管理员)

**认证**: 需要JWT token (管理员权限)

**响应格式**:
```json
{
  "data": {
    "total": 100,
    "unused": 50,
    "used": 50,
    "expired": 10
  }
}
```

---

## 通用响应格式

### 成功响应

成功响应格式分为以下几种：

1. **单个对象响应**：直接返回对象本身
```json
{
  "id": 1,
  "name": "string",
  "created_at": "2024-01-01T00:00:00Z"
}
```

2. **列表响应**：包装在data字段中
```json
{
  "data": [
    {
      "id": 1,
      "name": "string"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

3. **统计信息响应**：包装在data字段中
```json
{
  "data": {
    "total": 100,
    "unused": 50,
    "used": 50
  }
}
```

4. **消息响应**：直接返回消息对象
```json
{
  "message": "操作成功"
}
```

### 错误响应

所有错误响应都遵循以下格式：

```json
{
  "error": "错误信息描述"
}
```

### HTTP 状态码

- `200 OK`: 请求成功
- `201 Created`: 创建成功
- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误

### 认证头格式

对于需要认证的接口，请在请求头中添加：

```
Authorization: Bearer <JWT_TOKEN>
```

---

## 文件访问

上传的文件可以通过以下URL访问：

```
GET /uploads/{filename}
```

其中 `{filename}` 是文件在服务器上的存储名称。

---

## 注意事项

1. 所有时间字段都使用 ISO 8601 格式 (UTC)
2. 文件上传大小限制请参考服务器配置
3. 支持的图片格式：JPG, PNG, GIF, BMP
4. 支持的文档格式：PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX
5. 支持的视频格式：MP4, AVI, MOV, WMV
6. 支持的音频格式：MP3, WAV, FLAC, AAC

---

## 更新日志

- **v1.0.2**: 修正响应格式，确保与实际代码实现一致
- **v1.0.1**: 修正API路径，补充邀请码统计接口，更新响应格式
- **v1.0.0**: 初始版本，包含基础的用户、素材、标签、工作流管理功能
