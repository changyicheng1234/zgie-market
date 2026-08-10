import Component from "@glimmer/component";
import dAvatar from "discourse/ui-kit/helpers/d-avatar";
import { i18n } from "discourse-i18n";

export default class ZgieAnonymousAvatar extends Component {
  static shouldRender(args) {
    return args.post?.zgie_anonymous_author === true;
  }

  <template>
    <span
      class="zgie-anonymous-avatar"
      aria-label={{i18n "zgie_bbs.anonymous_author"}}
      title={{i18n "zgie_bbs.anonymous_author"}}
    >
      {{dAvatar
        @outletArgs.user
        extraClasses="main-avatar zgie-anonymous-avatar__image"
        imageSize=@outletArgs.size
        hideTitle=true
        loading="lazy"
      }}
    </span>
  </template>
}
