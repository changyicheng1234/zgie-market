# frozen_string_literal: true

desc "Apply the idempotent ZGIE BBS site settings and initial categories"
task "zgie_bbs:configure" => :environment do
  require_relative "../zgie_bbs/site_configurator"
  ZgieBbs::SiteConfigurator.call
end
