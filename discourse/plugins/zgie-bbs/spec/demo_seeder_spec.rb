# frozen_string_literal: true

RSpec.describe ZgieBbs::DemoSeeder do
  let(:env) do
    {
      "ZGIE_DEMO_EMAIL" => "demo-seeder@example.com",
      "ZGIE_DEMO_USERNAME" => "demo_seeder",
      "ZGIE_DEMO_PASSWORD" => "RSpec-VioletRiver!8246"
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
      expect(
        result.user.confirm_password?(env.fetch("ZGIE_DEMO_PASSWORD"))
      ).to eq(true)
      expect(result.topics.map(&:title)).to contain_exactly(
        *described_class::TOPICS.map { |topic| topic.fetch(:title) }
      )
      expect(replies.count).to eq(10)
      expect(replies.filter_map(&:reply_to_post_number)).to eq([2, 3, 5, 7])
      expect(rendered_elements).to all(be_present)
    end

    it "does not duplicate users, topics, or posts when rerun" do
      described_class.call(env:, output: StringIO.new)
      counts = [User.count, Topic.count, Post.count]

      described_class.call(env:, output: StringIO.new)

      expect([User.count, Topic.count, Post.count]).to eq(counts)
    end
  end
end
