# frozen_string_literal: true

require "net/http"

module ZgieBbs
  # 新主题发布时往企业微信群机器人推送一条带跳转链接的消息。
  # 纯提醒作用：webhook 未配置时静默跳过，发送失败也只记日志，不影响发帖主流程。
  module WecomNotify
    def self.enabled?
      webhook_url.present?
    end

    def self.webhook_url
      ENV["ZGIE_WECOM_WEBHOOK_URL"].to_s.strip
    end

    def self.notify_new_topic(post)
      return unless SiteSetting.zgie_bbs_enabled
      return unless post.post_number == 1
      return if post.topic.blank? || post.topic.archetype == Archetype.private_message
      return unless enabled?

      title = post.topic.title
      title = "[匿名] #{title}" if post.user&.anonymous?

      category_name = post.topic.category&.name || "未分类"
      topic_url = "#{Discourse.base_url}/t/#{post.topic.slug}/#{post.topic.id}"

      content =
        "📢 **新帖子发布**\n> 分区：#{category_name}\n> [#{title}](#{topic_url})"

      send_markdown(content)
    end

    def self.send_markdown(content)
      body = { msgtype: "markdown", markdown: { content: content } }.to_json

      uri = URI.parse(webhook_url)
      http = Net::HTTP.new(uri.host, uri.port)
      http.use_ssl = uri.scheme == "https"
      http.open_timeout = 5
      http.read_timeout = 5

      request = Net::HTTP::Post.new(uri.request_uri, "Content-Type" => "application/json; charset=utf-8")
      request.body = body

      response = http.request(request)
      unless response.is_a?(Net::HTTPSuccess)
        Rails.logger.warn("[zgie-bbs] 企业微信通知接口返回非成功状态: #{response.code}")
      end
    rescue => e
      Rails.logger.warn("[zgie-bbs] 发送企业微信新帖子通知失败: #{e.message}")
    end
  end
end
