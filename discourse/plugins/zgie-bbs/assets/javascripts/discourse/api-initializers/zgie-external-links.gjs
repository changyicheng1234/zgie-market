import { later } from "@ember/runloop";
import { apiInitializer } from "discourse/lib/api";
import { customPanels } from "discourse/lib/sidebar/custom-sections";

// Chat plugin's sidebar panel key (plugins/chat/.../lib/init-sidebar-state.js:
// `export const CHAT_PANEL = "chat"`). Registering here places our link right
// after chat's own last section for this panel ("开始新的DM" / chat-dms).
//
// The chat plugin creates this panel asynchronously (after loading the
// current user's chat preload data), so it may not exist yet when our own
// api-initializer runs. We poll briefly until it appears.
const CHAT_PANEL = "chat";
const POLL_INTERVAL_MS = 150;
const MAX_ATTEMPTS = 40; // ~6s upper bound

function chatPanelExists() {
  return customPanels?.some((panel) => panel.key === CHAT_PANEL);
}

function registerOpenISESection(api) {
  api.addSidebarSection(
    (BaseCustomSidebarSection, BaseCustomSidebarSectionLink) => {
      const OpenISESectionLink = class extends BaseCustomSidebarSectionLink {
        name = "openise";
        title = "OpenISE";
        text = "OpenISE";
        href = "https://openise.pages.dev";
        prefixType = "icon";
        prefixValue = "up-right-from-square";
      };

      return class OpenISESection extends BaseCustomSidebarSection {
        name = "openise-links";
        hideSectionHeader = true;
        title = "";

        get text() {
          return null;
        }

        get links() {
          return [new OpenISESectionLink()];
        }
      };
    },
    CHAT_PANEL
  );

  // Pushing into `panel.sections` (a plain array) doesn't invalidate the
  // sidebar's tracked getters on its own. Reassigning `sidebarState.mode`
  // to itself dirties the tracked field it's read from, forcing the
  // sidebar to recompute its section list and pick up the new one.
  const sidebarState = api.container.lookup("service:sidebar-state");
  if (sidebarState) {
    // eslint-disable-next-line no-self-assign
    sidebarState.mode = sidebarState.mode;
  }
}

function waitForChatPanelThenRegister(api, attempt = 0) {
  if (chatPanelExists()) {
    registerOpenISESection(api);
    return;
  }

  if (attempt >= MAX_ATTEMPTS) {
    return;
  }

  later(() => waitForChatPanelThenRegister(api, attempt + 1), POLL_INTERVAL_MS);
}

export default apiInitializer((api) => {
  waitForChatPanelThenRegister(api);
});
