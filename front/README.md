# AHSFNU 媒体云平台 - 前端

基于 Vue 3 + TypeScript + Element Plus 构建的现代化媒体管理平台前端应用。

## 🚀 技术栈

- **框架**: Vue 3 (Composition API)
- **语言**: TypeScript
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **构建工具**: Vite
- **HTTP 客户端**: Axios
- **代码规范**: ESLint + Prettier

## ✨ 功能特性

### 🔐 用户认证
- 用户登录/注册
- 验证码验证
- JWT 令牌管理
- 邀请码注册系统

### 📁 素材管理
- 图片/视频上传
- 素材预览和播放
- 批量操作支持
- 星标收藏功能
- 标签分类管理
- 工作流关联

### 🏷️ 标签系统
- 创建和管理标签
- 颜色自定义
- 标签筛选和搜索
- 素材标签关联

### 🔄 工作流管理
- 工作流创建和配置
- 成员权限管理
- 工作流类型分类
- JSON 配置支持

### 👥 用户管理 (管理员)
- 用户列表查看
- 角色权限管理
- 用户信息编辑
- 邀请码生成和管理

### 📊 数据统计
- 个人使用统计
- 素材数量统计
- 最近活动记录

## 🛠️ 开发环境

### 环境要求
- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 代码检查
```bash
npm run lint
```

### 代码格式化
```bash
npm run format
```

## 📁 项目结构

```
front/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口定义
│   │   └── index.ts       # API 客户端配置
│   ├── assets/            # 资源文件
│   │   └── main.css       # 全局样式
│   ├── components/        # 公共组件
│   ├── layouts/           # 布局组件
│   │   └── MainLayout.vue # 主布局
│   ├── router/            # 路由配置
│   │   └── index.ts       # 路由定义
│   ├── stores/            # 状态管理
│   │   └── auth.ts        # 认证状态
│   ├── types/             # TypeScript 类型定义
│   │   └── index.ts       # 全局类型
│   ├── views/             # 页面组件
│   │   ├── Login.vue      # 登录页
│   │   ├── Register.vue   # 注册页
│   │   ├── Dashboard.vue  # 仪表板
│   │   ├── Materials.vue  # 素材管理
│   │   ├── Tags.vue       # 标签管理
│   │   ├── Workflows.vue  # 工作流管理
│   │   ├── Profile.vue    # 个人资料
│   │   ├── Users.vue      # 用户管理
│   │   └── InviteCodes.vue # 邀请码管理
│   ├── App.vue            # 根组件
│   └── main.ts            # 应用入口
├── .eslintrc.cjs          # ESLint 配置
├── .prettierrc.json       # Prettier 配置
├── index.html             # HTML 模板
├── package.json           # 项目配置
├── tsconfig.app.json      # TypeScript 配置
├── vite.config.ts         # Vite 配置
└── README.md              # 项目说明
```

## 🎨 设计特色

### 现代化 UI 设计
- 采用 Element Plus 设计语言
- 响应式布局设计
- 优雅的动画效果
- 直观的用户交互

### 用户体验优化
- 加载状态提示
- 错误处理机制
- 表单验证反馈
- 操作确认对话框

### 性能优化
- 组件懒加载
- 图片懒加载
- 路由缓存
- 状态持久化

## 🔧 配置说明

### 开发环境配置
项目使用 Vite 作为构建工具，主要配置在 `vite.config.ts` 中：

- **端口**: 3000
- **API 代理**: 自动代理 `/api` 和 `/uploads` 到后端服务器
- **路径别名**: `@` 指向 `src` 目录

### API 配置
API 客户端配置在 `src/api/index.ts` 中：

- **基础 URL**: `/api/v1`
- **超时时间**: 10 秒
- **请求拦截器**: 自动添加 JWT 令牌
- **响应拦截器**: 统一错误处理

## 🚀 部署说明

### 构建生产版本
```bash
npm run build
```

构建完成后，`dist` 目录包含所有静态文件。

### 部署到 Web 服务器
将 `dist` 目录的内容部署到 Web 服务器的根目录。

### Nginx 配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /uploads {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📝 更新日志

### v1.0.0 (2024-01-01)
- ✨ 初始版本发布
- 🔐 用户认证系统
- 📁 素材管理功能
- 🏷️ 标签系统
- 🔄 工作流管理
- 👥 用户管理
- 📊 数据统计

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 项目地址: [GitHub Repository]
- 问题反馈: [Issues]
- 邮箱: [your-email@example.com]

---

**AHSFNU 媒体云平台** - 让媒体管理更简单、更高效！
