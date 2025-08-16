# 桌面端HMAC认证接口

## 概述

桌面端HMAC认证接口为桌面应用程序提供安全的认证机制，使用基于HMAC-SHA256的挑战-响应协议来验证客户端身份。

## 接口信息

- **接口地址**: `POST /api/v1/auth/desktop-auth`
- **请求方式**: POST
- **Content-Type**: application/json

## 认证流程

### 1. 客户端流程

1. 生成随机字符串（nonce）
2. 使用本地密钥和nonce计算HMAC-SHA256摘要
3. 将nonce和HMAC一起发送给服务端

### 2. 服务端流程

1. 接收客户端发送的nonce和HMAC
2. 使用相同的密钥和nonce计算HMAC-SHA256
3. 比较计算结果和客户端发送的HMAC
4. 如果相等，生成并返回auth_token

## 请求参数

```json
{
  "nonce": "随机字符串，建议16字节",
  "hmac": "HMAC-SHA256摘要，使用hex编码"
}
```

### 参数说明

- `nonce`: 随机字符串，用于防止重放攻击
- `hmac`: 使用HMAC-SHA256算法计算的摘要，必须使用hex编码

## 响应格式

### 成功响应 (200 OK)

```json
{
  "auth_token": "生成的认证token",
  "message": "桌面端认证成功"
}
```

### 错误响应

#### 400 Bad Request
```json
{
  "error": "请求参数错误: [具体错误信息]"
}
```

#### 401 Unauthorized
```json
{
  "error": "HMAC验证失败，请求可能被伪造"
}
```

#### 500 Internal Server Error
```json
{
  "error": "生成认证token失败"
}
```

## 配置要求

在环境变量或.env文件中设置HMAC密钥：

```bash
HMAC_SECRET=your-hmac-secret-key
```

## 安全特性

1. **防重放攻击**: 每次请求使用不同的nonce
2. **防伪造**: 只有拥有正确密钥的客户端才能生成有效的HMAC
3. **防篡改**: HMAC确保请求内容完整性
4. **时效性**: auth_token有过期时间限制

## 客户端实现示例

参考 `examples/desktop_client_example.go` 文件中的完整实现。

### 关键代码片段

```go
// 生成随机nonce
nonce, err := generateNonce()

// 计算HMAC
hmacValue := calculateHMAC(nonce, secretKey)

// 发送请求
reqBody := DesktopAuthRequest{
    Nonce: nonce,
    HMAC:  hmacValue,
}
```

## 注意事项

1. **密钥安全**: HMAC密钥必须妥善保管，不应在客户端代码中硬编码
2. **nonce唯一性**: 每次请求都应使用新的随机nonce
3. **HTTPS**: 生产环境建议使用HTTPS传输
4. **密钥轮换**: 定期更换HMAC密钥以提高安全性

## 与现有接口的区别

- **验证码接口**: 需要用户输入验证码，适合Web端
- **桌面端接口**: 基于密钥的自动认证，适合桌面应用
- **返回结果**: 都返回auth_token，可以用于后续的API调用
