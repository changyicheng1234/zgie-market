# Discourse 重构验证记录

验证日期：2026-08-06

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
| 上游 lint | Ruby、Prettier、Stylelint 全部通过 |
| 插件测试 | 3 examples，0 failures |

## 注册接口实测

通过真实 `/u.json` 注册接口提交了两次本地测试：

- 错误邀请码：`success=false`，返回“您输入的邀请码不正确”。
- 正确邀请码：`success=true`、`active=false`，创建待邮件激活账号并将激活邮件投递到 Mailpit。

这证明当前方案直接复用了 Discourse 的账号、密码、邮件激活与邀请码注册链路，没有自行仿写鉴权 API。

## 开发站当前状态

- 论坛：`http://localhost:4200`
- 邮件预览：`http://localhost:8025`
- 本地测试管理员：`admin@example.com`
- 管理员设置密码的邮件已进入 Mailpit，没有发送到外部。

`.env` 当前仍使用示例邀请码，只适合本机实验。共享或部署前必须替换 `ZGIE_INVITE_CODE`。

## 尚未覆盖

- 旧 MySQL 数据迁移尚未执行。
- `global_code` 是可重复使用的共享邀请码；若要求多个手输码逐个核销，需要继续扩展插件。
- 生产 SMTP、域名、TLS、备份与升级演练需要在部署主机验证。
