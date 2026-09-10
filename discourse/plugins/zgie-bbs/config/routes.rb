# frozen_string_literal: true

ZgieBbs::Engine.routes.draw { get "/hot-topics" => "hot_topics#index" }
