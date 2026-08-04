<!-- eslint-disable max-len -->
<!-- eslint-disable vue/no-use-v-if-with-v-for -->
<!-- eslint-disable vuejs-accessibility/click-events-have-key-events -->
<!-- eslint-disable vuejs-accessibility/alt-text -->
<template>
  <div style="margin-bottom: 115px">
    <van-overlay :show="showloading" duration="0.1" :custom-style="{background:'rgba(0,0,0,0.1)'}">
      <van-loading class="wrapper" size="24px" vertical>加载中...</van-loading>
    </van-overlay>
    <div class='container'>
      <h2>SSE_market</h2>
      <div class='icon-container'>
        <van-icon size='27' name='add' @click='goToPost'/>
      </div>
      <div class='search-container'>
        <van-search
          v-model='searchinfo'
          show-action
          show-action-slot
          placeholder='搜索  支持内容'
          @search='onSearch'
          shape='round'
        >
          <template #action>
            <div
                class='custom-search-button'
                @click='onSearch'
                @keydown.enter='onSearch'
                tabindex='0'
            >
              搜索
            </div>
          </template>
        </van-search>
      </div>
    </div>
    <div class='list'>
      <div
          class='item'
          v-for='(card, index) in cards'
          :key='index'
          @click='$refs.PostList.partitionBrowse(card.title)'
      >
        <div>
          <svg class='icon' aria-hidden='true'>
            <use v-bind:xlink:href='getGoodsHref(card.icon)'></use>
          </svg>
        </div>
        <!-- <p v-html='line'></p> -->
        <div class='title'>
          {{ card.title }}
        </div>
      </div>
    </div>
    <post-list ref="PostList"></post-list>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import PostList from '../../components/PostList.vue';

export default {
  components: {
    PostList,
  },
  data() {
    return {
      searchinfo: '',
      loading: false, // 是否正在加载帖子
      showloading: false, // 显示加载动画
      cards: [
        {
          title: '全部',
          icon: 'quanbu',
        },
        {
          title: '求职',
          icon: 'a-bianzu90',
        },
        {
          title: '日常吐槽',
          icon: 'pinglunqu',
        },
        {
          title: '打听求助',
          icon: 'wenhao',
        },
        {
          title: '二手闲置',
          icon: 'a-gouwugouwudai',
        },
        {
          title: '恋爱交友',
          icon: 'xin-qinglv',
        },
        {
          title: '学习交流',
          icon: 'wodexuexi',
        },
        {
          title: '其他',
          icon: 'qita',
        },
      ],
    };
  },
  created() {
    // this.showloading = true;
  },
  activated() {
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
  },
  methods: {

    getGoodsHref(iconName) {
      return `#icon-${iconName}`;
    },
    goToPost() {
      this.$router.push({ path: '/post' });
    },
    onSearch() {
      // 记得在进行切换加载操作的时候currentPage要恢复为1
      this.$refs.searchinfo = this.searchinfo;
      this.$refs.PostList.currentPage = 1;
      this.$refs.PostList.partitionBrowse('');
      // this.searchinfo = '';
    },
  },
};
</script>

<style scoped>
.list {
  padding: 20px;
  justify-items: center;
  display: grid;
  grid-gap: 40px;
  grid-template-columns: 1fr 1fr 1fr;
}

.list .item {
  justify-items: center;
  width: 150px;
  height: 150px;
  /* background: skyblue; */
}

.list .item .title {
  font-size: 25px;
}

.icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

</style>
