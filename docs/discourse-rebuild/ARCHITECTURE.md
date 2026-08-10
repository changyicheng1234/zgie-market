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
- 固定两层评论直接使用 Discourse 原生 Nested Replies，通过站点设置限制视觉深度。
- 如果未来确实需要独有业务，再以独立 Discourse 插件扩展。

## 注册模式

配置任务支持两种模式：

- `global_code`：开放本地注册，但必须输入全站邀请码。对应 Discourse 核心 `invite_code` 设置。
- `invite_links`：关闭普通注册入口，只接受管理员生成的邀请链接。对应 Discourse 核心 `invite_only` 设置。

两种模式不同时启用。旧系统“多个 9 位、单次使用、手工输入的 CDKey”不是 Discourse 核心能力；若最终必须逐码核销，再新增一个很小的注册插件。

## 匿名与互动隐私

匿名发布直接复用 Discourse 的 `AnonymousShadowCreator`：真实账号负责登录和资格校验，进入匿名模式后由关联的影子账号创建 Topic/Post。普通成员不能访问影子账号资料，页面也不会为匿名头像和用户名生成资料入口；管理员仍可追溯真实账号，满足社区治理需要。插件不新增匿名帖子表，也不把真实身份写入前端序列化结果。

点赞和反应仍使用 Discourse 的 PostAction 与 discourse-reactions 数据模型。插件只改变公开读取边界：页面保留 Emoji 和汇总计数，普通成员不能调用参与者名单接口，管理员保留审计读取权限。

右侧楼层进度是纯前端只读视图，数字来源于 Discourse 已有的 `post_number`、`posts_count` 和 `highest_post_number`，不创建另一套楼层编号或滚动状态。

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

`reply_to_post_number` 始终保留真实回复目标。站点将嵌套回复最大视觉深度设为 1，并开启深度封顶；因此回复二级留言时，数据、通知和“回复给谁”仍指向该二级留言，但主题页面把它压平为同一一级留言下的兄弟项。插件不新增评论表，也不接管 Discourse 的帖子 API。
