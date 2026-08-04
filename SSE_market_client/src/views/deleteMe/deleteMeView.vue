<template>
  <FormPageBase>
    <h3 style="text-align: center;">永久注销账号</h3>
    <b-form @keydown.enter="deleteMe">
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
      <b-form-group label='验证码'>
        <b-form-input
          v-model='user.valiCode'
          type='valiCode'
          placeholder='输入验证码'
        />
      </b-form-group>
      <b-button @click='deleteUser' variant='outline-danger' block>
        注销账号
      </b-button>
    </b-form>
  </FormPageBase>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import { required } from 'vuelidate/lib/validators';
import customValidator from '@/helper/validator';
import FormPageBase from '@/components/FormPageBase.vue';

export default {
  components: {
    FormPageBase,
  },
  data() {
    return {
      timeshow: false,
      time: 60,
      user: {
        phone: '',
        email: '',
        valiCode: '',
        mode: '',
      },
    };
  },
  validations: {
    user: {
      email: {
        required,
        email: customValidator.emailValidator,
      },
    },
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
  },
  methods: {
    ...mapActions('userModule', { userValidate: 'validateEmail' }),
    ...mapActions('userModule', ['deleteMe']),
    ...mapActions('userModule', ['logout']),
    validateState(name) {
      // 这里是es6 解构赋值
      const { $dirty, $error } = this.$v.user[name];
      return $dirty ? !$error : null;
    },
    ensureUser() {
      if (this.userInfo === null || this.userInfo === undefined) {
        this.$bvToast.toast('请先登录', {
          title: '系统提醒',
          variant: 'danger',
          solid: true,
        });
        setTimeout(() => this.$router.replace({ name: 'login' }), 2000);
        return false;
      }
      if (this.userInfo.email !== this.user.email) {
        this.$bvToast.toast('邮箱与登录用户不一致', {
          title: '系统提醒',
          variant: 'danger',
          solid: true,
        });
        return false;
      }
      this.$v.user.$touch();
      return !this.$v.user.$anyError;
    },
    validateEmail() {
      if (!this.ensureUser()) {
        return;
      }
      this.user.mode = 1;
      this.userValidate(this.user).then(() => {
        this.$bvToast.toast('已发送验证码，请将邮箱发送的验证码输入以完成注销验证', {
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
    deleteUser() {
      if (!this.ensureUser()) {
        return;
      }
      this.$confirm('此操作将永久注销用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(() => {
        this.user.phone = this.userInfo.phone;
        this.deleteMe(this.user).then(() => {
          this.$message({
            type: 'success',
            message: '注销成功!',
          });
          setTimeout(() => this.logout(), 2000);
        }).catch((err) => {
          if (err.response && err.response.data && err.response.data.msg) {
            this.$bvToast.toast(err.response.data.msg, {
              title: '数据验证错误',
              variant: 'danger',
              solid: true,
            });
          } else {
            console.error(err);
          }
        });
      }, () => {
        this.$message({
          type: 'info',
          message: '已取消注销',
        });
      }).catch(() => {
        this.$bvToast.toast({
          title: '出错了',
          variant: 'primary',
          solid: true,
        });
      });
    },
  },
};
</script>
