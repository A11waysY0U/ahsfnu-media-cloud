# API文档

## 基础信息
- 基础URL: `/api/v1`
- 静态文件访问: `/uploads`

## 通用响应格式

### 成功响应格式
```json
{
  "success": true,
  "message": "操作成功",
  "data": {
    // 具体数据内容
  }
}
```

### 错误响应格式
```json
{
  "success": false,
  "message": "错误信息",
  "error": "详细错误描述"
}
```

### 分页响应格式
```json
{
  "success": true,
  "message": "获取成功",
  "data": {
    "items": [
      // 数据列表
    ],
    "pagination": {
      "page": 1,
      "page_size": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

## 认证相关API

### 获取验证码
- **URL**: `/api/v1/auth/captcha`
- **方法**: GET
- **描述**: 获取图形验证码
- **认证要求**: 无
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "验证码生成成功",
    "data": {
      "captcha_id": "abc123def456",
      "captcha_image": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
      "expires_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 验证验证码
- **URL**: `/api/v1/auth/verify-captcha`
- **方法**: POST
- **描述**: 验证用户输入的验证码是否正确
- **认证要求**: 无
- **请求体**:
  ```json
  {
    "captcha_id": "验证码ID",
    "captcha_value": "用户输入的验证码"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "验证码验证成功",
    "data": {
      "valid": true,
      "verified_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 用户登录
- **URL**: `/api/v1/auth/login`
- **方法**: POST
- **描述**: 用户登录接口
- **认证要求**: 无
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "captcha_id": "验证码ID",
    "captcha_value": "验证码值"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "refresh_token_here",
      "user": {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com",
        "role": "user",
        "created_at": "2024-01-01T00:00:00Z",
        "last_login": "2024-01-01T12:00:00Z"
      },
      "expires_at": "2024-01-01T13:00:00Z"
    }
  }
  ```

### 用户注册
- **URL**: `/api/v1/auth/register`
- **方法**: POST
- **描述**: 用户注册接口
- **认证要求**: 无
- **请求体**:
  ```json
  {
    "username": "用户名",
    "password": "密码",
    "email": "电子邮箱",
    "invite_code": "邀请码",
    "captcha_id": "验证码ID",
    "captcha_value": "验证码值"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "注册成功",
    "data": {
      "user": {
        "id": 1,
        "username": "newuser",
        "email": "newuser@example.com",
        "role": "user",
        "created_at": "2024-01-01T12:00:00Z"
      },
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "refresh_token_here"
    }
  }
  ```

## 素材管理API

### 上传素材
- **URL**: `/api/v1/materials`
- **方法**: POST
- **描述**: 上传新的素材文件
- **认证要求**: 需要认证
- **请求体**: 表单数据，包含文件和元数据
  ```
  Content-Type: multipart/form-data
  
  file: 文件数据
  title: 素材标题
  description: 素材描述
  tags: 标签列表（可选）
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "素材上传成功",
    "data": {
      "id": 1,
      "title": "示例素材",
      "description": "这是一个示例素材",
      "filename": "example.jpg",
      "file_size": 1024000,
      "file_type": "image/jpeg",
      "file_path": "/uploads/example.jpg",
      "uploaded_by": 1,
      "created_at": "2024-01-01T12:00:00Z",
      "tags": [
        {
          "id": 1,
          "name": "图片",
          "description": "图片类型标签"
        }
      ]
    }
  }
  ```

### 更新素材
- **URL**: `/api/v1/materials/:id`
- **方法**: PUT
- **描述**: 更新指定ID的素材信息
- **认证要求**: 需要认证
- **URL参数**: 
  - id: 素材ID
- **请求体**:
  ```json
  {
    "title": "更新后的标题",
    "description": "更新后的描述",
    "tags": [1, 2, 3]
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "素材更新成功",
    "data": {
      "id": 1,
      "title": "更新后的标题",
      "description": "更新后的描述",
      "updated_at": "2024-01-01T12:00:00Z",
      "tags": [
        {
          "id": 1,
          "name": "图片"
        },
        {
          "id": 2,
          "name": "设计"
        }
      ]
    }
  }
  ```

### 获取单个素材
- **URL**: `/api/v1/materials/:id`
- **方法**: GET
- **描述**: 获取指定ID的素材详细信息
- **认证要求**: 需要认证
- **URL参数**:
  - id: 素材ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "id": 1,
      "title": "示例素材",
      "description": "这是一个示例素材",
      "filename": "example.jpg",
      "file_size": 1024000,
      "file_type": "image/jpeg",
      "file_path": "/uploads/example.jpg",
      "uploaded_by": {
        "id": 1,
        "username": "testuser"
      },
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "download_count": 5,
      "is_starred": true,
      "tags": [
        {
          "id": 1,
          "name": "图片",
          "description": "图片类型标签"
        }
      ]
    }
  }
  ```

### 删除素材
- **URL**: `/api/v1/materials/:id`
- **方法**: DELETE
- **描述**: 删除指定ID的素材
- **认证要求**: 需要认证
- **URL参数**:
  - id: 素材ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "素材删除成功",
    "data": {
      "deleted_id": 1,
      "deleted_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 搜索素材
- **URL**: `/api/v1/materials`
- **方法**: GET
- **描述**: 根据条件搜索素材
- **认证要求**: 需要认证
- **查询参数**:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10）
  - search: 搜索关键词
  - file_type: 文件类型过滤
  - tags: 标签过滤（逗号分隔）
  - uploaded_by: 上传者ID
  - sort_by: 排序字段（created_at, title, file_size）
  - sort_order: 排序方向（asc, desc）
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "搜索成功",
    "data": {
      "items": [
        {
          "id": 1,
          "title": "示例素材1",
          "description": "描述",
          "filename": "example1.jpg",
          "file_size": 1024000,
          "file_type": "image/jpeg",
          "file_path": "/uploads/example1.jpg",
          "uploaded_by": {
            "id": 1,
            "username": "testuser"
          },
          "created_at": "2024-01-01T12:00:00Z",
          "download_count": 5,
          "is_starred": true,
          "tags": [
            {
              "id": 1,
              "name": "图片"
            }
          ]
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total": 25,
        "total_pages": 3
      }
    }
  }
  ```

### 收藏/取消收藏素材
- **URL**: `/api/v1/materials/:id/star`
- **方法**: POST
- **描述**: 收藏或取消收藏指定素材
- **认证要求**: 需要认证
- **URL参数**:
  - id: 素材ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "收藏成功",
    "data": {
      "material_id": 1,
      "is_starred": true,
      "starred_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

## 邀请码管理API

### 生成邀请码
- **URL**: `/api/v1/invite_codes`
- **方法**: POST
- **描述**: 生成新的邀请码
- **认证要求**: 需要认证（可能需要管理员权限）
- **请求体**:
  ```json
  {
    "max_uses": 10,
    "expires_at": "2024-12-31T23:59:59Z",
    "description": "邀请码描述"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "邀请码生成成功",
    "data": {
      "id": 1,
      "code": "INVITE123456",
      "max_uses": 10,
      "current_uses": 0,
      "expires_at": "2024-12-31T23:59:59Z",
      "description": "邀请码描述",
      "created_by": 1,
      "created_at": "2024-01-01T12:00:00Z",
      "is_active": true
    }
  }
  ```

### 获取邀请码列表
- **URL**: `/api/v1/invite_codes`
- **方法**: GET
- **描述**: 获取邀请码列表
- **认证要求**: 需要认证（可能需要管理员权限）
- **查询参数**:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10）
  - is_active: 是否激活（true/false）
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "items": [
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
          "created_at": "2024-01-01T12:00:00Z",
          "is_active": true
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total": 5,
        "total_pages": 1
      }
    }
  }
  ```

## 用户管理API

### 获取个人资料
- **URL**: `/api/v1/profile`
- **方法**: GET
- **描述**: 获取当前登录用户的个人资料
- **认证要求**: 需要认证
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user",
      "avatar": "/uploads/avatar.jpg",
      "bio": "个人简介",
      "created_at": "2024-01-01T00:00:00Z",
      "last_login": "2024-01-01T12:00:00Z",
      "upload_count": 25,
      "star_count": 10
    }
  }
  ```

### 更新个人资料
- **URL**: `/api/v1/profile`
- **方法**: PUT
- **描述**: 更新当前登录用户的个人资料
- **认证要求**: 需要认证
- **请求体**:
  ```json
  {
    "email": "newemail@example.com",
    "bio": "新的个人简介",
    "avatar": "新的头像文件"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "个人资料更新成功",
    "data": {
      "id": 1,
      "username": "testuser",
      "email": "newemail@example.com",
      "role": "user",
      "avatar": "/uploads/new_avatar.jpg",
      "bio": "新的个人简介",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 修改密码
- **URL**: `/api/v1/profile/password`
- **方法**: PUT
- **描述**: 修改当前登录用户的密码
- **认证要求**: 需要认证
- **请求体**: 
  ```json
  {
    "old_password": "旧密码",
    "new_password": "新密码"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "密码修改成功",
    "data": {
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 获取用户列表
- **URL**: `/api/v1/users`
- **方法**: GET
- **描述**: 获取所有用户列表（管理功能）
- **认证要求**: 需要认证（需要管理员权限）
- **查询参数**:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10）
  - search: 搜索关键词
  - role: 角色过滤
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "items": [
        {
          "id": 1,
          "username": "testuser",
          "email": "test@example.com",
          "role": "user",
          "avatar": "/uploads/avatar.jpg",
          "created_at": "2024-01-01T00:00:00Z",
          "last_login": "2024-01-01T12:00:00Z",
          "upload_count": 25,
          "is_active": true
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total": 50,
        "total_pages": 5
      }
    }
  }
  ```

### 更新用户角色
- **URL**: `/api/v1/users/:id/role`
- **方法**: PUT
- **描述**: 管理员修改用户角色/权限
- **认证要求**: 需要认证（需要管理员权限）
- **URL参数**:
  - id: 用户ID
- **请求体**: 
  ```json
  {
    "role": "admin"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "用户角色更新成功",
    "data": {
      "user_id": 1,
      "old_role": "user",
      "new_role": "admin",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 删除用户
- **URL**: `/api/v1/users/:id`
- **方法**: DELETE
- **描述**: 管理员删除用户
- **认证要求**: 需要认证（需要管理员权限）
- **URL参数**:
  - id: 用户ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "用户删除成功",
    "data": {
      "deleted_user_id": 1,
      "deleted_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

## 标签管理API

### 获取标签列表
- **URL**: `/api/v1/tags`
- **方法**: GET
- **描述**: 获取所有标签列表
- **认证要求**: 需要认证
- **查询参数**:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10）
  - search: 搜索关键词
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "items": [
        {
          "id": 1,
          "name": "图片",
          "description": "图片类型标签",
          "created_at": "2024-01-01T00:00:00Z",
          "material_count": 15
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total": 20,
        "total_pages": 2
      }
    }
  }
  ```

### 创建标签
- **URL**: `/api/v1/tags`
- **方法**: POST
- **描述**: 创建新标签
- **认证要求**: 需要认证
- **请求体**: 
  ```json
  {
    "name": "标签名称",
    "description": "标签描述"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "标签创建成功",
    "data": {
      "id": 1,
      "name": "新标签",
      "description": "标签描述",
      "created_by": 1,
      "created_at": "2024-01-01T12:00:00Z",
      "material_count": 0
    }
  }
  ```

### 更新标签
- **URL**: `/api/v1/tags/:id`
- **方法**: PUT
- **描述**: 更新指定ID的标签
- **认证要求**: 需要认证
- **URL参数**:
  - id: 标签ID
- **请求体**:
  ```json
  {
    "name": "更新后的标签名称",
    "description": "更新后的标签描述"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "标签更新成功",
    "data": {
      "id": 1,
      "name": "更新后的标签名称",
      "description": "更新后的标签描述",
      "updated_at": "2024-01-01T12:00:00Z",
      "material_count": 15
    }
  }
  ```

### 删除标签
- **URL**: `/api/v1/tags/:id`
- **方法**: DELETE
- **描述**: 删除指定ID的标签
- **认证要求**: 需要认证
- **URL参数**:
  - id: 标签ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "标签删除成功",
    "data": {
      "deleted_tag_id": 1,
      "deleted_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 为素材添加标签
- **URL**: `/api/v1/tags/:id/materials/:materialId`
- **方法**: POST
- **描述**: 为指定素材添加标签
- **认证要求**: 需要认证
- **URL参数**:
  - id: 标签ID
  - materialId: 素材ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "标签添加成功",
    "data": {
      "tag_id": 1,
      "material_id": 1,
      "added_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 从素材移除标签
- **URL**: `/api/v1/tags/:id/materials/:materialId`
- **方法**: DELETE
- **描述**: 从指定素材移除标签
- **认证要求**: 需要认证
- **URL参数**:
  - id: 标签ID
  - materialId: 素材ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "标签移除成功",
    "data": {
      "tag_id": 1,
      "material_id": 1,
      "removed_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

## 工作流管理API

### 获取工作流列表
- **URL**: `/api/v1/workflows`
- **方法**: GET
- **描述**: 获取所有工作流列表
- **认证要求**: 需要认证
- **查询参数**:
  - page: 页码（默认1）
  - page_size: 每页数量（默认10）
  - status: 状态过滤（active, inactive）
  - created_by: 创建者ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "items": [
        {
          "id": 1,
          "name": "设计工作流",
          "description": "设计相关的工作流程",
          "status": "active",
          "created_by": {
            "id": 1,
            "username": "testuser"
          },
          "created_at": "2024-01-01T00:00:00Z",
          "member_count": 5,
          "material_count": 10
        }
      ],
      "pagination": {
        "page": 1,
        "page_size": 10,
        "total": 15,
        "total_pages": 2
      }
    }
  }
  ```

### 创建工作流
- **URL**: `/api/v1/workflows`
- **方法**: POST
- **描述**: 创建新的工作流
- **认证要求**: 需要认证
- **请求体**:
  ```json
  {
    "name": "工作流名称",
    "description": "工作流描述",
    "status": "active"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "工作流创建成功",
    "data": {
      "id": 1,
      "name": "设计工作流",
      "description": "设计相关的工作流程",
      "status": "active",
      "created_by": 1,
      "created_at": "2024-01-01T12:00:00Z",
      "member_count": 1,
      "material_count": 0
    }
  }
  ```

### 获取单个工作流
- **URL**: `/api/v1/workflows/:id`
- **方法**: GET
- **描述**: 获取指定ID的工作流详细信息
- **认证要求**: 需要认证
- **URL参数**:
  - id: 工作流ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "获取成功",
    "data": {
      "id": 1,
      "name": "设计工作流",
      "description": "设计相关的工作流程",
      "status": "active",
      "created_by": {
        "id": 1,
        "username": "testuser"
      },
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z",
      "members": [
        {
          "id": 1,
          "username": "testuser",
          "role": "owner",
          "joined_at": "2024-01-01T00:00:00Z"
        }
      ],
      "materials": [
        {
          "id": 1,
          "title": "设计素材",
          "filename": "design.jpg",
          "uploaded_by": {
            "id": 1,
            "username": "testuser"
          }
        }
      ]
    }
  }
  ```

### 更新工作流
- **URL**: `/api/v1/workflows/:id`
- **方法**: PUT
- **描述**: 更新指定ID的工作流
- **认证要求**: 需要认证
- **URL参数**:
  - id: 工作流ID
- **请求体**:
  ```json
  {
    "name": "更新后的工作流名称",
    "description": "更新后的工作流描述",
    "status": "inactive"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "工作流更新成功",
    "data": {
      "id": 1,
      "name": "更新后的工作流名称",
      "description": "更新后的工作流描述",
      "status": "inactive",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 删除工作流
- **URL**: `/api/v1/workflows/:id`
- **方法**: DELETE
- **描述**: 删除指定ID的工作流
- **认证要求**: 需要认证
- **URL参数**:
  - id: 工作流ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "工作流删除成功",
    "data": {
      "deleted_workflow_id": 1,
      "deleted_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 添加工作流成员
- **URL**: `/api/v1/workflows/:id/members`
- **方法**: POST
- **描述**: 向指定工作流添加成员
- **认证要求**: 需要认证
- **URL参数**:
  - id: 工作流ID
- **请求体**:
  ```json
  {
    "user_id": 2,
    "role": "member"
  }
  ```
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "成员添加成功",
    "data": {
      "workflow_id": 1,
      "user_id": 2,
      "role": "member",
      "added_at": "2024-01-01T12:00:00Z"
    }
  }
  ```

### 移除工作流成员
- **URL**: `/api/v1/workflows/:id/members/:userId`
- **方法**: DELETE
- **描述**: 从指定工作流移除成员
- **认证要求**: 需要认证
- **URL参数**:
  - id: 工作流ID
  - userId: 用户ID
- **响应示例**:
  ```json
  {
    "success": true,
    "message": "成员移除成功",
    "data": {
      "workflow_id": 1,
      "user_id": 2,
      "removed_at": "2024-01-01T12:00:00Z"
    }
  }
  ```
