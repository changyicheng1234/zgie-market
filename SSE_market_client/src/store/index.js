import Vue from 'vue';
import Vuex from 'vuex';
// eslint-disable-next-line import/no-cycle
import request from '@/utils/request';
import userModule from './module/user';
import postModule from './module/post';
import commentModule from './module/comment';

Vue.use(Vuex);

export default new Vuex.Store({
  strict: process.env.NODE_ENV !== 'production',
  state: {
    noticesNum: 0,
  },
  getters: {
  },
  mutations: {
    getNoticesNum(state) {
      request.get('/auth/getNoticeNum').then((response) => {
        state.noticesNum = response.data.unreadTotalNum;
      }).catch((error) => {
        console.error(error);
      });
    },
  },
  actions: {
  },
  modules: {
    userModule,
    postModule,
    commentModule,
  },
});
