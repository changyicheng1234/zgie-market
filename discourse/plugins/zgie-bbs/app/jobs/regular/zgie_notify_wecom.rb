# frozen_string_literal: true

module Jobs
  class ZgieNotifyWecom < ::Jobs::Base
    def execute(args)
      post = Post.find_by(id: args[:post_id])
      return if post.blank?

      ZgieBbs::WecomNotify.notify_new_topic(post)
    end
  end
end
