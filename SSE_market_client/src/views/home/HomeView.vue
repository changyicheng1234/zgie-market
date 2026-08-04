<template>
  <PostList />
</template>

<script>

import PostList from '@/components/PostList.vue';
import store from '@/store';

export default {
  components: {
    PostList,
  },
  created() {
    store.commit('getNotices');
  },
  beforeRouteLeave(to, from, next) {
    // 返回上一页面时清空本地缓存
    localStorage.removeItem('Partition');
    next();
  },
  methods: {
    // goback() {
    //   this.$router.replace({ name: 'partitions' });
    // },
  },
};

</script>
