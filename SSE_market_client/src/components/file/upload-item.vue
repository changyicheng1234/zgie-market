<template>
  <div>
    <div class="upload-box">
      <!-- 自定义内容区 -->
      <slot></slot>
      <!-- 文件上传区 :disabled="disUpload"-->
      <el-upload
        class="upload-demo"
        ref="upload"
        drag
        multiple
        with-credentials
        http-request
        :action="url"
        :limit="limit"
        :data="upOptions"
        :file-list="fileList"
        :auto-upload="autoUpload"
        :on-error="handleError"
        :on-exceed="handleExceed"
        :on-success="handleSuccess"
        >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">
          将文件拖到此处，或
          <em>点击上传</em>
        </div>
      </el-upload>
    </div>
  </div>
</template>

<script>
/**
 * auth:weilan
 * github: https://github.com/hql7
 * props: 输入
 * title:组件标题
 * templateCode:[String, Array]下载模板所需编码,多个模板时array
 * url:上传地址
 * options:上传附带参数
 * autoUpload:自动上传
 * uploadSuccess:父组件监听函数
 * reg:是否校验excel
 * path:是否显示文件路径
 * contract: 获取文件上传路径所需参数
 * defaultProps：路径树配置参数
 * view: 开启将在上传区下显示可预览文件列表区
 * iscom: 是否通用上传，仅区分通用上传使用内部通用逻辑还是采用外部逻辑
 * schedule-plan: 是否进度计划模块
 * emit: 输出
 * uploadSuccess：上传成功回调
 * length: 已上传文件列表长度
 * solt: 插槽
 *  templateMsg：下载模板的显示信息
 * --------> path模式为通用文件上传路径，文件夹的增删改均在内部完成。
 * --------> 同普通上传一致只需传入接口所需参数，输出均为上传成功回调。
 * --------> 手动上传模式均需要在父组件调用子组件方法，在父组件引入此组件时，给此组件加ref，然后this.$refs.[""].toUpload()即可!
 * --------> 此组件请使用v-if展示。使用范例：claim-apply.vue。
 */
import { uploadFile } from '@/api'; // 导入组装树函数、拼接路径函数
export default {
  data() {
    return {
      fileList: [], // 已上传文件
    };
  },
  props: {
    uploadPath:'course/',
    // 自动上传
    autoUpload: {
      type: Boolean,
      default: false
    },
    // 上传个数限制
    limit: Number,
  },
  methods: {
    // 手动上传
    toUpload() {
      this.$refs.upload.submit();
    },
    uploadfile(param){
      // 获取上传的文件
      var file = param.file
      //发送请求的参数格式为FormData
      const formData = new FormData();
      formData.append("file",file);
      formData.append("path",)
      uploadFile(formData);
    },
    // 上传成功回调
    handleSuccess(res, file,fileList) {

    },
    // 上传失败回调
    handleError(err) {
      this.$emit("uploadError", err);
    },
    // 文件超出限制
    handleExceed() {
      this.$message.error("超出文件上传个数限制，最大上传个数为：10！");
    },
  },
  computed: {
    // 上传参数
    upOptions() {
      return this.options;
    }
  },
};
</script>

<style lang="scss">
.upload-box {
  .el-upload,
  .el-upload-dragger {
    width: 100%;
  }
}
</style>
