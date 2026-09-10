import Component from "@glimmer/component";
import { getOwner } from "@ember/owner";
import { service } from "@ember/service";
import bodyClass from "discourse/helpers/body-class";
import { defaultHomepage } from "discourse/lib/utilities";
import dConcatClass from "discourse/ui-kit/helpers/d-concat-class";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import { i18n } from "discourse-i18n";

export default class ZgieMobileNavigation extends Component {
  @service currentUser;
  @service router;
  @service siteSettings;

  get enabled() {
    return (
      this.siteSettings.zgie_bbs_mobile_navigation &&
      this.currentUser &&
      !this.router.currentRouteName?.startsWith("admin")
    );
  }

  get bodyClasses() {
    return this.enabled ? "zgie-mobile-navigation-enabled" : "";
  }

  get chat() {
    // Chat is optional; importing its components would break sites without it.
    return getOwner(this).lookup("service:chat");
  }

  get items() {
    const route = this.router.currentRouteName || "";
    const path = (this.router.currentURL || "").split("?")[0];
    const userPath = this.currentUser.path;
    const messagesActive =
      route.startsWith("chat") ||
      path.toLowerCase().startsWith(`${userPath.toLowerCase()}/messages`);
    const personalActive =
      !messagesActive &&
      path.toLowerCase().startsWith(`${userPath.toLowerCase()}/`);
    const hotActive = route === "discovery.hot" || route === "discovery.top";
    const chatEnabled =
      this.chat?.userCanChat && this.chat?.userCanAccessDirectMessages;
    const chatUnread = chatEnabled
      ? getOwner(this)
          .lookup("service:chat-channels-manager")
          ?.directMessageChannels.reduce(
            (count, channel) => count + channel.tracking.unreadCount,
            0
          )
      : 0;
    const messageCount =
      chatUnread || this.currentUser.new_personal_messages_notifications_count;

    return [
      {
        key: "home",
        icon: "house",
        href: `/${defaultHomepage()}`,
        active: route.startsWith("discovery.") && !hotActive,
      },
      { key: "hot", icon: "fire", href: "/hot", active: hotActive },
      {
        key: "messages",
        icon: "envelope",
        href: chatEnabled ? "/chat/direct-messages" : `${userPath}/messages`,
        active: messagesActive,
        unread: messageCount,
      },
      {
        key: "profile",
        icon: "user",
        href: `${userPath}/summary`,
        active: personalActive,
        unread: this.currentUser.all_unread_notifications_count,
      },
    ].map((item) => ({
      ...item,
      label: i18n(`zgie_bbs.navigation.${item.key}`),
    }));
  }

  <template>
    {{bodyClass this.bodyClasses}}
    {{#if this.enabled}}
      <nav
        class="zgie-mobile-nav"
        aria-label={{i18n "zgie_bbs.navigation.label"}}
      >
        {{#each this.items as |item|}}
          <a
            class={{dConcatClass
              "zgie-mobile-nav__item"
              (if item.active "is-active")
            }}
            href={{item.href}}
            aria-current={{if item.active "page"}}
          >
            <span class="zgie-mobile-nav__icon">
              {{dIcon item.icon}}
              {{#if item.unread}}
                <span
                  class="zgie-mobile-nav__unread"
                  aria-label={{i18n "zgie_bbs.navigation.unread"}}
                ></span>
              {{/if}}
            </span>
            <span>{{item.label}}</span>
          </a>
        {{/each}}
      </nav>
    {{/if}}
  </template>
}
