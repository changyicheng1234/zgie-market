# frozen_string_literal: true

module ZgieBbs
  module PerPostAnonymity
    BOOLEAN_TYPE = ActiveModel::Type::Boolean.new

    def self.author_for(user, params)
      master_user = AnonymousShadowCreator.get_master(user) || user
      return master_user unless BOOLEAN_TYPE.cast(params[:zgie_anonymous])

      raise Discourse::InvalidAccess if private_message?(params)

      shadow_user =
        AnonymousShadowCreator.get(master_user) ||
          raise(Discourse::InvalidAccess)

      # AnonymousShadowCreator assigns TL1 immediately, while automatic group
      # membership may otherwise be refreshed asynchronously. The same request
      # must already be able to pass category and posting permissions.
      Group.refresh_automatic_groups_for_user!(shadow_user)
      shadow_user
    end

    def self.private_message?(params)
      return true if params[:archetype] == Archetype.private_message
      return false if params[:topic_id].blank?

      Topic.where(id: params[:topic_id], archetype: Archetype.private_message).exists?
    end

    def self.actor_for_owned_anonymous_content(user, content)
      return user if user.blank? || content.blank?
      return user unless SiteSetting.zgie_bbs_enabled

      anonymous_author = content.user
      anonymous_author&.master_user == user ? anonymous_author : user
    end
  end

  # Native anonymous mode changes the whole browser session. ZGIE anonymity is
  # selected in the composer for one new post instead. Existing anonymous
  # sessions can still use the endpoint once to return to their master account.
  module PerPostAnonymousSessionControllerExtension
    def toggle_anon
      return super unless SiteSetting.zgie_bbs_enabled
      return super if current_user&.anonymous?

      render json: failed_json, status: :forbidden
    end
  end

  # A member keeps moderation rights over posts created by their own shadow
  # account without exposing that relationship to other members or serializers.
  module AnonymousPostOwnershipGuardianExtension
    private

    def is_my_own?(object)
      return true if owns_zgie_anonymous_content?(object)

      super
    end

    def owns_zgie_anonymous_content?(object)
      return false unless SiteSetting.zgie_bbs_enabled
      return false unless object.is_a?(Post) || object.is_a?(Topic)
      return false if object.user_id.blank? || @user.blank? || @user.anonymous?

      zgie_anonymous_author_ids.include?(object.user_id)
    end

    def zgie_anonymous_author_ids
      @zgie_anonymous_author_ids ||=
        AnonymousUser.where(master_user_id: @user.id).pluck(:user_id).to_set
    end
  end

  # Editing or deleting as the master account would put the real user in the
  # public revision/deletion metadata. Once ownership is authorized above,
  # execute the mutation as that post's shadow author as well.
  module AnonymousPostRevisorActorExtension
    def revise!(editor, fields, opts = {}, &block)
      editor =
        PerPostAnonymity.actor_for_owned_anonymous_content(editor, @post)
      super(editor, fields, opts, &block)
    end
  end

  module AnonymousPostDestroyerActorExtension
    def initialize(user, post, opts = {})
      user = PerPostAnonymity.actor_for_owned_anonymous_content(user, post)
      super(user, post, opts)
    end
  end
end
