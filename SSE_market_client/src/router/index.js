import Vue from 'vue';
import VueRouter from 'vue-router';
import store from '@/store';
import userRoutes from './module/user';
import postRoutes from './module/post';

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
    path: '/filemanage',
    name: 'filemanage',
    meta: {
      auth: true,
      keepAlive: false,
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
  ...userRoutes, // 引入userRoutes中的路由
  ...postRoutes, // 引入postRoutes中的路由
];

// VueRouter是一个用于处理路由的插件。可以使用它来定义应用程序的路由映射，并在用户浏览时动态地更新视图，以此构建单页应用程序。
const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) { // 如果savedPosition存在，滚动条会自动跳到记录的值的地方
      return savedPosition;
    }
    return {
      x: 0,
      y: 0,
    };// savedPosition也是一个记录x轴和y轴位置的对象
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
      router.push({ name: 'login' });
    }
  } else {
    next();
  }
});

export default router;
