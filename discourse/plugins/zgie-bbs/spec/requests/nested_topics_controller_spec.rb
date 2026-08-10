# frozen_string_literal: true

RSpec.describe NestedTopicsController do
  fab!(:primary_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:peer_user) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:topic) { Fabricate(:topic, user: primary_user) }
  fab!(:op) { Fabricate(:post, topic:, user: primary_user, post_number: 1) }
  fab!(:root_reply) do
    Fabricate(:post, topic:, user: primary_user, raw: "甲的一级留言")
  end
  fab!(:child_reply) do
    Fabricate(
      :post,
      topic:,
      user: peer_user,
      raw: "乙回复甲的留言",
      reply_to_post_number: root_reply.post_number
    )
  end
  fab!(:reply_to_child) do
    Fabricate(
      :post,
      topic:,
      user: primary_user,
      raw: "甲继续回复乙",
      reply_to_post_number: child_reply.post_number
    )
  end

  before do
    ZgieBbs::SiteConfigurator.call(
      env: {
        "ZGIE_REGISTRATION_MODE" => "global_code",
        "ZGIE_INVITE_CODE" => "nested-controller-spec"
      },
      output: StringIO.new
    )
    sign_in(primary_user)
  end

  describe "#show" do
    it "keeps replies out of the root stream" do
      get "/n/#{topic.slug}/#{topic.id}.json", params: { sort: "old" }

      expect(response.status).to eq(200)
      expect(response.parsed_body["roots"].map { |root| root["id"] }).to eq(
        [root_reply.id]
      )
    end
  end

  describe "#children" do
    it "shows a reply to a child as its sibling under the same root" do
      get "/n/#{topic.slug}/#{topic.id}/children/#{root_reply.post_number}.json",
          params: {
            depth: 1,
            sort: "old"
          }

      expect(response.status).to eq(200)
      children = response.parsed_body["children"]
      expect(children.map { |child| child["id"] }).to eq(
        [child_reply.id, reply_to_child.id]
      )
      expect(children.map { |child| child["children"] }).to all(eq([]))
      expect(children.last["reply_to_post_number"]).to eq(
        child_reply.post_number
      )
      expect(children.last.dig("reply_to_user", "username")).to eq(
        peer_user.username
      )
      expect(children.first["zgie_replies_to_nested_reply"]).to eq(false)
      expect(children.last["zgie_replies_to_nested_reply"]).to eq(true)
    end
  end
end
