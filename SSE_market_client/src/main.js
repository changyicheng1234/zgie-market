import Vue from 'vue';
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue';
import Vuelidate from 'vuelidate';
import axios from 'axios';
import ElementUI from 'element-ui';
import VueParticles from 'vue-particles';
import 'element-ui/lib/theme-chalk/index.css';
import VueAxios from 'vue-axios';
import VueScrollTo from 'vue-scrollto';
import VueTypedJs from 'vue-typed-js';
import { Picker } from 'emoji-mart-vue';
import mouse from '@/utils/mouseClick';
import VueMarkdownEditor from '@kangc/v-md-editor';
import '@kangc/v-md-editor/lib/style/base-editor.css';
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress';
import '@kangc/v-md-editor/lib/theme/style/vuepress.css';
import Prism from 'prismjs';
import App from './App.vue';
import store from './store';
import './style/scss/index.scss';
import router from './router';

Vue.config.productionTip = false;

// Make BootstrapVue available throughout your project
Vue.use(BootstrapVue);
// Optionally install the BootstrapVue icon components plugin
Vue.use(IconsPlugin);
Vue.use(Vuelidate);
Vue.use(ElementUI);
Vue.use(mouse);
Vue.use(VueParticles);
Vue.use(VueAxios, axios);
Vue.use(VueScrollTo);
Vue.use(VueTypedJs);
Vue.use(Picker);
Vue.prototype.$axios = axios;

VueMarkdownEditor.use(vuepressTheme, {
  Prism,
});
Vue.use(VueMarkdownEditor);

new Vue({
  router,
  store,
  render: (h) => h(App),
  mounted() {
    this.$mouseClick();
  },
}).$mount('#app');
