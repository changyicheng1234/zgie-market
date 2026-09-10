# frozen_string_literal: true

module ZgieBbs
  class Engine < ::Rails::Engine
    engine_name PLUGIN_NAME
    isolate_namespace ZgieBbs
  end
end
