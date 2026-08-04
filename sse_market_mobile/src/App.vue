<template>
  <div id="app">
    <div id="main">
      <keep-alive>
      <router-view v-if="this.$route.meta.keepAlive"></router-view>
    </keep-alive>
    <router-view v-if="!this.$route.meta.keepAlive"></router-view>
    <navbar />
    </div>
  </div>
</template>

<script>
import Navbar from './views/layout/NavbarView.vue';

export default {
  components: { Navbar },
};
</script>

<style lang="scss">
// 对float/fixed元素要额外配置样式
.van-popup,.van-nav-bar,.van-tabbar {
max-width: 600px;
left:auto;
}
@media screen and (max-width: 600px) {
  .van-popup,
  .van-nav-bar {
    left: 0;
  }
}

// @media screen and (max-width: 550px) {
//   .van-popup,
//   .van-nav-bar {
//     left: auto;
//   }
// }
#app {
  display: flex;
  justify-content: center;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
#main {
  max-width: 600px;
  flex: 1;
}
nav {
  padding: 30px;
  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
