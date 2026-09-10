import { apiInitializer } from "discourse/lib/api";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import ZgieMobileNavigation from "../components/zgie-mobile-navigation";

export default apiInitializer((api) => {
  api.renderInOutlet("above-footer", ZgieMobileNavigation);
  api.renderInOutlet(
    "user-dropdown-notifications__before",
    <template>
      <span class="zgie-mobile-notifications-icon">{{dIcon "bell"}}</span>
    </template>
  );
});
