# frozen_string_literal: true

require_relative "page_objects/components/zgie_nested_replies"

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

  let(:nested_view) { PageObjects::Pages::NestedView.new }
  let(:composer) { PageObjects::Components::Composer.new }
  let(:reply_tree) { PageObjects::Components::ZgieNestedReplies.new }

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

  it "only lets second-level replies collapse and names the reply target" do
    nested_view.visit_nested(topic)

    expect(reply_tree).to have_no_collapse_control(root_reply)
    expect(reply_tree).to have_collapse_control(child_reply)
    expect(reply_tree).to have_reply_relation(
      child_reply,
      from: peer_user.username,
      to: primary_user.username
    )
    expect(reply_tree.vertical_gap(root_reply, child_reply)).to be >= 8

    reply_tree.collapse(child_reply)

    expect(reply_tree).to have_collapsed_reply(child_reply)
    expect(reply_tree).to have_reply_at_depth("一级留言内容", depth: 0)

    reply_tree.expand(child_reply)

    expect(reply_tree).to have_reply_at_depth("二级留言内容", depth: 1)
  end
end
