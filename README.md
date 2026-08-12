# 智工 BBS（Discourse 重构分支）

智工 BBS 是一个面向智能工程学院师生的邀请制、半匿名校园论坛。本分支使用完整 Discourse 作为论坛前端、后端和数据模型，只在独立插件中维护智工 BBS 的产品定制，不长期 fork 或复制上游核心源码。

当前分支：`codex/feature-discourse-bbs`。这是可以本地运行和测试的功能分支，尚未替换 `master` 中的旧 Vue / Go 实现。

## 项目组成

```text
zgie-market/
├── discourse/
│   ├── bin/dev                    # 本地开发统一命令
│   ├── compose.dev.yml            # Discourse 开发容器与持久卷
│   ├── deploy/app.yml.example     # 官方 discourse_docker 生产模板
│   └── plugins/zgie-bbs/          # 智工 BBS 的前后端定制、配置、任务和测试
├── docs/discourse-rebuild/
│   ├── PRODUCT.md                 # 产品定位、核心流程与非目标
│   ├── ARCHITECTURE.md            # 架构、匿名机制和旧数据映射
│   └── VALIDATION.md              # 已完成的接口、浏览器与插件验证
├── SSE_market_client/             # 旧 Vue 桌面前端，仅作对照/迁移来源
├── sse_market_mobile/             # 旧 Vue 移动前端，仅作对照/迁移来源
└── SSE_market_server/             # 旧 Go + MySQL 后端，仅作对照/迁移来源
```

启动脚本会把 `.env` 中指定分支的上游 Discourse 浅克隆到 `discourse/.cache/discourse`，并把 `zgie-bbs` 插件挂载进去。`.cache`、本地 `.env` 和运行数据都不会提交到 Git；本分支实际验证过的上游提交记录在 [`VALIDATION.md`](docs/discourse-rebuild/VALIDATION.md) 中。

## 当前已实现

- 完整 Discourse 论坛能力：主题、分类、搜索、最新/热门、富文本、图片、Emoji、收藏、通知、举报和管理后台。
- 本地账号鉴权：用户使用邮箱、用户名、密码和邀请码注册；内容默认仅登录成员可见。支持共享邀请码或 Discourse 邀请链接两种模式。
- 六个初始分类：日常吐槽、学习交流、打听求助、恋爱交友、课程专区和其他。
- 固定两层留言：回复一级留言时直接显示在其下；回复二级留言时显示 `用户 A > 用户 B`，但仍留在同一一级留言下，不追加到主题末尾。一级留言始终展开，只有二级留言可折叠。
- 单次匿名发布：发主题和回复时分别选择是否匿名，不在用户页提供全局匿名模式。匿名头像和用户名不能打开资料页，发布者仍可管理自己的匿名内容，管理员保留审计能力。
- 匿名草稿清理：匿名主题发布成功后同步清理真实账号草稿，并阻止旧的在途自动保存把它重新写回草稿箱。
- 互动隐私：点赞和反应只显示计数与 Emoji；普通成员不能展开或通过名单接口枚举参与者，管理员仍可审计。
- 稳定阅读进度：右侧显示当前内容在主题全部帖子中的阅读位置，位置按嵌套页面的实际顺序递增，不受二级留言原始楼层号影响。
- 内容限制：主题与回复只禁止空内容；正文上限提高到 Discourse 当前支持的 150,000 字符，标题保留数据库的 255 字符技术上限。
- 幂等演示数据：一个命令创建演示账号、10 个主题，以及包含富文本、图片和 10 条回复的综合测试主题。

当前尚未执行旧 MySQL 数据迁移，也未完成生产域名、SMTP、TLS、备份和升级演练。多枚手输邀请码逐码核销与 Tag 定制尚未实现。打分、二手交易、简历招聘、积分排行和统一身份认证不是本产品目标。

## 数据库与数据来源

新站运行时以 Discourse 原生 PostgreSQL 为唯一业务数据源：账号、主题、帖子、回复关系、分类、收藏、点赞、通知和站点设置均使用 Discourse 自己的表与 Rails 模型。插件不新增一套用户、主题或评论表。

| 组件 | 用途 | 是否为业务事实来源 |
| --- | --- | --- |
| PostgreSQL | 用户、内容、权限、互动、通知和站点配置 | 是 |
| Redis | 缓存、MessageBus 与后台任务协调 | 否，不存放论坛内容的权威副本 |
| Discourse uploads / `/shared` | 图片、附件、备份等文件数据 | 文件来源 |
| Mailpit | 只在本地接收激活和测试邮件 | 否 |
| 旧 MySQL | 未来一次性迁移的读取来源 | 不参与新站运行 |

如果迁移旧站数据，预期映射如下：

```text
User      -> Discourse User
Post      -> Topic + first Post
Pcomment  -> subsequent Post
Ccomment  -> subsequent Post with reply_to_post_number
Partition -> Category
Psave     -> Bookmark
Plike     -> PostAction(like)
```

本地演示内容不是从旧 MySQL 导入，而是由 `bin/rake zgie_bbs:seed_demo` 幂等生成；重复执行不会重复创建受控主题和回复。

## 本地启动

需要 Git、Docker（或兼容 Docker CLI 的 Podman）和 Docker Compose。首次运行会下载较大的 Discourse 开发镜像与上游源码。

```bash
git clone --branch codex/feature-discourse-bbs \
  https://github.com/changyicheng1234/zgie-market.git
cd zgie-market/discourse

cp .env.example .env
# 打开 .env，至少把 ZGIE_INVITE_CODE 改成自己的开发邀请码

./bin/dev bootstrap
./bin/dev seed
./bin/dev server
```

`server` 会在前台运行。启动完成后访问：

- 论坛：<http://localhost:4200>
- 本地邮件：<http://localhost:8025>
- 演示账号：`demo@example.com`
- 演示密码：`VioletRiver!8246-Campus`

后续启动不需要重复安装依赖：

```bash
cd discourse
./bin/dev up
./bin/dev server
```

常用维护命令：

```bash
./bin/dev status       # 查看容器状态
./bin/dev configure    # 幂等应用站点设置与分类
./bin/dev seed         # 创建或刷新本地演示数据
./bin/dev test         # 运行 zgie-bbs 插件测试
./bin/dev shell        # 进入 Discourse 开发容器
./bin/dev logs         # 查看容器日志
./bin/dev stop         # 停止容器并保留依赖与数据
```

完整开发说明见 [`discourse/README.md`](discourse/README.md)。生产环境不使用开发 Compose，而使用官方 `discourse_docker`；模板和步骤见 [`discourse/deploy/app.yml.example`](discourse/deploy/app.yml.example)。

## 进一步阅读

- [产品目标与范围](docs/discourse-rebuild/PRODUCT.md)
- [Discourse 重构架构](docs/discourse-rebuild/ARCHITECTURE.md)
- [测试与实机验证记录](docs/discourse-rebuild/VALIDATION.md)
