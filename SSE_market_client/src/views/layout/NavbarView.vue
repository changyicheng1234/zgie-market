<template>
  <div style="width: 100%;">
    <b-modal v-model='toLogin' title='登录' ok-only ok-title="取消登录"
             modal-class="custom-modal">
      <LoginView />
    </b-modal>
    <!--    电脑端头部（欢迎界面与搜索框）-->
    <div class="video-container" v-if="$route.name == 'home' && (!this.$route.query.partitions ||
      this.$route.query.partitions == '主页')" style="margin-bottom: -80px;">
      <div style="background-color: black; height: 400px;">
        <video autoplay loop muted playsinline>
          <source
            src="https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/vedio/Background.mp4"
            type="video/mp4">
        </video>
      </div>
      <div class="overlay">
        <vue-typed-js :showCursor="false"
                      :strings="['这是中山大学', '欢迎来到软工集市','今日火文等你看','SSE_market']">
          <h1 class="typing">SSE_market</h1>
        </vue-typed-js>
        <div class="search-box">
          <b-form-input placeholder="搜索  支持标签和内容" v-model="searchinfo"
                        style="height: 50px;border-radius: 20px;">
          </b-form-input>
          <b-button variant="dark" @click="refreshPageAndNavigate" class="d-flex text-nowrap"
                    style="margin-left:10px;padding:20px;border-radius: 20px;">
            <b-icon-search class="mr-2"></b-icon-search>
            搜索
          </b-button>
        </div>
      </div>
    </div>
    <!-- 导航栏-->
    <b-navbar fixed="top" v-if="(this.$route.name !== 'home' && this.$route.name !== 'register'
      && this.$route.name !== 'login' && this.$route.name !== 'modifyPassword'
      && this.$route.name !== 'deleteMe') || scrollPosition > 400
      || (this.$route.query.partitions && this.$route.query.partitions != '主页')"
      :style="{ 'background-color': isNightStyle ? 'rgb(246, 155, 10)' : 'rgb(17, 167, 226)'}"
      class="d-flex justify-content-between">
      <!-- SSE_MARKET文字-->
      <b-navbar-brand>
        <b-icon-shop class="mr-3"></b-icon-shop>
        SSE_market
      </b-navbar-brand>
      <!--      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>-->
      <!-- 搜索栏-->
      <div>
        <b-input-group v-if="this.$route.name == 'home' &&
            (this.$route.query.partitions == '主页' || !this.$route.query.partitions)">
          <b-form-input style="width: 600px; border-radius: 5px; height: 50px;
              z-index: 100;"
                        placeholder="搜索" v-model="searchinfo"></b-form-input>
          <b-input-group-append>
            <b-button style="margin-left:20px; border-radius: 5px; pointer-events: auto;"
                      variant="dark"
                      @click="refreshPageAndNavigate">
              <b-icon-search class="mr-2"></b-icon-search>
              搜索
            </b-button>
          </b-input-group-append>
        </b-input-group>
      </div>
      <!--夜间模式和个人信息-->
      <div class="d-flex">
        <div>
          <StyleButton/>
        </div>
        <!-- 个人登录信息 -->
        <div class="ml-2">
          <b-nav-item-dropdown right v-if="userInfo" class="d-flex">
            <template #button-content >
              <a style="color: #e4e7ed">{{ userInfo.name }}
                <b-icon-person-fill>
                </b-icon-person-fill></a>
            </template>
            <b-dropdown-item @click="logout">
              <b-icon-box-arrow-left class="mr-1"></b-icon-box-arrow-left>
              退出登录
            </b-dropdown-item>
          </b-nav-item-dropdown>
        </div>
      </div>
    </b-navbar>
  </div>
</template>

<script>

import { mapState, mapActions } from 'vuex';
import LoginView from '@/views/login/LoginView.vue';
import StyleButton from '@/components/StyleButton.vue';
import store from '@/store';

export default {
  // 获取在浏览器缓存中的包含用户信息的token，userInfo中包含用户的name和telephone
  // 所以要从前端返回用户信息时，一般采用telephone(因为name不唯一，id又无法从前端直接获取)
  components: {
    // PostForm,
    LoginView,
    StyleButton,
  },
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
      window.open('/modifyPassword', '_blank');
    },
    toDeleteUser() {
      window.open('/deleteMe', '_blank');
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

<style>
@import '../../style/css/NavbarView.css';

.pop {
  animation: pop 0.5s ease-in-out;
}
</style>
