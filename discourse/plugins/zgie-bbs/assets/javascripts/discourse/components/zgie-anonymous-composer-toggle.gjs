import Component from "@glimmer/component";
import { on } from "@ember/modifier";
import { action } from "@ember/object";
import { CREATE_TOPIC, REPLY } from "discourse/models/composer";
import DToggleSwitch from "discourse/ui-kit/d-toggle-switch";
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
    return this.model.zgieAnonymous === true;
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
      <DToggleSwitch
        @state={{this.enabled}}
        @translatedLabel={{this.label}}
        aria-label={{this.label}}
        data-zgie-anonymous-toggle
        {{on "click" this.toggle}}
      />
      <span class="zgie-anonymous-composer-toggle__hint">
        {{i18n "zgie_bbs.composer.anonymous_hint"}}
      </span>
    </div>
  </template>
}
