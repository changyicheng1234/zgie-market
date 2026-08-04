<template>
  <div>
    <v-md-editor
      v-model="editorValue"
      :disabled-menus="[]"
      height="calc(100vh - 150px)"
      @upload-image="handleUploadImage"
    />
  </div>
</template>

<script>
export default {
  name: 'MdEditor',
  model: {
    prop: 'content',
    event: 'update',
  },
  props: {
    // 接收父组件传递的值
    content: String,
  },
  data() {
    return {
      editorValue: this.content != null ? this.content : '',
    };
  },
  watch: {
    content(newContent) {
      this.editorValue = newContent;
      console.log(`editor:${this.editorValue}`);
    },
    editorValue(newValue) {
      this.$emit('update', newValue);
    },
  },
  methods: {
    // v-md-editor 文件上传
    handleUploadImage(event, insertImage, files) {
      // 上传
      files.forEach((file) => {
        this.crud.upload(file, 'image/vMdEditor/').then((res) => {
          // 获取返回数据
          const { url, name } = res.data.data;
          // 添加图片到内容
          insertImage({
            url,
            desc: name,
          });
        });
      });
    },
  },
};
</script>

<style>
</style>
