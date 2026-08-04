<template>
  <div>
    <van-overlay :show="showloading" duration="0.1" :custom-style="{background:'rgba(0,0,0,0.1)'}">
      <van-loading class="wrapper" size="24px" vertical>加载中...</van-loading>
    </van-overlay>
    <van-nav-bar title="课件资料"
                 left-text="返回"
                 right-text="刷新"
                 @click-left="goBack" @click-right="getFileList(now_path)"/>
    <van-cell>当前路径:{{ now_path }}</van-cell>
    <van-list>
      <van-cell class="border"
          v-for="(file, index) in fileList"
          :key="index"
          :title="file.Name"
          :is-link="file.Type === 1"
          @click="file.Type === 0 ? downloadFile(file.key) : openFolder(file.key)"
      >
        <template #label>
          <div>
            <span v-if="file.EditTime!==''">{{ formatDate(file.EditTime) }}</span>
            <span v-if="file.Size > 0"> | {{ formatFileSize(file.Size) }}</span>
          </div>
        </template>
      </van-cell>
    </van-list>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import request from '@/utils/request';

export default {
  data() {
    return {

      fileList: [],
      showloading: false,
      now_path: '',
    };
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
  },
  methods: {
    formatDate(date) {
      // 格式化日期时间
      const d = new Date(date);
      return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${String(
        d.getHours(),
      )
        .padStart(2, '0')}:${String(d.getMinutes())
        .padStart(2, '0')}`;
    },
    goBack() {
      this.$router.back();
    },
    downloadFile(key) {
      this.showloading = true;
      const formData = new FormData();
      formData.append('filename', key);
      return request.post('/auth/getObjectUrl', formData)
        .then((res) => {
          // console.log(res)
          // window.location.assign(res.data.data.url);
          window.location.href = res.data.data.url;
          // 用winodow.open()默认在苹果中会打不开
          // window.open(res.data.data.url);
          this.showloading = false;
        })
        .catch((err) => {
          this.$notify(`文件打开失败${err}`);
          this.showloading = false;
        });
      // 发起文件下载请求，key为文件的标识
      // 你可以使用浏览器的下载功能或调用后端接口来实现文件下载
      // 例如：window.location.href = `/api/download?key=${key}`;
    },
    openFolder(key) {
      this.getFileList(key);
    },
    getFileList(path) {
      console.log('getFileList()');
      this.showloading = true;
      const formData = new FormData();
      formData.append('path', path);
      request.post('/auth/getFileList', formData)
        .then((data) => {
          this.fileList = data.data || [];
          this.now_path = path;
          // this.$toast.success('文件列表获取成功');
          this.showloading = false;
        })
        .catch((error) => {
          console.error(`getFileList error${error}`);
          this.$toast.fail(`文件列表获取失败${error}`);
          this.showloading = false;
        });
    },
    formatFileSize(size) {
      const number = size * 1024;
      // 格式化文件大小，你可以根据实际需求编写该函数
      if (number < 1024) return `${number} B`;
      if (number < 1024 * 1024) return `${(number / 1024).toFixed(2)} KB`;
      return `${(number / (1024 * 1024)).toFixed(2)} MB`;
    },
  },
  mounted() {
    // 获取文件列表数据，你可以从后端API获取数据并填充到fileList
    // 例如：this.fileList = fetchDataFromApi();
    // this.getFileList(this.now_path);
  },
  created() {
    const parts = this.$route.params.email.split('@');
    if (parts[0].length === 0) {
      this.$toast.fail('请选择教师');
    }
    this.now_path = `course/${parts[0]}/`;
    this.getFileList(this.now_path);
  },
};
</script>

<style scoped>
/* 在这里添加页面样式 */
.border {
  border-top: rgb(128, 128, 128) 1px solid;

}
</style>
