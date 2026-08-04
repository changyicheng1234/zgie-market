<!-- eslint-disable max-len -->
<template>
  <div>
    <div>
      注册
    </div>

    <p v-html="line"></p>

    <div>
      <van-cell-group inset>
        <van-field
          v-model="user.name"
          label="用户名"
          placeholder="请输入用户名"
        />

        <van-field
          v-model="user.email"
          required
          label="邮箱"
          placeholder="请输入邮箱"
          :error-message="errorEmailMessage"
        />

        <van-field
          v-model="password"
          type="password"
          required
          label="密码"
          placeholder="请输入密码"
          :error-message="errorPasswordMessage"
        />

        <van-field
          v-model="password2"
          type="password"
          required
          label="再次输入密码"
          placeholder="请再次输入密码"
          :error-message="errorPassword2Message"
        />

        <van-field
          v-model="user.valiCode"
          required
          label="验证码"
          placeholder="请输入邮箱验证码"
        >
          <template #button>
            <van-button size="small" type="info" v-show="timeshow === true" disabled>{{time}}秒后重新获取</van-button>
            <van-button size="small" type="info" v-show="timeshow === false" @click="validateEmail">发送验证码</van-button>
          </template>
        </van-field>
        <van-field
          v-model="user.CDKey"
          required
          label="邀请码"
          placeholder="请输入邀请码"
        />
      </van-cell-group>
      <van-button size="normal" type="info" @click="register">
        注册
      </van-button>
    </div>
  </div>
</template>

<script>
import { required, minLength, maxLength } from 'vuelidate/lib/validators';
import { mapActions } from 'vuex';
import { len } from 'vuelidate/lib/validators/common';
// eslint-disable-next-line import/no-extraneous-dependencies
import CryptoJS from 'crypto-js';

import customValidator from '@/helper/validator';

export default {
  data() {
    return {
      time: 60,
      timeshow: false,
      line: '<br/>',
      errorEmailMessage: '',
      errorPassword2Message: '',
      errorPasswordMessage: '',
      errorNumMessage: '',
      errorPhoneMessage: '',
      emailCheck: false,
      password: '',
      password2: '',
      user: {
        phone: '',
        name: '',
        password: '',
        email: '',
        password2: '',
        valiCode: '',
        CDKey: '',
        mode: '',
      },
      key: '16bit secret key',
      validator: {
        name: {},
        password: {
          required,
          minLength: minLength(6),
        },
        password2: {
          required,
          minLength: minLength(6),
        },
        email: {
          required,
          email: customValidator.emailValidator,
        },
        num: {
          required,
          minLength: minLength(8),
          maxLength: maxLength(8),
        },
        valiCode: {},
        CDKey: {},
        mode: {},
      },
    };
  },

  methods: {
    ...mapActions('userModule', { userRegister: 'register' }),
    ...mapActions('userModule', { userValidate: 'validateEmail' }),

    setPassword(data, key) {
      const cypherKey = CryptoJS.enc.Utf8.parse(key);
      CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);
      const iv = CryptoJS.SHA256(key).toString();
      const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };
      return CryptoJS.AES.encrypt(data, cypherKey, cfg).toString();
    },

    validateEmail() {
      this.user.mode = 0;
      // console.error(this.user);
      if (this.emailCheck === true) {
        // 开始倒计时
        this.timeshow = true;
        this.time = 60;
        const setTimeoutS = setInterval(() => {
          this.time -= 1;
          if (this.time <= 0) {
            clearInterval(setTimeoutS);
            this.timeshow = false;
          }
        }, 1000);
        // console.error(this.user);
        this.userValidate(this.user)
          .then(() => {
            this.$toast.success('已发送验证码，请将邮箱发送的验证码输入以完成注册验证');
          })
          .catch((err) => {
            console.error(err);
            this.$toast.fail(err.response.data.msg);
            // 取消倒计时
            clearInterval(setTimeoutS);
            this.timeshow = false;
          });
      } else {
        this.$toast.fail('邮箱格式有误');
      }
    },

    register() {
      this.user.password = this.setPassword(this.password, this.key);
      this.user.password2 = this.setPassword(this.password2, this.key);
      // 请求
      this.userRegister(this.user)
        .then(() => {
          // 这里写从后端成功返回数据后的操作
          this.$toast.success('注册成功,请登录', {
            title: '系统提醒',
            variant: 'primary',
            solid: true,
          });
          setTimeout(() => {
            this.$router.push({ name: 'logintest' });
          });
        })
        .catch((err) => {
          // 这里写从后端返回数据失败后的操作
          this.$toast.fail(err.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            solid: true,
          });
        });
    },
  },

  watch: {
    password(password) {
      if (len(password) < 6) {
        this.errorPasswordMessage = '密码长度需要不小于6位';
      } else {
        this.errorPasswordMessage = '';
      }
    },
    password2(password2) {
      if (len(password2) < 6) {
        this.errorPassword2Message = '密码长度需要不小于6位';
      } else {
        this.errorPassword2Message = '';
      }
    },
    // eslint-disable-next-line func-names
    'user.email': function (email) {
      if (!customValidator.emailValidator(email)) {
        this.errorEmailMessage = '请输入格式正确的邮箱';
        this.emailCheck = false;
      } else {
        this.errorEmailMessage = '';
        this.emailCheck = true;
      }
    },
  },
};

</script>

<style>

</style>
