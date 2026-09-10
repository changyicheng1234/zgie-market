import Component from "@glimmer/component";
import { service } from "@ember/service";
import getURL from "discourse/lib/get-url";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import { i18n } from "discourse-i18n";
import { campusCategories } from "../lib/zgie-campus-categories";

export default class ZgieDiscoveryBanner extends Component {
  @service currentUser;
  @service site;
  @service siteSettings;
  @service router;

  searchURL = getURL("/search");

  categoriesURL = getURL("/categories");

  get enabled() {
    return (
      this.siteSettings.zgie_bbs_branded_layout &&
      (!this.site.desktopView ||
        this.router.currentRouteName === "discovery.custom") &&
      (this.currentUser || !this.siteSettings.login_required)
    );
  }

  get isHomepage() {
    if (this.site.desktopView) {
      return this.router.currentRouteName === "discovery.custom";
    }
    return ["discovery.latest", "discovery.categories"].includes(
      this.router.currentRouteName
    );
  }

  get title() {
    if (this.args.outletArgs?.category) {
      return this.args.outletArgs.category.name;
    }
    if (this.args.outletArgs?.tag) {
      return `#${this.args.outletArgs.tag.name}`;
    }
    if (this.isHomepage) {
      return i18n("zgie_bbs.banner.welcome", { site: this.siteSettings.title });
    }
    return i18n("zgie_bbs.banner.discussions");
  }

  get categories() {
    if (this.site.desktopView) {
      return campusCategories(this.site)
        .slice(0, 3)
        .map((category) => ({
          name: category.name,
          href: category.url,
          icon: "compass",
        }));
    }
    const available = this.site.categories.filter(
      (category) =>
        !category.parent_category_id &&
        category.id !== this.site.uncategorized_category_id
    );
    const preferred = ["daily", "study", "help"]
      .map((slug) => available.find((category) => category.slug === slug))
      .filter(Boolean);
    return [
      ...preferred,
      ...available.filter((category) => !preferred.includes(category)),
    ]
      .slice(0, 3)
      .map((category, index) => ({
        name: category.name,
        href: category.url,
        icon: ["comments", "book-open", "compass"][index],
      }));
  }

  <template>
    {{#if this.enabled}}
      <section
        class={{if
          this.isHomepage
          "zgie-banner"
          "zgie-banner zgie-banner--compact"
        }}
        aria-label={{this.title}}
      >
        <div class="zgie-banner__intro">
          <p class="zgie-banner__eyebrow">{{i18n "zgie_bbs.banner.eyebrow"}}</p>
          <h1>{{this.title}}</h1>
          {{#if this.isHomepage}}
            <p class="zgie-banner__description">{{i18n
                "zgie_bbs.banner.description"
              }}</p>
          {{/if}}
        </div>
        <form
          class="zgie-banner__search"
          role="search"
          action={{this.searchURL}}
          method="get"
        >
          <label class="sr-only" for="zgie-community-search">{{i18n
              "zgie_bbs.banner.search"
            }}</label>
          {{dIcon "magnifying-glass"}}
          <input
            id="zgie-community-search"
            name="q"
            type="search"
            placeholder={{i18n "zgie_bbs.banner.search"}}
          />
          <button
            type="submit"
            aria-label={{i18n "zgie_bbs.banner.submit"}}
          >{{dIcon "arrow-right"}}</button>
        </form>
        {{#if this.isHomepage}}
          <div class="zgie-banner__categories">
            {{#each this.categories as |category|}}
              <a class="zgie-banner__category" href={{category.href}}>
                <span class="zgie-banner__category-icon">{{dIcon
                    category.icon
                  }}</span>
                <span>{{category.name}}</span>
                <span class="zgie-banner__arrow" aria-hidden="true">{{dIcon
                    "arrow-right"
                  }}</span>
              </a>
            {{/each}}
          </div>
          <a class="zgie-banner__all" href={{this.categoriesURL}}>{{i18n
              "zgie_bbs.banner.all_categories"
            }}
            <span aria-hidden="true">→</span></a>
        {{/if}}
      </section>
    {{/if}}
  </template>
}
