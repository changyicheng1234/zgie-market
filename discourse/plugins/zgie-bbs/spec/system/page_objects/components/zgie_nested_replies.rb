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

      def has_no_reply_relation?(post)
        has_no_css?(wrapper_selector(post, ".zgie-reply-target"))
      end

      def has_no_jump_control?
        has_no_css?(".nested-post__continue-link")
      end

      def has_noninteractive_anonymous_identity?(post)
        wrapper = post_wrapper_selector(post)

        has_css?("#{wrapper} .zgie-anonymous-avatar") &&
          has_no_css?("#{wrapper} .topic-avatar a") &&
          has_css?("#{wrapper} .names .zgie-anonymous-name") &&
          has_no_css?("#{wrapper} .names [data-user-card]")
      end

      def has_topic_position?(current:, total:)
        has_css?(
          ".zgie-topic-position__count",
          text: %r{\A#{current}\s*/\s*#{total}\z}
        )
      end

      def has_topic_position_to_right_of_content?
        page.evaluate_script(<<~JS)
          (() => {
            const position = document.querySelector(".zgie-topic-position");
            const content = document.querySelector(".nested-view__op");
            return !position.classList.contains("--compact") &&
              position.getBoundingClientRect().left >= content.getBoundingClientRect().right + 16;
          })()
        JS
      end

      def scroll_post_to_reading_line(post)
        page.execute_script(<<~JS)
          const post = document.querySelector(
            ".nested-post__article[data-post-number='#{post.post_number}']"
          );
          const readingLine = Math.min(window.innerHeight * 0.36, 240);
          window.scrollTo(0, window.scrollY + post.getBoundingClientRect().top - readingLine + 20);
        JS
      end

      def has_read_only_reaction_summary?(post, reaction:)
        wrapper = post_wrapper_selector(post)
        counter = "#{wrapper} .discourse-reactions-counter"

        has_css?(
          "#{wrapper} #discourse-reactions-list-emoji-#{post.id}-#{reaction}"
        ) &&
          page.evaluate_script(
            "getComputedStyle(document.querySelector(#{counter.to_json})).pointerEvents"
          ) == "none"
      end

      def grow_content_above_reading_position
        page.execute_script(<<~JS)
          document.documentElement.style.overflowAnchor = "none";
          const delayedContent = document.createElement("div");
          delayedContent.style.height = "500px";
          document.querySelector(".nested-view__op-article").append(delayedContent);
        JS
        page.evaluate_async_script(<<~JS)
          const done = arguments[0];
          requestAnimationFrame(() => requestAnimationFrame(() => done(true)));
        JS
      end

      def return_to_top
        page.execute_script("window.scrollTo(0, 0)")
      end

      def force_click_reaction_summary(post)
        selector = "#{post_wrapper_selector(post)} .discourse-reactions-counter"
        page.execute_script(
          "document.querySelector(#{selector.to_json}).click()"
        )
      end

      def has_no_reaction_actor_popup?
        has_no_css?(".users-popup")
      end

      def load_more(post)
        find(wrapper_selector(post, ".nested-post-children__load-more")).click
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

      def post_wrapper_selector(post)
        ".nested-post:has(> .nested-post__main > [data-post-number='#{post.post_number}'])"
      end

      def wrapper_selector(post, child_selector)
        "#{post_wrapper_selector(post)} #{child_selector}"
      end
    end
  end
end
