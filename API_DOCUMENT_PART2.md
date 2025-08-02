# AHSFNU Media Cloud API 文档 - 标签和工作流管理

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

## 6. 工作流管理 API

### 6.1 获取工作流列表

**接口地址**: `GET /api/v1/workflows`

**认证要求**: 需要JWT Token

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
      "name": "工作流名称",
      "description": "工作流描述",
      "status": "active",
      "created_by": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "creator": {
        "id": 1,
        "username": "testuser"
      },
      "members": [
        {
          "id": 1,
          "user": {
            "id": 1,
            "username": "testuser"
          },
          "role": "owner"
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

### 6.2 创建工作流

**接口地址**: `POST /api/v1/workflows`

**认证要求**: 需要JWT Token

**请求格式**:
```json
{
  "name": "工作流名称",
  "description": "工作流描述",
  "members": [2, 3, 4]
}
```

**请求参数说明**:
- `name`: 工作流名称（必需）
- `description`: 工作流描述（可选）
- `members`: 成员用户ID列表（可选）

**响应格式**:
```json
{
  "id": 1,
  "name": "工作流名称",
  "description": "工作流描述",
  "status": "active",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  },
  "members": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "testuser"
      },
      "role": "owner"
    }
  ]
}
```

### 6.3 获取工作流详情

**接口地址**: `GET /api/v1/workflows/:id`

**认证要求**: 需要JWT Token

**URL参数**:
- `id`: 工作流ID

**响应格式**:
```json
{
  "id": 1,
  "name": "工作流名称",
  "description": "工作流描述",
  "status": "active",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  },
  "members": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "testuser"
      },
      "role": "owner"
    }
  ]
}
```

### 6.4 更新工作流

**接口地址**: `PUT /api/v1/workflows/:id`

**认证要求**: 需要JWT Token（工作流创建者或管理员）

**URL参数**:
- `id`: 工作流ID

**请求格式**:
```json
{
  "name": "更新后的工作流名称",
  "description": "更新后的工作流描述",
  "status": "inactive"
}
```

**响应格式**:
```json
{
  "id": 1,
  "name": "更新后的工作流名称",
  "description": "更新后的工作流描述",
  "status": "inactive",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 6.5 删除工作流

**接口地址**: `DELETE /api/v1/workflows/:id`

**认证要求**: 需要JWT Token（工作流创建者或管理员）

**URL参数**:
- `id`: 工作流ID

**响应格式**:
```json
{
  "message": "工作流删除成功"
}
```

### 6.6 添加工作流成员

**接口地址**: `POST /api/v1/workflows/:id/members`

**认证要求**: 需要JWT Token（工作流创建者或管理员）

**URL参数**:
- `id`: 工作流ID

**请求格式**:
```json
{
  "user_id": 2,
  "role": "member"
}
```

**请求参数说明**:
- `user_id`: 用户ID（必需）
- `role`: 角色（可选，默认member）

**响应格式**:
```json
{
  "id": 2,
  "workflow_id": 1,
  "user_id": 2,
  "role": "member",
  "user": {
    "id": 2,
    "username": "memberuser"
  }
}
```

### 6.7 移除工作流成员

**接口地址**: `DELETE /api/v1/workflows/:id/members/:userId`

**认证要求**: 需要JWT Token（工作流创建者或管理员）

**URL参数**:
- `id`: 工作流ID
- `userId`: 用户ID

**响应格式**:
```json
{
  "message": "成员移除成功"
}
```

---

## 错误码说明

- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 认证失败或token无效
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误 