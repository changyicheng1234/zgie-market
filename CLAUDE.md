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
- [x] **填入真实密钥**：`application.yml` 配置完成（COS/TMS/JWT/数据库）
- [x] **填入邮箱信息**：`.env` 配置完成（SYSU SMTP）

> ⚠️ 注意：`config/emailConst.go` 邮件模板中 logo 图片 URL 仍指向 `ssemarket.cn`，
> 等有自己的域名和 logo 后替换（搜索 `ssemarket.cn/new/android-chrome-192x192.png`）

### 阶段二：第三方服务开通（照搬软工集市方案）
- [x] 腾讯云子账号创建（zngc-cyc，AppId: 1463399285）
- [x] COS 存储桶：zngclt-1463399285（广州区，私有读写）
- [x] TMS 文本内容安全已开通
- [x] SMTP：中山大学官方邮箱 mail.sysu.edu.cn:465

### 阶段三：本地开发环境验证
- [ ] 安装 Go 1.23+
- [ ] 安装 MySQL 8.1，创建数据库（名称与 application.yml 的 `database` 字段一致）
- [ ] 安装 Redis
- [ ] 本地跑通 `go run main.go`
- [ ] 验证注册/登录/发帖流程

### 阶段四：服务器部署
- [x] 获取服务器（腾讯云 203.195.162.102）
- [x] 安装 Docker + Docker Compose
- [x] 修改 `compose.yml` 路径为 `/home/ubuntu/zgie-market/`
- [x] 安装 Nginx，配置静态文件服务 + `/api/` 反向代理
- [x] 首次部署，验证线上运行（PC: /pc/，移动: /mb/，API: /api/）
- [x] 域名：ise-market.duckdns.org（免费，DuckDNS）
- [x] 配置 HTTPS（acme.sh + DuckDNS DNS验证，证书已安装，自动续期每天3点）

---

## 人B — 前端 + 移动端

### 任务清单

#### 品牌替换（主要工作）
- [ ] PC端：`src/views/layout/NavbarView.vue` — 导航栏名称
- [ ] PC端：`src/views/login/LoginView.vue` — 登录页标题
- [ ] PC端：`src/utils/mouseClick.js` — 点击特效文字
- [ ] 移动端：全局搜索"软工"替换为"智工"
- [ ] 替换 logo、favicon、banner 图（替换 `public/` 目录下的图片文件）

#### 已完成（人A代劳）
- [x] `.env.production` 配置好了（API 地址指向服务器）
- [x] `.env.development` 配置好了（本地开发 API 指向线上服务器）
- [x] 首次 build + 上传服务器完成

---

### 人B 本地开发教程

> 目标：在自己电脑上改代码，立刻能看到效果，不依赖人A。

#### 第一步：拉代码

```bash
git clone https://github.com/[仓库地址]
```

#### 第二步：安装依赖

PC端：
```bash
cd SSE_market_client
npm install --legacy-peer-deps
```

移动端：
```bash
cd sse_market_mobile
yarn install --ignore-engines
```

> Node.js 版本要求：推荐 v18 或 v20（v24 有兼容性警告但能用）

#### 第三步：本地启动

PC端（访问 http://localhost:8080）：
```bash
cd SSE_market_client
npm run serve
```

移动端（访问 http://localhost:8081）：
```bash
cd sse_market_mobile
yarn serve
```

浏览器会自动打开。API 请求会自动转发到线上服务器 `https://ise-market.duckdns.org`，**登录/注册/发帖都能正常用**，不需要本地搭后端。

#### 第四步：改完提交

```bash
git add .
git commit -m "修改导航栏品牌名为智工集市"
git push
```

push 后告知人A，人A 执行一条命令重新 build 并部署到服务器。

---

### 人A 收到人B改动后的部署命令

```bash
# PC端
cd "D:/zstudy/研一/学院集市/SSE_market_client"
git pull && npx vue-cli-service build
scp -i ~/.ssh/zgie_server -r dist/. ubuntu@203.195.162.102:/home/ubuntu/zgie-market/frontend/pc/

# 移动端
cd "D:/zstudy/研一/学院集市/sse_market_mobile"
git pull && yarn build --ignore-engines
scp -i ~/.ssh/zgie_server -r dist/. ubuntu@203.195.162.102:/home/ubuntu/zgie-market/frontend/mb/
```

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
- 2026-08-05：后端 Docker Compose 部署完成，前端构建并上传，Nginx + HTTPS 配置完成。正式访问地址：https://ise-market.duckdns.org/pc/（PC端）、/mb/（移动端）
