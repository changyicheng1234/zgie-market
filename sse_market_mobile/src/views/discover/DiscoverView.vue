<template>
  <div>
    <div>
      <!-- <van-swipe :autoplay="3000">
        <van-swipe-item v-for="item in imagebox" :key="item.id">
          <a :href="getWebsiteURL(item.id)" aria-label="链接到图片">
            <img :src="item.idView" alt="" class="full-width-image"/>
          </a>
        </van-swipe-item>
      </van-swipe> -->
    </div>
    <div class="hotpost-topic">
      <van-icon name="fire-o" color="#fc0000" />
      24小时火文
    </div>
    <div>
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model="loading"
          :finished="finished"
          @load="onLoad"
        >
          <van-cell v-for="(hotpost, index) in hotposts" :key="index"
            :class="['custom-cell', { 'red-text': index < 3 }]"
            @click="() => { showDetails(hotpost); updatebrowse(hotpost) }"
          >
            <span :class="{ 'red-text': index < 3 }">{{ index + 1 }}. {{ hotpost.title }}</span>
            <template #right-icon>
              <van-icon name="fire-o" :badge="Math.floor(hotpost.heat)">
              </van-icon>
            </template>
          </van-cell>
        </van-list>
      </van-pull-refresh>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';

export default {
  computed: mapState({
    userInfo: (state) => state.userModule.userInfo,
  }),
  data() {
    return {
      posts: [],
      hotposts: [],
      userTelephone: '',
      postID: '',
      title: '',
      heat: '',

      loading: false,
      finished: false,
      refreshing: false,

      imagebox: [
        { id: 0, idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel1.png' },
        { id: 1, idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel2.png' },
        { id: 2, idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel3.png' },
        { id: 3, idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel4.png' },
      ],
    };
  },
  created() {
    this.calculateheat();
  },
  methods: {
    ...mapActions('postModule', { postHots: 'calculateheat' }),
    ...mapActions('postModule', { postUpdateBrowse: 'updatebrowse' }),
    updatebrowse(post) {
      // 切换浏览状态
      const updatedPost = { ...post, browse: post.browse };
      // 更新浏览次数
      updatedPost.browse = post.browse + 1;
      // 用更新后的 post 对象替换原先的 post 对象
      this.posts.splice(this.posts.indexOf(post), 1, updatedPost);
      this.userTelephone = this.userInfo.phone;
      this.postID = post.id;
      // 请求
      this.postUpdateBrowse({
        userTelephone: this.userTelephone,
        postID: this.postID,
      })
        .then(() => {})
        .catch((err) => {
          console.error(err);
        });
    },
    calculateheat() {
      // 向后端请求热度前10的帖子
      this.postHots({
        postID: this.postID,
        title: this.title,
        heat: this.heat,
      })
        .then((hots) => {
          this.hotposts = hots.data.map((hot) => ({
            id: hot.PostID,
            title: hot.Title,
            heat: hot.Heat,
          }));
        })
        .catch((err) => {
          console.error(err);
        });
    },
    onLoad() {
      setTimeout(() => {
        if (this.refreshing) {
          this.hotposts = [];
          this.refreshing = false;
        }

        this.hotposts.slice(0, 10);
        this.loading = false;

        if (this.hotposts.length <= 10) {
          this.finished = true;
        }
      }, 1000);
    },
    async onRefresh() {
      this.finished = false;
      await this.calculateheat(); // 调用该方法重新获取数据
      this.refreshing = false;
    },
    showDetails(hotpost) {
      this.$router.push({
        name: 'postDetails',
        params: { partition: this.partition },
        query: {
          id: hotpost.id, title: hotpost.title, before: this.$route.name, partition: this.partition,
        },
      });
    },
    getWebsiteURL(id) {
      if (id === 0) {
        return 'https://ztjy.sysu.edu.cn/';
      }
      if (id === 1) {
        return 'https://www.sysu.edu.cn/news/ztjj2/gdzt/zt2023ncjgzhy.htm';
      }
      if (id === 2) {
        return 'https://www.sysu.edu.cn/xxg/zdjj1.htm';
      }
      if (id === 3) {
        return 'https://ddhzt.sysu.edu.cn/';
      }
      return ''; // 添加默认的返回值
    },
  },
};
</script>

<style scoped>
.hotpost-topic {
  text-align: center;
  font-size: 40px;
  margin: 30px 0;
  color: #fc0000;
}
.custom-cell {
    text-align: left; /* Align cell content to the left */
}
.red-text {
  color: rgb(255, 128, 0);
}
.full-width-image {
  width: 100%;
  height: auto;
  margin: 30px 0;
}

</style>
