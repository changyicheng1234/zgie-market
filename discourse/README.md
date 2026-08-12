# 智工 BBS Discourse 实验

这里运行完整 Discourse，并把智工定制限制在 `plugins/zgie-bbs`。上游源码只下载到 `.cache/discourse`，不会被复制进本仓库。

## 本地启动

要求：Git、兼容 Docker CLI 的容器运行时，以及 Docker Compose。当前开发机使用 rootless Podman 的 Docker 兼容层。

```bash
cd /home/matrix/zgie-market-discourse/discourse
cp .env.example .env
# 修改 .env 中的 ZGIE_INVITE_CODE
./bin/dev bootstrap
./bin/dev server
```

浏览器访问 `http://localhost:4200`。邮件预览访问 `http://localhost:8025`。

`ZGIE_EXTERNAL_HOSTNAME` 与 `ZGIE_EXTERNAL_PORT` 必须和浏览器实际访问的地址一致，Discourse 才能为正文图片、Emoji 和邮件生成正确 URL。默认值分别为 `localhost` 与 `4200`。

本实验在 rootless Podman 上也可运行：Compose 会创建仅属于本项目的 PostgreSQL、Redis 和 Node.js 卷，并在容器内创建与开发进程匹配的本地 PostgreSQL 角色。该兼容设置只用于开发环境，不进入生产模板。

常用命令：

```bash
./bin/dev status
./bin/dev configure
./bin/dev seed
./bin/dev test
./bin/dev shell
./bin/dev logs
./bin/dev stop
```

`stop` 会保留已安装在开发容器内的 Ruby gems；再次执行 `up` 即可继续。`down` 会删除容器（但保留数据库和 Node.js 卷），下次需要重新执行 `bootstrap` 补齐 gems。

## 演示账号与渲染数据

开发站可用一个幂等任务创建演示账号与内容：

```bash
./bin/dev seed
```

默认登录信息：

- 邮箱：`demo@example.com`
- 密码：`VioletRiver!8246-Campus`

任务会创建 10 个分布在现有分类中的演示主题。其中 `[演示] 富文本与回复层级综合测试` 包含标题、粗体、斜体、删除线、引用、列表、代码块、表格、链接、图片与 Emoji，并带有 10 条回复。5 条回复通过 Discourse 原生 `reply_to_post_number` 指向其他回复，辅助演示用户 `zgiepeer` 用于呈现两位作者之间的回复关系。

站点全局使用 Discourse 原生嵌套回复，固定 `nested_replies_max_depth=1` 且开启深度封顶。效果是一级留言下面只有一层楼中楼：回复楼中楼时仍保留真实回复对象和通知关系，但视觉上作为同一一级留言下的同级项出现，不会追加到主题末尾。一级留言始终展开，只有二级留言可收起；二级留言头部以 `作者 > 被回复者` 显示真实用户名。默认按发布时间从旧到新排序。

演示链路中，第 1 条一级留言由 `zgiedemo` 发布，第 3 条由 `zgiepeer` 回复，第 10 条再由 `zgiedemo` 回复第 3 条。运行中的演示数据库会按自定义字段安全刷新这些受控夹具；其他主题和用户内容不会被修改。再次运行不会重复创建主题或回复。

这些凭据只用于本机开发；任务默认拒绝在生产环境运行。如确需在临时生产模式环境中生成，必须显式设置 `ZGIE_ALLOW_DEMO_SEED=true`。

## 匿名发布与互动隐私

所有已登录成员都可以在新主题或回复编辑器中打开“匿名发帖 / 匿名回复”。该选择只影响本次提交，不改变登录会话，也不会延续到下一条内容；头像菜单不再提供整段会话切换匿名身份的入口。匿名内容仍复用 Discourse 与真实账号关联的影子账号：普通成员看不到真实身份，匿名头像和用户名也不会打开资料页；站务人员仍可按需审计。发布者用真实账号登录时仍可管理自己的匿名内容。

主题标题、首帖和回复的产品最短字数门槛均已取消，只禁止空内容；正文上限提高到当前 Discourse 支持的 150,000 字符，标题仍受数据库 255 字符的技术上限约束。

点赞与反应区域只显示计数和 Emoji，不再打开参与者名单。普通成员对应的点赞及反应名单接口也会拒绝访问，避免只隐藏前端却仍能从接口枚举用户；管理员接口保留审计能力。

嵌套主题右侧显示只读的楼层进度，使用真实 `post_number / highest_post_number`，滚动时同步更新当前楼层。桌面空间足够时显示日期与竖向轨道，窄屏或右侧空间不足时缩为右下角计数浮块。

本机首次拉取大镜像若因网络出现 `unexpected EOF`，可启用 Podman 原生分层重试后再次执行 `bootstrap`：

```bash
podman pull --retry 20 --retry-delay 2s docker.io/discourse/discourse_dev:20260803-0122
```

第一次创建管理员：

```bash
./bin/dev shell
bin/rake admin:create
```

## 邀请注册

默认 `ZGIE_REGISTRATION_MODE=global_code`，Discourse 使用本地账号注册，并在原生注册表单中要求填写 `ZGIE_INVITE_CODE`。

如果改为：

```dotenv
ZGIE_REGISTRATION_MODE=invite_links
ZGIE_INVITE_CODE=
```

普通注册入口会关闭，只能通过 Discourse 管理员生成的邀请链接注册。

注意：`global_code` 是可重复使用的共享口令；`invite_links` 可以限制使用次数。如果最终要求“多个手输邀请码、每码一次”，需要在本插件中增加邀请码表和注册校验。

## 生产部署

生产环境不使用 `compose.dev.yml`，而使用官方 `discourse_docker`。参考 [`deploy/app.yml.example`](deploy/app.yml.example)：

1. 在服务器安装官方 `/var/discourse`。
2. 将示例复制为 `/var/discourse/containers/app.yml`。
3. 填写域名、管理员邮箱、SMTP 密码和邀请码。
4. 执行 `/var/discourse/launcher rebuild app`。
5. 首次启动后进入容器执行 `bin/rake zgie_bbs:configure`。
6. 将现有 Nginx 反向代理到 `127.0.0.1:20080`。

`app.yml` 含真实秘密，不能提交。

## 已验证结果

当前固定的上游提交、镜像、测试、HTTP 和注册验证结果见 [`../docs/discourse-rebuild/VALIDATION.md`](../docs/discourse-rebuild/VALIDATION.md)。
