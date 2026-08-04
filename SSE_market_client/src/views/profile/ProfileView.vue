<template>
  <div>
    <h4>个人信息</h4>
    <el-upload
        class="avatar-uploader"
        :action="uploadAvatarActionURL"
        :show-file-list="false"
        :on-success="handleAvatarSuccess"
        :before-upload="beforeAvatarUpload"
    >
      <img v-if="imageUrl" :src="imageUrl" class="avatar" alt="Uploaded Avatar">
      <i v-else class="el-icon-plus avatar-uploader-icon"></i>
    </el-upload>
    <div class="mt-5">
      <el-form :model="userInfo" @submit.native.prevent="updateUserInfo">
        <el-row>
          <el-col :span="5" class="field-col">
            <el-form-item label="用户名">
              <el-input v-model="userInfo.name" required></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="7" class="field-col">
            <el-form-item label="手机">
              <el-input v-model="userInfo.phone" required :disabled="true"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="11" class="field-col">
            <el-form-item label="邮箱">
              <el-input v-model="userInfo.email" required :disabled="true"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="用户等级积分">
          <el-input :value="userInfo.title + ' ' + userInfo.score + '分'"
          required :disabled="true" ></el-input>
        </el-form-item>
        <el-form-item label="个人简介">
          <el-input type="textarea" v-model="userInfo.intro"></el-input>
        </el-form-item>
        <el-button type="primary" native-type="submit">保存</el-button>
      </el-form>
    </div>
    <el-button type="text" @click="open">点击查看更多信息</el-button>
  </div>
</template>
<script>
import { mapActions, mapState } from 'vuex';
import request from '@/utils/request';

export default {
  computed: mapState({
    user: (state) => state.userModule.userInfo,
  }),
  data() {
    return {
      uploadAvatarActionURL: 'https://ssemarket.cn/api/auth/uploadAvatar',
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
    open() {
      if (this.userInfo.punishnum !== 0) {
        this.$alert(`解除封禁日期：${this.userInfo.ban}，违规次数：${this.userInfo.punishnum}次\n`, '违规信息', {
          confirmButtonText: '确定',
        });
      } else {
        this.$alert('您没有违规记录，请继续保持~', '违规信息', {
          confirmButtonText: '确定',
        });
      }
    },
    handleAvatarSuccess(res) {
      this.imageUrl = res.data.avatar_url;
      this.userInfo.avatarURL = this.imageUrl;
    },
    beforeAvatarUpload(file) {
      const isLt2M = file.size / 1024 / 1024 < 5;

      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过 5MB!');
      }
      return isLt2M;
    },
    ...mapActions('userModule', { getInfo: 'getInfo' }),
    updateUserInfo() {
      const updatedUserInfo = {
        userID: this.userInfo.userID,
        name: this.userInfo.name,
        num: this.userInfo.num,
        intro: this.userInfo.intro,
        avatarURL: this.userInfo.avatarURL,
      };

      request.post('auth/updateUserInfo', updatedUserInfo)
        .then((response) => {
          // 处理更新成功的逻辑
          console.log(response);
          console.log('用户信息更新成功');
          this.$bvToast.toast('用户信息更新成功', {
            title: '系统提醒',
            variant: 'primary',
            solid: true,
          });
          // Commit updated user information to Vuex
          this.updateLocalUserInfo(updatedUserInfo);
        })
        .catch((error) => {
          // 处理更新失败的逻辑
          console.error('用户信息更新失败', error);
        });
    },
    ...mapActions('userModule', { getInfo: 'getInfo', updateLocalUserInfo: 'updateLocalUserInfo' }),
  },
};
</script>

<style>
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409EFF;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
.field-col {
  margin-left: 10px;
}
</style>
