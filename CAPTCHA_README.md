# 验证码功能说明

本项目已集成人机验证功能，在登录和注册前需要先通过验证码验证。

## 功能特性

- ✅ 基于 base64Captcha 的数字验证码
- ✅ 验证码有效期：5分钟
- ✅ 认证token有效期：10分钟
- ✅ 验证成功后自动删除验证码，防止重复使用
- ✅ 支持用户名或邮箱登录

## 使用流程

### 1. 获取验证码

```bash
GET /api/v1/auth/captcha
```

**响应示例：**
```json
{
  "captcha_id": "abc123def456",
  "captcha_b64": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...",
  "auth_token": "eyJ0b2tlbiI6ImFiYzEyMyIsImV4cCI6MTYz..."
}
```

### 2. 验证验证码

```bash
POST /api/v1/auth/verify-captcha
Content-Type: application/json

{
  "captcha_id": "abc123def456",
  "captcha_code": "1234"
}
```

**响应示例：**
```json
{
  "auth_token": "eyJ0b2tlbiI6ImRlZjQ1NiIsImV4cCI6MTYz...",
  "message": "验证码验证成功"
}
```

### 3. 登录（需要认证token）

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123",
  "auth_token": "eyJ0b2tlbiI6ImRlZjQ1NiIsImV4cCI6MTYz..."
}
```

### 4. 注册（需要认证token）

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "invite_code": "INVITE123",
  "auth_token": "eyJ0b2tlbiI6ImRlZjQ1NiIsImV4cCI6MTYz..."
}
```

## 测试

运行测试脚本：

```bash
python test_captcha.py
```

测试脚本会：
1. 获取验证码并保存图片
2. 提示输入验证码
3. 验证验证码
4. 测试登录和注册功能

## 安全说明

- 验证码图片使用Base64编码返回
- 验证码答案存储在服务器内存中，5分钟后自动过期
- 认证token使用随机字节生成，10分钟后过期
- 验证码验证成功后立即删除，防止重复使用
- 登录和注册必须提供有效的认证token

## 错误处理

- 验证码错误：返回 400 状态码
- 认证token无效：返回 401 状态码
- 验证码过期：返回 400 状态码
- 服务器错误：返回 500 状态码

## 依赖

- `github.com/mojocn/base64Captcha` - 验证码生成库
- `crypto/rand` - 随机数生成
- `encoding/base64` - Base64编码
- `encoding/json` - JSON处理 