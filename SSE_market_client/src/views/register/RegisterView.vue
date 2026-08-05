<!-- <template>中放的是页面的html框架，可用html语言编写，也可调用bootstrap-vue等第三方提供的组件，但需要在main.js中引入  -->
<template>
  <FormPageBase>
    <h3 style="text-align: center;">用户注册</h3>
    <b-form @keydown.enter="register">
      <b-form-group label='用户名'>
        <!-- v-model里的内容是下面script中data()函数中return的属性，可实现数据连接 -->
        <b-form-input
          v-model='$v.user.name.$model'
          type='text'
          placeholder='输入您的用户名'
        />
      </b-form-group>
      <!-- <b-form-group label='账号'>
        <b-form-input
          v-model="$v.user.phone.$model"
          type="number"
          placeholder="输入账号,可以不填或乱填11位"
          :state="validateState('phone')"
        ></b-form-input>
        <b-form-invalid-feedback :state="validateState('phone')">
          手机号不符合要求
        </b-form-invalid-feedback>
      </b-form-group> -->
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
          v-model="$v.user.password.$model"
          type="password"
          placeholder="输入密码"
          :state="validateState('password')"
        />
        <b-form-invalid-feedback :state="validateState('password')">
          密码必须大于等于 6 位
        </b-form-invalid-feedback>
      </b-form-group>
      <b-form-group label='确认密码'>
        <b-form-input
          v-model="$v.user.password2.$model"
          type="password"
          placeholder="再次输入密码"
          :state="validateState('password2')"
        />
        <b-form-invalid-feedback :state="validateState('password2')">
          密码必须大于等于 6 位
        </b-form-invalid-feedback>
      </b-form-group>
      <b-form-group label='邮箱验证码'>
        <b-form-input
          v-model="$v.user.valiCode.$model"
          type="valiCode"
          placeholder="输入邮箱验证码"
        />
      </b-form-group>
      <b-form-group label='邀请码'>
        <b-form-input
          v-model="$v.user.CDKey.$model"
          type="text"
          placeholder="输入邀请码"
        />
      </b-form-group>
      <!-- 使用click添加事件，这里是定向到下面的register函数。click的另一个用户参见'views/layout/NavbarView.vue'文件 -->
      <b-button @click='register' variant='outline-primary' block>
        注册
      </b-button>
    </b-form>
  </FormPageBase>
</template>

<script>
// 引入一些需要的组件及库
import { required, minLength } from 'vuelidate/lib/validators';
import { mapActions } from 'vuex';
import customValidator from '@/helper/validator';

// eslint-disable-next-line import/no-extraneous-dependencies
import CryptoJS from 'crypto-js';
import FormPageBase from '@/components/FormPageBase.vue';

// export default用于导出模块中的一个默认值，函数一般放在其中
export default {
  components: {
    FormPageBase,
  },
  // data 函数用于定义组件的数据。它返回一个对象，其中包含组件的各种数据属性和初始值。
  data() {
    return {
      time: 60,
      timeshow: false,
      user: {
        name: '',
        phone: '',
        password: '',
        password2: '',
        email: '',
        valiCode: '',
        CDKey: '',
        mode: '',
      },
      key: '16bit secret key',
    };
  },
  // 用于定义数据的格式规则
  validations: {
    user: {
      name: {},
      phone: {
        phone: customValidator.telephoneValidator,
      },
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
      valiCode: {},
      CDKey: {},
      mode: {},
    },
  },
  // methods 属性用于定义组件的方法。它是一个对象，其中包含了一组可以在组件实例中使用的方法
  methods: {
    // 使用map将userRegister映射到'store/module'里的user函数，在下面使用时就写成userRegister
    ...mapActions('userModule', { userRegister: 'register' }),
    ...mapActions('userModule', { userValidate: 'validateEmail' }),
    // validators库里规定的，用于规则验证什么的，我也不是很了解
    validateState(name) {
      // 这里是es6 解构赋值
      const {
        $dirty,
        $error,
      } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    setPassword(data, key) {
      const cypherKey = CryptoJS.enc.Utf8.parse(key);
      CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);
      const iv = CryptoJS.SHA256(key).toString();
      const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };
      return CryptoJS.AES.encrypt(data, cypherKey, cfg).toString();
    },
    validateEmail() {
      this.user.mode = 0;
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        return;
      }
      this.$bvToast.toast('正在发送验证码请求', {
        title: '系统提醒',
        variant: 'primary',
        solid: true,
      });
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
      // 发送获取验证码请求
      this.userValidate(this.user).then(() => {
        this.$bvToast.toast('已发送验证码，请将邮箱发送的验证码输入以完成注册验证', {
          title: '系统提醒',
          variant: 'primary',
          solid: true,
        });
      }).catch((err) => {
        this.$bvToast.toast(err.response.data.msg, {
          title: '发送邮箱错误，请检查并重试',
          variant: 'danger',
          solid: true,
        });
        clearInterval(setTimeoutS);
        this.timeshow = false;
      });
    },
    register() {
      // 验证数据
      this.$v.user.$touch();
      if (this.$v.user.$anyError) {
        return;
      }
      // console.error(this.user);
      this.$bvToast.toast('已发送注册请求，请稍等', {
        title: '系统提醒',
        variant: 'primary',
        solid: true,
      });
      // 请求
      this.userRegister({
        ...this.user,
        password: this.setPassword(this.user.password, this.key),
        password2: this.setPassword(this.user.password2, this.key),
      }).then(() => {
        // 这里写从后端成功返回数据后的操作
        this.$bvToast.toast('注册成功，请登录', {
          title: '系统提醒',
          variant: 'primary',
          solid: true,
        });
        setTimeout(() => {
          this.$router.push({ name: 'login' });
        });
      }).catch((err) => {
        const msg = err.response ? err.response.data.msg : '网络错误，请检查连接后重试';
        this.$bvToast.toast(msg, {
          title: '注册失败',
          variant: 'danger',
          solid: true,
        });
      });
    },
  },
};
</script>
