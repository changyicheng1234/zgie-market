import storageService from '@/service/storageService';
// eslint-disable-next-line import/no-cycle
import userService from '@/service/userService';

const userModule = {
  namespaced: true,
  state: {
    // 从本地缓存中获取
    token: storageService.get(storageService.USER_TOKEN),
    userInfo: storageService.get(storageService.USER_INFO) ? JSON.parse(storageService.get(storageService.USER_INFO)) : null, //eslint-disable-line
  },

  mutations: {
    SET_TOKEN(state, token) {
      // 更新本地缓存
      storageService.set(storageService.USER_TOKEN, token);
      // 更新state
      state.token = token;
    },
    SET_USERINFO(state, userInfo) {
      // 更新本地缓存
      storageService.set(storageService.USER_INFO, JSON.stringify(userInfo));
      // 更新state
      state.userInfo = userInfo;
    },
  },

  actions: {
    updateLocalUserInfo(context, updatedUserInfo) {
      return new Promise((resolve) => {
        // merge the updated fields with the existing user information
        const userInfo = {
          ...context.state.userInfo,
          ...updatedUserInfo,
        };

        // Commit the mutation to update the state
        context.commit('SET_USERINFO', userInfo);
        resolve();
      });
    },
    // context是浏览器store上下文对象，用于保存token
    // { name, telephone, password }为传送给后端的数据的结构体变量
    register(context, {
      name,
      phone,
      password,
      password2,
      email,
      num,
      valiCode,
      CDKey,
    }) {
      return new Promise((resolve, reject) => {
        // userService.register()为发送网络请求的函数，用于向后端发送注册信息
        userService.register({
          name,
          phone,
          password,
          password2,
          email,
          num,
          valiCode,
          CDKey,
        }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    modifyPassword(context, {
      email, password, password2, valiCode,
    }) {
      return new Promise((resolve, reject) => {
        userService.modifyPassword({
          email, password, password2, valiCode,
        }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    deleteMe(context, { phone, email, valiCode }) {
      console.error(email);
      return new Promise((resolve, reject) => {
        userService.deleteUser({ phone, email, valiCode }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    // DEPRECATED
    identityValidate(context, { email, valiCode, mode }) {
      return new Promise((resolve, reject) => {
        userService.identityValidate({ email, valiCode, mode }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    validateEmail(context, { email, mode }) {
      return new Promise((resolve, reject) => {
        userService.validateEmail({ email, mode }).then((res) => {
          console.error('1');
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    login(context, { email, password }) {
      return new Promise((resolve, reject) => {
        userService.login({ email, password }).then((res) => {
          // 保存token
          context.commit('SET_TOKEN', res.data.data.token);
          return userService.info();
        }).then((res) => {
          console.log(res.data.data.user);
          // 保存用户信息
          context.commit('SET_USERINFO', res.data.data.user);
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },

    logout({ commit }) {
      // 清除token
      commit('SET_TOKEN', '');
      storageService.set(storageService.USER_TOKEN, '');
      // 清除用户信息
      commit('SET_USERINFO', '');
      storageService.set(storageService.USER_INFO, '');
      localStorage.removeItem('Partition');
      window.location.reload();
    },

    save(context, { userTelephone, postID, isSaved }) {
      return new Promise((resolve, reject) => {
        userService.save({ userTelephone, postID, isSaved }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
    getInfo(context, { phone }) {
      return new Promise((resolve, reject) => {
        userService.getInfo({ phone }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
      });
    },
  },
};

export default userModule;
