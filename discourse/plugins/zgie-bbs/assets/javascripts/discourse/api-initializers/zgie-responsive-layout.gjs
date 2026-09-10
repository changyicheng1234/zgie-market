import { apiInitializer } from "discourse/lib/api";
import dIcon from "discourse/ui-kit/helpers/d-icon";
import ZgieBrandShell from "../components/zgie-brand-shell";
import ZgieDiscoveryBanner from "../components/zgie-discovery-banner";
import ZgieMobileNavigation from "../components/zgie-mobile-navigation";

export default apiInitializer((api) => {
  api.renderInOutlet("above-footer", ZgieBrandShell);
  api.renderInOutlet("discovery-list-controls-above", ZgieDiscoveryBanner);
  api.registerValueTransformer(
    "welcome-banner-display-for-route",
    ({ value }) =>
      api.container.lookup("service:site-settings").zgie_bbs_branded_layout
        ? false
        : value
  );
  api.renderInOutlet("above-footer", ZgieMobileNavigation);
  api.renderInOutlet(
    "user-dropdown-notifications__before",
    <template>
      <span class="zgie-mobile-notifications-icon">{{dIcon "bell"}}</span>
    </template>
  );
});
