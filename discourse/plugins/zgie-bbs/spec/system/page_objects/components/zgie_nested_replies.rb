# frozen_string_literal: true

module PageObjects
  module Components
    class ZgieNestedReplies < PageObjects::Components::Base
      def has_reply_at_depth?(text, depth:)
        has_css?(
          ".nested-post.--depth-#{depth} > .nested-post__main > article.nested-post__article",
          text:
        )
      end

      def has_no_reply_at_depth?(text, depth:)
        has_no_css?(
          ".nested-post.--depth-#{depth} > .nested-post__main > article.nested-post__article",
          text:
        )
      end

      def has_collapse_control?(post)
        has_css?(
          wrapper_selector(
            post,
            "> .nested-post__gutter .nested-post__depth-line"
          )
        )
      end

      def has_no_collapse_control?(post)
        has_no_css?(
          wrapper_selector(
            post,
            "> .nested-post__gutter .nested-post__depth-line"
          )
        )
      end

      def has_reply_relation?(post, from:, to:)
        relation = /#{Regexp.escape(from)}\s*>\s*#{Regexp.escape(to)}/
        has_css?(wrapper_selector(post, ".names"), text: relation)
      end

      def vertical_gap(root_post, child_post)
        page.evaluate_script(<<~JS)
          (() => {
            const rootArticle = document.querySelector(
              "[data-post-number='#{root_post.post_number}']"
            );
            const childWrapper = document
              .querySelector("[data-post-number='#{child_post.post_number}']")
              .closest(".nested-post");
            const childAvatar = childWrapper.querySelector(
              ":scope > .nested-post__gutter .topic-avatar"
            );

            return childAvatar.getBoundingClientRect().top -
              rootArticle.getBoundingClientRect().bottom;
          })()
        JS
      end

      def collapse(post)
        find(
          wrapper_selector(
            post,
            "> .nested-post__gutter .nested-post__depth-line"
          )
        ).click
      end

      def expand(post)
        find(
          wrapper_selector(
            post,
            "> .nested-post__main .nested-post__collapsed-bar"
          )
        ).click
      end

      def has_collapsed_reply?(post)
        has_css?(
          wrapper_selector(
            post,
            "> .nested-post__main .nested-post__collapsed-bar"
          )
        )
      end

      private

      def wrapper_selector(post, child_selector)
        wrapper =
          ".nested-post:has(> .nested-post__main > [data-post-number='#{post.post_number}'])"
        "#{wrapper} #{child_selector}"
      end
    end
  end
end
