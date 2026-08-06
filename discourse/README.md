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

本实验在 rootless Podman 上也可运行：Compose 会创建仅属于本项目的 PostgreSQL、Redis 和 Node.js 卷，并在容器内创建与开发进程匹配的本地 PostgreSQL 角色。该兼容设置只用于开发环境，不进入生产模板。

常用命令：

```bash
./bin/dev status
./bin/dev configure
./bin/dev test
./bin/dev shell
./bin/dev logs
./bin/dev stop
```

`stop` 会保留已安装在开发容器内的 Ruby gems；再次执行 `up` 即可继续。`down` 会删除容器（但保留数据库和 Node.js 卷），下次需要重新执行 `bootstrap` 补齐 gems。

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
