import Component from "@glimmer/component";
import { tracked } from "@glimmer/tracking";
import { on } from "@ember/modifier";
import { action } from "@ember/object";
import didInsert from "@ember/render-modifiers/modifiers/did-insert";
import { service } from "@ember/service";
import TopicList from "discourse/components/topic-list/list";
import getURL from "discourse/lib/get-url";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import { i18n } from "discourse-i18n";
import ZgieDiscoveryBanner from "./zgie-discovery-banner";
import ZgieHotTopics from "./zgie-hot-topics";

export default class ZgieHomepage extends Component {
  @service store;
  @service site;
  @service siteSettings;
  @service currentUser;

  @tracked topics = [];
  @tracked loading = true;
  @tracked failed = false;

  allPostsURL = getURL("/latest");

  get enabled() {
    return (
      this.site.desktopView &&
      this.siteSettings.zgie_bbs_branded_layout &&
      (this.currentUser || !this.siteSettings.login_required)
    );
  }

  @action
  async load() {
    this.loading = true;
    this.failed = false;
    try {
      const list = await this.store.findFiltered("topicList", {
        filter: "latest",
        params: { per_page: 20 },
      });
      if (!this.isDestroying && !this.isDestroyed) {
        // A snapshot only: no load-more observer or live topic-list subscription.
        this.topics = list.topics.slice(0, 20);
      }
    } catch {
      if (!this.isDestroying && !this.isDestroyed) {
        this.failed = true;
      }
    } finally {
      if (!this.isDestroying && !this.isDestroyed) {
        this.loading = false;
      }
    }
  }

  <template>
    {{#if this.enabled}}
      <div class="zgie-home" {{didInsert this.load}}>
        <ZgieDiscoveryBanner />
        <h2 class="zgie-home__heading">{{i18n
            "zgie_bbs.campus.latest_posts"
          }}</h2>
        <div class="zgie-home__grid">
          <section
            class="zgie-home__feed"
            aria-label={{i18n "zgie_bbs.campus.latest_posts"}}
          >
            {{#if this.loading}}
              <p class="zgie-home__status" role="status">{{i18n "loading"}}</p>
            {{else if this.failed}}
              <p class="zgie-home__status" role="status">{{i18n
                  "zgie_bbs.campus.feed_error"
                }}</p>
              <button
                class="btn btn-default"
                type="button"
                {{on "click" this.load}}
              >{{i18n "zgie_bbs.campus.retry"}}</button>
            {{else if this.topics.length}}
              <TopicList
                @topics={{this.topics}}
                @showPosters={{true}}
                @canBulkSelect={{false}}
                @expandAllPinned={{false}}
                @expandGloballyPinned={{false}}
                @listContext="zgie-home"
              />
            {{else}}
              <p class="zgie-home__status">{{i18n "zgie_bbs.campus.empty"}}</p>
            {{/if}}
            <a class="zgie-home__all" href={{this.allPostsURL}}>{{i18n
                "zgie_bbs.campus.show_all"
              }}
              {{dIcon "arrow-right"}}</a>
          </section>
          <ZgieHotTopics />
        </div>
      </div>
    {{else}}
      <a href={{this.allPostsURL}}>{{i18n "zgie_bbs.campus.posts"}}</a>
    {{/if}}
  </template>
}
