<template>
  <div>
    <b-button :variant ="isNightStyle ? 'outline-warning':'primary'"
              class="back_button" style="margin-left: -100px"
              @click="goback">
      <b-icon-reply class="mr-2"></b-icon-reply>
      返回
    </b-button>
    <wl-Explorer
      ref="wl-explorer-cpt"
      :header-dropdown="headerHandle"
      :upload-options="uploadOptions"
      :columns="file_table_columns"
      :all-path="all_folder_list"
      :is-folder-fn="isFolderFn"
      :folderType="rource_type"
      :data="file_table_data"
      :props="explorer_prop"
      size="small"
      @handleFolder="handleFolder"
      @download="download"
      @del="fileDel"
      @closeFade="closeOtherLayout(fade)"
    >
    </wl-Explorer>
  </div>
</template>

<script>

import 'weilan-ui/lib/weilan-ui.css';
import { closeOtherLayout } from '@/utils/file'; // 导入关闭其他弹出类视图函数
import {
  getFileListApi, // 1获取文件夹列表
} from '@/api';
import submitBtn from '@/components/file/submit-btn.vue'; // 导入防抖组件
import fadeIn from '@/components/file/fade-in.vue'; // 引入滑入组件
import WlExplorer from '@/pages/WlExplorer';

export default {
  components: {
    // eslint-disable-next-line vue/no-unused-components
    fadeIn,
    // eslint-disable-next-line vue/no-unused-components
    submitBtn,
    WlExplorer,
  },
  data() {
    // eslint-disable-next-line no-underscore-dangle
    const _GB = 1024 * 1024;
    // const vm = this;
    return {
      load: {
        folder: false,
      }, // loading管理
      fade: {
        folder: false,
      }, // 弹出类视图管理
      headerHandle: [{
        name: '权限',
        command: 'auth',
      }], // 头部按钮更多操作-自定义权限
      file_table_columns: [
        {
          label: '名称',
          prop: 'Name',
          minWidth: 120,
        },
        {
          label: '修改日期',
          align: 'center',
          width: 120,
          formatter(row) {
            return row.EditTime.split('T')[0] || '-';
          },
        },
        {
          label: '类型',
          align: 'center',
          width: 90,
          formatter(row) {
            return row.Type === 1 ? '文件夹' : row.SuffixName;
          },
        },
        {
          label: '大小',
          minWidth: 90,
          align: 'center',
          formatter(row) {
            if (row.Size === null) return '-';
            if (row.Size < 1024) {
              // 1024以下显示kb
              return `${row.Size}KB`;
            }
            if (row.Size < _GB) {
              // 1024*1024
              // eslint-disable-next-line no-underscore-dangle
              const _mb = (row.Size / 1024).toFixed(2);
              return `${parseFloat(_mb)}MB`;
            }
            // eslint-disable-next-line no-underscore-dangle
            const _gb = (row.Size / _GB).toFixed(2);
            return `${parseFloat(_gb)}GB`;
          },
        },
      ], // 自定义表格列
      file_table_data: [], // 表格数据
      all_folder_list: [], // 所有文件夹列表
      tree_folder_list: [], // 树形文件夹列表
      type: {
        folder: 1,
        img: 2,
        video: 3,
        other: 4,
      }, // 文件类型
      rource_type: {
        self: 1, // 自建
      }, // 数据来源类型
      explorer_prop: {
        key: 'key', // 绝对请求路径
        name: 'Name',
        match: 'Name',
        splic: true,
        suffix: 'SuffixName',
        pathId: 'Id',
        pathPid: 'ParentId',
        pathName: 'Name',
        pathChildren: 'Children', // String 路径数据 children字段
        pathConnector: '\\', // String 路径父子数据拼接连接符,默认为'\'
        pathParents: 'Parents', // String 路径数据所有直系祖先节点自增长identityId逗号拼接
        pathIdentityId: 'IdentityId', // String 路径数据自增长id
      }, // 文件管理器配置项
      path: {}, // 路径数据
      folder_form: {
        ParentId: '',
        Name: '',
        preview: [],
        handle: [],
        Describe: '',
      }, // 文件夹表单
      folder_rules: {
        ParentId: [
          {
            required: true,
            message: '请选择文件路径',
            trigger: 'blur',
          },
        ],
        Name: [
          {
            required: true,
            message: '请填写文件夹名称',
            trigger: 'blur',
          },
        ],
      }, // 文件夹表单验证
      child_act_saved: {}, // 存储子组件数据，用于编辑更新
      tree_select_prop: {
        label: 'Name',
        children: 'Children',
      }, // 树形下拉框配置项
      uploadOptions: {
      }, // 上传文件附加参数
    };
  },
  methods: {
    // 返回
    goback() {
      // eslint-disable-next-line no-restricted-globals
      history.back();
    },
    closeOtherLayout,
    /**
     * @name 上传文件提交操作
     */
    download(data, func) {
      console.log(data, func);
    },
    /**
     * 根据关键词搜索文件
     * file: Object 文件属性
     * update: Boolean 数据是否需要更新（不需要表示已存在）
     */
    fileSearch(file, update) {
      if (update) {
        this.path = file;
        this.getFileList();
      }
    },
    // 获取文件夹列表
    getFileList(path) {
      console.log('getFileList()');
      getFileListApi(path)
        .then((data) => {
          this.file_table_data = data.data || [];
        }).catch((error) => {
          console.error(`getFileList error${error}`);
        });
    },
    /**
     * 编辑文件夹
     * act:Object 当前选中文件夹
     * type:String 操作类型 add添加edit编辑
     * file:Object 当前路径信息
     */
    handleFolder(act, type, file) {
      this.path = file;
      this.fade.folder = true;
      if (type === 'add') {
        this.fade.folder = true;
        this.$refs.folder_form.resetFields();
        this.folder_form.Id = '';
        this.folder_form.ParentId = file.id;
        return;
      }
      this.child_act_saved = act;
      this.folder_form = { ...act };
    },
    // 判断是否文件夹函数
    isFolderFn(row) {
      return row.Type === this.type.folder;
    },
  },
  created() {
    this.closeOtherLayout = closeOtherLayout;
    this.getFileList();
  },
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  position: absolute;
  top: 0;
  left: 0;
  padding: 25px;
  width: 100%;
  height: 100%;
  background: #f7f7f7;
  box-sizing: border-box;
}
</style>
