```txt
r2-gallery/
├── .env                     # 配置 Cloudflare R2 访问密钥、数据库配置等
├── main.go                  # 入口文件
├── go.mod                   # Go 依赖管理
├── config/
│   ├── config.go            # 读取 .env 配置
│   ├── r2.go                # Cloudflare R2 配置和初始化
│   ├── database.go          # 数据库初始化
├── models/
│   ├── user.go              # 用户模型
│   ├── image.go             # 图片模型
├── routes/
│   ├── auth.go              # 认证相关 API（登录、注册）
│   ├── image.go             # 图片管理 API（上传、查询、删除）
│   ├── user.go              # 用户管理 API
├── controllers/
│   ├── auth_controller.go   # 认证业务逻辑
│   ├── image_controller.go  # 图片上传、查询、删除业务逻辑
│   ├── user_controller.go   # 用户管理业务逻辑
├── middleware/
│   ├── auth_middleware.go   # JWT 身份验证中间件
├── services/
│   ├── r2_service.go        # Cloudflare R2 文件上传、删除、获取 URL
│   ├── jwt_service.go       # JWT 令牌生成和解析
├── utils/
│   ├── hash.go              # 密码哈希工具
│   ├── response.go          # 统一 API 响应格式
│   ├── logger.go            # 日志工具
└── docs/                    # API 文档（Swagger 或 Postman Collection）
```
