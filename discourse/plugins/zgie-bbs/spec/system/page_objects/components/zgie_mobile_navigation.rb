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

      def has_clear_topic_actions?
        return false unless page.has_css?(".nested-view__floating-reply")
        page.evaluate_script(<<~JS)
          (() => {
            const nav = document.querySelector('.zgie-mobile-nav').getBoundingClientRect();
            const reply = document.querySelector('.nested-view__floating-reply').getBoundingClientRect();
            const progress = document.querySelector('.zgie-topic-position').getBoundingClientRect();
            return reply.bottom <= nav.top && progress.bottom <= nav.top && progress.right < reply.left;
          })()
        JS
      end
    end
  end
end
