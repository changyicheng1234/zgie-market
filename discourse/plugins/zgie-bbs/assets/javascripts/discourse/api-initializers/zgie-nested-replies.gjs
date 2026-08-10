import Component from "@glimmer/component";
import { action } from "@ember/object";
import { apiInitializer } from "discourse/lib/api";
import DUserLink from "discourse/ui-kit/d-user-link";
import { i18n } from "discourse-i18n";

class ZgieReplyTarget extends Component {
  static shouldRender(args) {
    return (
      args.post.zgie_replies_to_nested_reply === true &&
      Boolean(args.post.reply_to_user?.username)
    );
  }

  <template>
    <span
      class="zgie-reply-target"
      title={{i18n
        "zgie_bbs.nested_reply_to"
        username=@outletArgs.post.reply_to_user.username
      }}
    >
      <span class="zgie-reply-target__separator" aria-hidden="true">&gt;</span>
      <DUserLink @user={{@outletArgs.post.reply_to_user}}>
        {{@outletArgs.post.reply_to_user.username}}
      </DUserLink>
    </span>
  </template>
}

export default apiInitializer((api) => {
  api.modifyClass(
    "component:nested/post",
    (Superclass) =>
      class extends Superclass {
        get zgieIsFirstLevelReply() {
          return (this.args.post.reply_to_post_number || 0) <= 1;
        }

        get effectiveExpanded() {
          return this.zgieIsFirstLevelReply ? true : super.effectiveExpanded;
        }

        get effectiveCollapsed() {
          return this.zgieIsFirstLevelReply ? false : super.effectiveCollapsed;
        }

        get showDepthLine() {
          if (this.zgieIsFirstLevelReply) {
            return false;
          }

          return this.atMaxDepth || super.showDepthLine;
        }

        get showDepthLineIcon() {
          return this.atMaxDepth ? true : super.showDepthLineIcon;
        }

        get depthLineCollapsed() {
          return this.atMaxDepth ? false : super.depthLineCollapsed;
        }

        @action
        toggleExpanded() {
          if (this.zgieIsFirstLevelReply) {
            return;
          }

          if (this.atMaxDepth) {
            this.collapsed = !this.collapsed;
            this.lineHighlighted = false;
            this.args.expansionState?.set(this.args.post.post_number, {
              expanded: this.expanded,
              collapsed: this.collapsed,
            });
            return;
          }

          return super.toggleExpanded();
        }

        @action
        handleDepthLine() {
          if (this.zgieIsFirstLevelReply) {
            return;
          }

          if (this.atMaxDepth) {
            this.toggleExpanded();
            return;
          }

          return super.handleDepthLine();
        }

        collapsePost() {
          if (this.zgieIsFirstLevelReply) {
            return;
          }

          return super.collapsePost();
        }
      }
  );

  api.renderAfterWrapperOutlet(
    "post-meta-data-poster-name",
    ZgieReplyTarget
  );
});
