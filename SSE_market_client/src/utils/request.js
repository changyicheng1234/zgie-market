import axios from 'axios';
import storageService from '@/service/storageService';
// eslint-disable-next-line import/no-cycle
import store from '@/store';

const service = axios.create({
  baseURL: process.env.VUE_APP_BASE_URL || 'https://ssemarket.cn/api/',
  timeout: 100000,
});

service.interceptors.request.use((config) => {
  Object.assign(config.headers, { Authorization: `Bearer ${storageService.get(storageService.USER_TOKEN)}` });
  return config;
}, (error) => {
  return Promise.reject(error);
});

// 添加响应拦截器
service.interceptors.response.use((response) => {
  return response;
}, async (error) => {
  const originalRequest = error.config;
  // 如果错误码为401且原始请求未重试，则说明token已经过期，需要重新获取
  // eslint-disable-next-line no-underscore-dangle
  if (error.response.status === 401 && !originalRequest._retry) {
    // this.showMessage();这里应该有提示消息，但没成功
    setTimeout(() => {
      // eslint-disable-next-line no-underscore-dangle
      originalRequest._retry = true;
      // 清除token
      store.commit('SET_TOKEN', '');
      storageService.set(storageService.USER_TOKEN, '');
      // 清除用户信息
      store.commit('SET_USERINFO', '');
      storageService.set(storageService.USER_INFO, '');
      localStorage.removeItem('Partition');
      window.location.reload();
    }, 2000);
  }
  return Promise.reject(error);
});
export default service;
