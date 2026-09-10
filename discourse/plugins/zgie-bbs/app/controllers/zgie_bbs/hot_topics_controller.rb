# frozen_string_literal: true

module ZgieBbs
  class HotTopicsController < ::ApplicationController
    requires_plugin PLUGIN_NAME
    before_action :ensure_logged_in, if: -> { SiteSetting.login_required }

    def index
      raise Discourse::NotFound unless SiteSetting.zgie_bbs_branded_layout

      topics = CampusTopicQuery.new(current_user).hot_topics

      render json: {
               topics:
                 topics.map { |topic|
                   {
                     id: topic.id,
                     title: topic.title,
                     url: topic.relative_url,
                     category_id: topic.category_id,
                     reply_count: [topic.posts_count - 1, 0].max
                   }
                 }
             }
    end
  end
end
