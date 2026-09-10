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
      configure_branding
      configure_urls
      configure_access
      configure_discussions
      create_categories
      @output.puts "智工 BBS configuration applied."
    end

    private

    def configure_identity
      SiteSetting.title = fetch("ZGIE_SITE_TITLE", "智工集市")
      SiteSetting.site_description =
        fetch("ZGIE_SITE_DESCRIPTION", "智能工程学院师生的邀请制交流社区")
      SiteSetting.default_locale = fetch("ZGIE_DEFAULT_LOCALE", "zh_CN")
      SiteSetting.enable_local_logins = true
      SiteSetting.enable_local_logins_via_email = true
      SiteSetting.allow_new_registrations = true
    end

    # Vectorized from a source PNG (potrace, traced per flat color layer) so
    # it stays crisp at favicon and header sizes; the wordmark text in
    # logo.svg is outlined to paths (Inkscape --export-text-to-path) so it
    # renders identically regardless of the client's installed fonts.
    # UploadCreator dedupes by sha1, so re-running this on an
    # already-configured site is a no-op.
    def configure_branding
      images_dir = File.expand_path("../../assets/images", __dir__)
      set_upload_setting(:favicon, File.join(images_dir, "favicon.svg"))
      set_upload_setting(:logo, File.join(images_dir, "logo-light.svg"))
      set_upload_setting(:logo_dark, File.join(images_dir, "logo.svg"))
      set_upload_setting(:logo_small, File.join(images_dir, "favicon.svg"))
    end

    def set_upload_setting(setting_name, path)
      return unless File.exist?(path)

      upload =
        UploadCreator.new(
          File.open(path),
          File.basename(path),
          for_site_setting: true,
          site_setting_name: setting_name.to_s
        ).create_for(Discourse::SYSTEM_USER_ID)

      SiteSetting.public_send("#{setting_name}=", upload) if upload.persisted?
    end

    def configure_urls
      hostname = fetch("ZGIE_EXTERNAL_HOSTNAME", "").strip
      port = fetch("ZGIE_EXTERNAL_PORT", "").strip

      SiteSetting.force_hostname = hostname if hostname.present?
      SiteSetting.port = port if port.present?

      # The default avatar proxy depends on avatars.discourse-cdn.com. Local
      # letter avatars keep an invite-only installation self-contained.
      SiteSetting.external_system_avatars_url = ""
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

    def configure_discussions
      # Empty content remains invalid, but ZGIE does not impose editorial
      # character-count thresholds on topics or replies. The maximum is the
      # largest value supported by this Discourse version.
      SiteSetting.min_topic_title_length = 1
      SiteSetting.max_topic_title_length = 255
      SiteSetting.min_first_post_length = 1
      SiteSetting.min_post_length = 1
      SiteSetting.max_post_length = 150_000
      SiteSetting.allow_anonymous_mode = true
      SiteSetting.allow_likes_in_anonymous_mode = true
      SiteSetting.anonymous_posting_allowed_groups =
        Group::AUTO_GROUPS[:trust_level_0].to_s
      SiteSetting.discourse_reactions_enabled = true
      SiteSetting.nested_replies_enabled = true
      SiteSetting.nested_replies_default = true
      SiteSetting.nested_replies_max_depth = 1
      SiteSetting.nested_replies_cap_nesting_depth = true
      SiteSetting.nested_replies_default_sort = "old"
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
