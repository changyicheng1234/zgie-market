# frozen_string_literal: true

RSpec.describe "ZGIE privacy" do
  fab!(:member) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:viewer) { Fabricate(:user, refresh_auto_groups: true) }
  fab!(:admin) { Fabricate(:admin, refresh_auto_groups: true) }
  fab!(:topic) { Fabricate(:topic, user: member) }
  fab!(:op) { Fabricate(:post, topic:, user: member, post_number: 1) }

  before do
    ZgieBbs::SiteConfigurator.call(
      env: {
        "ZGIE_REGISTRATION_MODE" => "global_code",
        "ZGIE_INVITE_CODE" => "privacy-spec"
      },
      output: StringIO.new
    )
    SiteSetting.fast_typing_threshold = "disabled"
  end

  it "uses a shadow only for selected topics and replies without changing the session" do
    sign_in(member)
    category = Category.find_by!(slug: "study")

    post "/posts.json",
         params: {
           title: "题",
           raw: "文",
           category: category.id,
           zgie_anonymous: "true"
         }

    expect(response.status).to eq(200), response.body
    topic_post = Post.find(response.parsed_body.fetch("id"))
    anonymous_user = topic_post.user
    expect(anonymous_user).to be_anonymous
    expect(anonymous_user.master_user).to eq(member)
    expect(session[:current_user_id]).to eq(member.id)

    post "/posts.json",
         params: {
           topic_id: topic_post.topic_id,
           raw: "匿",
           zgie_anonymous: true
         }

    expect(response.status).to eq(200)
    anonymous_reply = Post.find(response.parsed_body.fetch("id"))
    expect(anonymous_reply.user).to eq(anonymous_user)
    expect(session[:current_user_id]).to eq(member.id)
    expect(Guardian.new(member).can_edit_post?(anonymous_reply)).to eq(true)

    post "/posts.json",
         params: {
           topic_id: topic_post.topic_id,
           raw: "明",
           zgie_anonymous: "false"
         }

    expect(response.status).to eq(200)
    named_reply = Post.find(response.parsed_body.fetch("id"))
    expect(named_reply.user).to eq(member)
    expect(session[:current_user_id]).to eq(member.id)

    SiteSetting.editing_grace_period = 0
    put "/posts/#{anonymous_reply.id}.json",
        params: {
          post: {
            raw: "匿名修改",
            original_text: anonymous_reply.raw
          }
        }

    expect(response.status).to eq(200), response.body
    expect(anonymous_reply.reload.last_editor_id).to eq(anonymous_user.id)

    sign_in(viewer)
    get "/posts/#{anonymous_reply.id}/revisions/latest.json"
    expect(response.status).to eq(200)
    expect(response.parsed_body.fetch("username")).to eq(
      anonymous_user.username_lower
    )
    expect(response.body).not_to include(member.username, member.email)

    sign_in(member)
    delete "/posts/#{anonymous_reply.id}.json"
    expect(response.status).to eq(200)
    expect(anonymous_reply.reload).to be_user_deleted
    expect(anonymous_reply.last_editor_id).to eq(anonymous_user.id)
  end

  it "blocks entering native session-wide anonymity but lets old sessions exit" do
    sign_in(member)

    post "/u/toggle-anon.json"

    expect(response).to be_forbidden
    expect(session[:current_user_id]).to eq(member.id)

    anonymous_user = AnonymousShadowCreator.get(member)
    sign_in(anonymous_user)

    post "/u/toggle-anon.json"

    expect(response.status).to eq(200)
    expect(session[:current_user_id]).to eq(member.id)
  end

  it "marks shadow-authored posts and hides their profile from regular members" do
    anonymous_user = AnonymousShadowCreator.get(member)
    anonymous_reply =
      Fabricate(:post, topic:, user: anonymous_user, raw: "匿名一级留言")

    expect(Guardian.new(viewer).can_see_profile?(anonymous_user)).to eq(false)
    expect(Guardian.new(admin).can_see_profile?(anonymous_user)).to eq(true)

    sign_in(viewer)
    get "/u/#{anonymous_user.username}.json"
    expect(response.status).to eq(200)
    anonymous_profile = response.parsed_body.fetch("user")
    expect(anonymous_profile["profile_hidden"]).to eq(true)
    expect(anonymous_profile).not_to include(
      "bio_raw",
      "email",
      "location",
      "website"
    )

    get "/n/#{topic.slug}/#{topic.id}.json", params: { sort: "old" }

    expect(response.status).to eq(200)
    serialized_reply =
      response.parsed_body["roots"].find do |root|
        root["id"] == anonymous_reply.id
      end
    expect(serialized_reply["zgie_anonymous_author"]).to eq(true)
  end

  it "keeps like and reaction actors private from members but auditable by staff" do
    PostActionCreator.like(viewer, op)
    like_params = {
      id: op.id,
      post_action_type_id: PostActionType::LIKE_POST_ACTION_ID
    }
    reaction_paths = [
      "/discourse-reactions/posts/#{op.id}/reactions-users-list.json",
      "/discourse-reactions/posts/#{op.id}/reactions-users.json"
    ]

    sign_in(member)
    get "/post_action_users.json", params: like_params
    expect(response).to be_forbidden
    reaction_paths.each do |path|
      get path
      expect(response).to be_forbidden
    end

    sign_in(admin)
    get "/post_action_users.json", params: like_params
    expect(response.status).to eq(200)
    expect(
      response.parsed_body["post_action_users"].map { |user| user["id"] }
    ).to include(viewer.id)
    reaction_paths.each do |path|
      get path
      expect(response.status).to eq(200)
    end
  end
end
