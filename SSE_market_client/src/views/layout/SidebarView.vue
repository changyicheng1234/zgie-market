<template>
  <!--   侧边栏 -->
      <b-list-group flush class="list-group" :style="{
        'margin-top': '100px', 'position': 'sticky', 'top': '100px',
        'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
        'border': isNightStyle ? '1px solid rgb(50, 50, 50)' :
          '1px solid rgb(237, 235, 235)', 'height': 'fit-content'
      }" v-if="this.$route.name != 'register' && this.$route.name != 'modifyPassword'
                && this.$route.name != 'deleteMe' && this.$route.name !== 'login'">
        <!--     主页   -->
        <b-list-group-item to="/" :style="{ 'background-color': $route.path !== '/' ?
            isNightStyle ? 'rgb(70, 70, 70)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px; align-items: center;"
                           @click="togglePartitions()" class="click" ref="parent">
          <b-icon-house-fill class="mr-3"></b-icon-house-fill>
          主页
          <b-icon-caret-right-fill v-if="!showPartitions" style="margin-left: auto;">
          </b-icon-caret-right-fill>
          <b-icon-caret-down-fill v-if="showPartitions" style="margin-left: auto;">
          </b-icon-caret-down-fill>
        </b-list-group-item>
        <!--     分区   -->
        <b-list-group-item v-if="showPartitions" to="/partitions" class="click"
                           :style="{ 'background-color': $route.path !== '/partitions' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;">
          <b-icon-grid1x2-fill class="mr-3"></b-icon-grid1x2-fill>
          分区选择
        </b-list-group-item>
        <!-- 课程专区 -->
        <b-list-group-item
          to="/course"
          class="click"
          :style="{ 'background-color': $route.path !== '/course' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
          style="font-size: 18px;">
          <b-icon-people-fill class="mr-3"></b-icon-people-fill>
          课程专区
        </b-list-group-item>
        <!--     发帖   -->
        <b-list-group-item @click="showPostForm()"
                           :style="{ 'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
              'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;"
                           class="click">
          <b-icon-pencil-fill class="mr-3">
          </b-icon-pencil-fill>
          发帖
        </b-list-group-item>
        <b-modal v-model='PostFormVisible' title='发帖' ok-only ok-title="取消发帖"
                 modal-class="custom-modal">
          <PostForm :mode="'post'"/>
        </b-modal>
        <!--     通知   -->
        <b-list-group-item @click="toNotifications()"
                           :style="{ 'background-color': $route.path !== '/notice' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null }"
                           style="font-size: 18px;" class="click">
          <b-icon-bell-fill class="mr-3"></b-icon-bell-fill>
          通知
          <span v-if="noticesNum!=0"
                class="badge badge-danger ml-2 pop">{{ noticesNum }}</span>
        </b-list-group-item>
        <!--     反馈   -->
        <b-list-group-item @click="toFeedback()"
                           :style="{ 'background-color': $route.path !== '/feedback' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null
            : isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;" class="click">
          <b-icon-envelope-fill class="mr-3"></b-icon-envelope-fill>
          反馈
        </b-list-group-item>
        <b-modal v-model='FeedbackVisible' title='反馈' ok-only ok-title="取消反馈"
                 modal-class="custom-modal">
        </b-modal>
        <!--     我的   -->
        <b-list-group-item
          :style="{ 'font-size': '18px', 'align-items': 'center',
        'background-color': showProfiles ? 'rgb(245, 245, 245)' : '' ,
        'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
        'color': isNightStyle ? 'gray' : null}"
          @click="toggleProfiles" class="click">
          <b-icon-person-circle class="mr-3"></b-icon-person-circle>
          我的
          <b-icon-caret-right-fill v-if="!showProfiles" style="margin-left: auto;">
          </b-icon-caret-right-fill>
          <b-icon-caret-down-fill v-if="showProfiles" style="margin-left: auto;">
          </b-icon-caret-down-fill>
        </b-list-group-item>
        <!--     个人信息   -->
        <b-list-group-item class="click" v-if="showProfiles" to="/profile"
                           :style="{ 'background-color': $route.path !== '/profile' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;">
          <b-icon-table class="mr-3"></b-icon-table>
          个人信息
        </b-list-group-item>
        <!--我的收藏-->
        <b-list-group-item v-if="showProfiles" to="/save" class="click"
                           :style="{ 'background-color': $route.path !== '/save' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;">
          <b-icon-star-fill class="mr-3"></b-icon-star-fill>
          我的收藏
        </b-list-group-item>
        <!--发帖历史-->
        <b-list-group-item v-if="showProfiles" to="/history" class="click"
                           :style="{ 'background-color': $route.path !== '/history' ?
            isNightStyle ? 'rgb(50, 50, 50)' : null :
            isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)',
            'color': isNightStyle ? 'gray' : null}"
                           style="font-size: 18px;">
          <b-icon-clock-fill class="mr-3"></b-icon-clock-fill>
          发帖历史
        </b-list-group-item>
        <!--设置-->
        <b-list-group-item
          :style="{ 'font-size': '18px', 'align-items': 'center',
          'background-color': showSettings ? 'rgb(245, 245, 245)' : '' ,
          'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
          'color': isNightStyle ? 'gray' : null}"
          @click="toggleSettings" class="click">
          <b-icon-gear-fill class="mr-3"></b-icon-gear-fill>
          设置
          <b-icon-caret-right-fill v-if="!showSettings" style="margin-left: auto;">
          </b-icon-caret-right-fill>
          <b-icon-caret-down-fill v-if="showSettings" style="margin-left: auto;">
          </b-icon-caret-down-fill>
        </b-list-group-item>
        <!--修改密码-->
        <b-list-group-item v-if="showSettings" class="click"
                           @click="toModifyPassword()" style="font-size: 18px;"
                           :style="{ 'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
            'color': isNightStyle ? 'gray' : null}">
          <b-icon-lock-fill class="mr-3"></b-icon-lock-fill>
          修改密码
        </b-list-group-item>
        <!--注销账号-->
        <b-list-group-item v-if="showSettings" class="click"
                           @click="toDeleteUser()" style="font-size: 18px;"
                           :style="{ 'background-color': isNightStyle ? 'rgb(50, 50, 50)' : null,
            'color': isNightStyle ? 'gray' : null}">
          <b-icon-x-circle-fill class="mr-3"></b-icon-x-circle-fill>
          注销账号
        </b-list-group-item>
      </b-list-group>
</template>

<script>
import PostForm from '@/components/PostForm.vue';
import { mapActions, mapState } from 'vuex';
import store from '@/store';

export default {
  components: {
    PostForm,
  },
  name: 'SidebarView',
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
    noticesNum() {
      return store.state.noticesNum;
    },
    style() {
      return store.state.style;
    },
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
  },
  data() {
    return {
      PostFormVisible: false,
      FeedbackVisible: false,
      showPartitions: false,
      showSettings: false,
      showProfiles: false,
      searchinfo: '',
      scrollPosition: 0,
      toLogin: false,
    };
  },
  created() {
    if (typeof localStorage !== 'undefined') {
      if (!localStorage.getItem('Style')) {
        localStorage.setItem('Style', JSON.stringify('day'));
      }
    }
  },
  mounted() {
    window.addEventListener('scroll', this.handleScroll);
  },
  destroyed() {
    window.removeEventListener('scroll', this.handleScroll);
  },
  methods: {
    // 使用map将映射'store/module'里的logout函数
    ...mapActions('userModule', ['logout']),
    showPostForm() {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      this.PostFormVisible = true;
    },
    toNotifications() {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      if (this.$route.path !== '/notice') {
        this.$router.replace({ name: 'notice' });
      }
    },
    toFeedback() {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      if (this.$route.path !== '/feedback') {
        this.$router.replace({ name: 'feedback' });
      }
    },
    toModifyPassword() {
      const { href } = this.$router.resolve({ path: '/modifyPassword' });
      window.open(href, '_blank');
    },
    toDeleteUser() {
      const { href } = this.$router.resolve({ path: '/deleteMe' });
      window.open(href, '_blank');
    },
    handleScroll() {
      this.scrollPosition = window.scrollY;
    },
    toggleSettings() {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      this.showSettings = !this.showSettings;
    },
    toggleProfiles() {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      this.showProfiles = !this.showProfiles;
    },
    togglePartitions() {
      this.showPartitions = !this.showPartitions;
    },
    refreshPageAndNavigate() {
      this.$router.push({
        name: 'home',
        query: { searchinfo: this.searchinfo },
      });
      this.$router.go(0);
    },

  },
};
</script>

<style scoped>
@import '../../style/css/NavbarView.css';
</style>
