import { apiInitializer } from "discourse/lib/api";

export default apiInitializer((api) => {
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
    }
  );
});
