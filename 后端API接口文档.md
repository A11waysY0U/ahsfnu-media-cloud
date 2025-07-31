# 后端API接口文档

> 除 `/api/v1/public/materials` 外，所有接口均需登录（带 token）。

---

## 认证相关

| 方法 | 路径                | 说明         | 参数/Body         |
|------|---------------------|--------------|-------------------|
| POST | /api/v1/auth/login  | 用户登录     | username, password|
| POST | /api/v1/auth/register | 用户注册   | username, email, password |
| GET  | /api/v1/profile     | 获取当前用户 | 无                |
| GET  | /api/v1/users       | 获取用户列表 | 无                |

---

## 素材相关

| 方法 | 路径                        | 说明           | 参数/Body                |
|------|-----------------------------|----------------|--------------------------|
| GET  | /api/v1/public/materials    | 获取公开素材   | page, page_size, keyword, tag（标签名模糊）|
| GET  | /api/v1/materials           | 获取素材列表   | page, page_size, keyword, file_type, workflow_id, tags（如tags=1,2,3）|
| POST | /api/v1/materials           | 上传素材       | formData: file, workflow_id（可选）|
| GET  | /api/v1/materials/:id       | 获取素材详情   | id                       |
| PUT  | /api/v1/materials/:id       | 更新素材       | original_filename, is_starred, is_favorite, is_public, workflow_id |
| DELETE|/api/v1/materials/:id       | 删除素材       | id                       |
| POST | /api/v1/materials/:id/star  | 切换星标       | id                       |
| POST | /api/v1/materials/:id/favorite | 切换收藏    | id                       |
| POST | /api/v1/materials/:id/public  | 切换公开状态 | id                       |
| PUT  | /api/v1/materials/batch_workflow | 批量修改工作流 | material_ids, workflow_id |

---

## 标签相关

| 方法 | 路径                                         | 说明           | 参数/Body         |
|------|----------------------------------------------|----------------|-------------------|
| GET  | /api/v1/tags                                 | 获取标签列表   | 无                |
| POST | /api/v1/tags                                 | 创建标签       | name, color       |
| PUT  | /api/v1/tags/:id                             | 更新标签       | name, color       |
| DELETE|/api/v1/tags/:id                             | 删除标签       | id                |
| POST | /api/v1/tags/:id/materials/:materialId       | 给素材打标签   | id, materialId    |
| DELETE|/api/v1/tags/:id/materials/:materialId       | 移除素材标签   | id, materialId    |

---

## 工作流相关

| 方法 | 路径                                         | 说明           | 参数/Body         |
|------|----------------------------------------------|----------------|-------------------|
| GET  | /api/v1/workflows                            | 获取工作流列表 | page, page_size, keyword |
| POST | /api/v1/workflows                            | 创建工作流     | name, description |
| GET  | /api/v1/workflows/:id                        | 获取详情       | id                |
| PUT  | /api/v1/workflows/:id                        | 更新工作流     | name, description |
| DELETE|/api/v1/workflows/:id                        | 删除工作流     | id                |
| POST | /api/v1/workflows/:id/members                | 添加成员       | user_id           |
| DELETE|/api/v1/workflows/:id/members/:userId        | 移除成员       | userId            |

---

> 如需详细字段说明或示例请求，请联系开发者。 