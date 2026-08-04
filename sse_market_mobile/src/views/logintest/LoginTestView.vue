<template>
  <div class="login-view">
    <h2 class="login-title">SSE_market</h2>
    <div class="login-container">
        <div class="login">
          <van-row>
            <van-cell>
              <van-form>
                <van-field
                  v-model="user.email"
                  label="邮箱"
                  type="text"
                  placeholder="请输入邮箱"
                  required
                  :error-message="errorMessage"
                ></van-field>
                <van-field
                  v-model="password"
                  label="密码"
                  type="password"
                  placeholder="输入密码"
                  required
                  :error-message="errorInfo"
                ></van-field>
                <van-row style="margin: 10px" type="flex" justify="center">
                  <van-col>
                    <van-button round color="skyblue" @click="login">登录</van-button>
                  </van-col>
                </van-row>
                <van-row style="margin-top: 10px" type="flex" justify="center">
                  <van-col>
                    <router-link style="margin-top: 10px" to="/register">注册</router-link>
                    <router-link style="margin-top: 10px;margin-left:30px"
                      to="/modifyPassword">忘记密码?</router-link>
                  </van-col>
                </van-row>
              </van-form>
            </van-cell>
          </van-row>
        </div>
    <div class="footer"></div>
        </div>
  </div>
</template>

<script>

import { mapActions } from 'vuex';
import { len } from 'vuelidate/lib/validators/common';
// eslint-disable-next-line import/no-extraneous-dependencies
import CryptoJS from 'crypto-js';

import validator from '@/helper/validator';

export default {
  data() {
    return {
      password: '',
      user: {
        email: '',
        password: '',
      },
      showKeyboard: false,
      errorInfo: '', // 密码错误提示
      errorMessage: '', // 邮箱错误提示
      key: '16bit secret key',
      // item: require('@/assets/用户条款.txt'),
    };
  },
  methods: {
    ...mapActions('userModule', { userlogin: 'login' }),
    setPassword(data, key) {
      const cypherKey = CryptoJS.enc.Utf8.parse(key);
      CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);
      const iv = CryptoJS.SHA256(key).toString();
      const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };
      return CryptoJS.AES.encrypt(data, cypherKey, cfg).toString();
    },
    login() {
      this.user.password = this.setPassword(this.password, this.key);
      // console.error(this.user);
      this.userlogin(this.user)
        .then(() => {
          this.$toast.success('登录成功');
          this.$router.replace({ name: 'home' });
        })
        .catch((err) => {
          console.error(err); // 打印完整的错误对象
          if (err.response && err.response.data) {
            this.$toast.fail(err.response.data.msg);
          } else {
            this.$toast.fail('登录失败'); // 或者其他通用错误消息
          }
        });
    },
  },
  watch: {
    // eslint-disable-next-line func-names
    password(password) {
      if (len(password) < 6) {
        this.errorInfo = '密码长度>=6';
      } else {
        this.errorInfo = '';
      }
    },
    // eslint-disable-next-line func-names
    'user.email': function (email) {
      if (!validator.emailValidator(email)) {
        this.errorMessage = '请输入格式正确的邮箱';
      }
    },
  },
};
</script>

<style scoped>
.van-button--round {
  width: 200px;
}
.login-view {
  background-color: #edf6f9; /* 设置整体背景颜色为淡灰色 #e4e4e4*/
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  height: 100vh;
}

.login-title {
  padding: 20px 0;
  margin-top: 0px; /* 将标题向上移动 20px，以覆盖灰色部分 */
  margin-bottom: 0px;
  background-color: skyblue; /* 设置标题的背景颜色为蓝色 */
  width: 100%;
  box-sizing: border-box;
  text-align: center;
  color:black;
  size: 25px;
}

.login-container {
  background-color: white;
  padding: 20px;
}

.footer {
  flex-grow: 1;
}
</style>
