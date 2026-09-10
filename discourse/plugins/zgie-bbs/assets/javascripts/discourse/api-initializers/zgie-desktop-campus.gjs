import { service } from "@ember/service";
import { apiInitializer } from "discourse/lib/api";
import getURL from "discourse/lib/get-url";
import { i18n } from "discourse-i18n";
import ZgieHomepage from "../components/zgie-homepage";
import { campusCategories } from "../lib/zgie-campus-categories";

export default apiInitializer((api) => {
  api.renderInOutlet("custom-homepage", ZgieHomepage);
  const desktopHomeEnabled = () => {
    const site = api.container.lookup("service:site");
    const settings = api.container.lookup("service:site-settings");
    return (
      site.desktopView &&
      settings.zgie_bbs_branded_layout &&
      (api.getCurrentUser() || !settings.login_required)
    );
  };
  api.registerValueTransformer("home-logo-href", ({ value }) =>
    desktopHomeEnabled() ? getURL("/custom") : value
  );
  api.onPageChange((url) => {
    if (url === getURL("/") && desktopHomeEnabled()) {
      api.container.lookup("service:router").replaceWith("discovery.custom");
    }
  });
  api.addSidebarSection((BaseSection, BaseLink) => {
    class NavigationLink extends BaseLink {
      name;
      text;
      title;
      href;
      prefixValue;

      prefixType = "icon";

      constructor(name, text, href, icon) {
        super();
        this.name = name;
        this.text = text;
        this.title = text;
        this.href = href;
        this.prefixValue = icon;
      }
    }
    return class NavigationSection extends BaseSection {
      @service currentUser;
      @service site;
      @service siteSettings;

      name = "zgie-navigation";
      hideSectionHeader = true;
      text = "";
      title = "";

      get displaySection() {
        return (
          this.site.desktopView && this.siteSettings.zgie_bbs_branded_layout
        );
      }

      get links() {
        const links = [
          new NavigationLink(
            "zgie-home",
            i18n("zgie_bbs.campus.home"),
            getURL("/custom"),
            "house"
          ),
          new NavigationLink(
            "zgie-posts",
            i18n("zgie_bbs.campus.posts"),
            getURL("/latest"),
            "layer-group"
          ),
        ];
        if (this.currentUser) {
          links.push(
            new NavigationLink(
              "zgie-my-posts",
              i18n("zgie_bbs.campus.my_posts"),
              `${this.currentUser.path}/activity/topics`,
              "user"
            ),
            new NavigationLink(
              "zgie-my-messages",
              i18n("zgie_bbs.campus.my_messages"),
              `${this.currentUser.path}/messages`,
              "envelope"
            )
          );
        }
        return links;
      }
    };
  });
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
