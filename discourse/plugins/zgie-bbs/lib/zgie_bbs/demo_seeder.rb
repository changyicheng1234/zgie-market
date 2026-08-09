# frozen_string_literal: true

module ZgieBbs
  class DemoSeeder
    Result =
      Struct.new(
        :user,
        :peer_user,
        :topics,
        :showcase_topic,
        keyword_init: true
      )

    TOPIC_KEY_FIELD = "zgie_bbs_demo_topic_key"
    REPLY_KEY_FIELD = "zgie_bbs_demo_reply_key"
    DEFAULT_EMAIL = "demo@example.com"
    DEFAULT_USERNAME = "zgiedemo"
    DEFAULT_PASSWORD = "VioletRiver!8246-Campus"
    DEFAULT_PEER_EMAIL = "demo-peer@example.com"
    DEFAULT_PEER_USERNAME = "zgiepeer"
    DEFAULT_PEER_PASSWORD = "CedarBridge!3902-Peer"

    TOPICS = [
      {
        key: "rich-text-showcase",
        category_slug: "study",
        title: "[演示] 富文本与回复层级综合测试",
        raw: <<~'MARKDOWN'
          # 智工 BBS 富文本渲染测试

          这是一篇用于检查 **粗体**、*斜体*、~~删除线~~、`行内代码` 和 :rocket: Emoji 的综合帖子。

          > 好的校园社区不仅要能即时交流，也要让有价值的信息长期沉淀。

          ## 功能清单

          - Markdown 图文排版
          - 主题与回复关系
          - 代码、表格和引用

          1. 从最新或分类页进入主题
          2. 阅读内容并回复
          3. 点赞、收藏或举报

          ```ruby
          community = "智工 BBS"
          puts "欢迎来到 #{community}"
          ```

          | 渲染项 | 预期结果 |
          | --- | --- |
          | 标题与段落 | 层级清晰 |
          | 代码块 | 保留格式并高亮 |
          | 表格 | 正常显示两列 |

          [访问 Discourse 官方项目](https://github.com/discourse/discourse)

          ![Discourse 示例图标](/images/discourse-logo-sketch.png)

          ---

          欢迎在下面的十条回复中继续测试引用关系与富文本内容。
        MARKDOWN
      },
      {
        key: "canteen-recommendations",
        category_slug: "daily",
        title: "[演示] 今天食堂哪道菜最值得推荐？",
        raw: "大家可以分享窗口、菜名和用餐时段，给第一次来食堂的同学一份参考。"
      },
      {
        key: "calculus-review",
        category_slug: "study",
        title: "[演示] 高等数学复习资料与时间安排",
        raw: "准备期末复习时，你会怎样安排概念回顾、例题练习和模拟测试？欢迎分享资料与方法。"
      },
      {
        key: "network-troubleshooting",
        category_slug: "help",
        title: "[演示] 校园网偶尔断线，如何排查？",
        raw: "宿舍网络偶尔会断开。请把设备、系统、发生时间和已经尝试的方法写清楚，方便一起定位。"
      },
      {
        key: "badminton-partners",
        category_slug: "social",
        title: "[演示] 周末羽毛球搭子招募",
        raw: "计划周末约一场轻松的羽毛球活动，新手也欢迎。回复时可以注明方便的时间段。"
      },
      {
        key: "data-structures-roadmap",
        category_slug: "courses",
        title: "[演示] 数据结构课程学习路线",
        raw: "从线性表到图算法，怎样把课堂知识、代码练习和题目复盘串成一条有效的学习路线？"
      },
      {
        key: "developer-tools",
        category_slug: "other",
        title: "[演示] 大家常用哪些开发工具？",
        raw: "欢迎介绍编辑器、终端、版本控制和调试工具，也可以分享一两个真正提升效率的配置。"
      },
      {
        key: "study-spaces",
        category_slug: "daily",
        title: "[演示] 晚自习地点体验分享",
        raw: "哪些教室或公共空间适合晚自习？可以从开放时间、插座、网络和安静程度几个方面评价。"
      },
      {
        key: "project-teams",
        category_slug: "courses",
        title: "[演示] 课程项目组队与分工建议",
        raw: "课程项目怎样组队、拆任务和同步进度更顺畅？这里汇总大家踩过的坑和有效做法。"
      },
      {
        key: "freshman-faq",
        category_slug: "help",
        title: "[演示] 新生常见问题汇总",
        raw: "请把选课、实验室、校园服务等常见问题集中到这里，后续再整理成便于检索的答案。"
      }
    ].freeze

    REPLIES = [
      {
        key: "reply-01",
        author: :primary,
        raw: "甲的一级留言：下面的回复只会出现在这一条留言下，不会追加到主题末尾。"
      },
      { key: "reply-02", author: :primary, raw: <<~'MARKDOWN' },
          > 这是回复中的引用块。

          引用后继续写正文，并补充 **加粗结论** 与 *斜体提示*。
        MARKDOWN
      {
        key: "reply-03",
        author: :peer,
        parent_key: "reply-01",
        raw: "乙回复甲：这是第一条楼中楼回复，用于检查 `reply_to_post_number` 关系。"
      },
      { key: "reply-04", author: :primary, raw: <<~'MARKDOWN' },
          回复也可以包含代码块：

          ```javascript
          const forum = { name: "智工 BBS", inviteOnly: true };
          console.log(forum.name);
          ```
        MARKDOWN
      {
        key: "reply-05",
        author: :peer,
        parent_key: "reply-02",
        raw: "这是对第二条回复的二级回复，并附上[中山大学官网](https://www.sysu.edu.cn/)链接。"
      },
      { key: "reply-06", author: :primary, raw: <<~'MARKDOWN' },
          回复中的表格测试：

          | 功能 | 状态 |
          | --- | --- |
          | 发主题 | 正常 |
          | 写回复 | 正常 |
        MARKDOWN
      {
        key: "reply-07",
        author: :peer,
        parent_key: "reply-04",
        raw: "这是对代码块回复的二级回复，用于验证点击“回复”后能定位到对应楼层。"
      },
      {
        key: "reply-08",
        author: :peer,
        raw: "回复中的本地图片：![Discourse 示例图标](/images/discourse-logo-sketch.png)"
      },
      {
        key: "reply-09",
        author: :peer,
        parent_key: "reply-06",
        raw: "这是对表格回复的二级回复：列表、表格与楼层关系可以同时存在。"
      },
      {
        key: "reply-10",
        author: :primary,
        parent_key: "reply-03",
        raw: "甲回复乙：虽然这条回复实际指向乙，但视觉上仍与乙同级，只显示在甲的一级留言下。 :tada:"
      }
    ].freeze

    def self.call(env: ENV, output: $stdout)
      new(env:, output:).call
    end

    def initialize(env:, output:)
      @env = env
      @output = output
      @created_topics = 0
      @created_replies = 0
      @updated_replies = 0
    end

    def call
      prevent_production_seed!

      user =
        ensure_user(
          email: email,
          username: username,
          password: password,
          name: "智工演示用户"
        )
      peer_user =
        ensure_user(
          email: peer_email,
          username: peer_username,
          password: peer_password,
          name: "智工同学乙"
        )
      topics = TOPICS.map { |attributes| ensure_topic(user, attributes) }
      showcase_topic = topics.first
      ensure_replies({ primary: user, peer: peer_user }, showcase_topic)

      @output.puts "Demo login: #{user.email} / #{password}"
      @output.puts "Demo peer: #{peer_user.username}"
      @output.puts "Demo data: #{@created_topics} topics and #{@created_replies} replies created; " \
                     "#{@updated_replies} replies refreshed."
      @output.puts "Showcase topic: #{showcase_topic.relative_url}"

      Result.new(user:, peer_user:, topics:, showcase_topic:)
    end

    private

    def prevent_production_seed!
      return unless Rails.env.production?
      return if @env["ZGIE_ALLOW_DEMO_SEED"] == "true"

      raise "Refusing to create demo credentials in production without ZGIE_ALLOW_DEMO_SEED=true"
    end

    def ensure_user(email:, username:, password:, name:)
      user = User.find_by_email(email)
      user ||= User.new(email:, username:, name:)
      user.password = password unless user.persisted? &&
        user.confirm_password?(password)
      user.approved = true
      user.save!
      user.activate unless user.active?

      trust_level_one = TrustLevel[1]
      if user.trust_level < trust_level_one
        user.change_trust_level!(trust_level_one)
      end
      user
    end

    def ensure_topic(user, attributes)
      key = attributes.fetch(:key)
      topic = topic_by_key(key)
      return topic if topic

      category = Category.find_by!(slug: attributes.fetch(:category_slug))
      first_post =
        PostCreator.create!(
          user,
          title: attributes.fetch(:title),
          raw: attributes.fetch(:raw),
          category: category.id,
          custom_fields: {
            TOPIC_KEY_FIELD => key
          },
          post_alert_options: {
            skip_send_email: true
          },
          skip_validations: true
        )
      @created_topics += 1
      first_post.topic
    end

    def ensure_replies(authors, topic)
      posts_by_key = {}

      REPLIES.each do |attributes|
        key = attributes.fetch(:key)
        post = post_by_key(key)
        author = authors.fetch(attributes.fetch(:author))
        parent = posts_by_key.fetch(attributes[:parent_key]) if attributes[
          :parent_key
        ]

        if post && post.topic_id != topic.id
          raise "Demo reply #{key} belongs to another topic"
        end

        post = create_reply(author, topic, attributes, parent) unless post
        reconcile_reply(post, author, attributes, parent)

        posts_by_key[key] = post
      end
    end

    def reconcile_reply(post, author, attributes, parent)
      updated = false
      changes = {}
      raw = attributes.fetch(:raw)
      reply_to_post_number = parent&.post_number

      changes[:raw] = raw if post.raw.strip != raw.strip
      if post.reply_to_post_number != reply_to_post_number
        changes[:reply_to_post_number] = reply_to_post_number
      end

      if changes.present?
        revised =
          post.revise(
            Discourse.system_user,
            changes,
            bypass_bump: true,
            silent: true,
            skip_revision: true,
            skip_validations: true
          )
        if !revised && post.errors.present?
          raise "Could not refresh demo reply #{attributes.fetch(:key)}: " \
                  "#{post.errors.full_messages.join(", ")}"
        end

        updated = true if revised
      end

      if post.user_id != author.id
        PostOwnerChanger.new(
          post_ids: [post.id],
          topic_id: post.topic_id,
          new_owner: author,
          acting_user: Discourse.system_user,
          skip_revision: true
        ).change_owner!
        updated = true
      end

      @updated_replies += 1 if updated
      post.reload if updated
      post
    end

    def create_reply(user, topic, attributes, parent)
      options = {
        topic_id: topic.id,
        raw: attributes.fetch(:raw),
        custom_fields: {
          REPLY_KEY_FIELD => attributes.fetch(:key)
        },
        post_alert_options: {
          skip_send_email: true
        },
        skip_validations: true
      }
      options[:reply_to_post_number] = parent.post_number if parent

      PostCreator.create!(user, options).tap { @created_replies += 1 }
    end

    def topic_by_key(key)
      field = TopicCustomField.find_by(name: TOPIC_KEY_FIELD, value: key)
      Topic.find_by(id: field&.topic_id)
    end

    def post_by_key(key)
      field = PostCustomField.find_by(name: REPLY_KEY_FIELD, value: key)
      Post.find_by(id: field&.post_id)
    end

    def email
      @email ||= fetch("ZGIE_DEMO_EMAIL", DEFAULT_EMAIL).strip.downcase
    end

    def username
      @username ||= fetch("ZGIE_DEMO_USERNAME", DEFAULT_USERNAME).strip
    end

    def password
      @password ||= fetch("ZGIE_DEMO_PASSWORD", DEFAULT_PASSWORD)
    end

    def peer_email
      @peer_email ||=
        fetch("ZGIE_DEMO_PEER_EMAIL", DEFAULT_PEER_EMAIL).strip.downcase
    end

    def peer_username
      @peer_username ||=
        fetch("ZGIE_DEMO_PEER_USERNAME", DEFAULT_PEER_USERNAME).strip
    end

    def peer_password
      @peer_password ||= fetch("ZGIE_DEMO_PEER_PASSWORD", DEFAULT_PEER_PASSWORD)
    end

    def fetch(name, default)
      @env.fetch(name, default)
    end
  end
end
