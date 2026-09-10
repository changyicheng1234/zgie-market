# frozen_string_literal: true

module PageObjects
  module Components
    class ZgieMobileNavigation < PageObjects::Components::Base
      def has_visible_navigation?
        page.has_css?(".zgie-mobile-nav")
      end

      def has_no_visible_navigation?
        page.has_no_css?(".zgie-mobile-nav")
      end

      def select(key)
        find(".zgie-mobile-nav").click_link(
          I18n.t("js.zgie_bbs.navigation.#{key}")
        )
      end

      def has_active_item?(key)
        page.has_css?(
          ".zgie-mobile-nav a[aria-current='page']",
          text: I18n.t("js.zgie_bbs.navigation.#{key}")
        )
      end

      def has_no_horizontal_overflow?
        page.evaluate_script(
          "document.documentElement.scrollWidth <= window.innerWidth"
        )
      end
    end
  end
end
