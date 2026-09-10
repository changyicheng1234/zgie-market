# frozen_string_literal: true

module ZgieBbs
  class CampusTopicQuery < ::TopicQuery
    def hot_topics
      # TopicQuery applies Guardian's category permissions and excludes PMs,
      # deleted topics and category definitions. Rank before limiting to five.
      default_results(
        limit: false,
        skip_ordering: true,
        no_definitions: true,
        visible: true
      )
        .where(
          category_id:
            Category.where(slug: SiteConfigurator::CAMPUS_SLUGS).select(:id)
        )
        .joins(
          "LEFT JOIN topic_hot_scores ON topic_hot_scores.topic_id = topics.id"
        )
        .reorder(
          Arel.sql(
            "COALESCE(topic_hot_scores.score, 0) DESC, topics.bumped_at DESC, topics.id DESC"
          )
        )
        .limit(5)
    end
  end
end
