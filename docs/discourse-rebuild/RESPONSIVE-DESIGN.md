# 智工集市响应式界面

基线：`dbf2103`（2026-09-09 从 `origin/codex/feature-discourse-bbs` fast-forward）。开发分支：`codex/meta-mobile-navigation`。

## 移动端

`zgie_bbs_mobile_navigation` 控制 767px 及以下屏幕的四项导航：主页、热帖、私聊、个人。仅向已登录成员显示，管理后台保留原生界面。

- 主页进入站点配置的默认列表；热帖进入 `/hot`。
- 有聊天权限时私聊进入 Chat 的直接消息列表；关闭聊天或没有聊天权限时进入原生私信收件箱。
- 个人进入当前用户的资料概览；顶部使用铃铛呈现原有通知与账户菜单，保留通知、退出登录等操作。
- 底栏处理 iOS 安全区；原生编辑器打开或软键盘出现时隐藏，阅读进度浮块同步避让。
- Chat 内部的频道、直接消息与搜索放在顶部，底部始终为全站导航。
- 复用 Discourse 原生路由、编辑器、聊天和权限检查，不重新实现身份认证或发帖接口。

主要文件为 `zgie-mobile-navigation.gjs`、`zgie-responsive-layout.gjs`、`zgie-mobile.scss`。新增系统测试覆盖四项导航、聊天开关、回复编辑器与桌面隐藏。

## 本地验证

沿用 `discourse/README.md` 的开发容器。首次启动需要 `bundle install`、`pnpm install --frozen-lockfile`。系统测试还需要 `pnpm playwright-install` 与 `pnpm exec playwright install-deps chromium`。

```sh
RAILS_ENV=test LOAD_PLUGINS=1 bin/rails runner 'Stylesheet::Manager.recalculate_fs_asset_cachebuster!'
bin/rspec plugins/zgie-bbs/spec/system/use_mobile_navigation_spec.rb
```

修改 CSS 后应刷新测试样式摘要；Discourse 测试环境会复用磁盘上的样式清单，单独重跑 RSpec 可能仍引用旧 CSS。开发模式通过文件时间自动刷新，不受这个缓存问题影响。

所有开发操作均在本地进行；生产站点未发布此分支。可在管理后台关闭 `zgie_bbs_mobile_navigation`，恢复原生移动导航。
