import { apiInitializer } from "discourse/lib/api";

export default apiInitializer((api) => {
  const colors = api.container.lookup("service:interface-color");
  colors.ensureCorrectMode();
  // A first-time visitor may have neither a cookie nor an explicit preference.
  // Keep the native selector's accessible label valid before the first choice.
  if (colors.selectorAvailable && !colors.colorMode) {
    colors.colorMode = "auto";
  }
});
