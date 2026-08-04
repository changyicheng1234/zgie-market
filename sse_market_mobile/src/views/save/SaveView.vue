<!-- eslint-disable max-len -->
<!-- eslint-disable vue/no-use-v-if-with-v-for -->
<!-- eslint-disable vuejs-accessibility/click-events-have-key-events -->
<!-- eslint-disable vuejs-accessibility/alt-text -->
<template>
  <div margin='0'>
    <van-nav-bar
      title="我的收藏"
      left-text="返回"
      left-arrow
      @click-left="$router.push({name:'profile'})"
      fixed
      placeholder
      safe-area-inset-top
    />
    <post-list ref="PostList" searchsort="save"></post-list>
  </div>
</template>

<script>
import PostList from '@/components/PostList.vue';

export default {
  components: {
    PostList,
  },
};
</script>

<style scoped>
@import '@/style/css/PostList.css';
</style>
