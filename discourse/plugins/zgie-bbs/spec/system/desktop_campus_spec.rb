# frozen_string_literal: true

RSpec.describe "Desktop campus layout" do
  fab!(:user)
  let(:campus) { PageObjects::Components::ZgieDesktopCampus.new }

  before do
    SiteSetting.zgie_bbs_enabled = true
    SiteSetting.zgie_bbs_branded_layout = true
    ZgieBbs::SiteConfigurator.configure_campus_categories(output: StringIO.new)
    sign_in(user)
  end

  it "shows expanded campus categories and five working trending links" do
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
    visit("/latest")
    expect(campus).to have_expanded_categories
    expect(campus).to have_no_default_categories
    expect(campus).to have_five_topics
    expect(campus).to have_sidebar_beside_list
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
end
