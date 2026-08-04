<template>
  <div class='container'>
    <div class='row'>
      <div class='col-md-4' v-for='(card, index) in cards' :key='index'>
        <div class='card'
             :style="{ 'background-color': isNightStyle ? 'rgb(25,25,25)' : 'white'
             ,'color': 'grey'}">
          <img :src="card.imageUrl" class="card-img-top" alt=".."/>
          <div class='card-body'>
            <h5 class='card-title'>{{ card.title }}</h5>
            <p class='card-text'>{{ card.description }}</p>
            <a :href='card.linkUrl'
               :class="isNightStyle ? 'btn btn-outline-warning':'btn btn-primary'"
               @click="$router.push({ name: 'home', query: { partitions: card.title } })">
              {{ card.linkLabel }}</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      cards: [
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/daily.jpg',
          title: '日常吐槽',
          description: '言短情长，抒发心声。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/help.jpg',
          title: '打听求助',
          description: '疑虑相询，忧难互助。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/love.jpg',
          title: '恋爱交友',
          description: '情缘相连，同窗结友。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/study.jpg',
          title: '学习交流',
          description: '知识碰撞，学术思辨。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/secondhand.jpg',
          title: '二手闲置',
          description: '物尽其用，旧物新归。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/else.jpg',
          title: '其他',
          description: '畅所欲言，创意无限。',
          linkLabel: '跳转分区',
        },
        {
          imageUrl: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/recruit.jpg',
          title: '求职招募',
          description: '梦想职路，共创未来。若有需要投稿，请投稿到邮箱:zhouyt69@mail.sysu.edu.cn',
          linkLabel: '跳转分区',
        },
      ],
    };
  },
  methods: {
  },
  computed: {
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
  },
};
</script>

<style scoped>
.container {
  margin-top: 30px;
}

.card {
  margin-bottom: 30px;
}
</style>
