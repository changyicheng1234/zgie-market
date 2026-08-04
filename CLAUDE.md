# 智工集市项目开发记录

> 智能工程学院内部交流平台，基于 [软工集市](https://ssemarket.cn) 开源项目改造。
> 原项目仓库：https://gitee.com/yang-peiyue/SSE_market_server

## 项目结构

```
学院集市/
├── SSE_market_server/   # 后端 Go + Gin（人A负责）
├── SSE_market_client/   # PC端前端 Vue 2（人B负责）
└── sse_market_mobile/   # 移动端 Vue 2 + Vant（人B负责）
```

## 技术栈

| 模块 | 技术 |
|------|------|
| 后端 | Go 1.23 + Gin + GORM |
| 数据库 | MySQL 8.1 |
| 缓存 | Redis |
| 文件存储 | 腾讯云 COS |
| 内容审核 | 腾讯云 TMS |
| 前端 | Vue 2 + Bootstrap-Vue |
| 移动端 | Vue 2 + Vant 2 |
| 部署 | Docker Compose + Nginx |

---

## 人A — 后端 + 运维任务清单

### 阶段一：配置文件
- [x] 确定学院名称：**智能工程学院 / 智工集市**
- [x] 创建 `SSE_market_server/config/application.yml`（MySQL/Redis/COS/TMS/JWT 模板已创建）
- [x] 创建 `SSE_market_server/.env`（SMTP邮件模板已创建）
- [x] 修改 `config/emailConst.go` 品牌文字（"软工集市" → "智工集市" 全部替换完成）
- [x] TMS 密钥移入 application.yml 统一管理（原本硬编码）
- [ ] **填入真实密钥**：`application.yml` 中所有 TODO 项（COS/TMS/JWT/数据库密码）
- [ ] **填入邮箱信息**：`.env` 中所有 TODO 项（SMTP账号/授权码）

> ⚠️ 注意：`config/emailConst.go` 邮件模板中 logo 图片 URL 仍指向 `ssemarket.cn`，
> 等有自己的域名和 logo 后替换（搜索 `ssemarket.cn/new/android-chrome-192x192.png`）

### 阶段二：第三方服务开通（照搬软工集市方案）
- [ ] 腾讯云账号注册/登录：https://console.cloud.tencent.com
- [ ] 开通 COS 存储桶，获取 BucketName / AppId / SecretId / SecretKey
- [ ] 开通 TMS 文本内容安全，获取密钥（可复用 COS 同一子账号）
- [ ] 准备 SMTP 邮箱（QQ邮箱推荐），开启授权码

### 阶段三：本地开发环境验证
- [ ] 安装 Go 1.23+
- [ ] 安装 MySQL 8.1，创建数据库（名称与 application.yml 的 `database` 字段一致）
- [ ] 安装 Redis
- [ ] 本地跑通 `go run main.go`
- [ ] 验证注册/登录/发帖流程

### 阶段四：服务器部署
- [ ] 获取服务器（学院提供 or 云服务器，最低 2核4G，能跑 Docker）
- [ ] 安装 Docker + Docker Compose
- [ ] 修改 `compose.yml` 中 `/root/SSE_Market/` 路径为服务器实际路径
- [ ] 配置 Nginx + HTTPS（Let's Encrypt，原项目已有 certbot 续期脚本）
- [ ] 申请学院子域名 or 独立域名
- [ ] 首次部署，验证线上运行

---

## 人B — 前端 + 移动端任务清单

### 阶段一：品牌替换
- [ ] PC端：`src/views/layout/NavbarView.vue` — 导航栏名称
- [ ] PC端：`src/views/login/LoginView.vue` — 登录页标题
- [ ] PC端：`src/utils/mouseClick.js` — 点击特效文字
- [ ] 移动端：全局搜索"软工"替换为"智工"
- [ ] 替换 logo、favicon、banner 图

### 阶段二：环境配置
- [ ] PC端：修改 `.env.production` 的 `VUE_APP_BASE_URL` 为服务器后端地址
- [ ] 移动端：同上

### 阶段三：构建 & 部署
- [ ] `yarn build` 打包
- [ ] 将 dist/ 上传到服务器（Nginx 已配置静态路径）

---

## 待办：需要填入真实值的位置

| 文件 | 字段 | 说明 |
|------|------|------|
| `config/application.yml` | `datasource.password` | MySQL 密码，与 compose.yml 保持一致 |
| `config/application.yml` | `cos.bucketName/appid/secretId/secretKey` | 腾讯云 COS 信息 |
| `config/application.yml` | `tms.secretId/secretKey` | 腾讯云 TMS 密钥 |
| `config/application.yml` | `crypto.jwtKey` | 随机字符串，32位以上 |
| `.env` | 所有 SMTP_* 字段 | 邮箱账号和授权码 |
| `config/emailConst.go` | logo 图片 URL | 有自己域名后替换 `ssemarket.cn/new/android-chrome-192x192.png` |

---

## 问题记录

> 记录开发过程中遇到的坑和解决方案

---

## 进度日志

- 2026-08-04：项目启动，克隆三个原始仓库，完成技术栈调研
- 2026-08-04：确定方案（腾讯云全家桶 + CDKey注册），创建配置文件模板，完成品牌文字替换
