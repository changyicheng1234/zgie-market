<template>
  <div>
<!--    {{ this.options }}-->
<!--    {{ userInfo }}-->
<!--    value：{{ value }}-->
    <!-- 标题 -->
    <div class="jumbotron jumbotron-fluid">
      <div class="container">
        <h1 class="display-4">欢迎来到课程专区</h1>
        <p class="lead">课程专区，以老师为单位，以学生为主体</p>
        <p class="lead">请在下拉搜索选择框选择相应的教师</p>
        <div>
          <el-select v-model="value" filterable
                     @change="handlechange"
                     placeholder="请选择教师">
            <el-option
              v-for="item in options"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </div>
      </div>
    </div>
    <div class="card bg-light border-primary" style="width: fit-content">
      <div class="card-body">
        <div v-if="userInfo.identity=='teacher'">
          <a class="card-link" style="font-size: 30px" href="/filemanage">进入课件管理</a>
        </div>
        <div v-else>
          <a class="card-link"
             style="font-size: 30px"
             @keydown.enter="abc"
             @click="goToFileManage">课件资料</a>
        </div>
      </div>
    </div>
    <PostList
      :key="postListKey"
      :tag ="label"
      class="border-primary mt-5"/>
  </div>
</template>

<script>

import PostList from '@/components/PostList.vue';
import store from '@/store';
import { mapState } from 'vuex';
import request from '@/utils/request';

export default {
  data() {
    return {
      postListKey: 1,
      options: [],
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
      // this.label = this.getLabel(value);
      this.label = this.options.find((item) => item.value === value).label;
      this.$message.success(`切换成功to${this.label}`);
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
        });
    },

    // goback() {
    //   this.$router.replace({ name: 'partitions' });
    // },
    goToFileManage() {
      this.$router.push({
        name: 'filemanage',
        params: {
          email: this.value,
        },
      });
    },
  },
};

</script>
