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

移动端首阶段提交：`ea3db83`。后续跨页验证补齐了楼中楼浮动回复按钮和阅读进度对底栏的避让。

## Meta 风格视觉

`zgie_bbs_branded_layout` 独立控制紫色渐变、圆角内容面板、搜索横幅、分类入口与共用导航样式。保留智工集市品牌与原生页面结构；参考 Meta 的视觉布局，不复制其私有主题或品牌素材。

- 最新与分类页有完整横幅；其他发现页显示紧凑标题与搜索。
- 分类快捷入口从当前成员可见的分类中选取，优先日常吐槽、学习交流和打听求助，不硬编码分类 ID。
- 搜索直接进入 Discourse 原生搜索，列表排序、分类筛选、发帖、用户设置、私信与聊天仍由原生组件处理。
- 手机缩短横幅，平板调整卡片方向；桌面保持左侧导航、居中的内容面板和帖子阅读轨道。
- OpenISE 外链改用公开的主侧栏 API，移除依赖 Chat 面板出现的轮询。
- 深浅色共用 CSS 变量；新增 `logo-light.svg` 为浅色字标，原 `logo.svg` 继续用于深色模式。

关闭 `zgie_bbs_branded_layout` 可恢复原有视觉，底栏开关可以独立保留。管理后台及嵌入模式不启用这层视觉。

## 发布注意

此分支未部署到生产站点。部署插件后需要运行 `bin/rake zgie_bbs:configure` 才会上传并设置新的浅色字标；执行前确认生产容器已有正确的 `ZGIE_*` 环境变量，因为该任务也会重新应用仓库中定义的站点配置。已有生产配置的维护者也可只在外观设置中上传 `logo-light.svg` 为浅色 Logo，保留深色 Logo。

首次预览建议先使用默认 Foundation 主题。其他重度定制主题可能覆盖同一组件的样式，需要单独检查。

## 本地验证

沿用 `discourse/README.md` 的开发容器。首次启动需要 `bundle install`、`pnpm install --frozen-lockfile`。系统测试还需要 `pnpm playwright-install` 与 `pnpm exec playwright install-deps chromium`。

```sh
RAILS_ENV=test LOAD_PLUGINS=1 bin/rails runner 'Stylesheet::Manager.recalculate_fs_asset_cachebuster!'
bin/rspec plugins/zgie-bbs/spec/system/use_mobile_navigation_spec.rb
```

修改 CSS 后应刷新测试样式摘要；Discourse 测试环境会复用磁盘上的样式清单，单独重跑 RSpec 可能仍引用旧 CSS。开发模式通过文件时间自动刷新，不受这个缓存问题影响。

所有开发操作均在本地进行；生产站点未发布此分支。可在管理后台关闭 `zgie_bbs_mobile_navigation`，恢复原生移动导航。

## 验证记录（2026-09-10）

- 最终全插件回归：29 examples，0 failures（seed 18968），包含原有匿名发布、楼中楼、互动隐私、阅读进度，以及新增移动导航与横幅测试。
- 补充浮动回复避让后，导航与横幅系统测试：7 examples，0 failures。
- 浅色字标配置测试：5 examples，0 failures。
- 浏览器检查：320、360、390、768、1440px；浅色与深色；首页、分类、帖子、搜索、个人设置和 Chat 直接消息列表。
- 实际横幅搜索“富文本”返回本地演示帖子；手机帖子页浮动回复与阅读进度均位于底栏上方。
- ESLint、Prettier、Stylelint 与修改涉及的 Ruby 规范检查通过。

尚需真实设备补验：iOS Safari/PWA 的软键盘与安全区、安卓输入法弹出、真实推送到达后的未读提示。桌面浏览器尺寸模拟和系统测试不能完全替代真机行为。
