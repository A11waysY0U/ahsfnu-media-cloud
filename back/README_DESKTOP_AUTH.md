# 桌面端HMAC认证功能

## 新增功能概述

为Go后端新增了桌面端HMAC认证接口，提供与现有验证码接口相同的功能（返回auth_token），但使用更安全的HMAC-SHA256认证机制。

## 新增文件

### 1. 核心接口文件
- `internal/api/auth/desktop.go` - 桌面端认证接口实现

### 2. 配置文件更新
- `internal/config/config.go` - 新增HMAC配置结构

### 3. 路由配置更新
- `internal/api/routes.go` - 新增桌面端认证路由

### 4. 示例和文档
- `examples/desktop_client_example.go` - 客户端使用示例
- `docs/desktop_auth.md` - 详细接口文档
- `internal/api/auth/desktop_test.go` - 接口测试文件

## 接口信息

- **接口地址**: `POST /api/v1/auth/desktop-auth`
- **功能**: 桌面端HMAC认证，返回auth_token
- **认证方式**: HMAC-SHA256

## 配置要求

在环境变量中设置HMAC密钥：

```bash
HMAC_SECRET=your-hmac-secret-key
```

## 使用方法

### 1. 启动服务
```bash
go run cmd/main.go
```

### 2. 客户端调用
参考 `examples/desktop_client_example.go` 中的实现：

```go
// 生成随机nonce
nonce := generateNonce()

// 计算HMAC
hmacValue := calculateHMAC(nonce, secretKey)

// 发送认证请求
reqBody := DesktopAuthRequest{
    Nonce: nonce,
    HMAC:  hmacValue,
}
```

### 3. 运行测试
```bash
go test ./internal/api/auth/desktop_test.go ./internal/api/auth/desktop.go -v
```

## 安全特性

1. **防重放攻击**: 每次请求使用不同的nonce
2. **防伪造**: 只有拥有正确密钥的客户端才能生成有效的HMAC
3. **防篡改**: HMAC确保请求内容完整性
4. **时效性**: auth_token有过期时间限制

## 与现有接口的区别

| 特性 | 验证码接口 | 桌面端接口 |
|------|------------|------------|
| 认证方式 | 用户输入验证码 | HMAC-SHA256 |
| 适用场景 | Web端用户交互 | 桌面应用自动认证 |
| 安全性 | 中等 | 高 |
| 用户体验 | 需要用户操作 | 完全自动化 |

## 注意事项

1. **密钥安全**: HMAC密钥必须妥善保管
2. **HTTPS**: 生产环境建议使用HTTPS
3. **密钥轮换**: 定期更换HMAC密钥
4. **nonce唯一性**: 每次请求使用新的随机nonce

## 下一步

1. 在生产环境中设置安全的HMAC密钥
2. 考虑添加请求频率限制
3. 监控认证失败的情况
4. 定期审查和更新安全策略
