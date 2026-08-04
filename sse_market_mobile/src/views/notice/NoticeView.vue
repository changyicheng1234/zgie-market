<template>
  <div>
    <van-tabs v-model="noticesType" @click="changeType">
      <van-tab title="未读通知" name="unread"></van-tab>
      <van-tab title="已读通知" name="read"></van-tab>
    </van-tabs>
    <!-- 通知列表 -->
    <van-list v-if="notices.length > 0" class="notification-list" ref="totalGroup">
      <!-- 标记未读 -->
      <van-cell
        v-for="(notice, index) in notices"
        :key="index"
        :class="{ 'notification-unread': !notice.read }"
        @click="viewNotice(notice,index)"
      >
        <template #icon>
          <van-image
            v-if="notice.type==='pcomment'|| notice.type==='ccomment'|| notice.type==='post'"
            :src="notice.senderAvatar"
            width="48"
            height="48"
            round
          />
        </template>
        <van-col style="margin-outside: 5px">
          <van-row style="margin-top:5px">{{ formatNoticeTitle(notice) }}</van-row>
          <van-row style="margin-top: 5px">{{ formatNoticeContent(notice) }}</van-row>
          <van-row style="color: #888888;font-size: 14px">{{ formatDate(notice.time) }}</van-row>
        </van-col>
        <van-badge v-if="notice.read===false" :content="'New'"/>
      </van-cell>
      <div v-if="more">
        <van-cell class="loading" v-if="loading">正在加载...</van-cell>
        <van-button v-else block @click="loadMoreNotices">加载更多</van-button>
      </div>
      <div v-else>
        <van-empty description="已无更多通知"/>
      </div>
    </van-list>
    <van-empty v-else description="暂无通知"></van-empty>
    <!-- 查看通知弹窗 -->
    <van-popup v-if="currentNotice"
               v-model="modalVisible" position="center"
               overlay @close="markAsRead"
               safe-area-inset-bottom
               round>
      <div class="d-flex align-items-center" style="margin: 10px">
        <!-- 头像 -->
        <img
          v-if="currentNotice.type==='pcomment'
          || currentNotice.type==='ccomment'
          || currentNotice.type==='post'"
          :src="currentNotice.senderAvatar"
          width="48"
          height="48"
          class="avatar"
          alt=""
        >
        <div style="display: flex;flex-direction: column;">
          <!-- 标题 -->
          <van-col>
            <van-row v-if="currentNotice.type==='pcomment'
            || currentNotice.type==='ccomment'|| currentNotice.type==='post'"
                     class="title" STYLE="color: #409EFF">{{ currentNotice.senderName }}
            </van-row>
            <van-row v-else class="title" STYLE="color:saddlebrown">系统通知：</van-row>
            <van-row style="font-size:10px;color: #888888">
              {{ formatDate(currentNotice.time) }}
            </van-row>
          </van-col>
          <!-- 内容和类型 -->
          <div style="display: flex;flex-direction: column;">
            <span v-if="currentNotice.type === 'pcomment'" class="postjump"
                  @click="showDetails" @keydown.enter="showDetails">评论了你的帖子: </span>
            <span v-if="currentNotice.type === 'ccomment'" class="postjump"
                  @click="showDetails" @keydown.enter="showDetails">回复了你的评论: </span>
            <span v-if="currentNotice.type === 'punish'"
                  style="color: red;font-size: 20px">警告，你的账号出现违规： </span>
            <span v-if="currentNotice.type === 'sue'"
                  style="font-size: 20px">你的举报已得到处理： </span>
            <span v-if="currentNotice.type === 'feedback'"
                  style="font-size: 20px">你的反馈已得到处理回复： </span>
            <span v-if="currentNotice.type === 'post'"
                  style="font-size: 20px" class="postjump"
                  @click="showDetails" @keydown.enter="showDetails">在你的专区发帖：</span>
            <span class="preview mb-2"
                  style="font-size: 18px;color:black">{{ currentNotice.content }}</span>
          </div>
        </div>
      </div>
    </van-popup>
    <van-divider
      :style="{ color: '#1989fa', borderColor: '#1989fa', padding: '0 16px',height: '50px'}"
    >
      ヽ(ﾟ∀ﾟ)ﾒ(ﾟ∀ﾟ)ﾉ 到底啦(￣▽￣)~*
    </van-divider>
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
      notices: [],
      noticesType: 'unread', //
      more: true,
      noticesNum: {
        totalNum: 0,
        unreadNum: 0,
        readNum: 0,
      }, // 通知的总数
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
    this.getNoticeNum();
  },
  mounted() {
    this.getNotices();
  },
  methods: {
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
    formatDate(date) {
      // 格式化日期时间
      const d = new Date(date);
      return `${d.getFullYear()}年${
        d.getMonth() + 1
      }月${d.getDate()}日 ${String(d.getHours())
        .padStart(2, '0')}:${String(d.getMinutes())
        .padStart(2, '0')}`;
    },
    formatNoticeContent(notice, isTitle = false) {
      let content = '';
      switch (notice.type) {
        case 'pcomment':
          content += '评论了你的帖子：';
          break;
        case 'ccomment':
          content += '回复了你的评论：';
          break;
        case 'punish':
          content += '警告，你的账号出现违规：';
          break;
        case 'sue':
          content += '你的举报已得到处理：';
          break;
        case 'feedback':
          content += '你的反馈已得到处理回复：';
          break;
        case 'post':
          content += '在你的专区发帖：';
          break;
        default:
          content += '';
      }
      if (isTitle) {
        content += '\n';
        content += notice.content;
      } else {
        content += notice.content.length > 5 ? `${notice.content.slice(0, 5)}...` : `${notice.content}`;
      }
      return content;
    },
    formatNoticeTitle(notice) {
      let title = '';
      if (notice.type === 'pcomment' || notice.type === 'ccomment' || notice.type === 'post') {
        title += `${notice.senderName}：\n`;
      }
      return title;
    },
    //
    getNoticeNum() {
      request.get('/auth/getNoticeNum')
        .then((response) => {
          this.noticesNum.totalNum = response.data.totalNum;
          this.noticesNum.unreadNum = response.data.unreadTotalNum;
          this.noticesNum.readNum = response.data.readTotalNum;
        })
        .catch((error) => {
          this.$toast.fail('通知数量获取失败');
          console.error(error);
        });
    },
    // 获取通知列表
    // 我在这里突然想到了一个在非常小概率情况下会出现的bug，
    // 就是获取了通知数量后到获取通知这一极小的时间段内通知的数量又改变了，
    // 可见用getNoticeNum()的数量来标志获取通知结束是不合适的，所以看后端返回来的结果其实是最好的
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
    // ...
    async loadMoreNotices() {
      this.loading = true;
      setTimeout(() => {
        this.getNotices();
      }, 500);
      this.loading = false;
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
          const link = this.$router.resolve({
            name: 'postDetails',
            query: {
              id: nowNotice.postID,
              partition: this.partition,
              before: 'notice',
              pcommentID: nowNotice.target,
            },
          });
          // 在新页面加载后，可以通过监听 popstate 事件来处理返回操作
          // window.open(link.href, '_blank');
          window.location.href = link.href;
        } else if (nowNotice.type === 'ccomment') {
          const link = this.$router.resolve({
            name: 'postDetails',
            query: {
              id: nowNotice.postID,
              partition: this.partition,
              before: 'notice',
              pcommentID: nowNotice.pcommentID,
              ccommentID: nowNotice.target,
            },
          });
          window.open(link.href, '_blank');
          // window.location.href = link.href;
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
    // 查看通知
    viewNotice(notice, index) {
      this.currentNotice = notice;
      this.currentIndex = index;
      this.modalVisible = true;
    },
  },
};
</script>

<style lang="scss" scoped>
.notification-list {
  margin-top: 10px;
  margin-bottom: 10px;
}

.notification-item {
  cursor: pointer;
}

.notification-unread {
  background-color: #f7f7f7;
}

.notification-unread .van-cell__title {
  font-weight: bold;
}

p {
  color: #888;
}

.avatar {
  border-radius: 50%;
}

.postjump {
  font-size: 38px;
  cursor: pointer;
  color: #409eff;
  text-decoration: underline;
}

.preview {
  font-size: 14px;
  display: inline-block;
  max-width: 600px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: pre-wrap;
}

.notification-item {
  cursor: pointer;
}

.title {
  font-size: 38px; /* 调整为适合的字体大小 */
}

</style>
