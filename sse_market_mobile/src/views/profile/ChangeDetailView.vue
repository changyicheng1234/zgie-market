<template>
  <div>
    <van-uploader :after-read='afterRead' :max-count='1' accept='image/gif, image/jpeg, image/png'>
      <div class='avatar-container'>
        <img v-if='imageUrl' :src='imageUrl' class='profile_avatar' alt='' />
        <i v-else class='van-icon avatar'></i>
      </div>
    </van-uploader>

    <div class='information_box'>
      <van-cell-group>
        <van-field class='field-spacing' label='用户名' v-model='userInfo.name' />
        <van-field
          class='field-spacing'
          label='手机号'
          v-model='userInfo.phone'
          value='输入框已禁用'
          disabled
        />
        <van-field
          class='field-spacing'
          label='邮箱'
          v-model='userInfo.email'
          value='输入框已禁用'
          disabled
        />
        <van-field
          class='field-spacing'
          label='用户等级积分'
          :value="userInfo.title + ' ' + userInfo.score + '分'"
          disabled
        />
        <van-field
          class='field-spacing'
          v-model='userInfo.intro'
          rows='2'
          autosize
          label='简介'
          type='textarea'
          maxlength='100'
          placeholder='请输入留言'
          show-word-limit
        />
      </van-cell-group>
    </div>
    <van-button type='primary' @click='goback' style='margin-right: 20px'>返回主页</van-button>
    <van-button @click='updateUserInfo' type='info'>确定修改</van-button>
  </div>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import { Dialog } from 'vant';
// eslint-disable-next-line import/no-cycle
import request from '@/utils/request';

export default {
  computed: mapState({
    user: (state) => state.userModule.userInfo,
  }),
  data() {
    return {
      imageUrl: '',
      userTelephone: '',
      userInfo: {
        userID: '',
        phone: '',
        email: '',
        name: '',
        num: '',
        intro: '',
        score: '',
        title: '',
        ban: '',
        punishnum: '',
        avatarURL: '',
      },
    };
  },
  created() {
    this.getInfo({ phone: this.user.phone })
      .then((res) => {
        this.userInfo.userID = res.data.userID;
        this.userInfo.name = res.data.name;
        this.userInfo.num = res.data.num;
        this.userInfo.phone = res.data.phone;
        this.userInfo.email = res.data.email;
        this.userInfo.intro = res.data.intro;
        this.userInfo.score = res.data.score;
        this.userInfo.title = this.getUserTitle(res.data.score);
        this.userInfo.ban = res.data.ban;
        this.userInfo.punishnum = res.data.punishnum;
        this.userInfo.avatarURL = res.data.avatarURL;
        this.imageUrl = this.userInfo.avatarURL;
      })
      .catch((err) => {
        console.error(err);
      });
  },
  methods: {
    ...mapActions('userModule', ['logout']),
    ...mapActions('userModule', { getInfo: 'getInfo' }),
    afterRead(file) {
      const formData = new FormData();
      formData.append('file', file.file); // 'file' 必须与后端接收的字段名相匹配
      formData.append('phone', this.userInfo.phone);
      request
        .post('auth/updateAvatar', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
        .then((response) => {
          console.log('上传成功', response.data);
          window.location.reload();
        })
        .catch((error) => {
          console.error('上传失败', error);
          this.$bvToast.toast('上传失败', {
            title: '系统提醒',
            variant: 'danger',
            solid: true,
          });
          // 可以在此处处理上传失败的逻辑
        });
    },
    getUserTitle(userScore) {
      if (userScore < 100) {
        return '菜鸟';
      }
      if (userScore >= 100 && userScore < 300) {
        return '大虾';
      }
      if (userScore >= 300 && userScore < 600) {
        return '码农';
      }
      if (userScore >= 600 && userScore < 1000) {
        return '程序猿';
      }
      if (userScore >= 1000 && userScore < 2000) {
        return '工程师';
      }
      if (userScore >= 2000 && userScore < 3000) {
        return '大牛';
      }
      if (userScore >= 3000 && userScore < 4000) {
        return '专家';
      }
      if (userScore >= 4000 && userScore < 5000) {
        return '大神';
      }
      return '祖师爷';
    },
    updateUserInfo() {
      if (this.userInfo.name.trim() === '') {
        Dialog.alert({
          title: '保存失败',
          message: '用户名不能为空！',
          confirmButtonText: '确定',
        });
        return; // 如果用户名为空，则返回，不执行以下的保存逻辑
      }
      const updatedUserInfo = {
        userID: this.userInfo.userID,
        name: this.userInfo.name,
        num: this.userInfo.num,
        intro: this.userInfo.intro,
        avatarURL: this.userInfo.avatarURL,
      };
      request
        .post('auth/updateUserInfo', updatedUserInfo)
        .then((response) => {
          // 处理更新成功的逻辑
          console.log(response);
          console.log('用户信息更新成功');
          Dialog.alert({
            title: '保存成功',
            message: '用户信息已成功更新！',
            confirmButtonText: '确定',
          });
        })
        .catch((error) => {
          // 处理更新失败的逻辑
          console.error('用户信息更新失败', error);
        });
    },
    goback() {
      this.$router.push({ path: '/profile' });
    },
  },
};
</script>

<style>
.information_box {
  display: flex;
  flex-direction: row;
  padding-left: 10%;
  padding-right: 10%;
  border: none;
  margin-top: 15px;
  margin-bottom: 15px;
}
.avatar-container {
  width: 200px;
  height: 200px;
  border-radius: 10%;
  background-color: #eee;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  margin-top: 30px;
}

.profile_avatar, .van-icon.van-icon-add-o.avatar {
  width: 100%;
  height: 100%;
  object-fit: fill; /* 使图片填满容器 */
}

.field-spacing {
  margin-left: 10px;
  margin-right: 10px;
}
</style>
