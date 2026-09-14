# frozen_string_literal: true

module ZgieBbs
  class CampusTopicQuery < ::TopicQuery
    def hot_topics
      # Ranks across every category the user can see, not just the campus
      # ones -- CAMPUS_SLUGS still drives the desktop sidebar/banner
      # shortcuts elsewhere, but the trending list itself is site-wide.
      # TopicQuery applies Guardian's category permissions and excludes PMs,
      # deleted topics and category definitions. Rank before limiting to five.
      default_results(
        limit: false,
        skip_ordering: true,
        no_definitions: true,
        visible: true
      )
        # `no_definitions` filters on `categories.topic_id`, but nothing
        # upstream forces that table into the query as a real SQL join --
        # `Topic.includes(:category)` alone lets Rails silently fall back to
        # a separate preload query, so the raw `no_definitions` WHERE clause
        # blows up with "missing FROM-clause entry for table categories".
        # An explicit join guarantees it's actually there.
        .joins(:category)
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
