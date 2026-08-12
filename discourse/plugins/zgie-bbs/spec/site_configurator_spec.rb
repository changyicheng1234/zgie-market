# frozen_string_literal: true

RSpec.describe ZgieBbs::SiteConfigurator do
  let(:env) do
    {
      "ZGIE_SITE_TITLE" => "智工测试站",
      "ZGIE_SITE_DESCRIPTION" => "测试描述",
      "ZGIE_DEFAULT_LOCALE" => "zh_CN",
      "ZGIE_LOGIN_REQUIRED" => "true",
      "ZGIE_EXTERNAL_HOSTNAME" => "bbs.test.localhost",
      "ZGIE_EXTERNAL_PORT" => "4200",
      "ZGIE_REGISTRATION_MODE" => "global_code",
      "ZGIE_INVITE_CODE" => "test-code"
    }
  end

  describe ".call" do
    it "configures a private local-auth site with an invite code" do
      described_class.call(env:, output: StringIO.new)

      expect(SiteSetting.title).to eq("智工测试站")
      expect(SiteSetting.login_required).to eq(true)
      expect(SiteSetting.invite_only).to eq(false)
      expect(SiteSetting.invite_code).to eq("test-code")
      expect(SiteSetting.force_hostname).to eq("bbs.test.localhost")
      expect(SiteSetting.port).to eq("4200")
      expect(SiteSetting.external_system_avatars_url).to eq("")
      expect(Category.find_by(slug: "study")&.name).to eq("学习交流")
    end

    it "enforces one visual level of nested replies" do
      described_class.call(env:, output: StringIO.new)

      expect(
        [
          SiteSetting.nested_replies_enabled,
          SiteSetting.nested_replies_default,
          SiteSetting.nested_replies_max_depth,
          SiteSetting.nested_replies_cap_nesting_depth,
          SiteSetting.nested_replies_default_sort
        ]
      ).to eq([true, true, 1, true, "old"])
    end

    it "allows per-post anonymity and removes editorial length thresholds" do
      described_class.call(env:, output: StringIO.new)

      expect(
        [
          SiteSetting.min_topic_title_length,
          SiteSetting.max_topic_title_length,
          SiteSetting.min_first_post_length,
          SiteSetting.min_post_length,
          SiteSetting.max_post_length
        ]
      ).to eq([1, 255, 1, 1, 150_000])
      expect(SiteSetting.allow_anonymous_mode).to eq(true)
      expect(SiteSetting.allow_likes_in_anonymous_mode).to eq(true)
      expect(
        SiteSetting.anonymous_posting_allowed_groups_map
      ).to contain_exactly(Group::AUTO_GROUPS[:trust_level_0])
      expect(SiteSetting.discourse_reactions_enabled).to eq(true)
    end

    it "supports Discourse invite links without a global code" do
      env["ZGIE_REGISTRATION_MODE"] = "invite_links"
      env["ZGIE_INVITE_CODE"] = ""

      described_class.call(env:, output: StringIO.new)

      expect(SiteSetting.invite_only).to eq(true)
      expect(SiteSetting.invite_code).to eq("")
    end

    it "rejects a blank global invite code" do
      env["ZGIE_INVITE_CODE"] = ""

      expect {
        described_class.call(env:, output: StringIO.new)
      }.to raise_error(ArgumentError, /cannot be blank/)
    end
  end
end
