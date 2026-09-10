# frozen_string_literal: true

RSpec.describe "Use mobile navigation" do
  fab!(:user)
  fab!(:topic)
  fab!(:post) { Fabricate(:post, topic: topic) }

  let(:navigation) { PageObjects::Components::ZgieMobileNavigation.new }
  let(:topic_page) { PageObjects::Pages::Topic.new }
  let(:composer) { PageObjects::Components::Composer.new }

  before do
    SiteSetting.zgie_bbs_enabled = true
    SiteSetting.zgie_bbs_mobile_navigation = true
    SiteSetting.top_menu = "latest|hot|categories"
    SiteSetting.chat_enabled = false
    sign_in(user)
  end

  it "lets a mobile member switch between home, hot, messages and their profile",
     mobile: true do
    visit("/latest")
    expect(navigation).to have_active_item("home")
    expect(navigation).to have_no_horizontal_overflow

    navigation.select("hot")
    expect(page).to have_current_path("/hot")
    expect(navigation).to have_active_item("hot")

    navigation.select("profile")
    expect(page).to have_current_path("/u/#{user.username_lower}/summary")
    expect(navigation).to have_active_item("profile")

    navigation.select("messages")
    expect(page).to have_current_path("/u/#{user.username_lower}/messages")
    expect(navigation).to have_active_item("messages")

    page.refresh
    expect(navigation).to have_visible_navigation
    navigation.select("home")
    expect(page).to have_current_path("/latest")
  end

  it "gives the editor space when a mobile member replies", mobile: true do
    topic_page.visit_topic(topic)
    expect(navigation).to have_visible_navigation
    topic_page.click_post_action_button(post, :reply)
    expect(composer).to be_opened
    expect(navigation).to have_no_visible_navigation
    expect(navigation).to have_no_horizontal_overflow
  end

  it "keeps the desktop navigation clear of the mobile tab bar" do
    visit("/latest")
    expect(navigation).to have_no_visible_navigation
  end

  it "opens direct messages when chat is available", mobile: true do
    SiteSetting.chat_enabled = true
    SiteSetting.chat_allowed_groups = Group::AUTO_GROUPS[:everyone]
    SiteSetting.direct_message_enabled_groups = Group::AUTO_GROUPS[:everyone]
    user.user_option.update!(chat_enabled: true)
    visit("/latest")
    navigation.select("messages")
    expect(page).to have_current_path("/chat/direct-messages")
    expect(navigation).to have_active_item("messages")
    expect(navigation).to have_no_horizontal_overflow
  end
end
