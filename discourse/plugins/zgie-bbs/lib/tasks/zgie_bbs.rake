# frozen_string_literal: true

desc "Apply the idempotent ZGIE BBS site settings and initial categories"
task "zgie_bbs:configure" => :environment do
  require_relative "../zgie_bbs/site_configurator"
  ZgieBbs::SiteConfigurator.call
end

desc "Create the local ZGIE BBS demo account, topics, replies, and rich-text fixtures"
task "zgie_bbs:seed_demo" => %w[environment zgie_bbs:configure] do
  require_relative "../zgie_bbs/demo_seeder"
  ZgieBbs::DemoSeeder.call
end
