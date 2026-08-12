import { action } from "@ember/object";
import { apiInitializer } from "discourse/lib/api";
import ZgieAnonymousAvatar from "../components/zgie-anonymous-avatar";
import ZgieAnonymousComposerToggle from "../components/zgie-anonymous-composer-toggle";
import ZgieAnonymousName from "../components/zgie-anonymous-name";
import ZgieTopicPosition from "../components/zgie-topic-position";

function suppressPopup(event) {
  event?.preventDefault();
  event?.stopPropagation();
  event?.stopImmediatePropagation();
}

export default apiInitializer((api) => {
  api.serializeOnCreate("zgie_anonymous", "zgieAnonymous");
  api.serializeToDraft("zgie_anonymous", "zgieAnonymous");
  api.addTrackedPostProperties("zgie_anonymous_author");
  api.renderInOutlet("post-avatar", ZgieAnonymousAvatar);
  api.renderInOutlet("post-meta-data-poster-name-user-link", ZgieAnonymousName);
  api.renderInOutlet("topic-above-post-stream", ZgieTopicPosition);
  api.renderInOutlet(
    "composer-fields-below",
    ZgieAnonymousComposerToggle
  );

  api.modifyClass(
    "component:user-menu/profile-tab-content",
    (Superclass) =>
      class extends Superclass {
        get showToggleAnonymousButton() {
          return false;
        }
      }
  );

  api.modifyClass(
    "component:post/menu/liked-users-list",
    (Superclass) =>
      class extends Superclass {
        @action
        togglePopup(event) {
          suppressPopup(event);
        }
      }
  );

  api.modifyClass(
    "component:discourse-reactions-counter",
    (Superclass) =>
      class extends Superclass {
        @action
        click(event) {
          suppressPopup(event);
        }

        @action
        keyDown(event) {
          if (event.key === "Enter" || event.key === " ") {
            suppressPopup(event);
          }
        }
      },
    { ignoreMissing: true }
  );
});
