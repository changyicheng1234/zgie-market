import Component from "@glimmer/component";
import { tracked } from "@glimmer/tracking";
import { on } from "@ember/modifier";
import { action } from "@ember/object";
import didInsert from "@ember/render-modifiers/modifiers/did-insert";
import { service } from "@ember/service";
import { ajax } from "discourse/lib/ajax";
import getURL from "discourse/lib/get-url";
import dCategoryBadge from "discourse/ui-kit/helpers/d-category-badge";
import dFormatDate from "discourse/ui-kit/helpers/d-format-date";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import { i18n } from "discourse-i18n";

export default class ZgieHotTopics extends Component {
  @service site;
  @service currentUser;
  @service siteSettings;

  @tracked topics = [];
  @tracked loading = true;
  @tracked failed = false;
  allTopicsURL = getURL("/hot");

  get enabled() {
    return (
      this.siteSettings.zgie_bbs_branded_layout &&
      this.site.desktopView &&
      (this.currentUser || !this.siteSettings.login_required)
    );
  }

  @action
  async load() {
    this.loading = true;
    this.failed = false;
    try {
      const result = await ajax("/zgie-bbs/hot-topics.json");
      if (!this.isDestroying && !this.isDestroyed) {
        this.topics = result.topics.map((topic) => ({
          ...topic,
          href: getURL(topic.url),
          category: this.site.categories.find(
            (category) => category.id === topic.category_id
          ),
        }));
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
      <aside
        class="zgie-hot-topics"
        aria-label={{i18n "zgie_bbs.campus.hot"}}
        {{didInsert this.load}}
      >
        <h2>{{dIcon "fire"}} {{i18n "zgie_bbs.campus.hot"}}</h2>
        {{#if this.loading}}
          <p role="status">{{i18n "zgie_bbs.campus.loading"}}</p>
        {{else if this.failed}}
          <p role="status">{{i18n "zgie_bbs.campus.error"}}</p>
          <button
            class="btn btn-default"
            type="button"
            {{on "click" this.load}}
          >{{i18n "zgie_bbs.campus.retry"}}</button>
        {{else}}
          <ol class="zgie-hot-topics__list">
            {{#each this.topics key="id" as |topic|}}
              <li>
                <a
                  class="zgie-hot-topics__title"
                  href={{topic.href}}
                >{{topic.title}}</a>
                <div class="zgie-hot-topics__meta">
                  {{#if topic.category}}
                    {{dCategoryBadge topic.category link=true hideParent=true}}
                  {{/if}}
                  <span class="zgie-hot-topics__age">{{dFormatDate
                      topic.bumped_at
                      format="tiny"
                    }}</span>
                </div>
              </li>
            {{else}}
              <li class="zgie-hot-topics__empty">{{i18n
                  "zgie_bbs.campus.empty"
                }}</li>
            {{/each}}
          </ol>
        {{/if}}
        <a class="zgie-hot-topics__all" href={{this.allTopicsURL}}>{{i18n
            "zgie_bbs.campus.show_all"
          }}
          {{dIcon "arrow-right"}}</a>
      </aside>
    {{/if}}
  </template>
}
