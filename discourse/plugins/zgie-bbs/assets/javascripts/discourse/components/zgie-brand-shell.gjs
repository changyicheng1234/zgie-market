import Component from "@glimmer/component";
import { service } from "@ember/service";
import bodyClass from "discourse/helpers/body-class";
import EmbedMode from "discourse/lib/embed-mode";

export default class ZgieBrandShell extends Component {
  @service router;
  @service siteSettings;

  get enabled() {
    return (
      this.siteSettings.zgie_bbs_branded_layout &&
      !EmbedMode.enabled &&
      !this.router.currentRouteName?.startsWith("admin")
    );
  }

  <template>{{bodyClass (if this.enabled "zgie-branded")}}</template>
}
