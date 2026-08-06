# 智工 BBS：Discourse 重构实验

本分支验证一条与旧 Vue/Go 实现完全隔离的新路线：论坛前端、后端、权限、通知、搜索和管理后台全部使用完整 Discourse，仓库只维护智工 BBS 的站点配置与少量定制。

实验入口位于 [`discourse/`](discourse/README.md)。旧的 `SSE_market_client`、`sse_market_mobile` 和 `SSE_market_server` 暂时保留，只用于对照与后续数据迁移；它们不参与 Discourse 实验运行。实验已经通过运行验证，但在确认旧 MySQL 数据是否迁移前，不贸然删除旧运行时。

## 当前假设

- 社区面向智能工程学院师生。
- 内容默认仅登录用户可见。
- 使用 Discourse 本地账号（邮箱、用户名、密码）。
- 注册必须提供邀请码。
- 不建设打分、二手、简历和 OAuth2 应用平台。
- 接受 Discourse 原生的扁平话题回复流；“回复某人”通过引用关系表达，不保留旧后端的两层楼中楼数据模型。

## 实验成功标准

- [x] 能用官方 Discourse 开发镜像启动完整站点。
- [x] 默认语言为简体中文，站点品牌为“智工 BBS”。
- [x] 未登录用户无法阅读内容。
- [x] 用户可通过账号、密码和邀请码注册。
- [x] 首页具备 Latest、Hot、分类和搜索。
- [x] 发帖、回复、点赞、收藏、通知和举报由完整 Discourse 提供。
- [x] 定制代码作为插件独立存在，不修改上游 Discourse 核心源码。

实际运行结果见 [`docs/discourse-rebuild/VALIDATION.md`](docs/discourse-rebuild/VALIDATION.md)。
