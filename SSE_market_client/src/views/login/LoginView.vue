<template>
  <FormPageBase>
    <h3 style="text-align: center;">欢迎来到软工集市</h3>
    <b-form @keydown.enter='login'>
      <b-form-group label='邮箱'>
        <b-form-input
          v-model='$v.user.email.$model'
          type='email'
          placeholder='输入邮箱'
          :state="validateState('email')"
        />
        <b-form-invalid-feedback :state="validateState('email')">
          邮箱不符合要求
        </b-form-invalid-feedback>
      </b-form-group>
      <b-form-group label='密码'>
        <b-form-input
          v-model='$v.user.password.$model'
          type='password'
          placeholder='输入密码'
          :state="validateState('password')"
        />
        <b-form-invalid-feedback :state="validateState('password')">
          密码必须大于等于 6 位
        </b-form-invalid-feedback>
      </b-form-group>
      <b-form-group>
        <b-form-checkbox
          v-model="rememberMe"
          name="rememberMe"
        >
          记住我
        </b-form-checkbox>
      </b-form-group>
      <b-button @click='login' variant='outline-primary' block>
        登录
      </b-button>
    </b-form>
    <b-row>
      <b-col>
        <a href="#" @click="toRegister">立即注册！</a>
      </b-col>
      <b-col class="text-right">
        <a href="#" @click="toModifyPassword">忘记密码？</a>
      </b-col>
    </b-row>
  </FormPageBase>
</template>

<script>
import { required, minLength } from 'vuelidate/lib/validators';
import { mapActions } from 'vuex';
import customValidator from '@/helper/validator';

// eslint-disable-next-line import/no-extraneous-dependencies
import CryptoJS from 'crypto-js';
import FormPageBase from '@/components/FormPageBase.vue';

export default {
  components: {
    FormPageBase,
  },
  data() {
    return {
      user: {
        email: '',
        password: '',
      },
      rememberMe: false,
    };
  },
  validations: {
    user: {
      email: {
        required,
        email: customValidator.emailValidator,
      },
      password: {
        required,
        minLength: minLength(6),
      },
    },
  },

  mounted() {
    // Load saved login information if rememberMe is true
    if (localStorage.rememberMe === 'true') {
      this.user.email = localStorage.email || '';
      this.user.password = localStorage.password || '';
      this.rememberMe = true;
    }
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
    validateState(name) {
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    login() {
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        return;
      }
      // Save login information if rememberMe is true
      if (this.rememberMe) {
        localStorage.rememberMe = true;
        localStorage.email = this.user.email;
        localStorage.password = this.user.password;
      } else {
        // Clear saved login information if rememberMe is false
        localStorage.removeItem('rememberMe');
        localStorage.removeItem('email');
        localStorage.removeItem('password');
      }
      this.userlogin({
        email: this.user.email,
        password: this.setPassword(this.user.password, '16bit secret key'),
      }).then(() => {
        this.$router.replace({ name: 'home' });
      }).catch((err) => {
        this.$bvToast.toast(err.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
      });
    },
    toRegister() {
      const { href } = this.$router.resolve({ path: '/register' });
      window.open(href, '_blank');
    },
    toModifyPassword() {
      const { href } = this.$router.resolve({ path: '/modifyPassword' });
      window.open(href, '_blank');
    },
  },
};
</script>
