# frozen_string_literal: true

RSpec.describe "Explore branded layout" do
  fab!(:user)
  fab!(:category) do
    Fabricate(:category, name: "Campus study", slug: "courses")
  end
  let(:banner) { PageObjects::Components::ZgieDiscoveryBanner.new }

  before do
    SiteSetting.zgie_bbs_enabled = true
    SiteSetting.zgie_bbs_branded_layout = true
    SiteSetting.top_menu = "latest|hot|categories"
    sign_in(user)
  end

  it "lets a member explore a category and search from the welcome banner" do
    visit("/custom")
    expect(banner).to have_branded_layout
    expect(banner).to have_category(category)
    banner.select_category(category)
    expect(page).to have_current_path("/c/#{category.slug}/#{category.id}")

    visit("/custom")
    banner.search("campus")
    expect(page).to have_current_path("/search?q=campus")
    expect(banner).to have_branded_layout
    expect(banner).to have_no_banner
  end

  it "lets an administrator disable the visual layer without disabling navigation",
     mobile: true do
    SiteSetting.zgie_bbs_branded_layout = false
    visit("/latest")
    expect(banner).to have_no_branded_layout
    expect(banner).to have_no_banner
    expect(
      PageObjects::Components::ZgieMobileNavigation.new
    ).to have_visible_navigation
  end
end
