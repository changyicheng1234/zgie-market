import Vue from 'vue';
import VueRouter from 'vue-router';
import { Toast } from 'vant';
import store from '@/store';

Vue.use(VueRouter);

const routes = [
  {
    path: '/',
    name: 'home',
    meta: {
      auth: true,
      keepAlive: true,
    },
    component: () => import('../views/home/HomeView.vue'),
  },
  {
    path: '/welcome',
    name: 'welcome',
    component: () => import('../views/welcome/WelcomeView.vue'),
  },
  {
    path: '/discover',
    name: 'discover',
    meta: {
      auth: true,
    },
    component: () => import('../views/discover/DiscoverView.vue'),
  },
  {
    path: '/notice',
    name: 'notice',
    meta: {
      auth: true,
    },
    component: () => import('../views/notice/NoticeView.vue'),
  },
  {
    path: '/profile',
    name: 'profile',
    meta: {
      auth: true,
    },
    component: () => import('../views/profile/ProfileView.vue'),
  },
  {
    path: '/logintest',
    name: 'logintest',
    component: () => import('../views/logintest/LoginTestView.vue'),
  },
  {
    path: '/postDetails',
    name: 'postDetails',
    meta: {
      auth: true,
    },
    component: () => import('../views/postDetails/PostDetailsView.vue'),
  },
  {
    path: '/post',
    name: 'post',
    meta: {
      auth: true,
    },
    component: () => import('../views/post/PostView.vue'),
  },
  {
    path: '/save',
    name: 'save',
    meta: {
      auth: true,
      keepAlive: true,
    },
    component: () => import('../views/save/SaveView.vue'),
  },
  {
    path: '/history',
    name: 'history',
    meta: {
      auth: true,
      keepAlive: true,
    },
    component: () => import('../views/history/HistoryView.vue'),
  },
  {
    path: '/feedback',
    name: 'feedback',
    meta: {
      auth: true,
    },
    component: () => import('../views/feedback/FeedbackView.vue'),
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/register/RegisterView.vue'),
  },
  {
    path: '/changeDetail',
    name: 'changeDetail',
    meta: {
      auth: true,
    },
    component: () => import('../views/profile/ChangeDetailView.vue'),
  },
  {
    path: '/modifyPassword',
    name: 'modifyPassword',
    component: () => import('../views/modifyPassword/ModifyPasswordView.vue'),
  },
  {
    path: '/filemanage',
    name: 'filemanage',
    meta: {
      auth: true,
    },
    component: () => import('../views/file/FileManageView.vue'),
  },
  {
    path: '/course',
    name: 'course',
    meta: {
      auth: true,
      keepAlive: true,
    },
    component: () => import('../views/course/CourseView.vue'),
  },
];

const router = new VueRouter({
  // mode: 'history',
  mode: 'hash',
  base: process.env.BASE_URL,
  scrollBehavior(to, from, savePosition) {
    if (savePosition) {
      return savePosition;
    }
    return { x: 0, y: 0 };
  },
  routes,
});
router.beforeEach((to, from, next) => {
  if (to.meta.auth) { // 判断是否需要登录
    // 判断用户是否登录
    if (store.state.userModule.token) {
      // 这里还要判断token 的有效性 比如有没有过期 需要后台发放token 的时候 带上token 的有效期，
      // 如果 token 无效 需要 请求token
      // 每次跳转路由都更新一下通知数
      store.commit('getNoticesNum');
      next();
    } else {
      // 跳转登录
      Toast.fail('请先登录');
      if (from.name !== 'logintest') router.push({ name: 'logintest' });
    }
  } else {
    next();
  }
});
export default router;
