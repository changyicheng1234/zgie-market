<template>
  <!-- 文件管理器 auth：weilan time：2019-10-25 github：https://github.com/hql7 -->
  <div class="wl-explorer" @contextmenu.prevent>
    <b-button @click="fresh">刷新</b-button>
    <!-- 头部按钮区 -->
    <el-form
      v-if="this.userInfo.identity==='teacher'"
      class="wl-header-btn"
      :inline="true"
      :size="size"
      @submit.native.prevent
    >
      <el-form-item>
        <el-button type="primary" @click="showAdd()"
        >新增文件夹
        </el-button
        >
        <submit-btn
          type="danger"
          :size="size"
          @btn="handleDel"
          :status="load.del"
        >删除
        </submit-btn
        >
        <el-button @click="showUpload">上传文件</el-button>
        <!-- solt自定义头部按钮区 -->
        <slot name="header-btn"></slot>
      </el-form-item>
      <el-form-item v-show="uploading.ing">
        <span>正在上传：</span>
        <span class="c-blue u-uploading-name">{{ uploading.name }}</span>
        <span class="c-blue">({{ uploading.percentage }}%)</span>
      </el-form-item>
      <el-form-item class="u-right">
        <i
          class="iconfont icon-wl-list file-show-type"
          v-show="layout.show_list"
          @click="layout.show_list = !layout.show_list"
        ></i>
        <i
          class="iconfont icon-wl-grid file-show-type"
          v-show="!layout.show_list"
          @click="layout.show_list = !layout.show_list"
        ></i>
      </el-form-item>
    </el-form>
    <!--文件路径操作区-->
    <el-form
      :inline="true"
      :size="size"
      :model="file"
      class="wl-header-file"
      @submit.native.prevent
    >
      <!--      <el-form-item class="file-search-box">-->
      <!--        <el-input-->
      <!--            v-model="file.key"-->
      <!--            placeholder="请输入关键字搜索"-->
      <!--            @keyup.enter.native="fileSearch()"-->
      <!--        >-->
      <!--          <el-button-->
      <!--              slot="append"-->
      <!--              icon="el-icon-search file-search"-->
      <!--              @click="fileSearch()"-->
      <!--          ></el-button>-->
      <!--        </el-input>-->
      <!--      </el-form-item>-->
      <div class="grid-content bg-purple-dark">当前路径:{{ now_path }}</div>
    </el-form>
    <!-- 主内容区 -->
    <el-scrollbar class="wl-main-scroll">
      <!-- 文件列表区 -->
      <div class="wl-main-list">
        <!-- 表格型文件列表 -->
        <el-table
          v-show="layout.show_list"
          @selection-change="filrChecked"
          highlight-current-row
          :border="showBorder"
          :data="self_data"
          class="wl-table"
          ref="wl-table"
        >
          <el-table-column
            v-if="showCheckbox"
            align="center"
            type="selection"
            width="55"
          ></el-table-column>
          <el-table-column
            v-if="showIndex"
            align="center"
            type="index"
            label="序号"
            width="55"
          ></el-table-column>
          <slot name="table-column-top"></slot>
          <el-table-column
            v-for="i of selfColumns"
            show-overflow-tooltip
            :key="i._id"
            :prop="i.prop"
            :width="i.width"
            :label="i.label"
            :fixed="i.fixed"
            :align="i.align"
            :sort-by="i.sortBy"
            :sortable="i.sortable"
            :min-width="i.minWidth"
            :formatter="i.formatter"
            :column-key="i.columnKey"
            :class-name="i.className"
            :sort-method="i.sortMethod"
            :header-align="i.headerAlign"
            :render-header="i.renderHeader"
            :label-class-name="i.labelClassName"
          >
            <template slot-scope="scope">
              <!-- 非名称列 -->
              <template v-if="i.prop !== selfProps.name">
                {{
                  i.formatter
                    ? i.formatter(
                      scope.row,
                      scope.column,
                      scope.row[i.prop],
                      scope.$index
                    )
                    : scope.row[i.prop]
                }}
              </template>
              <!-- 名称列 -->
              <div
                v-else
                @click="enterTheLower(scope.row, scope.row[selfIsFolder])"
                class="wl-name-col wl-is-folder"
              >
                <!-- 不同文件类型图标区 -->
                <div class="namecol-iconbox">
                  <img
                    :src="fileTypeIcon(scope.row)"
                    class="name-col-icon"
                    alt="文件类型图标"
                  />
                </div>
                <!-- 不同文件类型 显示内容-->
                <div class="namecol-textbox">
                  {{
                    i.formatter
                      ? i.formatter(
                        scope.row,
                        scope.column,
                        scope.row[i.prop],
                        scope.$index
                      )
                      : scope.row[i.prop]
                  }}
                </div>
              </div>
            </template>
          </el-table-column>
          <slot name="table-column-bottom"></slot>
        </el-table>
        <!-- 列表型文件列表 -->
        <ul class="wl-list" v-show="!layout.show_list">
          <li
            class="wl-list-item wl-is-folder"
            v-for="(i, idx) in self_data"
            :key="i.Id"
          >
            <el-checkbox
              class="wl-checkbox"
              :class="{ 'wl-checkbox-checked': i._checked }"
              @change="listItemCheck($event, i)"
              v-model="i._checked"
            ></el-checkbox>
            <div @click="enterTheLower(i, i[selfIsFolder])">
              <img
                :src="fileTypeIcon(i)"
                class="name-col-icon"
                alt="文件类型图标"
              />
              <p class="list-item-name" :title="i[selfProps.name]">
                {{
                  i.formatter
                    ? i.formatter(i, null, i[selfProps.name], idx)
                    : i[selfProps.name]
                }}
              </p>
            </div>
          </li>
        </ul>
      </div>
    </el-scrollbar>
    <!-- 文件预览区 -->
    <template v-if="usePreview">
      <file-view
        v-show="layout.view"
        ref="file-view"
        class="file-view-components"
        :previewType="previewType"
        :previewOptions="previewOptions"
        @close="layout.view = false"
      ></file-view>
    </template>
    <!-- 文件上传区 -->
    <template v-if="useUpload">
      <fade-in v-show="layout.upload">
        <h3 class="edit-header">上传文件</h3>
        <el-scrollbar class="scroll">
          <el-form
            :size="size"
            ref="template_form"
            label-position="top"
            class="template_form rule-form"
          >
            <el-form-item label="文件路径">
              {{ now_path }}
            </el-form-item>
            <el-form-item label="导入文件">
              <el-upload
                class="upload-demo"
                ref="upload"
                drag
                multiple
                :http-request="uploadFile"
                action=""
                limit=10
                :auto-upload="false"
                :file-list="fileList"
                :on-exceed="handleExceed"
              >
                <i class="el-icon-upload"></i>
                <div class="el-upload__text">
                  将文件拖到此处，或
                  <em>点击上传</em>
                </div>
              </el-upload>
            </el-form-item>
          </el-form>
        </el-scrollbar>
        <!-- 按钮区 -->
        <div class="submit-btn-box">
          <el-button style="margin-left: 10px;"
                     size="small"
                     type="success"
                     @click="submitUpload">
            上传到服务器
          </el-button>
          <el-button :size="size" @click="layout.upload = false"
          >取消
          </el-button
          >
        </div>
      </fade-in>
    </template>
    <!--    文件夹新增区  -->
    <template>
      <fade-in v-show="layout.add_folder">
        <h3 class="edit-header">
          <p>新增文件夹</p>
        </h3>
        <el-scrollbar class="scroll">
          <el-form
            size="medium"
            ref="folder_form"
            label-position="top"
            class="folder_form rule-form"
            @keyup.enter.native="handleAdd()"
          >
            <el-form-item label="文件夹路径" prop="ParentId">
              {{ now_path }}
            </el-form-item>
            <el-form-item label="文件夹名称 " prop="Name">
              <el-input v-model="folder_form.Name" placeholder="请输入文件夹名称"></el-input>
            </el-form-item>
          </el-form>
        </el-scrollbar>
        <div class="submit-btn-box">
          <submit-btn @btn="handleAdd()"
                      :status="load.folder"></submit-btn>
          <el-button size="medium" @click="closeAdd">取消</el-button>
        </div>
      </fade-in>
    </template>
  </div>
</template>

<script>
import submitBtn from '@/components/file/submit-btn.vue'; // 导入防抖组件
import fadeIn from '@/components/file/fade-in.vue'; // 引入滑入组件
import uploadItem from '@/components/file/upload-item'; // 导入导入组件
import { addFolderApi, getFileListApi, uploadFile, getUrl, delFileApi, delFolderApi } from '@/api';
import { mapState } from 'vuex';
import { closeOtherLayout } from '@/utils/file'; // 导入组装树函数、拼接路径函数
export default {
  name: 'wlExplorer',
  components: {
    submitBtn,
    fadeIn,
    uploadItem
  },
  data() {
    return {
      loading: false,
      formData: '',
      fileList: [],
      load: {
        del: false, // 删除
        move: false, // 移动
        upload: false, // 上传
      }, // loading状态
      uploading: {
        name: '',
        percentage: 0,
        ing: false,
      }, // 当前上传文件状态
      layout: {
        show_list: true, // 文件展示形式
        edit_path: false, // 是否编辑路径
        view: false, // 预览视图
        upload: false, // 上传视图
        add_folder: false, // 新增文件夹
      }, // 视图管理
      folder_form: {
        ParentId: '',
        Name: '',
        preview: [],
        handle: [],
      }, // 文件夹表单
      file: {
        path: '', // 文件路径
        key: '', // 关键字
      }, // 文件相关参数
      self_data: [], // 当前数据
      now_path: 'course/',//当前路径
      file_checked_data: [], // 列表多选数据
      upload_data: {
        name: null,
      }, // 上传提交操作抛出的信息
    };
  },
  props: {
    /**
     * 头部更多操作自定义内容
     * 需要包含内容：
     * name: String 每条数据的名字
     * command: Function 每条数据的指令
     * disabled: Boolean 每条数据的禁用
     * divided: Boolean 每条数据的显示分割线
     * icon: String 每条数据的图标类名
     */
    headerDropdown: Array,
    // 是否显示复选框
    showCheckbox: {
      type: Boolean,
      default: true,
    },
    // 是否显示顺序号
    showIndex: {
      type: Boolean,
      default: true,
    },
    // 表格是否显示边框
    showBorder: {
      type: Boolean,
      default: true,
    },
    // 文件表格数据
    data: Array,
    // 文件表头数据【[参数：所有el-Table-column Attributes] (https://element.eleme.cn/#/zh-CN/component/table)】
    columns: Array,
    /**
     * 配置项
     * isFolder: Boolean
     * name: String 表示名称列的字段
     */
    props: Object,
    // 校验是否文件夹函数，（row）参数为当前行数据，用于复杂类型，当isFolderFn优先使用计算结果，不存在时使用props配置内的isFolder字段
    isFolderFn: Function,
    // 是否使用自带上传组件
    useUpload: {
      type: Boolean,
      default: true,
    },
    // 是否使用自带预览组件
    usePreview: {
      type: Boolean,
      default: true,
    },
    // 预览文件类型
    previewType: {
      type: String,
      default: 'img',
    },
    // 预览文件地址或配置项
    previewOptions: Object,
    size: {
      type: String,
      default: 'medium',
    },
  },
  methods: {
    // 删除文件
    handleDel() {
      if (this.file_checked_data.length === 0) {
        this.$message({
          showClose: true,
          message: '请选择要删除的文件或文件夹',
          type: 'error',
        });
        return;
      }
      this.$confirm('是否确认删除选中数据？', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      })
        .then(() => {
          let loadingInstance = this.$loading({fullscreen:true});
          this.file_checked_data.forEach((item) => {
            if (item.isFolder) {
              delFolderApi(item.key)
                .then(() => {
                  this.$message.success('删除成功');
                })
                .catch(() => {
                  this.$message.error('删除失败');
                });
            } else {
              delFileApi(item.key)
                .then(() => {
                  this.$message.success('删除成功');
                })
                .catch(() => {
                  this.$message.error('删除失败');
                });
            }
          });
          loadingInstance.close();
        });
        this.fresh();
    },
    // 文件超出限制
    handleExceed() {
      this.$message.error('超出文件上传个数限制，最大上传个数为：10！');
    },
    // 重写el-upload的上传
    uploadFile(file) {
      this.formData.append('file', file.file);
    },
    // 提交文件上传
    submitUpload() {
      this.$message('正在上传，请稍等...');
      this.formData = new FormData(); //创建FormData对象
      this.$refs.upload.submit(); // file将添加到this.formData
      this.formData.append('path', this.now_path); // 上传文件路径
      console.log(this.formData);
      let loadingInstance = this.$loading({fullscreen:true});
      uploadFile(this.formData)
        .then(() => {
          this.$message.success('上传成功');
          this.fileList = [];
          this.formData = '';
          loadingInstance.close();
        })
        .catch(() => {
          this.$message.error('上传失败');
          loadingInstance.close();
        });
      this.fresh();
    },
    // 新建文件夹
    handleAdd() {
      let loadingInstance = this.$loading({fullscreen:true});
      addFolderApi(this.now_path + this.folder_form.Name + '/')
        .then(() => {
          this.$message({
            message: '新建文件夹成功',
            type: 'success',
          });
          this.folder_form.Name = '';
          loadingInstance.close();
        })
        .catch((err) => {
          this.$message({
            message: '新建文件夹失败' + err,
            type: 'error',
          });
          loadingInstance.close();
        });
      this.fresh();
    },
    /**
     * 文件夹时进入下一级-文件时预览文件
     * row: Object 行数据
     * isFolder: Boolean 行是否文件夹
     */
    enterTheLower(row, isFolder) {
      if (!isFolder) {
        getUrl(row['key'])
          .then((res) => {
            // console.log(res)
            this.$message('文件地址' + res.data.data.url);
            // open()方法在苹果系统中可能会被截止而打不开
            // window.open(res.data.data.url);
            window.location.href = res.data.data.url;

          })
          .catch((err) => {
            this.$message.error('文件打开失败' + err);
          });

        // this.previewFile(row);
        return;
      }
      this.now_path = row[this.selfProps.key];
      this.getFileList(row[this.selfProps.key]);
    },
    // 获取文件夹列表
    getFileList(path) {
      console.log('getFileList()');
      let loadingInstance = this.$loading({fullscreen:true});
      getFileListApi(path)
        .then((data) => {
          this.data = data.data || [];
          loadingInstance.close();
        })
        .catch((error) => {
          console.error(`getFileList error${error}`);
          loadingInstance.close();
        });
    },
    // 刷新
    fresh(){
      setTimeout(()=>{
        this.getFileList(this.now_path);
      },1500);
    },
    // 显示上传界面
    showUpload() {
      this.upload_selected = this.file.id;
      this.upload_data = {
        parentPathId: this.file.pid,
        pathId: this.file.id,
        isCurrentFolder: true,
      };
      closeOtherLayout(this.layout,'upload');
      // this.layout.upload = true;
      },
    // 显示新增文件夹页面
    showAdd() {
      this.layout.add_folder = true;
    },
    closeAdd() {
      this.layout.add_folder = false;
    },
    // 关闭上传界面
    closeUpload() {
      this.layout.upload = false;
    },
    // 文件上传失败回调
    uploadError(err) {
      this.$emit('uploadError', err);
      this.load.upload = false;
    },
    // 根据文件类型显示图标
    fileTypeIcon(row) {
      let _path = '';
      // 文件夹
      if (row[this.selfIsFolder]) {
        _path = row[this.selfIsLock]
          ? require('./images/file_automatic@3x.png')
          : require('./images/folder@3x.png');
        return _path;
      }
      // 其他根据后缀类型
      let _suffix = row[this.selfProps.suffix];
      if (!_suffix) {
        _path = require('./images/file_none@3x.png');
        return _path;
      }
      if (['jpg', 'jpeg', 'png', 'gif', 'bmp'].includes(_suffix)) {
        // 图片
        _path = require('./images/file_img@3x.png');
      } else if (['zip', 'rar', '7z'].includes(_suffix)) {
        _path = require('./images/file_zip@3x.png');
      } else if (
        ['avi', 'mp4', 'rmvb', 'flv', 'mov', 'm2v', 'mkv'].includes(_suffix)
      ) {
        _path = require('./images/file_video@3x.png');
      } else if (['mp3', 'wav', 'wmv', 'wma'].includes(_suffix)) {
        _path = require('./images/file_mp3@3x.png');
      } else if (['xls', 'xlsx'].includes(_suffix)) {
        _path = require('./images/file_excel@3x.png');
      } else if (['doc', 'docx'].includes(_suffix)) {
        _path = require('./images/file_docx@3x.png');
      } else if ('pdf' == _suffix) {
        _path = require('./images/file_pdf@3x.png');
      } else if ('ppt' == _suffix) {
        _path = require('./images/file_ppt@3x.png');
      } else if ('txt' == _suffix) {
        _path = require('./images/file_txt@3x.png');
      } else {
        _path = require('./images/file_none@3x.png');
      }
      return _path;
    },
    // 记录多选列表数据
    filrChecked(val) {
      this.self_data.forEach((i) => (i._checked = false));
      val.forEach((i) => (i._checked = true));
      this.file_checked_data = val;
    },
    // 列表模式记录多选数据
    listItemCheck(check, val) {
      this.$refs['wl-table'].toggleRowSelection(val);
    },
    // 预览文件
    previewFile(row) {
      this.$emit('preview', row, this.showPreview);
    },
    // 打开预览组件
    showPreview() {
      this.layout.view = true;
    },
    // 处理数据变动
    handleDataChange(val) {
      let _data = val || [];
      if (this.isFolderFn) {
        _data.forEach((i) => {
          i.isFolder = this.isFolderFn(i);
        });
      }
      this.self_data = _data;
    },
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
    // 自身表头数据
    selfColumns() {
      let _data = this.columns || [];
      _data.forEach((i, idx) => {
        i._id = `_col_${idx}`;
      });
      return _data;
    },
    // 自身配置项
    selfProps() {
      return {
        isFolder: 'isFolder', // Boolean 用于有布尔值字段表示数据是否文件夹类型的情况,当使用isFolderFn函数时，此参数会被忽略
        name: 'name', // String 用于显示名称列的字段
        suffix: 'suffix', // String 用于判断后缀或显示文件类型列的字段
        pathName: 'name', // String 路径数据 显示名称字段
        pathDisabled: 'disabled', // String 路径数据 禁用字段
        ...this.props,
      };
    },
    // 自身移动 下拉框树 配置项
    selfMoveProps() {
      return {
        label: this.selfProps.pathName,
        children: this.selfProps.pathChildren,
        disabled: this.selfProps.pathDisabled,
      };
    },
    // 将是否文件夹的两种判断方式合并返回
    selfIsFolder() {
      return this.isFolderFn ? 'isFolder' : this.selfProps.isFolder;
    },
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
  },
  watch: {
    // 检测data数据更新列表
    data(val) {
      this.handleDataChange(val);
    },
  },
  created() {
    if (this.data && this.data.length > 0) {
      this.handleDataChange(this.data);
    }
    if (this.userInfo.identity == 'teacher') {
      let parts = this.userInfo.email.split('@');
      this.now_path='course/'+parts[0]+'/';
    }
    else{
      let parts = this.$route.params.email.split('@');
      this.now_path='course/'+parts[0]+'/';
    }
    this.fresh();
  },
};
</script>

<style lang="scss">
@import "./css/index.css";
@import "./css/clear.css";
@import "./icons/iconfont.css";
</style>
