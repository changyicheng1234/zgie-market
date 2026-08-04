<template>
  <div>
    <!-- 标题 -->
    <div class="jumbotron jumbotron-fluid">
      <div class="container">
        <h2>欢迎来到课程专区</h2>
        <p class="content-text">课程专区，以老师为单位，以学生为主体</p>
        <p class="content-text">请在下拉搜索选择框选择相应的教师</p>
        <van-field
          readonly
          clickable
          label="教师"
          :value="label"
          placeholder="请选择教师"
          @click="showPicker = true"
        />
        <van-popup v-model="showPicker" round position="bottom">
          <van-picker
            show-toolbar
            :columns="column"
            @cancel="showPicker = false"
            @confirm="handlechange"
          />
        </van-popup>
      </div>
    </div>
      <van-row type="flex" justify="center" style="margin: 5px">
<!--          <van-col v-if="userInfo.identity === 'teacher'">-->
<!--              <van-button type="primary"-->
<!--                          style="font-size: 30px" href="/filemanage">进入课件管理</van-button>-->
<!--          </van-col>-->
          <van-col>
              <van-button type="primary"
                          style="font-size: 30px" @keydown.enter="abc"
                          @click="goToFileManage">课件资料</van-button>
          </van-col>
      </van-row>
      <post-list :key="postListKey" :tag="label" class="border-primary mt-5" />
  </div>
</template>

<script>

import { mapState } from 'vuex';
import store from '@/store';
import request from '@/utils/request';
import PostList from '@/components/PostList.vue';

export default {
  data() {
    return {
      showPicker: false,
      postListKey: 1,
      options: [],
      column: [],
      value: '',
      label: '',
    };
  },
  components: {
    PostList,
  },
  created() {
    store.commit('getNotices');
    localStorage.setItem('Partition', JSON.stringify('课程专区'));
    this.getTag();
    if (this.userInfo.identity === 'teacher') {
      this.value = this.userInfo.email;
    }
  },
  beforeRouteLeave(to, from, next) {
    // 返回上一页面时清空本地缓存
    localStorage.removeItem('Partition');
    next();
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
  },
  methods: {
    handlechange(value) {
      this.label = value;
      this.$notify({ type: 'success', message: `切换成功to${this.label}` });
      this.showPicker = false;
      this.postListKey += 1;
    },
    // 下面的函数不知道为什么算不出来，看来大道至简
    // 昨晚搞到将近凌晨1点，然后在床上还想这个东西，突然悟到了就是forEach循环里面的return其实并不能传到给getLabel函数，
    // 所以forEach函数里面的return其实并没有起到作用，getLabel还是返回的‘’
    getLabel(value) {
      this.options.forEach(
        // eslint-disable-next-line consistent-return
        (item) => {
          if (item.value === value) {
            // this.$message.success(item.label);
            return item.label;
          }
        },
      );
      return '';
    },
    getTag() {
      request.get('/auth/getTags?type=course')
        .then((res) => {
          this.options = res.data.data.tags.map((tag) => ({
            label: tag.Name,
            value: tag.Value,
          }));
          this.column = this.options.map((option) => option.label);
        });
    },

    // goback() {
    //   this.$router.replace({ name: 'partitions' });
    // },
    goToFileManage() {
      this.$router.push({
        name: 'filemanage',
        params: {
          email: this.options.find((option) => option.label === this.label).value,
        },
      });
    },
  },
};

</script>
<style>
.content-text{
  font-size:25px
}
</style>
