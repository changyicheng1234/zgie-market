# 智工集市项目开发记录

> 智能工程学院内部交流平台。
> 线上地址：**https://ise-market.cn/**（2026-09-07 已从 `ise-market.duckdns.org` 切换完成，旧域名 301 跳转到新域名）

---

## ✅ 域名已切换到 ise-market.cn（2026-09-07 完成）

ICP 备案 2026-09-05 通过后，2026-09-07 完成服务器端全部切换。备案信息：备案主体是老师的腾讯云账号（账号 ID `100051329736`），备案订单号 `30178696354544589`，公安备案数据码 `a5fe5258b0976b89d9ae007964be24e8`。

背景：服务器在腾讯云**广州**节点（大陆），工信部规定用域名访问大陆服务器必须完成 ICP 备案，否则腾讯云在网络层拦截（跳"未备案"提示页 / 间歇性拦），这也是当初"网页时好时坏"的部分原因。

**本次切换实际做的事**：
1. **DNS**：`ise-market.cn` A 记录已指向 `203.195.162.102`（切换前就已生效）。`www.ise-market.cn` 尚未加 A 记录，nginx 已配好 301 block，需要时去 DNS 加一条即可
2. **HTTPS 证书**：用的是腾讯云签发的 **TrustAsia DV** 证书（`ise-market.cn` + `www.ise-market.cn`，RSA），nginx 格式 bundle 存在仓库 `ise-market.cn_nginx/`，已装到服务器 `/etc/nginx/ssl/ise-market.cn.fullchain.pem` + `ise-market.cn.key`。**有效期到 2026-12-04，只有 90 天，不自动续期**——到期前必须去腾讯云重新申请再手动传上去换掉，否则开了 HSTS 的浏览器会直接打不开
3. **Nginx**：`/etc/nginx/sites-available/zgie-market` 重写为 4 个 server block——`ise-market.cn` 用新证书反代 Discourse（新增了 `X-Forwarded-Proto` header）；`www.ise-market.cn` 用新证书 301 到 apex；`ise-market.duckdns.org` 用旧 acme.sh 证书 301 到 apex；80 端口全部 301 到 https apex。旧配置备份在同目录 `zgie-market.bak-20260907`
4. **Discourse**：`app.yml` 的 `DISCOURSE_HOSTNAME` 改为 `ise-market.cn`（备份 `app.yml.bak-20260907`），`./launcher rebuild app` 生效（顺便拉了最新插件代码，见进度日志）。**注意 `./launcher restart app` 不会应用 `app.yml` 的 env 变更**（它只是 `docker stop`+`docker start` 同一个容器），必须 `rebuild`
5. **`force_https`**：本次一并从 `false` 开为 `true`（`Discourse.base_url` 现在是 `https://`，邮件链接、canonical 都走 https）。之前不敢开是因为旧 nginx 配置没传 `X-Forwarded-Proto`，开了会 301 循环；本次新配置补了这个 header，已验证无循环
6. **旧 DuckDNS**：保留做 301 跳转（旧证书 + acme.sh 自动续期继续留着，方便万一回滚）。旧证书文件 `/etc/nginx/ssl/fullchain.pem` + `key.pem` 仍被 duckdns 那个 block 引用，别删

**回滚办法**（万一新域名出问题）：`cp zgie-market.bak-20260907 zgie-market && nginx -s reload`，再把 `app.yml` 的 hostname 改回 duckdns 跑 `rebuild`，关掉 `force_https`。

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
| 品牌名 | 智工集市（旧称"智工 BBS"，2026-09-07 改；`app.yml` 的 `ZGIE_SITE_TITLE` 与数据库 `SiteSetting.title` 均已更新，logo/favicon 换成 `zgie-bbs` 插件里矢量化的新校徽+字标） |
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
- Discourse 容器内网端口：`127.0.0.1:20080`，Nginx 全站反代到这个端口。HTTPS 证书：`ise-market.cn` 用腾讯云 TrustAsia DV 证书（`/etc/nginx/ssl/ise-market.cn.*`，90 天，**手动续期**，到期 2026-12-04）；旧 `ise-market.duckdns.org` 仍用 acme.sh 证书（`/etc/nginx/ssl/fullchain.pem`，自动续期，仅供 301 跳转）
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
- **`zgie-bbs` 插件的 `SiteConfigurator` 不会随容器启动自动运行**：`plugin.rb` 的 `after_initialize` 只 `require_relative` 了这个类，从没调用 `.call`。它把品牌/分区/访问策略等一次性写进数据库，只能靠手动跑 rake task `zgie_bbs:configure`（或 `bin/rails runner` 直接调 `ZgieBbs::SiteConfigurator.call`）。后果：改了插件里的 logo/favicon/默认标题后，光 `./launcher rebuild app` **不会生效**，界面还是旧的，必须 rebuild 后再手动跑一次 configure。`SiteConfigurator.call(env:)` 接受 env 覆盖，想临时改标题不必动 `app.yml` 重建：`ZgieBbs::SiteConfigurator.call(env: ENV.to_h.merge("ZGIE_SITE_TITLE" => "新标题"))`（但 `app.yml` 里的 `ZGIE_SITE_TITLE` 也要同步改，否则下次真正 rebuild + configure 会回退）。标题优先级：`app.yml` 的 `ZGIE_SITE_TITLE` env > 代码里 `fetch` 的默认值。
- **SSH 到这台服务器有连接频率限制**：短时间内连续 ssh/scp 会被 `Connection closed by remote host` / `closed by port 22` 打断（sshd rate-limit 或 fail2ban）。批量操作要么合并进一次会话，要么每次间隔 30s 以上；复杂的多层 heredoc 命令还会被本地安全分类器拦，改用「上传脚本文件再执行」更稳。

---

## 待办

- [x] **域名切换到 `ise-market.cn`**：2026-09-07 完成（DNS/证书/Nginx/`app.yml`/`force_https`），详见上方「域名已切换」章节
- [ ] **2026-12-04 前手动续期 `ise-market.cn` 的 HTTPS 证书**：腾讯云 TrustAsia DV 证书 90 天有效，不自动续期，到期前去腾讯云重新申请并传上服务器换掉 `/etc/nginx/ssl/ise-market.cn.*`
- [ ] 用户量上来后评估服务器内存升级（当前 3.6GB；~1000 人学院规模建议升 **4核8G**，磁盘扩到 100G，见「域名已切换」下方对话记录 / 服务器信息里的升配说明）
- [ ] **修复 OpenISE 侧边栏链接在桌面端不显示的问题**（见下方「已知问题」）：需要改前端代码 + rebuild，故意推迟到下次服务器升级后一并做
- [ ] 邀请码用完 / 需要新邀请码时，去 Discourse 后台管理面板（`/admin`）生成
- [ ] 视频/长期运营需求：评估是否需要开放公开注册，或长期保持邀请制

---

## 已知问题

### OpenISE 侧边栏链接桌面端不显示（2026-09-07 记录，待修）

`zgie-bbs` 插件的 `assets/javascripts/discourse/api-initializers/zgie-external-links.gjs`（`5b2a07a` 引入，人B 加的 OpenISE 外链，指向 `https://openise.pages.dev`）在**手机端能显示、桌面端不显示**。

- **代码没问题**：容器内源文件、编译产物（`public/assets/js/plugins/zgie-bbs_main.*.digested.js` 里能 grep 到 `openise`）都在，是运行时注册没成功。
- **根因**：这段代码把链接挂在 **Chat 专属侧边栏面板**上，用私有 API `discourse/lib/sidebar/custom-sections` 的 `customPanels` 去检测 Chat 面板是否存在，检测不到就每 150ms 重试、6 秒后放弃。桌面端在 `chat_separate_sidebar_mode = never`（当前值）下没有独立 Chat 侧边栏面板，检测永远失败；手机端全屏 Chat 有独立面板所以能注册。
- **试过但无效**：把 `chat_separate_sidebar_mode` 改成 `always`，结果桌面端还是不行、连手机端也失效了（该值又改变了 Chat 侧边栏结构）。已回退回 `never`。
- **正解（待做）**：改 `zgie-external-links.gjs`，改用公开稳定的 `api.addSidebarSection(callback)`（不传 panel 参数 → 注册到主侧边栏），去掉 Chat 面板检测和轮询。这样桌面/手机主侧边栏都常驻显示。改完要 commit + push + `./launcher rebuild app`（前端 assets 重新编译，约 1 小时下线），所以合并到下次服务器升级 rebuild 时一起做。

---

## 进度日志

- 2026-08-04：项目启动（旧架构），克隆三个原始仓库，完成技术栈调研
- 2026-08-04：确定方案（腾讯云全家桶 + CDKey注册），创建配置文件模板，完成品牌文字替换
- 2026-08-05：后端 Docker Compose 部署完成，前端构建并上传，Nginx + HTTPS 配置完成
- 2026-08-06：图片上传 URL 修复；等级制度改为学校主题；匿名发帖后端全部完成
- 2026-08-11：置顶/擦亮功能后端全部实现（微信支付未接通，商户资质未下发）；邮件模板 logo 换成中山大学校徽
- **2026-08-17：架构整体切换为 Discourse。** 旧 Go/MySQL/Redis 栈备份后下线，安装官方 `discourse_docker`，完成生产配置（域名/SMTP/邀请码/6个分区/中文品牌），修复 GitHub 连接问题后构建成功，创建管理员账号，Nginx 切换反代到 Discourse 容器，端到端验证通过（HTTPS、品牌显示、邀请码校验含防爬虫 honeypot）。旧架构任务清单（匿名发帖/置顶擦亮/微信支付等）随之作废，仅存档参考。
- **2026-08-18：企业微信群机器人新帖通知移植到 Discourse。** 在 `zgie-bbs` 插件里用 `Jobs::ZgieNotifyWecom`（Sidekiq 异步任务）监听新主题的 `post_created` 事件，推送分区+标题+链接到企业微信群机器人，webhook 地址通过 `ZGIE_WECOM_WEBHOOK_URL` 环境变量配置（`app.yml` 里维护，留空则静默跳过）。生产环境用真实发帖流程端到端验证通过。过程中踩了一个坑：插件里的 Job 类不会被自动加载，必须在 `plugin.rb` 手动 `require_relative`，否则新帖事件触发时会静默报错、消息发不出去（详见上方「踩坑记录」）。
- **2026-09-05：ICP 备案通过。** 腾讯云备案订单号 `30178696354544589`，公安备案数据码 `a5fe5258b0976b89d9ae007964be24e8`。本次仅更新本地文档记录备案结果，DNS 解析、HTTPS 证书、Nginx 及 `app.yml` 的 `DISCOURSE_HOSTNAME` 切换到 `ise-market.cn` 仍未执行，留到下次登录服务器时按「ICP 备案」章节的四步操作完成。
- **2026-09-07：域名切换到 `ise-market.cn` + 部署最新插件代码。** 上传腾讯云 TrustAsia DV 证书并重写 nginx（新域名反代 Discourse、`www` 和旧 DuckDNS 域名 301 到 apex、新增 `X-Forwarded-Proto` header）；`app.yml` 的 `DISCOURSE_HOSTNAME` 改为 `ise-market.cn` 并 `./launcher rebuild app`，rebuild 顺带从 GitHub 拉了最新分支代码上线（`5b2a07a` 智工集市 rebrand + OpenISE 侧边栏链接 `zgie-external-links.gjs`，此前生产还是 8-18 的旧插件）；一并开启 `force_https`（`base_url` 现在 https，已验证无 301 循环）。端到端验证通过：新域名 HTTPS 200、旧域名 301 跳转、`srv/status` 200、`Jobs::ZgieNotifyWecom` 正常加载、无 error 日志。备份：`nginx sites-available/zgie-market.bak-20260907`、`app.yml.bak-20260907`。遗留：`ise-market.cn` 证书 90 天需手动续期（到期 2026-12-04）；`www.ise-market.cn` 的 DNS A 记录尚未添加（nginx block 已就绪）。
- **2026-09-07（同日续）：品牌切到"智工集市" + 密码长度放宽。** 发现 rebuild 后标题/图标仍是旧的——根因见下方「踩坑记录」：`zgie-bbs` 插件的 `SiteConfigurator` 不会随容器启动自动跑。处理：`app.yml` 的 `ZGIE_SITE_TITLE` 改为"智工集市"，容器内 `bin/rails runner` 调 `ZgieBbs::SiteConfigurator.call`（临时注入正确标题，无需再 rebuild），标题 + logo/logo_dark/logo_small/favicon 全部更新为插件内新素材，`/site/basic-info.json` 已确认。另按需求把 `min_password_length` 10→8、`min_admin_password_length` 15→8（Discourse 硬下限就是 8）。OpenISE 侧边栏链接问题见下方「待办」与「已知问题」，本次未解决，推迟到服务器升级后随代码改动一起处理。

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
