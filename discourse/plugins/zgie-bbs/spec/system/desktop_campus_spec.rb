# frozen_string_literal: true

RSpec.describe "Desktop campus layout" do
  fab!(:user)
  let(:campus) { PageObjects::Components::ZgieDesktopCampus.new }

  before do
    SiteSetting.zgie_bbs_enabled = true
    SiteSetting.zgie_bbs_branded_layout = true
    SiteSetting.default_locale = "zh_CN"
    ZgieBbs::SiteConfigurator.configure_campus_categories(output: StringIO.new)
    sign_in(user)
  end

  it "shows expanded campus categories and five working trending links" do
    user.user_option.update!(interface_color_mode: UserOption::DARK_MODE)
    category = Category.find_by!(slug: "internships")
    6.times do |index|
      Fabricate(
        :post,
        topic:
          Fabricate(
            :topic,
            category: category,
            title: "Campus internship opportunity #{index}"
          )
      )
    end
    page.current_window.resize_to(1440, 1000)
    visit("/custom")
    expect(campus).to have_expanded_categories
    expect(campus).to have_no_default_categories
    expect(campus).to have_five_topics
    expect(campus).to have_sidebar_beside_list
    expect(campus).to have_matching_dark_cards
    campus.open_first_topic
    expect(page).to have_current_path(%r{^/t/})
  end

  it "keeps desktop additions out of mobile navigation", mobile: true do
    visit("/latest")
    expect(campus).to have_no_desktop_additions
    expect(
      PageObjects::Components::ZgieMobileNavigation.new
    ).to have_visible_navigation
  end

  it "separates the finite homepage from the full post list and orders navigation" do
    category = Category.find_by!(slug: "daily")
    23.times { Fabricate(:post, topic: Fabricate(:topic, category: category)) }
    page.current_window.resize_to(1440, 1000)
    visit("/")
    expect(page).to have_current_path("/custom")
    expect(campus).to have_twenty_home_posts
    expect(campus).to have_ordered_navigation
    campus.open_all_posts
    expect(page).to have_current_path("/latest")
    expect(campus).to have_plain_post_list
  end
end
