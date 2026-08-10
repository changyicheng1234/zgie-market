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
  end

  it "switches any signed-in member to a shadow account that can create topics and replies" do
    sign_in(member)

    post "/u/toggle-anon.json"

    expect(response.status).to eq(200)
    anonymous_user = User.find(session[:current_user_id])
    expect(anonymous_user).to be_anonymous
    # Automatic groups are populated out of band in the test environment.
    anonymous_user.send(:trigger_user_automatic_group_refresh)
    anonymous_user.reload
    category = Category.find_by!(slug: "study")
    anonymous_guardian = Guardian.new(anonymous_user)

    expect(anonymous_guardian.can_create_post?(nil)).to eq(true)
    expect(anonymous_user.belonging_to_group_ids).to include(
      Group::AUTO_GROUPS[:trust_level_0]
    )
    expect(
      Category.topic_create_allowed(anonymous_guardian).pluck(:id)
    ).to include(category.id)
    expect(anonymous_guardian.can_create_topic?(nil)).to eq(true)
    expect(anonymous_guardian.can_create_topic_on_category?(category)).to eq(
      true
    )

    topic_post =
      PostCreator.create!(
        anonymous_user,
        title: "这是一条用于验证匿名发帖功能的测试主题",
        raw: "这是一段用于确认匿名主题能够成功发布的测试正文内容。",
        category: category.id
      )
    reply =
      PostCreator.create!(
        anonymous_user,
        topic_id: topic_post.topic_id,
        raw: "这是一段用于确认匿名留言能够成功发布的测试回复内容。"
      )

    expect(topic_post.user).to eq(anonymous_user)
    expect(reply.user).to eq(anonymous_user)
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
