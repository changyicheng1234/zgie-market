# 智工集市项目开发记录

> 智能工程学院内部交流平台。
> 线上地址：**https://ise-market.duckdns.org/**（临时，即将切换到 `ise-market.cn`，见下方「ICP 备案」）

---

## ⚠️ ICP 备案进行中（2026-08-17 发现，待完成）

服务器在腾讯云**广州**节点（大陆），根据工信部规定，用域名访问大陆服务器上的网站必须完成 **ICP 备案**，否则腾讯云会在网络层拦截该域名的访问（表现为跳转到"未备案"提示页，或间歇性拦不稳定）。这也是之前排查"网页时好时坏"时的部分真实原因。

- **备案主体**：服务器是老师的腾讯云账号买的，备案需要账号实名认证人配合（大概率是老师本人），具体谁做主体待和老师确认
- **新域名**：已在腾讯云购买正式域名 `ise-market.cn`（DuckDNS 免费域名不被备案系统接受，只能作废/过渡用）
- **当前状态**：域名已买，**尚未提交备案申请**，也**尚未**把 `ise-market.cn` 解析到服务器（备案审核期间故意不解析，避免"未备案先上线"）
- **备案通过后要做的技术切换**（届时找 Claude 继续）：
  1. DNS：`ise-market.cn` A 记录指向 `203.195.162.102`
  2. 用腾讯云免费证书服务给新域名重新申请 HTTPS 证书
  3. 改 Nginx 配置 + `/var/discourse/containers/app.yml` 里的 `DISCOURSE_HOSTNAME` 为 `ise-market.cn`，`./launcher rebuild app` 生效
  4. 旧 DuckDNS 地址视情况保留跳转或弃用
- 备案周期一般 1~20 个工作日，期间新域名无法通过大陆网络正常访问

---

## ⚠️ 架构变更（2026-08-17）

原本基于 [软工集市](https://ssemarket.cn)（Go/Gin + Vue2）改造的自研前后端架构，**已于 2026-08-17 整体替换为 Discourse**（人B 主导的重构方案，代号"智工 BBS"）。

- 旧的 `SSE_market_server`（Go 后端）/ `SSE_market_client`（PC 前端）/ `sse_market_mobile`（移动端）三个目录中的代码**已不再部署上线**，仅作历史存档保留在本仓库中，不再需要维护或修复其中的 bug（包括本文档下方"后端待完成清单"里记录的置顶/擦亮、匿名发帖、微信支付对接等，均已随旧架构一起下线）。
- 旧栈的 MySQL 数据已于下线前全库导出备份，未丢失，见下方「服务器信息」。
- 如果之后要恢复或参考旧架构的实现细节（比如匿名发帖的脱敏逻辑、置顶擦亮的定价模型），本文档最下方仍保留完整的旧开发记录供查阅。

## 当前架构：Discourse

| 项目 | 内容 |
|------|------|
| 平台 | [Discourse](https://www.discourse.org/)（官方 `discourse_docker` 部署） |
| 品牌名 | 智工 BBS |
| 分区 | 日常吐槽 / 学习交流 / 打听求助 / 恋爱交友 / 课程专区 / 其他 |
| 注册方式 | 邀请码模式（非公开注册），当前邀请码：`44d734ae4fd49c1e20b2` |
| 邮件 | 中山大学官方 SMTP（mail.sysu.edu.cn） |
| 管理员账号 | 邮箱 `yic19250@gmail.com`，用户名 `user1`（登录名不变，昵称已改为"智工第一深情"，因为 `unicode_usernames` 关闭中文只能填昵称字段）。密码未留存明文记录，忘记密码走邮箱重置，或 SSH 上服务器用 Discourse 的 rake task 改密码 |
| 插件来源分支 | `codex/feature-discourse-bbs` |

**邀请码机制说明**：`invite_code` 是 Discourse 的**全局共享站点设置**（后台 site settings 搜 "invite code"），不是每人一个的邀请链接（Discourse 原生 Invite Links 功能这里没用，`invites` 表是空的）。所有人共用同一个码、无使用次数限制，一旦泄露到目标学生群体之外任何人都能用它注册。要改码/失效旧码直接去后台改这一项设置即可立即生效。

**登录排障经验**：
- 用邮箱登录务必输入完整邮箱（包括 `.com` 等后缀），浏览器自动填充有时会截断，之前遇到过 `yic19250@gmail` 少了 `.com` 导致登录失败，报错其实是正常的"账号不存在"，不是 bug
- 若登录点击后一直转圈没反应，先看浏览器控制台是否有 `net::ERR_CONNECTION_CLOSED`（尤其是 `/session/csrf` 请求）——这通常是大陆网络链路本身不稳定，不是 Discourse 或本项目的 bug；`/srv/pv`、`/srv/se` 返回 403 是正常的匿名统计信标失败，可以忽略

### 服务器信息

- 服务器：腾讯云 203.195.162.102
- SSH：`ssh -i ~/.ssh/zgie_server ubuntu@203.195.162.102`
- Discourse 安装目录：`/var/discourse`（官方 `discourse_docker` 标准布局）
- Discourse 容器内网端口：`127.0.0.1:20080`，Nginx 全站反代到这个端口，复用原有 HTTPS 证书（acme.sh + DuckDNS，自动续期）
- 旧栈 MySQL 全库备份：`/home/ubuntu/old-stack-backup/mysql-all-20260817.sql`
- 服务器规格：腾讯云 2核 / 3.6GB 内存 / 59GB 磁盘，**单容器全家桶**（Postgres+Redis+Sidekiq+Unicorn 全跑在同一个 `app` 容器里，没有拆分数据库/对象存储）
- 2026-08-18 资源快照（4个注册用户，几乎零负载下的基线）：内存用 2.1G/3.6G（app 容器占 1.335G，即 37%），load average 0.1~0.2（2核基本空闲），磁盘用 15G/59G（数据库本身仅 41MB，uploads 仅 1.1MB）。**结论：内存会比磁盘更早成为瓶颈**——Postgres/Redis/Rails 的基础占用不会因为内容少而低多少，并发用户一多内存先顶不住；磁盘增长速度慢，42GB 剩余空间够撑很久
- Unicorn worker 数（`UNICORN_WORKERS`，决定同时能处理几个请求）是 `discourse_docker` 在**构建镜像时**按当时内存自动算出来的（现在是 3 个），服务器升级内存后不会自动变多，必须重新 `./launcher rebuild app` 才会按新内存重算
- **公网 IP 大概率不是弹性 IP（EIP）**：查了实例元数据接口 `metadata.tencentyun.com` 的 `eip-address` 字段返回 404（未绑定 EIP 的典型表现，但没有控制台权限做不到 100% 确认）。这意味着关机类操作（比如以后配置升级）理论上公网 IP 可能会变，DNS（DuckDNS 或以后的 `ise-market.cn`）都硬编码解析到 `203.195.162.102`。**已跟人A确认过：不转弹性 IP，嫌麻烦，风险接受**——EIP 绑定在运行中的实例上其实免费（只有闲置不绑定才收费），但转不转是个可选项，不是必须做的事。真要升配的话，升级完检查一下 IP 有没有变，变了就去 DNS 那边手动改一下解析，几分钟能搞定，不是大问题
- **升配（加 CPU/内存）不会丢数据、不用重新部署**：腾讯云的"变更配置"只是给同一台机器分配更多硬件资源，系统盘/数据盘原封不动，Discourse、数据库、上传文件、Nginx、SSL 证书都不用动。因为是包年包月，走"补差价"流程，不用退款重买。会有几分钟停机（腾讯云要求关机才能变配）。当前判断：**200~500 注册用户、活跃期同时在线 20~30 人左右**是该考虑升级的信号，建议目标配置 **4核8G**（够撑学院级规模，不用短期内二次升级）

### 常用运维命令

```bash
# 登录服务器
ssh -i ~/.ssh/zgie_server ubuntu@203.195.162.102

# 进入 Discourse 容器
cd /var/discourse
./launcher enter app

# 重建容器（改配置后，比如 app.yml 变动）
cd /var/discourse
./launcher rebuild app

# 重启容器（不需要重建时）
./launcher restart app

# 查看日志
./launcher logs app
```

### 踩坑记录

- **容器构建时连 GitHub 随机卡死 / 连接被重置**：服务器国际线路问题，解决方法是给容器内 git 配置 `ghfast.top` 代理镜像。构建耗时约1小时（含两次重试）。
- **`launcher rebuild app` 不能挂在不稳定的 SSH 会话前台跑**：SSH 断线会连带杀掉触发它的 shell，导致旧容器已停、新容器没建完就中断，网站直接下线。必须用 `sudo bash -c 'cd /var/discourse && setsid nohup ./launcher rebuild app > /home/ubuntu/rebuild.log 2>&1 < /dev/null &'` 之类的方式让它彻底脱离终端，再轮询日志/`docker ps` 确认完成。
- **Discourse 插件里的 `Jobs::*` 任务类不会被自动加载**：不管放在 `<plugin>/jobs/regular/` 还是 `<plugin>/app/jobs/regular/`，都得在 `plugin.rb` 的 `after_initialize` 里手动 `require_relative` 一下（参考官方 `discourse-zendesk-plugin` 的写法），否则 `Jobs.enqueue` 在真正运行时会报 `uninitialized constant`，而且这个错误在 `on(:post_created)` 钩子里会被静默吞掉，表现为"代码部署了但通知就是发不出去"，很难第一时间发现。
- **`docker exec` 里直接跑 `bin/rails runner` 报 `NoDatabaseError`**：默认以 root 身份执行会连不上 Postgres（走 socket 的 peer 认证对不上），要加 `-u discourse`（容器内运行 Discourse 的系统用户）。

---

## 待办

- [ ] **ICP 备案**：确定备案主体人选（老师 or 授权你操作），用新域名 `ise-market.cn` 提交腾讯云备案，通过后按上方步骤切换域名
- [ ] 用户量上来后评估服务器内存升级（当前 3.6GB，建议 4GB+）
- [ ] 邀请码用完 / 需要新邀请码时，去 Discourse 后台管理面板（`/admin`）生成
- [ ] 视频/长期运营需求：评估是否需要开放公开注册，或长期保持邀请制

---

## 进度日志

- 2026-08-04：项目启动（旧架构），克隆三个原始仓库，完成技术栈调研
- 2026-08-04：确定方案（腾讯云全家桶 + CDKey注册），创建配置文件模板，完成品牌文字替换
- 2026-08-05：后端 Docker Compose 部署完成，前端构建并上传，Nginx + HTTPS 配置完成
- 2026-08-06：图片上传 URL 修复；等级制度改为学校主题；匿名发帖后端全部完成
- 2026-08-11：置顶/擦亮功能后端全部实现（微信支付未接通，商户资质未下发）；邮件模板 logo 换成中山大学校徽
- **2026-08-17：架构整体切换为 Discourse。** 旧 Go/MySQL/Redis 栈备份后下线，安装官方 `discourse_docker`，完成生产配置（域名/SMTP/邀请码/6个分区/中文品牌），修复 GitHub 连接问题后构建成功，创建管理员账号，Nginx 切换反代到 Discourse 容器，端到端验证通过（HTTPS、品牌显示、邀请码校验含防爬虫 honeypot）。旧架构任务清单（匿名发帖/置顶擦亮/微信支付等）随之作废，仅存档参考。
- **2026-08-18：企业微信群机器人新帖通知移植到 Discourse。** 在 `zgie-bbs` 插件里用 `Jobs::ZgieNotifyWecom`（Sidekiq 异步任务）监听新主题的 `post_created` 事件，推送分区+标题+链接到企业微信群机器人，webhook 地址通过 `ZGIE_WECOM_WEBHOOK_URL` 环境变量配置（`app.yml` 里维护，留空则静默跳过）。生产环境用真实发帖流程端到端验证通过。过程中踩了一个坑：插件里的 Job 类不会被自动加载，必须在 `plugin.rb` 手动 `require_relative`，否则新帖事件触发时会静默报错、消息发不出去（详见上方「踩坑记录」）。

---

<details>
<summary>旧架构开发记录存档（Go/Gin + Vue2，2026-08-04 ~ 2026-08-17，已下线）</summary>

原项目仓库：https://gitee.com/yang-peiyue/SSE_market_server

### 项目结构

```
学院集市/
├── SSE_market_server/   # 后端 Go + Gin（人A负责）
├── SSE_market_client/   # PC端前端 Vue 2（人B负责）
└── sse_market_mobile/   # 移动端 Vue 2 + Vant（人B负责）
```

### 技术栈

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

### 功能开发进度（旧架构，均已随下线作废）

- 匿名发帖：`model.Post.IsAnonymous`，发帖/列表/详情接口均已支持，脱敏逻辑（匿名时 UserName 替换为"匿名用户"等）在 `controller/postController.go`
- 置顶/擦亮：`BoostOrder` 表 + `service/boostService.go` 定价表 + 微信支付 APIv3 对接（`api/wechatpay.go`），置顶用 `IsTop`/`TopExpireAt` 排序，擦亮用热度公式临时加成，微信支付商户资质始终未下发，功能从未真实跑通过
- 等级制度改造：6个前端文件，等级名称改为学校风格，积分阈值降低
- 新帖子企业微信群通知：`service/wecomNotify.go`，webhook 一直未配置，功能从未触发过

### 问题记录

**bcrypt hash 在 SSH 命令里被 shell 转义破坏**：直接把 bcrypt hash（含 `$2b$10$...`）拼进 SSH 命令，shell 会把 `$` 符号解释为变量，导致 hash 损坏。解决方法：用 Python 的 subprocess 传参（不经过 shell），或用 heredoc + `$'...'` 语法。

</details>
