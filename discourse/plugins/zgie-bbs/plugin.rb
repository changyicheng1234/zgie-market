# frozen_string_literal: true

# name: zgie-bbs
# about: Branding and reproducible site configuration for the ZGIE invite-only BBS.
# version: 0.1.0
# authors: ZGIE Market Team
# url: https://github.com/changyicheng1234/zgie-market
# required_version: 3.5.0

enabled_site_setting :zgie_bbs_enabled

register_asset "stylesheets/common/zgie-bbs.scss"

module ::ZgieBbs
  PLUGIN_NAME = "zgie-bbs"
end

after_initialize { require_relative "lib/zgie_bbs/site_configurator" }
