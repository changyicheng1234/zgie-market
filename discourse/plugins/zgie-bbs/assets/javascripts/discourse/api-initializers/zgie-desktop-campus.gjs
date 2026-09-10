import { service } from "@ember/service";
import { apiInitializer } from "discourse/lib/api";
import { i18n } from "discourse-i18n";
import ZgieHotTopics from "../components/zgie-hot-topics";
import { campusCategories } from "../lib/zgie-campus-categories";

export default apiInitializer((api) => {
  api.renderInOutlet("discovery-above", ZgieHotTopics);
  api.addSidebarSection((BaseSection, BaseLink) => {
    class CampusLink extends BaseLink {
      prefixType = "icon";

      prefixValue = "folder";

      constructor(category) {
        super();
        this.category = category;
      }

      get name() {
        return `zgie-campus-${this.category.slug}`;
      }

      get text() {
        return this.category.name;
      }

      get title() {
        return this.category.name;
      }

      get href() {
        return this.category.url;
      }
    }
    return class CampusSection extends BaseSection {
      @service site;
      @service siteSettings;

      name = "zgie-campus";
      collapsedByDefault = false;

      get text() {
        return i18n("zgie_bbs.campus.categories");
      }

      get title() {
        return this.text;
      }

      get displaySection() {
        return (
          this.siteSettings.zgie_bbs_branded_layout && this.site.desktopView
        );
      }

      get links() {
        return campusCategories(this.site).map(
          (category) => new CampusLink(category)
        );
      }
    };
  });
});
