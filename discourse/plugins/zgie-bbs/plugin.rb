# frozen_string_literal: true

# name: zgie-bbs
# about: Branding and reproducible site configuration for the ZGIE invite-only BBS.
# version: 0.1.0
# authors: ZGIE Market Team
# url: https://github.com/changyicheng1234/zgie-market
# required_version: 3.5.0

enabled_site_setting :zgie_bbs_enabled

register_asset "stylesheets/common/zgie-bbs.scss"
register_asset "stylesheets/common/zgie-mobile.scss"
register_asset "stylesheets/common/zgie-branded.scss"
register_svg_icon "fire"
register_svg_icon "book-open"
register_svg_icon "compass"

module ::ZgieBbs
  PLUGIN_NAME = "zgie-bbs"

  module PostReplyTargetExtension
    def zgie_replies_to_nested_reply?
      if defined?(@zgie_replies_to_nested_reply)
        return @zgie_replies_to_nested_reply
      end

      parent = reply_to_post
      @zgie_replies_to_nested_reply =
        parent.present? && parent.reply_to_post_number.present? &&
          parent.reply_to_post_number != 1
    end
  end

  module NestedReplyTargetPreloaderExtension
    def prepare(posts)
      preload_zgie_reply_target_levels(posts)
      super
    end

    private

    def preload_zgie_reply_target_levels(posts)
      posts_by_number = posts.index_by(&:post_number)
      target_numbers =
        posts
          .filter_map do |post|
            target = post.reply_to_post_number
            target if target.present? && target != 1
          end
          .uniq
      missing_target_numbers = target_numbers - posts_by_number.keys
      missing_target_levels =
        if missing_target_numbers.empty?
          {}
        else
          Post
            .where(topic_id: @topic.id, post_number: missing_target_numbers)
            .pluck(:post_number, :reply_to_post_number)
            .to_h
        end

      posts.each do |post|
        target_number = post.reply_to_post_number
        target = posts_by_number[target_number]
        target_parent_number =
          target&.reply_to_post_number || missing_target_levels[target_number]

        post.instance_variable_set(
          :@zgie_replies_to_nested_reply,
          target_number.present? && target_number != 1 &&
            target_parent_number.present? && target_parent_number != 1
        )
      end
    end
  end
end

on(:post_created) do |post, options, author|
  ZgieBbs::PerPostAnonymity.clear_master_draft!(post, options, author)
  Jobs.enqueue(:zgie_notify_wecom, post_id: post.id)
end

after_initialize do
  require_relative "lib/zgie_bbs/demo_seeder"
  require_relative "lib/zgie_bbs/per_post_anonymity"
  require_relative "lib/zgie_bbs/privacy"
  require_relative "lib/zgie_bbs/site_configurator"
  require_relative "lib/zgie_bbs/wecom_notify"
  require_relative "jobs/regular/zgie_notify_wecom"

  add_permitted_post_create_param :zgie_anonymous
  register_modifier(:posts_controller_create_user) do |user, create_params|
    ZgieBbs::PerPostAnonymity.author_for(user, create_params)
  end

  reloadable_patch do
    Guardian.prepend(ZgieBbs::AnonymousProfileGuardianExtension)
    Guardian.prepend(ZgieBbs::AnonymousPostOwnershipGuardianExtension)
    PostDestroyer.prepend(ZgieBbs::AnonymousPostDestroyerActorExtension)
    PostRevisor.prepend(ZgieBbs::AnonymousPostRevisorActorExtension)
    UsersController.prepend(ZgieBbs::PerPostAnonymousSessionControllerExtension)
    PostActionUsersController.prepend(
      ZgieBbs::HiddenLikeActorsControllerExtension
    )
    Post.prepend(ZgieBbs::PostReplyTargetExtension)
    NestedReplies::PostPreloader.prepend(
      ZgieBbs::NestedReplyTargetPreloaderExtension
    )

    reaction_controller =
      "DiscourseReactions::CustomReactionsController".safe_constantize
    if reaction_controller
      reaction_controller.prepend(
        ZgieBbs::HiddenReactionActorsControllerExtension
      )
    end
  end

  add_to_serializer(
    :post,
    :zgie_replies_to_nested_reply,
    include_condition: -> { topic&.nested_view? }
  ) { object.zgie_replies_to_nested_reply? }

  add_to_serializer(
    :post,
    :zgie_anonymous_author,
    include_condition: -> { SiteSetting.allow_anonymous_mode }
  ) { object.user&.anonymous? || false }
end
