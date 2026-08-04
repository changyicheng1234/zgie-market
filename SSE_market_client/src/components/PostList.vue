<template>
  <div class='page-container'>
    <b-col class="ml-4">
      <p>当前分区:{{ partition }}</p>
      <div class='home-view'>
        <b-modal v-model='toLogin' title='登录' ok-only ok-title="取消登录"
                 modal-class="custom-modal">
          <LoginView />
        </b-modal>
        <div class="login-section" v-if="this.$route.name == 'home' && partition == '主页'
          && !userInfo">
          <b-button :variant="isNightStyle ? 'outline-warning':'primary'" @click="toLogin = true"
                    style="margin-top:100px;width: 150px;border-radius: 20px;">
            立即登录
          </b-button>
          <div class="register-section" style="margin-top:40px;color: white;">
            <span>还没有账号？</span>
            <a href="#" onclick="window.open('/register', '_blank');">立即注册！</a>
          </div>
        </div>
        <b-button :variant="isNightStyle ? 'outline-warning':'primary'"
                  v-if="this.partition !== '主页' && this.partition !== '课程专区'"
                  class="back_button"
                  style="margin-left: -100px"
                  @click="goback">
          <b-icon-reply class="mr-2"></b-icon-reply>
          返回
        </b-button>
        <div v-if="posts.length > 0">
          <b-row>
            <b-col v-for='post in posts' :key='post.id' cols='12' md='12' lg='12' class='mb-3'>
              <b-card class='px-3 py-2 card-shadow'
                      @click="() => { showDetails(post); updatebrowse(post) }"
                      :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
              'color': isNightStyle ? 'gray' : null }"
                      style="max-width:900px;min-width: 700px">
                <div class='text-muted float-right' @click.stop>
                  <b-icon icon='three-dots-vertical' @click.stop='toggleMenu(post)'></b-icon>
                </div>
                <b-list-group
                  v-if='post.showMenu'
                  class="float-right"
                  style='
                width: 100px;
                height: 1.25rem;
                margin-top: -20px;
                font-size: 0.9rem;
              '
                  @click.stop
                >
                  <b-list-group-item
                    :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null, 'z-index': 9999}">
                    <b-icon class="mr-2" :icon="post.isSaved ? 'star-fill' : 'star'"
                            @click.stop="save(post)"
                            :class="{ 'text-warning': post.isSaved }"></b-icon>
                    收藏
                  </b-list-group-item>
                  <b-list-group-item
                    :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null, 'z-index': 9999 }"
                    v-if='post.authorTelephone !== userInfo.phone'
                    @click.stop='showReportModal = true'
                  >
                    <b-icon-exclamation-triangle class='mr-2'></b-icon-exclamation-triangle>
                    举报
                  </b-list-group-item>
                  <b-modal
                    v-model='showReportModal'
                    title='举报'
                    @hidden='clearReportReason'
                    @ok='submitReport(post)'
                    ok-title='Submit'
                  >
                    <b-form-textarea v-model='reportReason' placeholder='请输入举报原因' rows='8'>
                    </b-form-textarea>
                  </b-modal>
                  <b-list-group-item
                    v-if='post.authorTelephone === userInfo.phone'
                    :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null, 'z-index': 9999}"
                    @click.stop='showDeleteModal = true'
                  >
                    <b-icon-trash class='mr-2'></b-icon-trash>
                    删除
                  </b-list-group-item>
                  <b-modal
                    v-model='showDeleteModal'
                    title='确认删除'
                    ok-title='Confirm'
                    @ok='postdelete(post)'
                  >
                    <p>你确定要删除这个帖子吗？</p>
                  </b-modal>
                </b-list-group>
                <b-row class='mt-0'>
                  <b-col md='4' class='mb-2'>
                    <b-avatar :src="post.authorAvatar" size="4rem" class="mr-2"></b-avatar>
                    <div class="author-box" @click.stop
                      :style="{ 'background-color': isNightStyle ?
                      'rgb(246, 155, 10)' : 'rgb(17, 167, 226)' }">
                      {{ post.author }}
                      <span v-if="post.authorIdentity==='teacher'"
                            class="badge badge-primary font-weight-light">教师</span>
                    </div>
                    <span
                    :class="{
                    'badge': true, 'badge-pill': true,
                    'badge-primary': post.authorTitle === '菜鸟',
                    'badge-success': post.authorTitle === '大虾',
                    'badge-danger': post.authorTitle === '码农',
                    'badge-custom': post.authorTitle === '程序猿',
                    'badge-warning': post.authorTitle === '工程师',
                    'badge-info': post.authorTitle === '大牛',
                    'badge-custom2': post.authorTitle === '专家',
                    'badge-secondary': post.authorTitle === '大神',
                    'badge-dark': post.authorTitle === '祖师爷' // 添加一个自定义样式类
                    }"
                    style="margin-left: 10px;">
                    {{ post.authorTitle }}
                    </span>
                  </b-col>
                </b-row>
                <b-card-title>{{ post.title }}</b-card-title>
                <b-card-text class="pre word-limit">{{ post.content }}</b-card-text>
                <div v-if="fileListGet.length > 0" class="photo-viewer">
                  <div class="thumbnail-container">
                    <template v-if="fileListGet(post).length === 4">
                      <div>
                        <img :src="fileListGet(post)[0]"
                             width="270"
                             height="270"
                             @click="handlePictureCardPreview(fileListGet(post)[0])"
                             @keyup.enter="handlePictureCardPreview(fileListGet(post)[0])"
                             alt="Post Photo"/>
                        <img :src="fileListGet(post)[1]"
                             width="270"
                             height="270"
                             style="margin-top:20px"
                             @click="handlePictureCardPreview(fileListGet(post)[1])"
                             @keyup.enter="handlePictureCardPreview(fileListGet(post)[1])"
                             alt="Post Photo"/>
                      </div>
                      <div>
                        <img :src="fileListGet(post)[2]"
                             width="270"
                             height="270"
                             @click="handlePictureCardPreview(fileListGet(post)[2])"
                             @keyup.enter="handlePictureCardPreview(fileListGet(post)[2])"
                             alt="Post Photo"/>
                        <img :src="fileListGet(post)[3]"
                             width="270"
                             height="270"
                             style="margin-top:20px"
                             @click="handlePictureCardPreview(fileListGet(post)[3])"
                             @keyup.enter="handlePictureCardPreview(fileListGet(post)[3])"
                             alt="Post Photo"/>
                      </div>
                    </template>
                    <template v-else>
                      <div v-for="(file, index) in fileListGet(post)" :key="index">
                        <img :src="file"
                             width="270"
                             height="270"
                             @click="handlePictureCardPreview(file)"
                             @keyup.enter="handlePictureCardPreview(file)"
                             alt="Post Photo"/>
                      </div>
                    </template>
                  </div>
                </div>
                <div class="d-flex justify-content-between">
                  <div class='d-flex justify-content-between'>
                    <small class='text-muted'>{{ formatDate(post.postTime) }}</small>
                  </div>
                  <div class="tag-group">
                    <span class="tag-group__title"></span>
                    <el-tag v-for="tag in post.tag" :key="tag.label" :type="tag.type"
                            effect="plain" size="mini">{{ tag.label }}
                    </el-tag>
                  </div>
                </div>
                <div class='d-flex justify-content-between align-items-center mt-3'>
                  <div class='text-muted'>
                    <b-icon :icon="post.isLiked ? 'heart-fill' : 'heart'"
                            @click.stop="like(post)"
                            :class="{ 'text-danger': post.isLiked }"></b-icon>
                    {{ post.like }}
                  </div>
                  <div class='text-muted'>
                    <b-icon-eye-fill></b-icon-eye-fill>
                    {{ post.browse }}
                  </div>
                  <div class='text-muted'>
                    <b-icon icon='chat-dots-fill'></b-icon>
                    {{ post.comment }}
                  </div>
                </div>
              </b-card>
            </b-col>
          </b-row>
        </div>
        <div v-else>
          <el-empty description="没有符合要求的帖子哦"></el-empty>
        </div>
        <el-pagination
          class="is-background"
          style="margin-left: 300px"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="currentPage"
          :page-sizes="[5, 10, 25, 50]"
          :page-size="pageSize"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalItems">
        </el-pagination>
      </div>
    </b-col>
    <b-col>
      <div class="audio-section" :style="{ marginTop: userInfo ? '0px' : '250px' }"
           v-if="this.$route.name == 'home' && partition == '主页'">
        <audio ref="audio" :src="music_path" style="margin-top: 140px; z-index: 1000;" controls/>
      </div>
      <div class="video-section" v-if="this.$route.name == 'home' && partition == '主页'"
           :style="{ marginTop: userInfo ? '250px' : '500px' }">
        <video ref="videoPlayer" :src="video_path" style="margin-top: 10px; z-index: 1000;"
               controls></video>
      </div>
      <div class='hots-bar' :style="{ marginTop: userInfo ? '500px' : '750px'}"
           v-if="this.$route.name == 'home' && partition == '主页'">
        <b-card class="px-3 py-2 card-shadow" style="width: 300px; height: 100%;"
                :style="{ 'background-color': isNightStyle ? 'rgb(70, 70, 70)' : null }">
          <div>
            <b-card-header style="font-weight: bold;"
                           :style="{ 'color': isNightStyle ? 'gray' : null }">
              24小时火文
              <b-button class="bi bi-arrow-clockwise"
                        :style="{ 'background-color': isNightStyle ? 'rgb(70, 70, 70)' : 'white',
            'color': isNightStyle ? 'gray' : 'black' }"
                        @click="calculateheat()">
                <b-icon-arrow-clockwise></b-icon-arrow-clockwise>
                刷新
              </b-button>
            </b-card-header>
          </div>
          <b-list-group-item v-for="(hotpost, index) in hotposts" :key="index"
                             :style="{ 'background-color': isNightStyle ? 'rgb(70, 70, 70)' : null,
            'color': isNightStyle ? 'gray' : null }"
                             @click="() => { showDetails(hotpost); updatebrowse(hotpost) }">
            <div>
              <template v-if="index<3">
                <span style="color: red; font-weight: bold;">{{ index + 1 }}</span>
              </template>
              <template v-else>
              <span style="font-weight: bold;"
                    :style="{ 'color': isNightStyle ? 'gray' : 'black'}">
                {{ index + 1 }}</span>
              </template>
              <template v-if="hotpost.title.length > 8">
                {{ hotpost.title.substring(0, 8) }}
              </template>
              <template v-else>
                {{ hotpost.title }}
              </template>
              <span style="float: right;" class="badge badge-danger ml-2 pop">
              {{ Math.floor(hotpost.heat) }}
            </span>
            </div>
          </b-list-group-item>
        </b-card>
        <div class="px-3 py-2 card-shadow block" style="width: 300px; height: 100%;"
             :style="{ 'background-color': isNightStyle ? 'rgb(70, 70, 70)' : null }">
          <el-carousel height="120px">
            <el-carousel-item v-for="item in imagebox" :key="item.id">
              <a :href="getWebsiteURL(item.id)">
                <span>Some text</span>
                <img :src="item.idView" class="image" alt="">
              </a>
            </el-carousel-item>
          </el-carousel>
        </div>
      </div>
    </b-col>
  </div>
</template>

<script>

import { mapState, mapActions } from 'vuex';
import LoginView from '@/views/login/LoginView.vue';

export default {
  components: {
    LoginView,
  },
  activated() {
    console.log('PostList activated!');
    if (localStorage.getItem('Partition')) {
      this.partition = JSON.parse(localStorage.getItem('Partition'));
    } else if (this.$route.query.partitions && this.$route.query.partitions !== '主页') {
      // 检测partition是否有变
      if (this.partition !== this.$route.query.partitions) {
        this.totalItems = 0;
        this.currentPage = 0;
        this.posts = [];
      }
      this.partition = this.$route.query.partitions;
      // 将partition保存在本地缓存中
      localStorage.setItem('Partition', JSON.stringify(this.$route.query.partitions));
    } else {
      if (this.$route.name === 'course') this.partition = '课程专区';
      else this.partition = '主页';
      this.totalItems = 0;
      this.currentPage = 0;
      this.posts = [];
    }
    // 返回到页面的时候刷新一下
    this.PostNum();
    this.browsePosts();
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
    }),
    tagTypeMap() {
      return {
        大厂: '',
        高工资: 'success',
        实习: 'danger',
      };
    },
    fileListGet() {
      return (post) => {
        if (!post.photos || post.photos === '') return [];
        return post.photos.split('|');
      };
    },
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
  },
  data() {
    return {
      totalItems: 0,
      pageSize: 5,
      currentPage: 1,
      music_path: '/山高水长.mp3',
      video_path: '/宣传片.mp4',
      posts: [],
      partition: '',
      fileList: [],
      hotposts: [],
      imagebox: [
        {
          id: 0,
          idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel1.png',
        },
        {
          id: 1,
          idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/assets/image/caraousel2.png',
        },
        {
          id: 2,
          idView: 'https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com//assets/image/caraousel3.png',
        },
        // { id: 3, idView: require('../assets/image/caraousel4.png') },
      ],
      dialogImageUrl: '',
      dialogVisible: false,
      userTelephone: '',
      postID: '',
      isSaved: '',
      isLiked: '',
      browse: '',
      title: '',
      heat: '',
      searchinfo: '',
      showDeleteModal: false,
      showReportModal: false,
      reportReason: '',
      toLogin: false,
    };
  },
  props: {
    tag: String,
    searchsort: {
      type: String,
      default: 'home',
    }, // 分表查询时用于区分,包含home,history,save
  },
  created() {
    if (typeof localStorage !== 'undefined') {
      if (!localStorage.getItem('Style')) {
        localStorage.setItem('Style', JSON.stringify('day'));
      }
    }
    if (localStorage.getItem('Partition')) {
      this.partition = JSON.parse(localStorage.getItem('Partition'));
    } else if (this.$route.query.partitions && this.$route.query.partitions !== '主页') {
      this.partition = this.$route.query.partitions;
      // 将partition保存在本地缓存中
      localStorage.setItem('Partition', JSON.stringify(this.$route.query.partitions));
    } else {
      this.partition = '主页';
    }
    this.searchinfo = this.$route.query.searchinfo;
    // 在页面创建时默认加载主页帖子列表
    this.PostNum();
    this.browsePosts();
    this.calculateheat();
  },
  methods: {
    ...mapActions('postModule', { postBrowse: 'browse' }),
    ...mapActions('postModule', { getPostNum: 'getPostNum' }),
    ...mapActions('postModule', { postLike: 'like' }),
    ...mapActions('userModule', { postSave: 'save' }),
    ...mapActions('postModule', { postUpdateBrowse: 'updatebrowse' }),
    ...mapActions('postModule', { postHots: 'calculateheat' }),
    ...mapActions('postModule', { deletepost: 'deletepost' }),
    ...mapActions('postModule', { submitreport: 'submitreport' }),
    urlToLink(str) { // 链接识别和显示
      if (!str) {
        return '';
      }
      let link = str;
      const urlPattern = /(https?:\/\/|www\.)[a-zA-Z_0-9\-@]+(\.\w[a-zA-Z_0-9\-:]+)+(\/[()~#&\-=?+%/.\w]+)?/g;
      link = str.replace(urlPattern, (match) => {
        let href = match;
        if (match.indexOf('http') === -1) {
          href = `http://${match}`;
        }
        return `<a target="_blank" href="${href}">${match}</a>`;
      });
      return link;
    },
    handleSizeChange(val) {
      this.pageSize = val;
      this.currentPage = 1; // 如果改变每页显示的数量，回到第一页
      this.browsePosts();
    },
    handleCurrentChange(val) {
      this.currentPage = val;
      this.browsePosts();
    },
    goback() {
      this.$router.replace({ name: 'partitions' });
    },
    toRegister() {
      window.open('/register', '_blank');
    },
    showDetails(post) {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      // const routeLink = this.$router.resolve({
      //   name: 'postDetails',
      //   params: { partition: this.partition },
      //   query: {
      //     id: post.id, title: post.title, before: this.$route.name, partition: this.partition,
      //   },
      // });
      // window.open(routeLink.href, '_blank');
      this.$router.push({
        name: 'postDetails',
        params: {
          id: post.id,
          title: post.title,
          before: this.$route.name,
          partition: this.partition,
        },
        query: {
          id: post.id,
          title: post.title,
          before: this.$route.name,
          partition: this.partition,
        },
      });
    },
    getUserTitle(userScore) {
      if (userScore < 100) {
        return '菜鸟';
      }
      if (userScore >= 100 && userScore < 300) {
        return '大虾';
      }
      if (userScore >= 300 && userScore < 600) {
        return '码农';
      }
      if (userScore >= 600 && userScore < 1000) {
        return '程序猿';
      }
      if (userScore >= 1000 && userScore < 2000) {
        return '工程师';
      }
      if (userScore >= 2000 && userScore < 3000) {
        return '大牛';
      }
      if (userScore >= 3000 && userScore < 4000) {
        return '专家';
      }
      if (userScore >= 4000 && userScore < 5000) {
        return '大神';
      }
      return '祖师爷';
    },
    // 查询满足要求的帖子数量
    async PostNum() {
      try {
        if (this.userInfo) {
          this.userTelephone = this.userInfo.phone;
        } else {
          // 游客访问
          this.userTelephone = '00000000000';
        }
        const { data } = await this.getPostNum({
          userTelephone: this.userTelephone,
          partition: this.partition,
          searchinfo: this.searchinfo,
          searchsort: this.searchsort,
        });
        this.totalItems = data.Postcount;
      } catch (error) {
        console.error(error);
      }
    },
    // 分页加载帖子
    async browsePosts() {
      // 从后端返回一个结构体变量的方法
      try {
        // 向后端发送请求并获取帖子列表
        const { data } = await this.postBrowse({
          userTelephone: this.userTelephone,
          partition: this.partition,
          searchinfo: this.searchinfo,
          tag: this.$props.tag,
          searchsort: this.searchsort,
          limit: this.pageSize,
          offset: (this.currentPage - 1) * this.pageSize,
        });
        console.log(data);
        // 将获取到的帖子列表数据赋值给 posts 变量
        this.posts = data
          .map((post) => ({
            id: post.PostID,
            author: post.UserName,
            authorTitle: this.getUserTitle(post.UserScore),
            authorTelephone: post.UserTelephone,
            authorAvatar: post.UserAvatar,
            authorIdentity: post.UserIdentity,
            title: post.Title,
            content: post.Content,
            like: post.Like,
            comment: post.Comment,
            postTime: post.PostTime,
            isSaved: post.IsSaved,
            isLiked: post.IsLiked,
            browse: post.Browse,
            heat: post.Heat,
            photos: post.Photos,
            tag: post.Tag ? post.Tag.split(',')
              .map((tagText) => ({
                type: this.tagTypeMap[tagText.trim()], // 使用 this.tagTypeMap
                label: tagText.trim(),
              })) : [],
            showMenu: false,
          }))
          .sort((a, b) => new Date(b.postTime) - new Date(a.postTime));
      } catch (error) {
        console.error(error);
      }
    },
    formatDate(date) {
      // 格式化日期时间
      const d = new Date(date);
      return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日 ${String(
        d.getHours(),
      )
        .padStart(2, '0')}:${String(d.getMinutes())
        .padStart(2, '0')}:${String(
        d.getSeconds(),
      )
        .padStart(2, '0')}`;
    },
    toggleMenu(post) {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      const updatedShowMenu = {
        ...post,
        showMenu: !post.showMenu,
      };
      this.posts.splice(this.posts.indexOf(post), 1, updatedShowMenu);
    },
    save(post) {
      // 切换收藏状态
      const updatedPost = {
        ...post,
        isSaved: !post.isSaved,
      };
      // 用更新后的 post 对象替换原先的 post 对象
      this.posts.splice(this.posts.indexOf(post), 1, updatedPost);
      this.userTelephone = this.userInfo.phone;
      this.postID = post.id;
      this.isSaved = post.isSaved;
      // 请求
      this.postSave({
        userTelephone: this.userTelephone,
        postID: this.postID,
        isSaved: this.isSaved,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
    },
    like(post) {
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      // 切换点赞状态
      const updatedPost = {
        ...post,
        isLiked: !post.isLiked,
      };
      // 更新点赞数
      updatedPost.like += post.isLiked ? -1 : 1;
      // 用更新后的 post 对象替换原先的 post 对象
      this.posts.splice(this.posts.indexOf(post), 1, updatedPost);
      this.userTelephone = this.userInfo.phone;
      this.postID = post.id;
      this.isLiked = post.isLiked;
      // 请求
      this.postLike({
        userTelephone: this.userTelephone,
        postID: this.postID,
        isLiked: this.isLiked,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
    },
    updatebrowse(post) {
      // 切换浏览状态
      const updatedPost = {
        ...post,
        browse: post.browse,
      };
      // 更新浏览次数
      updatedPost.browse = post.browse + 1;
      // 用更新后的 post 对象替换原先的 post 对象
      this.posts.splice(this.posts.indexOf(post), 1, updatedPost);
      this.userTelephone = this.userInfo.phone;
      this.postID = post.id;
      // 请求
      this.postUpdateBrowse({
        userTelephone: this.userTelephone,
        postID: this.postID,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
    },
    calculateheat() {
      // 向后端请求热度前10的帖子
      this.postHots({
        postID: this.postID,
        title: this.title,
        heat: this.heat,
      })
        .then((hots) => {
          this.hotposts = hots.data.map((hot) => ({
            id: hot.PostID,
            title: hot.Title,
            heat: hot.Heat,
          }));
        })
        .catch((err) => {
          console.error(err);
        });
    },
    postdelete(post) {
      this.postID = post.id;
      this.deletepost({
        postID: this.postID,
      })
        .then(() => {
          this.$bvToast.toast('帖子删除成功', {
            title: '系统提醒',
            variant: 'primary',
            solid: true,
          });
          setTimeout(() => {
            this.$router.go(0);
          }, 2000);
        })
        .catch((err) => {
          console.log(err);
          this.$bvToast.toast('未找到该帖子', {
            title: '数据异常',
            variant: 'danger',
            solid: true,
          });
          setTimeout(() => {
            this.$router.go(0);
          }, 2000);
        });
    },
    submitReport(post) {
      this.postID = post.id;
      this.userTelephone = this.userInfo.phone;
      this.submitreport({
        TargetID: this.postID,
        Targettype: 'post',
        userTelephone: this.userTelephone,
        Reason: this.reportReason,
      })
        .then(() => {
          this.$bvToast.toast('举报发送成功', {
            title: '系统提醒',
            variant: 'primary',
            solid: true,
          });
        })
        .catch((err) => {
          this.$bvToast.toast(err.response.data.msg, {
            title: '数据验证错误',
            variant: 'danger',
            solid: true,
          });
        });
    },
    clearReportReason() {
      this.reportReason = '';
    },
    handlePictureCardPreview(file) {
      this.dialogImageUrl = file;
      this.dialogVisible = true;
    },
    getWebsiteURL(id) {
      if (id === 0) {
        return 'https://ztjy.sysu.edu.cn/';
      }
      if (id === 1) {
        return 'https://www.sysu.edu.cn/news/ztjj2/gdzt/zt2023ncjgzhy.htm';
      }
      if (id === 2) {
        return 'https://www.sysu.edu.cn/xxg/zdjj1.htm';
      }
      if (id === 3) {
        return 'https://ddhzt.sysu.edu.cn/';
      }
      return ''; // 添加默认的返回值
    },
  },
};
</script>

<style lang='scss' scoped>
@import '@/style/css/HomeView.css';
.thumbnail-container {
  display: flex;
  flex-wrap: wrap;
}
.thumbnail-container div {
  width: calc(100% / 3);
  padding: 10px;
  box-sizing: border-box;
}

.b-card-title:hover {
  color: rgba(56, 104, 225, 0.9);
}

.page-container {
  display: flex;
}

.home-view {
  flex: 1;
  margin-left: 0px; /* Adjust the margin as needed */
  z-index: 1000;
}
.hots-bar {
  flex: 0;
  //position: absolute;
  //margin-left: 950px;
}

.el-pagination.is-background .el-pager li.active {
  background-color: #409EFF;
  color: #fff;
}

.el-pagination.is-background .el-pager li:not(.active) {
  background-color: #e4e7ed;
}

.block {
  position: relative; /* 设置父容器为相对定位，为绝对定位的图片提供参考 */
}

.block .image {
  position: absolute; /* 将图片设置为绝对定位 */
  top: 0;
  left: 0;
  width: 100%; /* 图片宽度占满父容器 */
  height: 100%; /* 图片高度占满父容器 */
  object-fit: cover; /* 按比例缩放图片，以覆盖整个容器 */
}

/* 添加一个自定义样式类 */
.badge-custom {
  background-color: #ff6600; /* 橙色背景 */
  color: #fff; /* 白色文字颜色 */
}

.badge-custom2 {
  background-color: #c428eb; /* 橙色背景 */
  color: #fff; /* 白色文字颜色 */
}

.pre {
  white-space: pre-wrap;
  word-wrap: anywhere;
}

.word-limit {
 display: -webkit-box;
 -webkit-box-orient: vertical;
 -webkit-line-clamp: 3;
 overflow: hidden;
}
</style>
