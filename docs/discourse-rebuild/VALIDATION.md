# Discourse 重构验证记录

验证日期：2026-08-10

## 固定版本

- 实验分支：`experiment/discourse-rebuild`
- 独立 worktree：`/home/matrix/zgie-market-discourse`
- Discourse：`cb5b167333b1f8af31b8d5dd31cd1962ba18bf92`
- 开发镜像：`docker.io/discourse/discourse_dev:20260803-0122`
- 镜像 ID：`09fd84d99572886ec9d2e1239b6ce0c62cf78016cc9c77a9b3a6d852e7b3e88f`
- 运行时：rootless Podman 5.8.4，Docker Compose 兼容层

## 验证结果

| 项目 | 结果 |
| --- | --- |
| PostgreSQL | `accepting connections` |
| Redis | `PONG` |
| Mailpit | HTTP 200，成功捕获本地邮件 |
| Rails + Ember | 宿主 `http://127.0.0.1:4200/` 返回 HTTP 200 |
| 品牌与语言 | `智工 BBS`、`zh_CN` |
| 内容权限 | `login_required=true` |
| 本地注册 | `enable_local_logins=true`、`allow_new_registrations=true` |
| 邀请码 | `require_invite_code=true`、`invite_only=false` |
| 初始分类 | `daily`、`study`、`help`、`social`、`courses`、`other`，共 6 个 |
| 幂等配置 | 第二次运行配置任务后仍为 6 个分类 |
| 演示账号 | `demo@example.com`，已激活，密码登录有效 |
| 演示内容 | 10 个主题；综合主题含 10 条回复 |
| 回复关系 | 5 条回复使用原生 `reply_to_post_number`；其中一条回复二级留言 |
| 固定两层留言 | 原生 Nested Replies，全局启用、最大视觉深度 1、深度封顶 |
| 折叠边界 | 一级留言始终展开且无折叠入口；二级留言可独立收起和展开 |
| 层级布局 | 一级正文与二级回复区使用独立边界；桌面和移动端均无重叠或横向溢出 |
| 回复对象 | 二级留言头部显示真实的 `作者 > 被回复者` 用户名关系 |
| 匿名发布 | 所有已登录成员可切换到原生匿名影子账号，并创建主题与留言 |
| 匿名身份入口 | 匿名头像、用户名不含资料链接；普通成员不能访问匿名资料 |
| 点赞与反应隐私 | 保留计数与 Emoji；不弹出用户列表，普通成员名单 API 返回 403，管理员可审计 |
| 右侧楼层进度 | 桌面显示日期、轨道及真实 `当前楼层 / 最高楼层`；滚动后同步更新 |
| 富文本 cooked HTML | 标题、粗体、斜体、引用、代码块、表格、链接与图片均已生成 |
| 图片资源 | 本地头像、Emoji、正文图片均返回 HTTP 200，烘焙 URL 使用 `localhost:4200` |
| 上游 lint | Ruby、Prettier、Stylelint 全部通过 |
| 插件测试 | 17 examples，0 failures（含真实浏览器系统测试） |

## 注册接口实测

通过真实 `/u.json` 注册接口提交了两次本地测试：

- 错误邀请码：`success=false`，返回“您输入的邀请码不正确”。
- 正确邀请码：`success=true`、`active=false`，创建待邮件激活账号并将激活邮件投递到 Mailpit。

这证明当前方案直接复用了 Discourse 的账号、密码、邮件激活与邀请码注册链路，没有自行仿写鉴权 API。

## 演示数据接口实测

执行 `./bin/dev seed` 首次创建 10 个主题和 10 条回复；第二次执行显示 `0 topics and 0 replies created`，证明任务不会重复导入内容。

通过真实 `/session.json` 使用 `demo@example.com` 和本地演示密码登录后：

- `/session/current.json` 返回当前用户 `zgiedemo`。
- `/latest.json` 返回全部 10 个 `[演示]` 主题。
- `/t/topic/13.json` 返回 `[演示] 富文本与回复层级综合测试`、`posts_count=11` 和 `is_nested_view=true`。
- 回复关系为楼层 `4 → 2`、`6 → 3`、`8 → 5`、`10 → 7`、`11 → 4`。
- `/n/topic/13.json?sort=old` 的根留言仅为楼层 `2`、`3`、`5`、`7`、`9`，二级回复不会进入根留言流。
- `/n/topic/13/children/2.json?depth=1&sort=old` 将楼层 `4` 与 `11` 作为同级项返回；楼层 `11` 仍保留 `reply_to_post_number=4` 和回复目标 `zgiepeer`。
- 首帖 cooked HTML 中存在 `h1`、`strong`、`em`、`blockquote`、`pre`、`table`、`a` 和 `img` 元素。
- 首帖图片与 Emoji 已重新烘焙为 `//localhost:4200/...`；`zgiedemo`、`zgiepeer` 使用站内 `/letter_avatar/...`，不依赖外部头像 CDN。

浏览器系统测试实际执行了以下流程：登录用户进入嵌套主题，点击二级留言的回复按钮并提交回复，页面最终在视觉深度 1 找到新内容，并确认视觉深度 0 不存在该内容；随后确认一级留言没有折叠控件、二级留言显示真实回复对象，并实际完成二级留言的收起与重新展开。测试还创建匿名留言，确认头像和用户名均无资料入口；创建反应后确认 Emoji 与计数可见、名单不可展开；滚动到指定楼层后确认右侧进度从首帖更新为对应的真实 `post_number`。

开发站浏览器实测结果：主题初始显示 `1 / 11`，滚动到楼层 6 后显示 `6 / 11`；进度轨道位于正文右侧；反应区域渲染两个 Emoji 摘要，`pointer-events=none`，强制触发点击也没有出现参与者弹窗。

## 开发站当前状态

- 论坛：`http://localhost:4200`
- 邮件预览：`http://localhost:8025`
- 本地测试管理员：`admin@example.com`
- 本地演示账号：`demo@example.com` / `VioletRiver!8246-Campus`
- 富文本综合主题：`http://localhost:4200/t/topic/13`
- 管理员设置密码的邮件已进入 Mailpit，没有发送到外部。

`.env` 当前仍使用示例邀请码，只适合本机实验。共享或部署前必须替换 `ZGIE_INVITE_CODE`。

## 尚未覆盖

- 旧 MySQL 数据迁移尚未执行。
- `global_code` 是可重复使用的共享邀请码；若要求多个手输码逐个核销，需要继续扩展插件。
- 生产 SMTP、域名、TLS、备份与升级演练需要在部署主机验证。
