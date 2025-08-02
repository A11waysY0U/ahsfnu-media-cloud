# AHSFNU Media Cloud API 完整文档

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

## API 接口总览

### 1. 认证相关 API (无需认证)

| 方法 | 接口地址 | 描述 |
|------|----------|------|
| GET | `/api/v1/auth/captcha` | 获取验证码 |
| POST | `/api/v1/auth/verify-captcha` | 验证验证码 |
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/register` | 用户注册 |

### 2. 素材管理 API (需要认证)

| 方法 | 接口地址 | 描述 |
|------|----------|------|
| POST | `/api/v1/materials` | 上传素材 |
| PUT | `/api/v1/materials/:id` | 更新素材 |
| GET | `/api/v1/materials/:id` | 获取素材详情 |
| DELETE | `/api/v1/materials/:id` | 删除素材 |
| GET | `/api/v1/materials` | 搜索素材 |
| POST | `/api/v1/materials/:id/star` | 收藏/取消收藏素材 |

### 3. 用户管理 API (需要认证)

| 方法 | 接口地址 | 描述 | 权限要求 |
|------|----------|------|----------|
| GET | `/api/v1/profile` | 获取个人资料 | 普通用户 |
| PUT | `/api/v1/profile` | 更新个人资料 | 普通用户 |
| PUT | `/api/v1/profile/password` | 修改密码 | 普通用户 |
| GET | `/api/v1/users` | 获取用户列表 | 管理员 |
| PUT | `/api/v1/users/:id/role` | 更新用户角色 | 管理员 |
| DELETE | `/api/v1/users/:id` | 删除用户 | 管理员 |

### 4. 邀请码管理 API (需要认证)

| 方法 | 接口地址 | 描述 | 权限要求 |
|------|----------|------|----------|
| POST | `/api/v1/invite_codes` | 生成邀请码 | 管理员 |
| GET | `/api/v1/invite_codes` | 获取邀请码列表 | 管理员 |

### 5. 标签管理 API (需要认证)

| 方法 | 接口地址 | 描述 |
|------|----------|------|
| GET | `/api/v1/tags` | 获取标签列表 |
| POST | `/api/v1/tags` | 创建标签 |
| PUT | `/api/v1/tags/:id` | 更新标签 |
| DELETE | `/api/v1/tags/:id` | 删除标签 |
| POST | `/api/v1/tags/:id/materials/:materialId` | 为素材添加标签 |
| DELETE | `/api/v1/tags/:id/materials/:materialId` | 从素材移除标签 |

### 6. 工作流管理 API (需要认证)

| 方法 | 接口地址 | 描述 | 权限要求 |
|------|----------|------|----------|
| GET | `/api/v1/workflows` | 获取工作流列表 | 普通用户 |
| POST | `/api/v1/workflows` | 创建工作流 | 普通用户 |
| GET | `/api/v1/workflows/:id` | 获取工作流详情 | 普通用户 |
| PUT | `/api/v1/workflows/:id` | 更新工作流 | 创建者/管理员 |
| DELETE | `/api/v1/workflows/:id` | 删除工作流 | 创建者/管理员 |
| POST | `/api/v1/workflows/:id/members` | 添加工作流成员 | 创建者/管理员 |
| DELETE | `/api/v1/workflows/:id/members/:userId` | 移除工作流成员 | 创建者/管理员 |

---

## 认证说明

### JWT Token 使用

在登录或注册成功后，客户端会收到一个JWT token。在后续的API请求中，需要在请求头中包含此token：

```
Authorization: Bearer <your_jwt_token>
```

### 用户角色说明

- **user**: 普通用户，可以管理自己的素材、标签和工作流
- **admin**: 管理员，拥有所有权限，可以管理用户和邀请码

---

## 错误码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 请求成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 认证失败或token无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 使用示例

### 1. 用户注册流程

```bash
# 1. 获取验证码
curl -X GET http://localhost:8080/api/v1/auth/captcha

# 2. 验证验证码
curl -X POST http://localhost:8080/api/v1/auth/verify-captcha \
  -H "Content-Type: application/json" \
  -d '{
    "captcha_id": "abc123",
    "captcha_code": "1234"
  }'

# 3. 用户注册
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "invite_code": "INVITE123456",
    "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

### 2. 上传素材

```bash
# 上传素材文件
curl -X POST http://localhost:8080/api/v1/materials \
  -H "Authorization: Bearer <your_jwt_token>" \
  -F "file=@/path/to/image.jpg" \
  -F "is_public=true"
```

### 3. 创建标签

```bash
# 创建新标签
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "设计素材",
    "color": "#409EFF"
  }'
```

### 4. 创建工作流

```bash
# 创建工作流
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "设计项目工作流",
    "description": "用于管理设计项目的素材",
    "members": [2, 3]
  }'
```

---

## 注意事项

1. **文件上传限制**: 支持图片和视频文件，文件大小限制请参考服务器配置
2. **权限控制**: 用户只能操作自己创建的素材，管理员可以操作所有内容
3. **邀请码**: 新用户注册需要有效的邀请码
4. **工作流协作**: 工作流成员可以共享和协作管理素材
5. **标签系统**: 支持为素材添加多个标签进行分类管理

---

## 相关文档

- [API_DOCUMENT.md](./API_DOCUMENT.md) - 详细API文档（认证、素材、用户、邀请码）
- [API_DOCUMENT_PART2.md](./API_DOCUMENT_PART2.md) - 标签和工作流管理API文档 