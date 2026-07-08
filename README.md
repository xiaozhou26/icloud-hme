# iCloud Hide My Email 本地管理工具

[English](#english) | 中文

通过逆向 iCloud Web 接口和 IMAP 邮件协议，实现 Apple iCloud 隐藏邮箱别名的创建、列出和邮件收取功能。

## 功能特性

- ✅ **创建 HME 别名** — 自动生成 iCloud 隐藏邮箱地址
- ✅ **列出所有别名** — 查看账号下的所有 HME 别名
- ✅ **收取邮件** — 通过 IMAP 读取发到 HME 别名的邮件
- ✅ **多账号管理** — 支持多个 iCloud 账号并行管理
- ✅ **双认证模式** — 支持 Cookie 和 App Password 两种认证方式

## 快速开始

### 1. 安装

```bash
# 前置要求: Go 1.26+
go version  # 确认 Go 版本

# 克隆项目
git clone <your-repo-url>
cd icloud-hme

# 编译
go build -o icloud-hme.exe .
```

### 2. 配置账号

在项目根目录创建 `accounts.json`:

```json
{
  "accounts": [
    {
      "id": "acc_1",
      "name": "主号",
      "cookies": [
        {
          "domain": ".icloud.com",
          "name": "x-apple-session-token",
          "value": "YOUR_SESSION_TOKEN_HERE"
        }
      ],
      "app_passwords": [
        {
          "icloud_email": "your_email@icloud.com",
          "password": "YOUR_APP_PASSWORD_HERE"
        }
      ]
    }
  ]
}
```

### 3. 启动服务

```bash
./icloud-hme.exe

# 服务默认监听 :8080
# 可通过环境变量 PORT 修改端口
PORT=9090 ./icloud-hme.exe
```

## API 接口

### 核心接口

#### 创建 HME 别名

```bash
POST /api/create

# 请求体
{
  "account_id": "acc_1",      # 必填: 账号 ID
  "label": "注册某网站"        # 可选: 别名标签
}

# 响应
{
  "success": true,
  "data": {
    "email": "xyz123@icloud.com",
    "label": "注册某网站",
    "created_at": "2024-01-15T10:30:00Z",
    "account_id": "acc_1"
  }
}
```

#### 读取邮件

```bash
GET /api/inbox?account_id=acc_1&alias=xyz123@icloud.com&limit=20&days=7

# 参数说明:
#   account_id - 必填: 账号 ID
#   alias      - 可选: 只读取发到该别名的邮件
#   limit      - 可选: 返回邮件数量 (默认 20)
#   days       - 可选: 查找最近几天的邮件 (默认 7)

# 响应
{
  "success": true,
  "data": {
    "account_id": "acc_1",
    "alias": "xyz123@icloud.com",
    "count": 5,
    "messages": [
      {
        "uid": 12345,
        "subject": "欢迎注册",
        "from": "noreply@example.com",
        "to": ["xyz123@icloud.com"],
        "date": "2024-01-15T10:35:00Z",
        "body": "感谢您的注册..."
      }
    ]
  }
}
```

### 账号管理接口

#### 列出所有账号

```bash
GET /api/accounts

# 响应
{
  "success": true,
  "data": [
    {"id": "acc_1", "name": "主号"},
    {"id": "acc_2", "name": "副号"}
  ]
}
```

#### 添加账号

**简化版（cookies 可选）:**

```bash
POST /api/accounts

# 请求体
{
  "name": "新账号",
  "host": "icloud.com",           # 可选
  "proxy": "http://..."           # 可选
}

# 响应 - 状态为 pending,需登录
{
  "success": true,
  "data": {
    "id": "acc_xxx",
    "name": "新账号",
    "status": "pending"
  }
}
```

**完整版（带 Cookie）:**

```bash
POST /api/accounts

# 请求体
{
  "name": "新账号",
  "cookies": "{\"x-apple-session-token\":\"token_value\"}",  # JSON 或 Header 格式
  "host": "icloud.com",           # 可选
  "proxy": "http://..."           # 可选
}

# 响应
{
  "success": true,
  "data": {
    "id": "acc_3",
    "name": "新账号",
    "status": "active"
  }
}
```

#### 账号登录（获取 Cookie）

```bash
POST /api/accounts/:id/login

# 请求体
{
  "password": "用户的常规iCloud密码",  # 不是 App Password
  "otp_code": "123456"                  # 可选,2FA 验证码
}

# 响应
{
  "success": true,
  "data": {
    "id": "acc_1",
    "cookies": {
      "x-apple-session-token": "...",
      "X-APPLE-WEBAUTH-TOKEN": "..."
    }
  }
}
```

#### 删除账号

```bash
DELETE /api/accounts/:id

# 响应
{
  "success": true,
  "data": {"id": "acc_3"}
}
```

#### 设置 App Password

```bash
POST /api/accounts/:id/password

# 请求体
{
  "icloud_email": "your_email@icloud.com",
  "app_password": "xxxx-xxxx-xxxx-xxxx"
}

# 响应
{
  "success": true,
  "data": {
    "id": "acc_1",
    "icloud_email": "your_email@icloud.com"
  }
}
```

### 别名管理接口

#### 列出所有别名

```bash
GET /api/aliases?account_id=acc_1

# 响应
{
  "success": true,
  "data": {
    "account_id": "acc_1",
    "count": 15,
    "aliases": [
      {
        "email": "xyz123@icloud.com",
        "label": "注册某网站",
        "created_at": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

#### 停用别名

```bash
POST /api/aliases/:id/deactivate

# 请求体
{
  "account_id": "acc_1"
}

# 响应
{
  "success": true,
  "data": {
    "anonymous_id": "abc123",
    "success": true
  }
}
```

#### 激活别名

```bash
POST /api/aliases/:id/reactivate

# 请求体
{
  "account_id": "acc_1"
}

# 响应
{
  "success": true,
  "data": {
    "anonymous_id": "abc123",
    "success": true
  }
}
```

#### 删除别名

```bash
DELETE /api/aliases/:id

# 请求体
{
  "account_id": "acc_1"
}

# 响应
{
  "success": true,
  "data": {
    "anonymous_id": "abc123"
  }
}
```

## 认证方式

### 方式一: Cookie 认证 (推荐)

Cookie 认证可以创建别名和读取邮件。

**获取 Cookie:**

1. 使用浏览器登录 [iCloud.com](https://www.icloud.com)
2. 打开浏览器开发者工具 (F12)
3. 进入 Application → Cookies → icloud.com
4. 找到以下 Cookie:
   - `x-apple-session-token` (必需)
   - `X-APPLE-WEBAUTH-TOKEN` (可选)
   - `X-APPLE-WEBAUTH-USER` (可选)

**完整 Cookie JSON 格式:**

```json
[
  {
    "domain": ".icloud.com",
    "expirationDate": 1735689600,
    "hostOnly": false,
    "httpOnly": true,
    "name": "x-apple-session-token",
    "path": "/",
    "sameSite": "unspecified",
    "secure": true,
    "session": false,
    "storeId": "0",
    "value": "YOUR_SESSION_TOKEN_VALUE"
  }
]
```

### 方式二: App Password 认证

App Password 认证仅用于读取邮件 (IMAP)。

**生成 App Password:**

1. 登录 [appleid.apple.com](https://appleid.apple.com)
2. 进入 "登录和安全" → "App 专用密码"
3. 生成新密码，用于此工具

## 项目架构

```
icloud-hme/
├── main.go                 # 入口: 加载配置、初始化管理器、启动服务
├── accounts.json           # 账号配置文件 (自动生成)
├── go.mod
└── internal/
    ├── account/
    │   └── manager.go      # 多账号管理器 (持久化、客户端工厂)
    ├── hme/
    │   └── client.go       # iCloud HME Web 客户端 (Cookie + App Password)
    ├── mail/
    │   └── client.go       # IMAP 邮件客户端 (邮件读取、搜索)
    └── server/
        └── server.go       # HTTP API (Gin 路由 + 请求处理)
```

### 核心模块

- **account.Manager**: 管理多个 iCloud 账号，负责配置持久化和客户端创建
- **hme.Client**: 封装 iCloud HME Web API，支持 Cookie 和 App Password 双认证
- **mail.Client**: IMAP 邮件客户端，负责连接 iCloud 邮箱、搜索和读取邮件
- **server.Server**: HTTP API 服务，提供 RESTful 接口

## 技术栈

- **Go 1.26+**
- **Gin** — HTTP 框架
- **go-imap** — IMAP 协议实现
- **tls-client** — TLS 指纹模拟 (绕过 iCloud 反爬)

## 常见问题

### Q: 创建别名返回 401/403 错误?

**A:** Cookie 已过期，需要重新获取。iCloud Cookie 有效期通常为 24 小时。

### Q: 读取邮件返回超时?

**A:** 检查网络连接，确保可以访问 `imap.mail.me.com:993`。

### Q: 如何查看某个别名收到了哪些邮件?

**A:** 调用 `GET /api/inbox?account_id=acc_1&alias=your_alias@icloud.com`

### Q: 支持同时管理多个 iCloud 账号吗?

**A:** 支持，在 `accounts.json` 中配置多个账号即可，每个账号有独立的 `id`。

## 开发指南

### 本地开发

```bash
# 安装依赖
go mod download

# 运行 (开发模式，带日志)
go run main.go

# 编译
go build -o icloud-hme.exe .

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o icloud-hme .
```

### 代码规范

- 代码注释使用中文
- 错误信息返回给用户时使用中文
- API 响应格式统一: `{success: bool, data: any, message: string}`

## 许可证

MIT License

---

## English

A local management tool for Apple iCloud Hide My Email (HME) aliases, supporting creation, listing, and email reading through reverse-engineered iCloud Web API and IMAP protocol.

### Features

- Create HME aliases automatically
- List all aliases for an account
- Read emails sent to HME aliases via IMAP
- Manage multiple iCloud accounts
- Dual authentication: Cookie and App Password

### Quick Start

```bash
# Build
go build -o icloud-hme.exe .

# Create accounts.json with your credentials
# Run
./icloud-hme.exe

# API endpoints
# POST /api/create     - Create HME alias
# GET  /api/inbox      - Read emails
# GET  /api/aliases    - List aliases
```

See [API Documentation](#api-接口) for detailed usage.
