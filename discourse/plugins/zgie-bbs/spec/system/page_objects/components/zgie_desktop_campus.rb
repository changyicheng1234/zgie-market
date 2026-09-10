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
            const list = document.querySelector('.zgie-home__feed').getBoundingClientRect();
            return sidebar.left >= list.right && Math.abs(sidebar.top - list.top) <= 1;
          })()
        JS
      end

      def open_first_topic
        find(".zgie-hot-topics__title", match: :first).click
      end

      def has_matching_dark_cards?
        page.evaluate_script(<<~JS)
          [...document.querySelectorAll('.zgie-home__feed, .zgie-hot-topics')].every(element => {
            const style = getComputedStyle(element);
            return style.backgroundColor === 'rgb(26, 26, 26)' && style.borderRadius === '16px';
          }) && [...document.querySelectorAll('.zgie-home .topic-list-item')].every(element =>
            getComputedStyle(element).backgroundColor === 'rgb(26, 26, 26)'
          )
        JS
      end

      def has_twenty_home_posts?
        page.has_css?(".zgie-home__feed .topic-list-item", count: 20)
      end

      def has_ordered_navigation?
        links = all('[data-section-name="zgie-navigation"] a').map(&:text)
        links == %w[首页 帖子 我的帖子 我的消息] && page.evaluate_script(<<~JS)
          document.querySelector('[data-section-name="zgie-navigation"]').getBoundingClientRect().bottom <=
          document.querySelector('[data-section-name="zgie-campus"]').getBoundingClientRect().top
        JS
      end

      def open_all_posts
        find(".zgie-home__all").click
      end

      def has_plain_post_list?
        page.has_css?(".topic-list-item", minimum: 23) &&
          page.has_no_css?(".zgie-hot-topics") &&
          page.has_no_css?(".zgie-banner")
      end

      def has_no_desktop_additions?
        page.has_no_css?(".zgie-hot-topics") &&
          page.has_no_css?('[data-section-name="zgie-campus"]')
      end
    end
  end
end
