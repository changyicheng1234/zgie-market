<template>
<div>
    <div>
        修改密码
    </div>

    <p v-html="line"></p>

    <div>
      <van-form inset>
        <van-field
          v-model="user.email"
          label="邮箱"
          required
          type="email"
          placeholder="请输入邮箱"
          :error-message="errorEmailInfo"
        />
        <van-field
          v-model="user.password"
          type="password"
          required
          label="密码"
          placeholder="请输入密码"
          :error-message="errorPasswordInfo"
        />

        <van-field
          v-model="user.password2"
          type="password"
          required
          label="再次输入密码"
          placeholder="请再次输入密码"
          :error-message="errorPassword2Info"
        />

        <van-field
          v-model="user.valiCode"
          label="验证码"
          required
          placeholder="请输入邮箱验证码"
        >
          <template #button>
            <van-button size="small" type="info" v-show="timeshow === true" disabled>{{time}}秒后重新获取</van-button>
            <van-button size="small" type="info" v-show="timeshow === false" @click="validateEmail">发送验证码</van-button>
          </template>
        </van-field>
      </van-form>

      <van-button size="normal" type="info" @click="modifyPassword">
        修改密码
      </van-button>
    </div>
</div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import { required, minLength } from 'vuelidate/lib/validators';
import { len } from 'vuelidate/lib/validators/common';
// eslint-disable-next-line import/no-extraneous-dependencies
import CryptoJS from 'crypto-js';

import customValidator from '@/helper/validator';

export default {
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),
  data() {
    return {
      time: 60,
      timeshow: false,
      line: '<br/>',
      errorEmailInfo: '',
      errorPasswordInfo: '',
      errorPassword2Info: '',
      user: {
        mode: '',
        email: '',
        password: '',
        password2: '',
        valiCode: '',
      },
      key: '16bit secret key',
    };
  },
  validations: {
    user: {
      password: {
        required,
        password: minLength(6),
      },
      password2: {
        required,
        password2: minLength(6),
      },
      email: {
        required,
        email: customValidator.emailValidator,
      },
    },
  },
  methods: {
    ...mapActions('userModule', { userModify: 'modifyPassword' }),
    ...mapActions('userModule', { userValidate: 'validateEmail' }),

    setPassword(data, key) {
      const cypherKey = CryptoJS.enc.Utf8.parse(key);
      CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);
      const iv = CryptoJS.SHA256(key).toString();
      const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };
      return CryptoJS.AES.encrypt(data, cypherKey, cfg).toString();
    },

    modifyPassword() {
      if (len(this.user.password) < 6
        || this.user.password !== this.user.password2
        || !customValidator.emailValidator(this.user.email)
        || len(this.user.valiCode) !== 6
      ) {
        this.$toast.fail('请检查输入信息');
        return;
      }
      this.userModify({
        ...this.user,
        password: this.setPassword(this.user.password, this.key),
        password2: this.setPassword(this.user.password2, this.key),
      }).then(() => {
        this.$toast.success('修改密码成功');
      }).catch((err) => {
        this.$toast.fail(err.response.data.msg);
      });
    },

    validateEmail() {
      if (!customValidator.emailValidator(this.user.email)) {
        this.$toast.fail('请输入格式正确的邮箱');
        return;
      }
      this.user.mode = 1;
      this.userValidate(this.user).then(() => {
        this.timeshow = true;
        this.time = 60;
        this.$toast.success('已发送验证码，请将邮箱发送的验证码输入以完成注册验证');
        const setTimeoutS = setInterval(() => {
          this.time -= 1;
          if (this.time <= 0) {
            clearInterval(setTimeoutS);
            this.timeshow = false;
          }
        }, 1000);
      }).catch((err) => {
        this.$toast.fail(err.response.data.msg);
      });
    },
  },

  watch: {
    // eslint-disable-next-line func-names
    'user.password': function (password) {
      if (len(password) < 6) {
        this.errorPasswordInfo = '密码长度需要不小于6位';
      } else {
        this.errorPasswordInfo = '';
      }
    },
    // eslint-disable-next-line func-names
    'user.password2': function (password2) {
      if (len(password2) < 6) {
        this.errorPassword2Info = '密码长度需要不小于6位';
      } else {
        this.errorPassword2Info = '';
      }
    },
    // eslint-disable-next-line func-names
    'user.email': function (email) {
      if (!customValidator.emailValidator(email)) {
        this.errorEmailInfo = '请输入格式正确的邮箱';
      } else {
        this.errorEmailInfo = '';
      }
    },
  },
};
</script>
