# Discourse 实验架构

## 运行结构

```text
Browser
   |
   v
Discourse (Ember + Rails)
   |-- PostgreSQL
   |-- Redis
   |-- Sidekiq
   `-- zgie-bbs plugin
```

开发环境使用 Discourse 官方 `discourse/discourse_dev` 镜像。上游源码被浅克隆到 `discourse/.cache/discourse`，该目录不进入本仓库；智工定制插件通过独立 bind mount 加载。

生产环境使用官方 `discourse_docker` standalone 模板，由外部 Nginx 负责 TLS 和反向代理。站点持久数据只保存在 Discourse 的 PostgreSQL、Redis 和 `/shared` 卷中。

## 为什么不修改核心源码

- 保留官方安全更新和升级路径。
- 避免维护 Rails、Ember 与序列化协议的长期 fork。
- 品牌和样式放在插件资产中。
- 一次性站点设置和分类放在幂等配置任务中。
- 如果未来确实需要独有业务，再以独立 Discourse 插件扩展。

## 注册模式

配置任务支持两种模式：

- `global_code`：开放本地注册，但必须输入全站邀请码。对应 Discourse 核心 `invite_code` 设置。
- `invite_links`：关闭普通注册入口，只接受管理员生成的邀请链接。对应 Discourse 核心 `invite_only` 设置。

两种模式不同时启用。旧系统“多个 9 位、单次使用、手工输入的 CDKey”不是 Discourse 核心能力；若最终必须逐码核销，再新增一个很小的注册插件。

## 与旧系统的关系

旧 Vue/Go/MySQL 代码不参与新站运行。若存在需要保留的数据，MySQL 仅作为一次性迁移源：

```text
User     -> Discourse User
Post     -> Topic + first Post
Pcomment -> subsequent Post
Ccomment -> subsequent Post with reply_to_post_number
Partition -> Category
Psave    -> Bookmark
Plike    -> PostAction(like)
```
