# frozen_string_literal: true

require_relative "page_objects/components/zgie_nested_replies"

RSpec.describe "Fixed nested replies" do
  fab!(:primary_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:peer_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:topic) { Fabricate(:topic, user: primary_user) }
  fab!(:op) { Fabricate(:post, topic:, user: primary_user, post_number: 1) }
  fab!(:root_reply) do
    Fabricate(:post, topic:, user: primary_user, raw: "甲发布的一级留言")
  end
  fab!(:child_reply) do
    Fabricate(
      :post,
      topic:,
      user: peer_user,
      raw: "乙在楼中楼回复甲",
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

    expect(reply_tree).to have_reply_at_depth("乙在楼中楼回复甲", depth: 1)
    nested_view.click_reply_on_post(child_reply)
    composer.fill_content(reply_text)
    composer.submit

    expect(composer).to be_closed
    expect(reply_tree).to have_reply_at_depth(reply_text, depth: 1)
    expect(reply_tree).to have_no_reply_at_depth(reply_text, depth: 0)
  end
end
