/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable import/no-duplicates */
/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable import/no-duplicates */
/* eslint-disable import/no-extraneous-dependencies */
// eslint-disable-next-line import/no-extraneous-dependencies
/* eslint-disable import/no-duplicates */
/* eslint-disable import/no-extraneous-dependencies */
/* eslint-disable import/extensions */
import Vue from 'vue';
import Vant from 'vant';
import 'vant/lib/index.css';
import axios from 'axios';
import VueAxios from 'vue-axios';
import preview from 'vue-photo-preview';
import VueScrollTo from 'vue-scrollto';
import App from './App.vue';
import './utils/amfe-flexible';
import './assets/css/index.css';
// eslint-disable-next-line
import router from './router';
import store from './store';
import 'vue-photo-preview/dist/skin.css';
// eslint-disable-next-line import/extensions
import './assets/iconfont/iconfont.js';

const options = {
  fullscreenEl: false, // 控制是否显示右上角全屏按钮
  closeEl: false, // 控制是否显示右上角关闭按钮
  tapToClose: true, // 点击滑动区域应关闭图库
  shareEl: false, // 控制是否显示分享按钮
  zoomEl: false, // 控制是否显示放大缩小按钮
  counterEl: false, // 控制是否显示左上角图片数量按钮
  tapToToggleControls: true, // 点击应切换控件的可见性
  clickToCloseNonZoomable: true, // 点击图片应关闭图库，仅当图像小于视口的大小时
  indexIndicatorSep: ' / ', // 图片数量的分隔符
};

Vue.config.productionTip = false;
Vue.use(Vant);
Vue.use(VueAxios, axios);
Vue.use(preview, options);
Vue.use(VueScrollTo);
Vue.prototype.$axios = axios;

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
