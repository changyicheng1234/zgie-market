# frozen_string_literal: true

module PageObjects
  module Components
    class ZgieDiscoveryBanner < PageObjects::Components::Base
      def has_category?(category)
        page.has_css?(".zgie-banner__category", text: category.name)
      end

      def select_category(category)
        find(".zgie-banner__categories").click_link(category.name)
      end

      def search(term)
        find(".zgie-banner__search").fill_in(
          "zgie-community-search",
          with: term
        )
        find(".zgie-banner__search button[type='submit']").click
      end

      def has_no_banner?
        page.has_no_css?(".zgie-banner")
      end

      def has_branded_layout?
        page.has_css?("body.zgie-branded")
      end

      def has_no_branded_layout?
        page.has_no_css?("body.zgie-branded")
      end
    end
  end
end
