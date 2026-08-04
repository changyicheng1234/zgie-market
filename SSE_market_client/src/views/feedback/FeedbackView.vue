<template>
  <div>
    <el-input
        type="textarea"
        placeholder="请输入反馈，如果有文件需要上传，请在下方上传"
        v-model="textarea"
        maxlength="50"
        show-word-limit
        style="margin-left: 10px;"
    >
    </el-input>
    <div style="padding-bottom: 15px;" />
    <el-upload
        class="upload-demo"
        style="margin-left: 10px;"
        :action="uploadZipActionURL"
        :on-remove="handleRemove"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        :before-upload="beforeUpload"
        multiple
        :limit="1"
        :on-exceed="handleExceed"
        :file-list="fileList"
        accept=".zip"
    >
      <el-button size="small" type="primary" >上传附件</el-button>
      <div slot="tip" class="el-upload__tip">只能上传zip文件,且不超过10MB</div>
    </el-upload>
    <div style="padding-bottom: 15px;" />
    <el-button type="primary" style="margin-left: 10px;" @click="submitfeedback">提交反馈</el-button>
  </div>
</template>

<script>
import { mapActions } from 'vuex';

export default {
  data() {
    return {
      uploadZipActionURL: 'https://ssemarker.cn:8080/api/auth/uploadZip',
      textarea: '',
      fileList: [],
      feedback: {
        ftext: '',
        attachment: null,
      },
    };
  },
  methods: {
    ...mapActions('postModule', { SubmitFeedback: 'submitFeedback' }),
    handleUploadSuccess(response, file, fileList) {
      this.fileList.push(response.zipURL);
      console.log(file, fileList);
    },
    handleUploadError() {
      this.dialogVisible = true;
    },
    // 不要有beforeRemove，要不会导致大文件不能自动删除
    handleRemove(file, fileList) {
      const index = this.fileList.indexOf(file);
      if (index > -1) {
        this.fileList.splice(index, 1);
      }
      console.log(file, fileList);
    },
    handleExceed(files, fileList) {
      console.log(files, fileList);
      this.$message.warning(`当前限制选择 1 个文件，你选择了 ${files.length} 个文件`);
    },
    beforeUpload(file) {
      if (file.size / 1024 / 1024 > 10) {
        this.$message.error('文件大小不能超过10MB');
        return false;
      }
      return true;
    },
    submitfeedback() {
      if (this.textarea === '') {
        this.$message.error('输入框内容不能为空');
        return;
      }
      this.feedback.ftext = this.textarea;
      if (this.fileList.length > 0) {
        [this.feedback.attachment] = this.fileList;
      } else {
        this.feedback.attachment = null;
      }
      this.SubmitFeedback(this.feedback)
        .then(() => {
          this.$message({
            message: '感谢您的反馈！',
            type: 'success',
          });
          this.$router.push({ path: '/' });
        });
    },
  },
};
</script>
