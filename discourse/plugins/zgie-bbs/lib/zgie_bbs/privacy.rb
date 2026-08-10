# frozen_string_literal: true

module ZgieBbs
  module AnonymousProfileGuardianExtension
    def can_see_profile?(profile_user)
      return false if profile_user&.anonymous? && !is_staff?

      super
    end
  end

  module HiddenLikeActorsControllerExtension
    def index
      if params[:post_action_type_id].to_i ==
           PostActionType::LIKE_POST_ACTION_ID && !guardian.is_staff?
        raise Discourse::InvalidAccess
      end

      super
    end
  end

  module HiddenReactionActorsControllerExtension
    def reactions_users_list
      ensure_zgie_can_view_reaction_actors!
      super
    end

    def post_reactions_users
      ensure_zgie_can_view_reaction_actors!
      super
    end

    private

    def ensure_zgie_can_view_reaction_actors!
      raise Discourse::InvalidAccess unless guardian.is_staff?
    end
  end
end
