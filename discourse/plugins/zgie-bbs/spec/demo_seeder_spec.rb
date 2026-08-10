# frozen_string_literal: true

RSpec.describe ZgieBbs::DemoSeeder do
  let(:env) do
    {
      "ZGIE_DEMO_EMAIL" => "demo-seeder@example.com",
      "ZGIE_DEMO_USERNAME" => "demo_seeder",
      "ZGIE_DEMO_PASSWORD" => "RSpec-VioletRiver!8246",
      "ZGIE_DEMO_PEER_EMAIL" => "demo-peer-seeder@example.com",
      "ZGIE_DEMO_PEER_USERNAME" => "demo_peer_seeder",
      "ZGIE_DEMO_PEER_PASSWORD" => "RSpec-CedarBridge!3902"
    }
  end

  before do
    ZgieBbs::SiteConfigurator.call(
      env: {
        "ZGIE_REGISTRATION_MODE" => "global_code",
        "ZGIE_INVITE_CODE" => "demo-seeder-spec"
      },
      output: StringIO.new
    )
  end

  describe ".call" do
    it "creates a login and the complete rendering dataset" do
      result = described_class.call(env:, output: StringIO.new)
      replies =
        result.showcase_topic.posts.where("post_number > 1").order(:post_number)
      rendered =
        Nokogiri::HTML5.fragment(result.showcase_topic.first_post.cooked)
      rendered_elements =
        %w[h1 strong em blockquote pre table a img].map do |name|
          rendered.at_css(name)
        end

      expect(result.user).to be_active
      expect(result.peer_user).to be_active
      expect(result.user.avatar_template).to start_with("/letter_avatar/")
      expect(
        result.user.confirm_password?(env.fetch("ZGIE_DEMO_PASSWORD"))
      ).to eq(true)
      expect(result.topics.map(&:title)).to contain_exactly(
        *described_class::TOPICS.map { |topic| topic.fetch(:title) }
      )
      expect(replies.count).to eq(10)
      expect(rendered_elements).to all(be_present)
      expect(
        described_class::REPLIES.map { |reply| reply.fetch(:raw) }.join
      ).not_to match(/甲回复乙|乙回复甲/)
    end

    it "creates a reply-to-child chain with two distinct authors" do
      result = described_class.call(env:, output: StringIO.new)
      replies =
        described_class::REPLIES.to_h do |attributes|
          field =
            PostCustomField.find_by(
              name: described_class::REPLY_KEY_FIELD,
              value: attributes.fetch(:key)
            )
          [attributes.fetch(:key), Post.find(field.post_id)]
        end
      root = replies.fetch("reply-01")
      child = replies.fetch("reply-03")
      reply_to_child = replies.fetch("reply-10")

      expect(child.reply_to_post_number).to eq(root.post_number)
      expect(child.user).to eq(result.peer_user)
      expect(reply_to_child.reply_to_post_number).to eq(child.post_number)
      expect(reply_to_child.reply_to_user_id).to eq(result.peer_user.id)
      expect(reply_to_child.user).to eq(result.user)
    end

    it "does not duplicate users, topics, or posts when rerun" do
      first_result = described_class.call(env:, output: StringIO.new)
      counts = [User.count, Topic.count, Post.count]
      first_result.showcase_topic.first_post.update_columns(
        raw: "stale fixture",
        cooked: "<p>stale fixture</p>"
      )

      second_result = described_class.call(env:, output: StringIO.new)

      expect([User.count, Topic.count, Post.count]).to eq(counts)
      expect(second_result.showcase_topic.first_post.raw.strip).to eq(
        described_class::TOPICS.first.fetch(:raw).strip
      )
      expect(second_result.showcase_topic.first_post.cooked).not_to include(
        "stale fixture"
      )
    end
  end
end
