<template>
  <div>
    <b-nav tabs >
      <b-nav-item :active="noticesType==='unread'" @click="changeType('unread')">未读通知
      </b-nav-item>
      <b-nav-item :active="noticesType==='read'" @click="changeType('read')">已读通知</b-nav-item>
    </b-nav>
    <!-- 通知列表 -->
    <b-list-group v-if="notices.length > 0" class="notification-list">
      <transition-group name="list" tag="div" ref="listGroup">
        <b-list-group-item
          v-for="(notice, index) in notices"
          :key="index"
          @click="viewNotice(notice,index)"
          class="notification-item"
          :class="{ 'notification-unread': !notice.read }"
          :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null}">
          <div class="d-flex align-items-center">
            <!-- 头像 -->
            <b-avatar v-if="notice.type==='pcomment'
            || notice.type==='ccomment'
            || notice.type==='post'"
                      :src="notice.senderAvatar" size="2rem" class="mr-2"></b-avatar>
            <!-- 发送者名字和通知内容 -->
            <div class="d-flex flex-column justify-content-center">
              <small class="text-muted">{{ formatDate(notice.time) }}</small>
              <div class="d-flex align-items-center justify-content-between">
              <span v-if="notice.type==='pcomment'
              || notice.type==='ccomment' || notice.type==='post'"
                    STYLE="color: #409EFF" class="mr-2">{{ notice.senderName }}</span>
                <span v-else STYLE="color:saddlebrown" class="mr-2">系统通知：</span>
                <!-- 内容和类型 -->
                <span v-if="notice.type === 'pcomment'">评论了你的帖子:</span>
                <span v-if="notice.type === 'ccomment'">回复了你的评论:</span>
                <span v-if="notice.type === 'punishment'"
                      style="color: red">警告，你的账号出现违规：</span>
                <span v-if="notice.type === 'sue'">你的举报已得到处理：</span>
                <span v-if="notice.type === 'feedback'">你的反馈已得到处理回复：</span>
                <span v-if="notice.type === 'post'">在你的专区发帖：</span>
                <span style="color: gray">
                {{
                    notice.content.length > 5 ?
                      notice.content.slice(0, 5) + '...' : notice.content + ''
                  }}</span>
                <!-- 标记未读 -->
                <transition name="fade">
                  <span v-if="!notice.read" class="badge badge-danger ml-2 pop">New</span>
                </transition>
              </div>
            </div>
          </div>
        </b-list-group-item>
      </transition-group>
      <div v-if="more">
        <b-card class="loading" v-if="loading">正在加载...</b-card>
        <b-button v-else block @click="loadMoreNotices">加载更多</b-button>
      </div>
      <div style="font-size: 20px;font-family: 'Teko', sans-serif" v-else>
        ∑(っ°Д°;)っ已无更多通知
      </div>
    </b-list-group>
    <p v-else>没有通知</p>

    <!-- 查看通知弹窗 -->
    <b-modal v-if="currentNotice" v-model="modalVisible" hide-footer @hide="markAsRead">
      <div class="d-flex align-items-center">
        <!-- 头像 -->
        <b-avatar v-if="currentNotice.type==='pcomment'
        || currentNotice.type==='ccomment'|| currentNotice.type==='post'"
                  :src="currentNotice.senderAvatar" size="2rem" class="mr-2"></b-avatar>
        <div class="flex-grow-1">
          <!-- 标题 -->
          <div class="d-flex justify-content-between align-items-center">
            <span v-if="currentNotice.type==='pcomment'
            || currentNotice.type==='ccomment'||currentNotice.type==='post'"
                  STYLE="color: #409EFF">{{ currentNotice.senderName }}</span>
            <span v-else STYLE="color:saddlebrown">系统通知：</span>
            <small class="text-muted">{{ formatDate(currentNotice.time) }}</small>
          </div>
          <!-- 内容和类型 -->
          <div>
            <span v-if="currentNotice.type === 'pcomment'" @keydown.enter="showDetails"
                  class="postjump" @click="showDetails">评论了你的帖子: </span>
            <span v-if="currentNotice.type === 'ccomment'" @keydown.enter="showDetails"
                  class="postjump" @click="showDetails">回复了你的评论: </span>
            <span v-if="currentNotice.type === 'punishment'"
                  style="color: red">警告，你的账号出现违规： </span>
            <span v-if="currentNotice.type === 'sue'">你的举报已得到处理： </span>
            <span v-if="currentNotice.type === 'feedback'">你的反馈已得到处理回复： </span>
            <span v-if="currentNotice.type === 'post'"
                  style="font-size: 20px" class="postjump"
                  @click="showDetails" @keydown.enter="showDetails">在你的专区发帖：</span>
            <span class="preview mb-1" style="color:gray">{{ currentNotice.content }}</span>
          </div>
        </div>
      </div>
    </b-modal>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import request from '@/utils/request';

export default {
  data() {
    return {
      loading: false,
      page: 1,
      pagesize: 10,
      notices: [],
      more: true,
      noticesType: 'unread',
      currentNotice: null,
      currentIndex: 0,
      modalVisible: false, // 弹出窗口是否可见
      // 用来存通知分页查询的相关参数
      getparams: {
        pageSize: 5, // 每页查询条数
        requireID: 0, // 查询ID，查询在此noticeID之前的通知，由于ID是越新越大的，也就是不管新来的通知
        read: 0, // 获取未读的通知,read=0;获取已读的通知,read=1
      },
    };
  },
  mounted() {
    this.getNotices();
    window.addEventListener('scroll', this.handleScroll);
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
  created() {
    if (typeof localStorage !== 'undefined') {
      if (!localStorage.getItem('Style')) {
        localStorage.setItem('Style', JSON.stringify('day'));
      }
    }
    if (this.$route.params.before) {
      this.before = this.$route.params.before;
    }
  },
  methods: {
    handleScroll() {
      this.$nextTick(() => {
        const listGroup = this.$refs.listGroup.$el;
        if (listGroup.getBoundingClientRect().bottom <= window.innerHeight) {
          this.loadMoreNotices();
        }
      });
    },
    async loadMoreNotices() {
      this.loading = true;
      setTimeout(() => {
        this.getNotices();
      }, 500);
      this.loading = false;
    },
    // 获取通知列表
    getNotices() {
      request.get('/auth/getNotice', {
        params: {
          requireID: this.getparams.requireID,
          pageSize: this.getparams.pageSize,
          read: this.getparams.read,
        },
      })
        .then((response) => {
          // 注意这里接受的notice的字段是小写开头的驼峰,下面是对应的go结构体
          // type NoticeResponse struct {
          //   NoticeID     int    `json:"noticeID"`
          //   ReceiverName string `json:"receiverName"`
          //   SenderName   string `json:"senderName"`
          //   SenderAvatar string `json:"senderAvatar"`
          //   Type         string `json:"type"`
          //   Content      string `json:"content"`
          //   Read         bool   `json:"read"`
          //   PostID       int    `json:"postID"`
          //   Target       int    `json:"target"`
          //   PcommentID   int    `json:"pcommentID"`
          // }
          if (response.data.code === 201) {
            this.more = false;
            return;
          }
          this.getparams.requireID = response.data.lastID;
          const newNotices = response.data.noticeList.sort((a, b) => {
            // 先比较是否已读
            if (a.read && !b.read) {
              return 1;
            }
            if (!a.read && b.read) {
              return -1;
            }
            // 如果两个通知都已读或都未读，则比较时间
            const timeA = Date.parse(a.time);
            const timeB = Date.parse(b.time);
            if (timeA < timeB) {
              return 1;
            }
            if (timeA > timeB) {
              return -1;
            }
            return 0;
          });
          this.notices.push(...newNotices);
          if (response.data.noticeList.length < 5) this.more = false;
        })
        .catch((error) => {
          console.error(error);
        });
    },
    // 查看通知
    viewNotice(notice, index) {
      this.currentNotice = notice;
      this.currentIndex = index;
      this.modalVisible = true;
    },
    // 标记为已读
    markAsRead() {
      if (this.currentNotice.read === false) {
        const { noticeID } = this.currentNotice;
        request.patch(`/auth/readNotice/${noticeID}`, { read: true })
          .then(() => {
            this.notices[this.currentIndex].read = true;
            this.currentIndex = 0;
          })
          .catch((error) => {
            console.error(error);
          });
      }
      this.currentNotice = null;
    },
    showDetails() {
      // 标为已读
      const nowNotice = this.currentNotice;
      // const nowComment = this.currentComment;
      this.markAsRead();
      if (!this.userInfo) {
        this.toLogin = true;
        return;
      }
      // 跳转到对应的帖子
      setTimeout(() => {
        if (nowNotice.type === 'pcomment') {
          // const link = this.$router.resolve({
          //   name: 'postDetails',
          //   query: {
          //     id: nowNotice.postID,
          //     partition: this.partition,
          //     before: 'notice',
          //     pcommentID: nowNotice.target,
          //   },
          // });
          // window.open(link.href, '_blank');
          this.$router.push({
            name: 'postDetails',
            query: {
              id: nowNotice.postID,
              partition: this.partition,
              before: 'notice',
              pcommentID: nowNotice.target,
            },
          });
        } else if (nowNotice.type === 'ccomment') {
          // const link = this.$router.resolve({
          //   name: 'postDetails',
          //   query: {
          //     id: nowNotice.postID,
          //     partition: this.partition,
          //     before: 'notice',
          //     pcommentID: nowNotice.pcommentID,
          //     ccommentID: nowNotice.target,
          //   },
          // });
          // window.open(link.href, '_blank');
          this.$router.push({
            name: 'postDetails',
            query: {
              id: nowNotice.postID,
              partition: this.partition,
              before: 'notice',
              pcommentID: nowNotice.pcommentID,
              ccommentID: nowNotice.target,
            },
          });
        } else if (nowNotice.type === 'post') {
          const link = this.$router.resolve({
            name: 'postDetails',
            query: {
              id: nowNotice.postID,
              partition: this.partition,
              before: 'notice',
            },
          });
          window.open(link.href, '_blank');
        }
      }, 100);
    },
    formatDate(date) {
      // 格式化日期时间
      const d = new Date(date);
      return `${d.getFullYear()}年${
        d.getMonth() + 1
      }月${d.getDate()}日 ${String(d.getHours())
        .padStart(2, '0')}:${String(d.getMinutes())
        .padStart(2, '0')}:${String(d.getSeconds())
        .padStart(2, '0')}`;
    },
    changeType(name) {
      this.notices = [];
      this.noticesType = name;
      this.getparams.requireID = 0;
      this.more = true;
      if (name === 'unread') {
        this.getparams.read = 0;
      } else {
        this.getparams.read = 1;
      }
      this.getNotices();
    },
  },
};
</script>
<style>
.notification-item {
  cursor: pointer;
}

.postjump:hover {
  color: cornflowerblue;
  cursor: pointer;
}

@keyframes pop {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.2);
    opacity: 1;
  }
  100% {
    transform: scale(0);
    opacity: 0;
  }
}

.pop {
  animation: pop 0.5s ease-in-out;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.7s cubic-bezier(0.5, 0, 0.5, 1);
}

.list-enter,
.list-leave-to {
  opacity: 0;
  transform: rotate(180deg);
}
</style>
