# frozen_string_literal: true

module ZgieBbs
  class SiteConfigurator
    CATEGORIES = [
      { name: "日常吐槽", slug: "daily", color: "E0891B" },
      { name: "学习交流", slug: "study", color: "0B6CF0" },
      { name: "打听求助", slug: "help", color: "2F855A" },
      { name: "恋爱交友", slug: "social", color: "D53F8C" },
      { name: "课程专区", slug: "courses", color: "3A4B8F" },
      { name: "其他", slug: "other", color: "718096" }
    ].freeze

    TRUE_VALUES = %w[1 true yes on].freeze
    REGISTRATION_MODES = %w[global_code invite_links].freeze

    def self.call(env: ENV, output: $stdout)
      new(env:, output:).call
    end

    def initialize(env:, output:)
      @env = env
      @output = output
    end

    def call
      configure_identity
      configure_access
      create_categories
      @output.puts "智工 BBS configuration applied."
    end

    private

    def configure_identity
      SiteSetting.title = fetch("ZGIE_SITE_TITLE", "智工 BBS")
      SiteSetting.site_description =
        fetch("ZGIE_SITE_DESCRIPTION", "智能工程学院师生的邀请制交流社区")
      SiteSetting.default_locale = fetch("ZGIE_DEFAULT_LOCALE", "zh_CN")
      SiteSetting.enable_local_logins = true
      SiteSetting.enable_local_logins_via_email = true
      SiteSetting.allow_new_registrations = true
    end

    def configure_access
      SiteSetting.login_required = truthy?(fetch("ZGIE_LOGIN_REQUIRED", "true"))
      SiteSetting.must_approve_users = false
      SiteSetting.invite_allowed_groups = "1|2"

      mode = fetch("ZGIE_REGISTRATION_MODE", "global_code")
      if REGISTRATION_MODES.exclude?(mode)
        raise ArgumentError,
              "ZGIE_REGISTRATION_MODE must be one of: #{REGISTRATION_MODES.join(", ")}"
      end

      if mode == "global_code"
        invite_code = fetch("ZGIE_INVITE_CODE", "").strip
        if invite_code.empty?
          raise ArgumentError,
                "ZGIE_INVITE_CODE cannot be blank in global_code mode"
        end

        SiteSetting.invite_only = false
        SiteSetting.invite_code = invite_code
      else
        SiteSetting.invite_only = true
        SiteSetting.invite_code = ""
      end

      @output.puts "Registration mode: #{mode}"
    end

    def create_categories
      CATEGORIES.each do |attributes|
        category = Category.find_or_initialize_by(slug: attributes.fetch(:slug))
        next unless category.new_record?

        category.name = attributes.fetch(:name)
        category.color = attributes.fetch(:color)
        category.text_color = "FFFFFF"
        category.user = Discourse.system_user
        category.save!
        @output.puts "Created category: #{category.name}"
      end
    end

    def fetch(name, default)
      @env.fetch(name, default)
    end

    def truthy?(value)
      TRUE_VALUES.include?(value.to_s.downcase)
    end
  end
end
