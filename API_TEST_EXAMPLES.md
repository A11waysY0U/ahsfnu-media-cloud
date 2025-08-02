# AHSFNU Media Cloud API 测试示例

## 环境准备

### 1. 启动服务器
```bash
# 编译并运行服务器
go build -o main.exe main.go
./main.exe
```

### 2. 测试工具
推荐使用以下工具进行API测试：
- **Postman**: 图形化API测试工具
- **curl**: 命令行工具
- **Python requests**: 编程方式测试

---

## 完整测试流程

### 步骤1: 获取验证码

```bash
curl -X GET http://localhost:8080/api/v1/auth/captcha
```

**预期响应**:
```json
{
  "captcha_id": "abc123def456",
  "captcha_b64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 步骤2: 验证验证码

```bash
curl -X POST http://localhost:8080/api/v1/auth/verify-captcha \
  -H "Content-Type: application/json" \
  -d '{
    "captcha_id": "abc123def456",
    "captcha_code": "1234"
  }'
```

**预期响应**:
```json
{
  "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "验证码验证成功"
}
```

### 步骤3: 用户注册

```bash
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

**预期响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 步骤4: 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "auth_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

**预期响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 步骤5: 上传素材

```bash
# 保存JWT token
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 上传素材
curl -X POST http://localhost:8080/api/v1/materials \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.jpg" \
  -F "is_public=true"
```

**预期响应**:
```json
{
  "id": 1,
  "filename": "abc123.jpg",
  "original_filename": "image.jpg",
  "file_path": "/uploads/abc123.jpg",
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
  "is_public": true,
  "thumbnail_path": "/uploads/thumbnails/abc123.jpg",
  "uploader": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 步骤6: 创建标签

```bash
curl -X POST http://localhost:8080/api/v1/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "设计素材",
    "color": "#409EFF"
  }'
```

**预期响应**:
```json
{
  "id": 1,
  "name": "设计素材",
  "color": "#409EFF",
  "created_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "creator": {
    "id": 1,
    "username": "testuser"
  }
}
```

### 步骤7: 为素材添加标签

```bash
curl -X POST http://localhost:8080/api/v1/tags/1/materials/1 \
  -H "Authorization: Bearer $TOKEN"
```

**预期响应**:
```json
{
  "message": "标签添加成功"
}
```

### 步骤8: 创建工作流

```bash
curl -X POST http://localhost:8080/api/v1/workflows \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "设计项目工作流",
    "description": "用于管理设计项目的素材",
    "members": []
  }'
```

**预期响应**:
```json
{
  "id": 1,
  "name": "设计项目工作流",
  "description": "用于管理设计项目的素材",
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

### 步骤9: 搜索素材

```bash
curl -X GET "http://localhost:8080/api/v1/materials?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

**预期响应**:
```json
{
  "data": [
    {
      "id": 1,
      "filename": "abc123.jpg",
      "original_filename": "image.jpg",
      "file_path": "/uploads/abc123.jpg",
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
      "is_public": true,
      "thumbnail_path": "/uploads/thumbnails/abc123.jpg",
      "uploader": {
        "id": 1,
        "username": "testuser"
      },
      "material_tags": [
        {
          "id": 1,
          "tag": {
            "id": 1,
            "name": "设计素材",
            "color": "#409EFF"
          }
        }
      ]
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 1
  }
}
```

---

## Python 测试脚本

```python
import requests
import json

# 基础配置
BASE_URL = "http://localhost:8080"
API_BASE = f"{BASE_URL}/api/v1"

class MediaCloudAPI:
    def __init__(self):
        self.session = requests.Session()
        self.token = None
    
    def get_captcha(self):
        """获取验证码"""
        response = self.session.get(f"{API_BASE}/auth/captcha")
        return response.json()
    
    def verify_captcha(self, captcha_id, captcha_code):
        """验证验证码"""
        data = {
            "captcha_id": captcha_id,
            "captcha_code": captcha_code
        }
        response = self.session.post(f"{API_BASE}/auth/verify-captcha", json=data)
        return response.json()
    
    def register(self, username, email, password, invite_code, auth_token):
        """用户注册"""
        data = {
            "username": username,
            "email": email,
            "password": password,
            "invite_code": invite_code,
            "auth_token": auth_token
        }
        response = self.session.post(f"{API_BASE}/auth/register", json=data)
        result = response.json()
        if "token" in result:
            self.token = result["token"]
        return result
    
    def login(self, username, password, auth_token):
        """用户登录"""
        data = {
            "username": username,
            "password": password,
            "auth_token": auth_token
        }
        response = self.session.post(f"{API_BASE}/auth/login", json=data)
        result = response.json()
        if "token" in result:
            self.token = result["token"]
        return result
    
    def upload_material(self, file_path, is_public=False):
        """上传素材"""
        headers = {"Authorization": f"Bearer {self.token}"}
        files = {"file": open(file_path, "rb")}
        data = {"is_public": is_public}
        response = self.session.post(f"{API_BASE}/materials", headers=headers, files=files, data=data)
        return response.json()
    
    def create_tag(self, name, color="#409EFF"):
        """创建标签"""
        headers = {"Authorization": f"Bearer {self.token}"}
        data = {"name": name, "color": color}
        response = self.session.post(f"{API_BASE}/tags", headers=headers, json=data)
        return response.json()
    
    def search_materials(self, page=1, page_size=10):
        """搜索素材"""
        headers = {"Authorization": f"Bearer {self.token}"}
        params = {"page": page, "page_size": page_size}
        response = self.session.get(f"{API_BASE}/materials", headers=headers, params=params)
        return response.json()

# 使用示例
if __name__ == "__main__":
    api = MediaCloudAPI()
    
    # 1. 获取验证码
    captcha_data = api.get_captcha()
    print("验证码数据:", captcha_data)
    
    # 2. 验证验证码（需要用户输入验证码）
    captcha_code = input("请输入验证码: ")
    verify_result = api.verify_captcha(captcha_data["captcha_id"], captcha_code)
    print("验证结果:", verify_result)
    
    # 3. 注册用户
    register_result = api.register(
        username="testuser",
        email="test@example.com",
        password="password123",
        invite_code="INVITE123456",
        auth_token=verify_result["auth_token"]
    )
    print("注册结果:", register_result)
    
    # 4. 创建标签
    tag_result = api.create_tag("设计素材", "#409EFF")
    print("标签创建结果:", tag_result)
    
    # 5. 搜索素材
    materials = api.search_materials()
    print("素材列表:", materials)
```

---

## 常见问题

### 1. 验证码错误
- 确保验证码ID和验证码匹配
- 验证码区分大小写

### 2. 认证失败
- 检查JWT token是否有效
- 确保token格式正确：`Bearer <token>`

### 3. 文件上传失败
- 检查文件格式是否支持
- 确认文件大小在限制范围内
- 验证文件路径是否正确

### 4. 权限不足
- 确认用户角色是否有相应权限
- 检查是否为资源所有者或管理员

---

## 测试检查清单

- [ ] 验证码获取和验证
- [ ] 用户注册和登录
- [ ] 素材上传和管理
- [ ] 标签创建和管理
- [ ] 工作流创建和管理
- [ ] 权限控制测试
- [ ] 错误处理测试
- [ ] 分页功能测试 