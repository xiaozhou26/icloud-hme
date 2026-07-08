# iCloud Hide My Email API 文档

## 概述

HTTP JSON API，所有接口返回统一格式：

```json
{
  "success": true,
  "data": {},
  "message": ""
}
```

**错误响应:**
- `400 Bad Request` — 参数错误
- `401 Unauthorized` — 会话失效
- `404 Not Found` — 账号不存在
- `502 Bad Gateway` — iCloud 服务错误

---

## 核心接口

### 1. 创建 HME 别名

```http
POST /api/create
Content-Type: application/json

{
  "account_id": "acc_1",
  "label": "注册某网站"
}
```

**响应:**
```json
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

**参数说明:**
- `account_id` (必填) — 账号 ID
- `label` (可选) — 别名标签，默认为 "Created YYYY-MM-DD HH:mm"

**错误情况:**
- `401` — Cookie 过期，需更新
- `502` — iCloud 服务错误，会自动重试 5 次

---

### 2. 读取邮件

```http
GET /api/inbox?account_id=acc_1&alias=xyz123@icloud.com&limit=20&days=7
```

**响应:**
```json
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

**参数说明:**
- `account_id` (必填) — 账号 ID
- `alias` (可选) — 只返回发到该别名的邮件
- `limit` (可选) — 返回邮件数量，默认 20
- `days` (可选) — 查找最近几天的邮件，默认 7

**邮件搜索逻辑:**
- 原生 IMAP SEARCH: `OR (TO alias) (CC alias)`
- 本地兜底: 过滤 To/CC/BCC 包含别名的邮件

---

## 账号管理接口

### 3. 列出所有账号

```http
GET /api/accounts
```

**响应:**
```json
{
  "success": true,
  "data": [
    {
      "id": "acc_1",
      "name": "主号",
      "host": "imap.mail.me.com"
    }
  ]
}
```

**注意:** 响应中不包含敏感信息（cookies、app_passwords）

---

### 4. 添加账号

**简化版（cookies 可选）:**
```http
POST /api/accounts
Content-Type: application/json

{
  "name": "新账号",
  "host": "icloud.com",
  "proxy": "http://user:pass@host:port"
}
```

**完整版（包含 Cookie）:**
```http
POST /api/accounts
Content-Type: application/json

{
  "name": "新账号",
  "cookies": "{\"x-apple-session-token\":\"token_value\"}",
  "host": "icloud.com",
  "proxy": "http://user:pass@host:port"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "id": "acc_3",
    "name": "新账号",
    "host": "icloud.com",
    "status": "pending"
  }
}
```

**参数说明:**
- `name` (必填) — 账号名称
- `cookies` (可选) — Cookie 字符串,支持两种格式:
  - JSON: `"{\"name\":\"value\"}"`
  - Header: `"name1=value1; name2=value2"`
- `host` (可选) — iCloud 域名,默认 `icloud.com`
- `proxy` (可选) — HTTP/SOCKS5 代理

**注意:** 不传 cookies 时,账号状态为 `pending`,需通过 `/login` 接口登录获取 Cookie

---

### 5. 账号密码登录（获取 Cookie）

```http
POST /api/accounts/:id/login
Content-Type: application/json

{
  "password": "用户的常规iCloud密码",
  "otp_code": "123456"  // 可选,2FA 验证码
}
```

**参数说明:**
- `:id` (路径参数) — 账号 ID
- `password` (必填) — iCloud 账号的常规密码(**不是** App Password)
- `otp_code` (可选) — 双重认证验证码

**响应:**
```json
{
  "success": true,
  "data": {
    "id": "acc_1",
    "cookies": {
      "x-apple-session-token": "...",
      "X-APPLE-WEBAUTH-TOKEN": "...",
      "X-APPLE-WEBAUTH-USER": "..."
    }
  }
}
```

**注意事项:**
- 密码是登录 appleid.apple.com 的**常规账号密码**,不是 App 专用密码
- 登录前账号必须已设置 `icloud_email` 字段
- 登录成功后 Cookie 会自动保存到 accounts.json
- 启用 2FA 时,第一次请求会被拒绝,需要带 `otp_code` 重试

---

### 6. 删除账号

**Cookie 格式:**
```json
[
  {
    "domain": ".icloud.com",
    "name": "x-apple-session-token",
    "value": "YOUR_TOKEN"
  },
  {
    "domain": ".icloud.com",
    "name": "X-APPLE-WEBAUTH-TOKEN",
    "value": "YOUR_TOKEN"
  }
]
```

---

### 5. 删除账号

```http
DELETE /api/accounts/:id
```

**响应:**
```json
{
  "success": true,
  "data": {
    "id": "acc_3"
  }
}
```

**错误情况:**
- `404` — 账号不存在

---

### 6. 设置 App Password

```http
POST /api/accounts/:id/password
Content-Type: application/json

{
  "icloud_email": "your_email@icloud.com",
  "app_password": "xxxx-xxxx-xxxx-xxxx"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "id": "acc_1",
    "icloud_email": "your_email@icloud.com"
  }
}
```

**参数说明:**
- `icloud_email` (必填) — iCloud 邮箱地址
- `app_password` (必填) — App 专用密码

**用途:** App Password 用于 IMAP 邮件读取，生成方式见 [appleid.apple.com](https://appleid.apple.com)

---

## 别名管理接口

### 7. 列出所有别名

```http
GET /api/aliases?account_id=acc_1
```

**响应:**
```json
{
  "success": true,
  "data": {
    "account_id": "acc_1",
    "count": 15,
    "aliases": [
      {
        "email": "xyz123@icloud.com",
        "anonymousId": "abc123",
        "label": "注册某网站",
        "active": true,
        "createdAt": "2024-01-15T10:30:00Z"
      }
    ]
  }
}
```

**参数说明:**
- `account_id` (必填) — 账号 ID

**别名字段:**
- `email` — HME 邮箱地址
- `anonymousId` — 别名唯一标识（用于停用/激活/删除）
- `label` — 用户定义的标签
- `active` — 是否激活
- `createdAt` — 创建时间

---

### 8. 停用别名

```http
POST /api/aliases/:id/deactivate
Content-Type: application/json

{
  "account_id": "acc_1"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "anonymous_id": "abc123",
    "success": true
  }
}
```

**参数说明:**
- `:id` (路径参数) — 别名的 `anonymousId`
- `account_id` (必填) — 账号 ID

**说明:** 停用后别名不再接收邮件，但可随时激活恢复

---

### 9. 激活别名

```http
POST /api/aliases/:id/reactivate
Content-Type: application/json

{
  "account_id": "acc_1"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "anonymous_id": "abc123",
    "success": true
  }
}
```

**参数说明:**
- `:id` (路径参数) — 别名的 `anonymousId`
- `account_id` (必填) — 账号 ID

**说明:** 激活已停用的别名，恢复邮件接收

---

### 10. 删除别名

```http
DELETE /api/aliases/:id
Content-Type: application/json

{
  "account_id": "acc_1"
}
```

**响应:**
```json
{
  "success": true,
  "data": {
    "anonymous_id": "abc123"
  }
}
```

**参数说明:**
- `:id` (路径参数) — 别名的 `anonymousId`
- `account_id` (必填) — 账号 ID

**注意:** 删除不可恢复！如果直接删除失败，会先停用再删除

---

## 使用示例

### curl 示例

```bash
# 创建别名
curl -X POST http://localhost:8080/api/create \
  -H "Content-Type: application/json" \
  -d '{"account_id": "acc_1", "label": "GitHub"}'

# 读取邮件
curl "http://localhost:8080/api/inbox?account_id=acc_1&alias=xyz123@icloud.com&limit=10"

# 列出别名
curl "http://localhost:8080/api/aliases?account_id=acc_1"

# 停用别名
curl -X POST http://localhost:8080/api/aliases/abc123/deactivate \
  -H "Content-Type: application/json" \
  -d '{"account_id": "acc_1"}'

# 删除别名
curl -X DELETE http://localhost:8080/api/aliases/abc123 \
  -H "Content-Type: application/json" \
  -d '{"account_id": "acc_1"}'
```

### Python 示例

```python
import requests

BASE_URL = "http://localhost:8080/api"

# 创建别名
resp = requests.post(f"{BASE_URL}/create", json={
    "account_id": "acc_1",
    "label": "Netflix"
})
print(resp.json())

# 读取邮件
resp = requests.get(f"{BASE_URL}/inbox", params={
    "account_id": "acc_1",
    "alias": "xyz123@icloud.com",
    "limit": 10
})
print(resp.json())

# 列出别名
resp = requests.get(f"{BASE_URL}/aliases", params={"account_id": "acc_1"})
for alias in resp.json()["data"]["aliases"]:
    print(f"{alias['email']} - {alias['label']} (active: {alias['active']})")
```

---

## 认证说明

### Cookie 认证 (高级功能)

用于：创建别名、列出别名、停用/激活/删除别名

**获取方式:**
1. 浏览器登录 [icloud.com](https://www.icloud.com)
2. F12 → Application → Cookies → icloud.com
3. 复制 `x-apple-session-token` 值

**有效期:** 约 24 小时

### App Password 认证 (IMAP)

用于：读取邮件

**获取方式:**
1. 登录 [appleid.apple.com](https://appleid.apple.com)
2. 登录和安全 → App 专用密码
3. 生成新密码

---

## 错误处理

### 会话失效 (401)

```json
{
  "success": false,
  "message": "iCloud 会话失效，请更新 Cookie: HTTP 401"
}
```

**解决:** 更新 `accounts.json` 中的 Cookie

### iCloud 服务错误 (502)

```json
{
  "success": false,
  "message": "创建邮箱失败: HTTP 429"
}
```

**说明:** 429 错误会自动重试最多 5 次

### 参数错误 (400)

```json
{
  "success": false,
  "message": "参数错误: account_id 必填"
}
```

---

## 限制

- **创建频率**: iCloud 限制别名创建频率，过快会返回 429
- **Cookie 有效期**: 约 24 小时，需定期更新
- **邮件读取**: 依赖 IMAP 连接，超时默认 30 秒
