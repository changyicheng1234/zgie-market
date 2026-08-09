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
    end
  end
end
