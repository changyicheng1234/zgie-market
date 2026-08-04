<template>
  <FormPageBase>
    <h3 style="text-align: center;">修改密码</h3>
    <b-form @keydown.enter="modifyPassword">
      <b-form-group label='邮箱'>
        <b-input-group>
          <b-form-input
            v-model="$v.user.email.$model"
            type="email"
            placeholder="输入邮箱"
            :state="validateState('email')"
          />
          <b-input-group-append width="150px">
            <b-button v-show="timeshow===true" disabled variant='outline-primary'>
              {{ time }}秒后重新获取
            </b-button>
            <b-button @click='validateEmail' v-show="timeshow===false" variant='outline-primary'>
              发送验证码
            </b-button>
          </b-input-group-append>
        </b-input-group>
        <b-form-invalid-feedback :state="validateState('email')">
          邮箱格式不符合要求
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
      <b-form-group label='确认密码'>
        <b-form-input
          v-model='$v.user.password2.$model'
          type='password'
          placeholder='再次输入密码'
          :state="validateState('password2')"
        />
        <b-form-invalid-feedback :state="validateState('password2')">
          密码必须大于等于 6 位
        </b-form-invalid-feedback>
      </b-form-group>
      <b-form-group label='验证码'>
        <b-form-input
          v-model='user.valiCode'
          type='valiCode'
          placeholder='输入验证码'
        />
      </b-form-group>

      <b-button @click='modifyPassword' variant="outline-primary" block>
        修改密码
      </b-button>
      <b-button @click="$router.replace({ name: 'login' })" variant="outline-primary" block>
        返回登录
      </b-button>
    </b-form>
  </FormPageBase>
</template>

<script>
import { mapActions } from 'vuex';
import { required, minLength } from 'vuelidate/lib/validators';
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
      time: 60,
      timeshow: false,
      user: {
        mode: 1,
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
      email: {
        required,
        email: customValidator.emailValidator,
      },
      password: {
        required,
        password: minLength(6),
      },
      password2: {
        required,
        password2: minLength(6),
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
    validateState(name) {
      // 这里是es6 解构赋值
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    validateEmail() {
      this.user.mode = 1;
      this.$v.user.$touch();
      if (this.$v.user.email.$anyError) {
        return;
      }
      this.userValidate(this.user).then(() => {
        this.$bvToast.toast('已发送验证码，请将邮箱发送的验证码输入以完成注册验证', {
          title: '系统提醒',
          variant: 'primary',
          solid: true,
        });
        this.timeshow = true;
        this.time = 60;
        const setTimeoutS = setInterval(() => {
          this.time -= 1;
          if (this.time <= 0) {
            clearInterval(setTimeoutS);
            this.timeshow = false;
          }
        }, 1000);
      }).catch((err) => {
        this.$bvToast.toast(err.response.data.msg, {
          title: '发送邮箱错误',
          variant: 'danger',
          solid: true,
        });
      });
    },
    modifyPassword() {
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        return;
      }
      this.userModify({
        ...this.user,
        password: this.setPassword(this.user.password, this.key),
        password2: this.setPassword(this.user.password2, this.key),
      }).then(() => {
        this.$bvToast.toast('修改密码成功', {
          title: '系统提醒',
          variant: 'primary',
          solid: true,
        });
      }).catch((err) => {
        this.$bvToast.toast(err.response.data.msg, {
          title: '数据验证错误',
          variant: 'danger',
          solid: true,
        });
      });
    },
  },
};
</script>
