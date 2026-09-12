# frozen_string_literal: true

require_relative "page_objects/components/zgie_nested_replies"
require_relative "page_objects/components/zgie_anonymous_composer_toggle"

RSpec.describe "Fixed nested replies" do
  fab!(:primary_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:peer_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:topic) { Fabricate(:topic, user: primary_user) }
  fab!(:op) { Fabricate(:post, topic:, user: primary_user, post_number: 1) }
  fab!(:root_reply) do
    Fabricate(:post, topic:, user: primary_user, raw: "一级留言内容")
  end
  fab!(:child_reply) do
    Fabricate(
      :post,
      topic:,
      user: peer_user,
      raw: "二级留言内容",
      reply_to_post_number: root_reply.post_number
    )
  end
  fab!(:reply_to_child) do
    Fabricate(
      :post,
      topic:,
      user: primary_user,
      raw: "一级作者继续回复二级作者",
      reply_to_post_number: child_reply.post_number
    )
  end

  let(:nested_view) { PageObjects::Pages::NestedView.new }
  let(:topic_page) { PageObjects::Pages::Topic.new }
  let(:drafts_page) { PageObjects::Pages::UserActivityDrafts.new }
  let(:composer) { PageObjects::Components::Composer.new }
  let(:anonymous_toggle) do
    PageObjects::Components::ZgieAnonymousComposerToggle.new
  end
  let(:reply_tree) { PageObjects::Components::ZgieNestedReplies.new }
  let(:user_menu) { PageObjects::Components::UserMenu.new }

  before do
    ZgieBbs::SiteConfigurator.call(
      env: {
        "ZGIE_REGISTRATION_MODE" => "global_code",
        "ZGIE_INVITE_CODE" => "nested-system-spec"
      },
      output: StringIO.new
    )
    sign_in(primary_user)
  end

  it "keeps a reply to a child alongside that child under the same root" do
    reply_text = "甲回复乙后仍留在同一个楼中楼"
    nested_view.visit_nested(topic)

    expect(reply_tree).to have_reply_at_depth("二级留言内容", depth: 1)
    nested_view.click_reply_on_post(child_reply)
    composer.fill_content(reply_text)
    composer.submit

    expect(composer).to be_closed
    expect(reply_tree).to have_reply_at_depth(reply_text, depth: 1)
    expect(reply_tree).to have_no_reply_at_depth(reply_text, depth: 0)
  end

  it "only lets second-level replies collapse and selectively names the reply target" do
    nested_view.visit_nested(topic)

    expect(reply_tree).to have_no_collapse_control(root_reply)
    expect(reply_tree).to have_collapse_control(child_reply)
    expect(reply_tree).to have_no_reply_relation(child_reply)
    reply_tree.load_more(root_reply)
    expect(reply_tree).to have_reply_relation(
      reply_to_child,
      from: primary_user.username,
      to: peer_user.username
    )
    expect(reply_tree).to have_no_jump_control

    reply_tree.collapse(child_reply)

    expect(reply_tree).to have_collapsed_reply(child_reply)
    expect(reply_tree).to have_reply_at_depth("一级留言内容", depth: 0)

    reply_tree.expand(child_reply)

    expect(reply_tree).to have_reply_at_depth("二级留言内容", depth: 1)
  end

  it "renders anonymous authors without profile links and shows topic position" do
    anonymous_user = AnonymousShadowCreator.get(peer_user)
    anonymous_reply =
      Fabricate(
        :post,
        topic:,
        user: anonymous_user,
        raw: "这是一条不会暴露用户资料入口的匿名留言",
        reply_to_post_number: op.post_number
      )
    6.times do |index|
      Fabricate(
        :post,
        topic:,
        user: primary_user,
        raw: "用于验证滚动楼层定位的后续留言 #{index + 1}",
        reply_to_post_number: op.post_number
      )
    end

    resize_window(height: 700) do
      nested_view.visit_nested(topic)

      expect(reply_tree).to have_noninteractive_anonymous_identity(
        anonymous_reply
      )
      expect(
        page.evaluate_script(
          'getComputedStyle(document.querySelector(".zgie-topic-position")).position'
        )
      ).to eq("fixed")
      expect(reply_tree).to have_topic_position_to_right_of_content
      reply_tree.scroll_post_to_reading_line(anonymous_reply)
      expect(reply_tree).to have_topic_position(
        current: 4,
        total: topic.reload.highest_post_number
      )
    end
  end

  it "keeps the reading position moving forward past a late nested reply" do
    progress_topic = Fabricate(:topic, user: primary_user)
    progress_op =
      Fabricate(
        :post,
        topic: progress_topic,
        user: primary_user,
        post_number: 1
      )
    early_root =
      Fabricate(
        :post,
        topic: progress_topic,
        user: primary_user,
        raw: "较早的一级留言"
      )
    later_root =
      Fabricate(:post, topic: progress_topic, user: peer_user, raw: "后续一级留言")
    late_child =
      Fabricate(
        :post,
        topic: progress_topic,
        user: peer_user,
        raw: "稍后发布但插回前文的二级留言",
        reply_to_post_number: early_root.post_number
      )
    6.times do |index|
      Fabricate(
        :post,
        topic: progress_topic,
        user: primary_user,
        raw: "撑开滚动区域的后续留言 #{index + 1}",
        reply_to_post_number: progress_op.post_number
      )
    end

    resize_window(height: 700) do
      nested_view.visit_nested(progress_topic)

      reply_tree.scroll_post_to_reading_line(late_child)
      expect(reply_tree).to have_topic_position(current: 3, total: 10)

      reply_tree.scroll_post_to_reading_line(later_root)
      expect(reply_tree).to have_topic_position(current: 4, total: 10)

      reply_tree.grow_content_above_reading_position
      expect(reply_tree).to have_topic_position(current: 4, total: 10)

      reply_tree.return_to_top
      expect(reply_tree).to have_topic_position(current: 1, total: 10)
      reply_tree.scroll_post_to_reading_line(early_root)
      expect(reply_tree).to have_topic_position(current: 2, total: 10)
    end
  end

  it "chooses anonymity per submission and removes the user-menu switch" do
    nested_view.visit_nested(topic)
    nested_view.click_reply_on_post(root_reply)

    expect(anonymous_toggle).to have_label("匿名回复")
    expect(anonymous_toggle).to have_state(false)

    anonymous_toggle.click_visible_label
    expect(anonymous_toggle).to have_state(true)
    composer.fill_content("匿")
    composer.submit

    expect(composer).to be_closed
    anonymous_reply = Post.find_by!(topic:, raw: "匿")
    expect(anonymous_reply.user).to be_anonymous

    nested_view.visit_nested(topic)
    expect(reply_tree).to have_noninteractive_anonymous_identity(
      anonymous_reply
    )

    nested_view.click_reply_on_post(root_reply)
    expect(anonymous_toggle).to have_state(false)
    composer.fill_content("明")
    composer.submit
    expect(Post.find_by!(topic:, raw: "明").user).to eq(primary_user)

    visit "/new-topic"
    expect(anonymous_toggle).to have_label("匿名发帖")

    visit topic.url
    user_menu.open.click_profile_tab
    expect(page).to have_no_css(
      "#quick-access-profile .enable-anonymous, " \
        "#quick-access-profile .disable-anonymous"
    )
  end

  it "keeps a successfully published anonymous topic out of the user's drafts" do
    title = "匿名发布后不保留草稿"
    visit "/new-topic"
    composer.fill_title(title)
    composer.fill_content("发布成功后，这段正文不能继续出现在草稿箱中。")
    anonymous_toggle.click_visible_label

    try_until_success(reason: "wait for the composer autosave") do
      expect(
        Draft
          .where(user: primary_user)
          .where("draft_key LIKE ?", "#{Draft::NEW_TOPIC}_%")
          .exists?
      ).to eq(true)
    end

    composer.submit

    expect(topic_page).to have_topic_title(title)
    drafts_page.visit(primary_user)
    expect(drafts_page).to have_no_drafts
  end

  it "keeps reaction emoji and counts visible without opening actor lists" do
    DiscourseReactions::ReactionManager.new(
      reaction_value: "heart",
      user: peer_user,
      post: root_reply
    ).toggle!

    nested_view.visit_nested(topic)

    expect(reply_tree).to have_read_only_reaction_summary(
      root_reply,
      reaction: "heart"
    )
    reply_tree.force_click_reaction_summary(root_reply)
    expect(reply_tree).to have_no_reaction_actor_popup
  end
end
