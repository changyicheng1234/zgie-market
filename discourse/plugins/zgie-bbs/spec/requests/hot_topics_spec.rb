# frozen_string_literal: true

RSpec.describe "Campus trending topics" do
  fab!(:member, :user)
  fab!(:category) { Fabricate(:category, slug: "courses") }

  before do |example|
    SiteSetting.zgie_bbs_enabled = true
    SiteSetting.zgie_bbs_branded_layout = true
    sign_in(member) unless example.metadata[:anonymous]
  end

  it "ranks campus topics by hot score before limiting to five" do
    topics =
      7.times.map do |index|
        topic = Fabricate(:topic, category: category, posts_count: 3)
        TopicHotScore.create!(topic_id: topic.id, score: index + 1)
        topic
      end
    other = Fabricate(:topic)
    TopicHotScore.create!(topic_id: other.id, score: 999)

    get "/zgie-bbs/hot-topics.json"
    expect(response.status).to eq(200)
    result = response.parsed_body.fetch("topics")
    expect(result.pluck("id")).to eq(topics.reverse.first(5).map(&:id))
    expect(result.first).to include(
      "url" => topics.last.relative_url,
      "reply_count" => 2
    )
  end

  it "excludes private categories, PMs, deleted and unlisted topics, and category definitions" do
    private_category =
      Fabricate(
        :private_category,
        slug: "internships",
        group: Fabricate(:group)
      )
    hidden = Fabricate(:topic, category: private_category)
    pm = Fabricate(:private_message_topic, category: category)
    deleted = Fabricate(:topic, category: category, deleted_at: Time.current)
    unlisted = Fabricate(:topic, category: category, visible: false)
    definition = Fabricate(:topic, category: category)
    category.update!(topic_id: definition.id)
    [hidden, pm, deleted, unlisted, definition].each do |topic|
      TopicHotScore.create!(topic_id: topic.id, score: 999)
    end
    allowed = Fabricate(:topic, category: category)

    get "/zgie-bbs/hot-topics.json"
    expect(response.status).to eq(200)
    expect(response.parsed_body.fetch("topics").pluck("id")).to eq([allowed.id])
  end

  it "requires authentication when the site requires login", anonymous: true do
    SiteSetting.login_required = true
    get "/zgie-bbs/hot-topics.json"
    expect(response.status).to eq(403)
  end

  it "returns an empty list when no campus posts exist" do
    get "/zgie-bbs/hot-topics.json"
    expect(response.status).to eq(200)
    expect(response.parsed_body.fetch("topics")).to eq([])
  end
end
