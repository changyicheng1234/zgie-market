# frozen_string_literal: true

module PageObjects
  module Components
    class ZgieDesktopCampus < PageObjects::Components::Base
      def has_expanded_categories?
        page.has_css?(
          '[data-section-name="zgie-campus"].sidebar-section--expanded'
        ) &&
          %w[实习推荐 考研保研分享 二手交易 课程经验分享 日常吐槽].all? do |name|
            page.has_css?('[data-section-name="zgie-campus"] a', text: name)
          end
      end

      def has_no_default_categories?
        page.has_no_css?('.sidebar-section[data-section-name="categories"]')
      end

      def has_five_topics?
        page.has_css?(".zgie-hot-topics__title", count: 5)
      end

      def has_sidebar_beside_list?
        page.evaluate_script(<<~JS)
          (() => {
            const sidebar = document.querySelector('.zgie-hot-topics').getBoundingClientRect();
            const list = document.querySelector('.list-container').getBoundingClientRect();
            return sidebar.left >= list.right && sidebar.top < list.bottom;
          })()
        JS
      end

      def open_first_topic
        find(".zgie-hot-topics__title", match: :first).click
      end

      def has_no_desktop_additions?
        page.has_no_css?(".zgie-hot-topics") &&
          page.has_no_css?('[data-section-name="zgie-campus"]')
      end
    end
  end
end
