import Component from "@glimmer/component";
import { on } from "@ember/modifier";
import { action } from "@ember/object";
import { CREATE_TOPIC, REPLY } from "discourse/models/composer";
import { i18n } from "discourse-i18n";

export default class ZgieAnonymousComposerToggle extends Component {
  static shouldRender(args) {
    return (
      !args.model?.privateMessage &&
      [CREATE_TOPIC, REPLY].includes(args.model?.action)
    );
  }

  get model() {
    return this.args.outletArgs.model;
  }

  get enabled() {
    return this.model.get("zgieAnonymous") === true;
  }

  get label() {
    const key =
      this.model.action === REPLY
        ? "zgie_bbs.composer.anonymous_reply"
        : "zgie_bbs.composer.anonymous_topic";

    return i18n(key);
  }

  @action
  toggle() {
    this.model.set("zgieAnonymous", !this.enabled);
  }

  <template>
    <div class="zgie-anonymous-composer-toggle">
      <button
        class="zgie-anonymous-composer-toggle__control"
        type="button"
        role="switch"
        aria-checked={{if this.enabled "true" "false"}}
        aria-label={{this.label}}
        aria-describedby="zgie-anonymous-composer-toggle-hint"
        data-zgie-anonymous-toggle
        {{on "click" this.toggle}}
      >
        <span
          class="zgie-anonymous-composer-toggle__slider"
          aria-hidden="true"
        ></span>
        <span class="zgie-anonymous-composer-toggle__label">
          {{this.label}}
        </span>
      </button>
      <span
        id="zgie-anonymous-composer-toggle-hint"
        class="zgie-anonymous-composer-toggle__hint"
      >
        {{i18n "zgie_bbs.composer.anonymous_hint"}}
      </span>
    </div>
  </template>
}
