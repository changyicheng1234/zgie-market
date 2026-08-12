# frozen_string_literal: true

module PageObjects
  module Components
    class ZgieAnonymousComposerToggle < PageObjects::Components::Base
      SELECTOR = ".zgie-anonymous-composer-toggle"

      def has_label?(text)
        has_css?(SELECTOR, text: text)
      end

      def has_state?(enabled)
        has_css?(
          "#{SELECTOR} [data-zgie-anonymous-toggle]" \
            "[aria-checked='#{enabled}']"
        )
      end

      def click_visible_label
        find("#{SELECTOR}__label").click
      end
    end
  end
end
