<template>
  <div>
    <van-overlay :show="showloading" duration="0.1" :custom-style="{background:'rgba(0,0,0,0.1)'}">
      <van-loading class="wrapper" size="24px" vertical>加载中...</van-loading>
    </van-overlay>
    <van-nav-bar
        :title="post.title"
        left-text="返回"
        right-text="更多"
        left-arrow
        @click-left="goback"
        @click-right="toggleMenu"
        fixed
        placeholder
        safe-area-inset-top
    />
    <!--帖子详情-->
    <div class='postDetails' style="margin-left:10px">
      <van-cell class='mx-auto my-5' style="max-width: 750px;"
                :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
          'color': isNightStyle ? 'gray' : null}">
        <van-popup v-model="post.showMenu" position="bottom" round
                   :style="{ height: '40' }" @click.stop>
          <div
              :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null,
                    margin:'20px'}">
            <van-icon name='star-o' size='20' @click='save(post)' v-if='!post.isSaved'/>
            <van-icon v-else color='rgb(255,220,0)' name='star' size='20' @click='save(post)'/>
            收藏
          </div>
          <hr>
          <van-row
              v-if="this.post.authorTelephone !== userInfo.phone"
              :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null,
                    margin:'20px'}"
              @click.stop="showReportModal = true">
            <van-icon size="20px" name="failure" class="mr-2"></van-icon>
            举报
          </van-row>
          <van-dialog v-model='showReportModal'
                      title='举报'
                      @confirm='submitReport("post",post.postID)'
                      @cancel='clearReportReason'
                      @hidden='clearReportReason'
                      confirm-button-text="提交"
                      show-cancel-button>
            <van-field v-model='reportReason' type="textarea" placeholder='请输入举报原因' rows='8'>
            </van-field>
          </van-dialog>
          <van-row
              v-if="this.post.authorTelephone === userInfo.phone"
              :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
                    'color': isNightStyle ? 'gray' : null,
                    margin:'20px'}"
              @click.stop="showDeleteModal = true">
            <van-icon size="20px" name="delete-o" class="mr-2"></van-icon>
            删除
          </van-row>
          <van-dialog v-model="showDeleteModal" message="你确定要删除这个帖子吗？"
                      showCancelButton
                      @confirm="postdelete(post)">
          </van-dialog>
        </van-popup>
        <div class="van-row--flex">
          <van-image :src="post.authorAvatar"
                     width="48"
                     height="48"
                     round class="mr-3"></van-image>
          <van-col>
            <div style='margin-top: 10px;margin-left: 5px;' class='username-container'>
              <span style='margin-top: 10px;' class='username'>{{ post.author }}</span>
            </div>
          </van-col>
          <van-col style="margin-top: 10px;margin-left: 5px;">
            <van-tag v-if="post.authorIdentity==='teacher'" color="#7232dd" plain>
              教师
            </van-tag>
            <van-tag
                :color="getTagColor(post.authorTitle)">{{ post.authorTitle }}
            </van-tag>
          </van-col>
        </div>
        <van-row class="post_title">{{ post.title }}
        </van-row>
        <van-row class="post_content"
                 style="margin-bottom: 10px">
          {{ post.content }}
        </van-row>
        <div v-if="fileListGet.length > 0" class="photo-viewer van-row">
          <div class="thumbnail-container">
            <template v-if="fileListGet.length === 4">
              <div>
                <img :src="fileListGet[0]"
                     class="photo"
                     @click="handlePictureCardPreview(fileListGet[0])"
                     @keyup.enter="handlePictureCardPreview(fileListGet[0])"
                     @loadeddata="handlePictureCardPreview(fileListGet[0])"
                     alt="" preview-text="Post Photo"
                />
                <img :src="fileListGet[1]"
                     class="photo"
                     style="margin-top:5px"
                     @click="handlePictureCardPreview(fileListGet[1])"
                     @keyup.enter="handlePictureCardPreview(fileListGet[1])"
                     @loadeddata="handlePictureCardPreview(fileListGet[1])"
                     alt=""
                />
              </div>
              <div>
                <img :src="fileListGet[2]"
                     class="photo"
                     @click="handlePictureCardPreview(fileListGet[2])"
                     @keyup.enter="handlePictureCardPreview(fileListGet[2])"
                     @loadeddata="handlePictureCardPreview(fileListGet[2])"
                     alt=""
                />
                <img :src="fileListGet[3]"
                     class="photo"
                     style="margin-top:5px"
                     @click="handlePictureCardPreview(fileListGet[3])"
                     @keyup.enter="handlePictureCardPreview(fileListGet[3])"
                     @loadeddata="handlePictureCardPreview(fileListGet[3])"
                     alt=""/>
              </div>
            </template>
            <template v-else>
              <div v-for="(file, index) in fileListGet" :key="index">
                <img :src="file"
                     class="photo"
                     @click="handlePictureCardPreview(file)"
                     @keyup.enter="handlePictureCardPreview(file)"
                     @loadeddata="handlePictureCardPreview(file)"
                     alt=""/>
              </div>
            </template>
          </div>
        </div>
        <div>
          <div class='d-flex justify-content-between'>
            <span class='post_content' style="color:rgb(128,128,128);
            margin-left:25px;">{{ formatDate(post.postTime) }}</span>
          </div>
          <div class="tag-group" v-if="post.tag">
            <span class="tag-group__title"></span>
            <van-tag v-for="tag in post.tag" :key="tag.label" :type="tag.type"
                     effect="plain" size="mini">{{ tag.label }}
            </van-tag>
          </div>
        </div>
        <div class='van-row--flex' style="margin-bottom: 5px;margin-top: 20px;">
          <div class="text-muted">
            <van-icon size="20px" :name="post.isLiked ? 'good-job' : 'good-job-o'"
                      :color="post.isLiked ? '#ee0a24' : ''"
                      @click.stop="like()" :class="{ 'text-danger': post.isLiked }"></van-icon>
            <span class="post-like">{{ post.like }}</span>
          </div>
          <div class="text-muted" style="margin-left:50px">
            <van-icon size="22px" name="eye-o" style="vertical-align: -1px;"></van-icon>
            <span class="browse-like">{{ post.browse }}</span>
          </div>
          <div class='text-muted' style="margin-left:50px; ">
            <van-icon size="20px" name='comment-o' style="vertical-align: -2px;"></van-icon>
            <span class="comment-like">{{ commentsNum }}</span>
          </div>
          <van-button class="div" @click="post.showCommentForm
            = !post.showCommentForm" type="default" plain
                      style="margin-left:20px;border:none;color:blue;margin-top:-1px;
                      text-decoration: underline;text-underline-offset:0.2em;">
            {{ post.showCommentForm ? '隐藏评论' : '评论' }}
          </van-button>
        </div>
      </van-cell>

      <!--显示帖子评论窗口-->
      <van-popup v-model="post.showCommentForm" position="bottom" :style="{ height: '50%' }">
        <div>
          <van-field v-model="pcomment.content"
                     rows="3"
                     autosize
                     maxlength="1000"
                     show-word-limit
                     type="textarea"
                     ref="pcommentTextarea"
                     id="pcommentInput"
                     placeholder="请写下你的精彩评论...">
          </van-field>
        </div>
        <div class="van-row" style="margin-right: 5px">
          <!-- <van-button style="margin-right: 2px"
                      type='default' size="small" plain
                      round native-type="button"
                      @click="showEmojiStatus()">😀
          </van-button> -->
          <div v-if="showEmoji">
            <picker
                :include="['people']"
                :showSearch="false"
                :showPreview="false"
                :showCategories="false"
                @select="addEmojiToPcomment"
            />
          </div>
          <van-button
              @click="pcommentPost"
              round type="primary" :loading="buttonLoading">提交评论
          </van-button>
        </div>

      </van-popup>
    </div>

    <!--显示和发表帖子评论、评论的评论-->
    <div class='comment'>
      <!--帖子评论排序按钮-->
      <div class="van-row--flex van-row--justify-center">
        <van-button @click="sortkind='Date';comments=sortcomments(comments)"
                    plain hairline icon="descending"
                    type="primary" size="small" style="margin-left: 190px">时间
        </van-button>
        <van-button @click="sortkind='heat';comments=sortcomments(comments)"
                    plain hairline icon="descending"
                    type="primary" size="small" style="margin-left: 2px">热度
        </van-button>
      </div>
      <!-- 帖子评论-->
      <van-list name="comment-list" tag="div">
        <div v-for="(comment, index) in visibleComments"
             :key="index" ref="commentRef" :id="`comment-${comment.pcommentID}`">
          <div class="van-hairline--top-bottom"
               :style="{ 'background-color': isNightStyle ? 'rgb(50,50,50)' : 'white',
          'color': isNightStyle ? 'gray' : null }">
            <!-- 显示评论详情-->
            <div style="margin-bottom: 0.5px;margin-left: 20px;">
              <!--              评论头像-->
              <div class="van-col">
                <van-image :src="comment.authorAvatar"
                           width="30"
                           height="30"
                           round style="margin-right: 20px"></van-image>
              </div>
              <!--        评论作者和评论内容-->
              <div class="van-row" style="margin-top: 20px">
                <div class="van-row">
                  <div class="van-col comment-author-name" style="margin-top:5px;margin-left: -15px;
                font-size:13px;">{{ comment.author }}
                  </div>
                  <div class="van-col" style="margin-top: 11px;margin-left: 7px;font-size:11px;">
                    <van-tag v-if="comment.authorIdentity==='teacher'" color="#7232dd" plain>
                      教师
                    </van-tag>
                    <van-tag :color="getTagColor(comment.authorTitle)">
                      {{ comment.authorTitle }}
                    </van-tag>
                  </div>
                </div>
                <div class="comment-content" style="margin-top:20px;">
                  {{ comment.content }}
                </div>
                <div>
                  <div class="post_content" style="color:rgb(128,128,128);">
                    {{ formatDate(comment.commentTime) }}
                  </div>
                </div>
                <!--显示每个评论的点赞和回复数，点赞和回复图片对应点赞和回复功能-->
                <div class='van-row--flex justify-content-between align-items-center'>
                  <div class="text-muted">
                    <van-icon size="20px" :name="comment.isLiked ? 'good-job' : 'good-job-o'"
                              :color="comment.isLiked ? '#ee0a24' : ''"
                              @click.stop="pclike(index)"
                              :class="{ 'text-danger': comment.isLiked }">
                    </van-icon>
                    <span style="font-size: 14px;vertical-align: 1.2px;margin: 0.8px;">
                      {{ comment.likeNum }}</span>
                  </div>
                  <van-icon size="20px" name="comment-o"
                            style="vertical-align: middle; margin-top: 13px;
                            margin-left: 30px;"
                            @click.stop="comment.showReplyForm
            = !comment.showReplyForm">
                  </van-icon>
                  <div class='text-muted' @click.stop>
                    <van-icon size="20px" name='ellipsis'
                              @click.stop="comment.showMenu =
              !comment.showMenu" style="margin-left: 33px;margin-top: 1px;"></van-icon>
                  </div>
                </div>
                <!--对帖子评论的更多功能选择菜单：举报和删除-->
                <van-popup
                    v-model='comment.showMenu'
                    position="bottom"
                    round
                    :style="{height:'10%'}"
                    @click.stop
                >
                  <!--                  举报评论按钮-->
                  <div class="van-row--flex"
                       v-if='comment.authorTelephone !== userInfo.phone'
                       @click.stop='showReportModal = true'
                       @keydown="handleKeyboardEvent"
                  >
                    <van-icon name="failure" class='mr-2' style="font-size:25px;
                    margin-left: 10px;margin-top: 30px;">
                    </van-icon>
                    <span style="font-size:20px;
                    margin-left: 8px;margin-top: 29px;">举报</span>
                  </div>
                  <!--                  提交举报弹窗-->
                  <van-dialog v-model='showReportModal'
                              title='举报'
                              @confirm='submitReport("pcomment",comment.pcommentID)'
                              @cancel='clearReportReason'
                              @hidden='clearReportReason'
                              confirm-button-text="提交"
                              show-cancel-button>
                    <van-field v-model='reportReason'
                               type="textarea" placeholder='请输入举报原因' rows='8'>
                    </van-field>
                  </van-dialog>
                  <!--                  删除评论按钮-->
                  <div
                      v-if='comment.authorTelephone === userInfo.phone'
                      @click.stop='showDeleteModal = true'
                      @keydown.stop="handleKeyboardEvent"
                  >
                    <van-icon size="20px" name="delete-o" class='mr-2'></van-icon>
                    <span style="font-size:20px;margin-left: 5px;">删除</span>
                  </div>
                  <!--              确认删除弹窗    -->
                  <van-dialog
                      v-model='showDeleteModal'
                      show-cancel-button
                      @confirm='pcommentdelete(comment)'
                  >
                    <p>你确定要删除这条评论吗？</p>
                  </van-dialog>
                </van-popup>
                <!--如果点击了评论，将显示评论窗口-->
                <van-popup
                    round
                    position="bottom" :style="{ height: '50%' }"
                    v-model="comment.showReplyForm">
                  <van-form @submit="ccommentPost(index)">
                    <van-field v-model="ccomment.content"
                               rows="3"
                               maxlength="100"
                               show-word-limit
                               autosize
                               type="textarea"
                               placeholder="请写下你的精彩评论...">
                    </van-field>
                    <!--                    表情选择器和提交评论按钮-->
                    <div>
                      <!-- <van-button style="margin-right: 2px"
                                  type="default" native-type="button"
                                  plain size="small"
                                  @click="showEmojiStatus()">😀
                      </van-button> -->
                      <div v-if="showEmoji">
                        <picker
                            :include="['people']"
                            :showSearch="false"
                            :showPreview="false"
                            :showCategories="false"
                            @select="addEmojiToCcomment"
                            style="height: 100px"
                        />
                      </div>
                      <van-button type="primary"
                                  plain native-type="submit" :loading="buttonLoading">
                        提交评论
                      </van-button>
                    </div>
                  </van-form>
                </van-popup>
              </div>
            </div>
            <!-- 查看回复列表按钮-->
            <van-button v-if="comment.subComments.length > 0"
                        @click="showRepliesModal=true;showcommentsindex=index"
                        size="small"
                        type="default"
                        plain>
              查看回复共{{ len(comment.subComments) }}条
            </van-button>
          </div>
        </div>
      </van-list>
      <!-- 评论的评论（下弹窗） -->
      <van-popup
          position="bottom"
          round
          :style="{ height: '80%' }"
          v-model="showRepliesModal">
        <!--显示评论的评论/回复列表-->
        <van-list name="comment-list">
          <div v-for="(subComment, subIndex) in visibleSubComments(showcommentsindex)"
               :key="subIndex" :id="`ccomment-${subComment.ccommentID}`" tabindex="0">
            <!--   用线来间隔-->
            <hr>
            <div class="van-row">
              <!-- 头像-->
              <div class="van-col" style="margin-right: 5px">
                <van-image width="30"
                           height="30" round class="mr-3"
                           :src="subComment.authorAvatar"></van-image>
              </div>
              <!-- 作者名称和回复内容、时间&点赞、举报和回复图标-->
              <div class="van-row">
                <!--作者名称和回复内容-->
                <div class="van-row">
                  <van-row>
                    <div class="van-col comment-author-name" style="margin-top: 6px;
                  margin-left: 4px; font-size: 13px;">{{ subComment.author }}
                    </div>
                    <van-col style="margin-top: -8px;margin-left: 5px;">
                      <van-tag v-if="subComment.authorIdentity==='teacher'" color="#7232dd" plain>
                        教师
                      </van-tag>
                      <van-tag
                          :color="getTagColor(subComment.authorTitle)">{{ subComment.authorTitle }}
                      </van-tag>
                    </van-col>
                  </van-row>
                  <div
                      class="van-row--flex"
                      v-if="subComment.userTargetName !== ''">
                            <span
                                class="comment-content"
                                style="color: cadetblue;font-size: 13px;margin-left: 10px;">
                              回复@{{ subComment.userTargetName }}:</span>
                  </div>
                  <div
                      class="comment-content" style="font-size: 13px; margin-left: 10px;">
                    {{ subComment.content }}
                  </div>
                </div>
                <!--回复时间&点赞、举报和回复图标-->
                <div class="van-row--flex">
                  <!--回复时间-->
                  <div></div>
                  <div class="text-muted" style="font-size: 10px; margin-left: 10px;">
                    {{ formatDate(subComment.commentTime) }}
                  </div>
                  <!--点赞图标-->
                  <div class="text-muted">
                    <van-icon style="font-size: 16px;margin-left: 60px;"
                              :name="subComment.isLiked ? 'good-job' : 'good-job-o'"
                              :color="subComment.isLiked ? '#ee0a24' : ''"
                              @click.stop="cclike(showcommentsindex,subIndex)"
                              :class="{ 'text-danger': subComment.isLiked }">
                    </van-icon>
                    <span style="font-size: 10px;margin-left: -1px;">
                      {{ subComment.likeNum }}</span>
                  </div>
                  <!--举报/删除图标-->
                  <div class='text-muted'>
                    <!--举报图标-->
                    <div v-if="subComment.authorTelephone !== userInfo.phone">
                      <van-icon style="font-size: 16px;margin-left: 10px;" name='failure'
                                @click.stop='showReportModal = true'></van-icon>
                      <van-dialog v-model='showReportModal'
                                  title='举报'
                                  @confirm='submitReport("ccomment",subComment.ccommentID)'
                                  @cancel='clearReportReason'
                                  @hidden='clearReportReason'
                                  confirm-button-text="提交"
                                  show-cancel-button>
                        <van-field autosize type="textarea"
                            v-model='reportReason' placeholder='请输入举报原因' rows='8'>
                        </van-field>
                      </van-dialog>
                    </div>
                    <!--删除图标-->
                    <div v-else>
                      <van-icon style="font-size: 16px;margin-left: 10px;"
                                name='delete-o' @click.stop='showDeleteModal = true'></van-icon>
                      <van-dialog
                          v-model='showDeleteModal'
                          confirm-button-text='确认'
                          show-cancel-button
                          @confirm='ccommentdelete(subComment)'
                      >
                        <p>你确定要删除这条评论吗？</p>
                      </van-dialog>
                    </div>
                  </div>
                  <!--回复按钮，点击后跳出评论的评论的回复窗口-->
                  <div class="text-muted">
                    <van-icon style="font-size: 16px;margin-left: 10px;margin-top: 2px;"
                              name="comment-o"
                              @click="replyshow = !replyshow; nowReplyComment=subComment">
                      <span style="font-size: 10px;margin-left:5px;vertical-align: 3px;">回复</span>
                    </van-icon>
                  </div>

                </div>
              </div>
            </div>
          </div>
        </van-list>
        <!-- 弹出评论回复窗口-->
        <van-popup position="bottom" round
                   :style="{height:'50%'}"
                   v-model="replyshow">
          <van-form @submit=
                        "ccommentPost(showcommentsindex,
                          nowReplyComment.author,
                          nowReplyComment.ccommentID)">
            <van-field v-model="ccomment.content"
                       rows="2"
                       autosize
                       maxlength="100"
                       show-word-limit
                       type="textarea"
                       id="ccommentInput"
                       :placeholder="'回复@'+nowReplyComment.author">
            </van-field>
            <!--表情选择器-->
            <div>
              <!-- <van-button style="margin-right: 2px"
                          type='primary' size="small" plain
                          native-type="button"
                          @click="showEmojiStatus()">😀
              </van-button> -->
              <div v-if="showEmoji">
                <picker
                    :include="['people']"
                    :showSearch="false"
                    :showPreview="false"
                    :showCategories="false"
                    @select="addEmojiToCcomment"
                />
              </div>
              <van-button type="primary" plain native-type="submit" :loading="buttonLoading">
                提交回复
              </van-button>
            </div>
          </van-form>
        </van-popup>
      </van-popup>
      <!--加载更多帖子评论-->
      <van-button v-if="visibleComments.length < comments.length" @click="showAllComments()">
        {{ allComments ? '折叠评论' : '展开全部评论' }}
      </van-button>
    </div>
  </div>

</template>

<script>

import { mapState, mapActions } from 'vuex';
import { len } from 'vuelidate/lib/validators/common';
import { Picker } from 'emoji-mart-vue';
import { ImagePreview } from 'vant';

export default {
  components: {
    Picker,
    [ImagePreview.Component.name]: ImagePreview.Component,
  },
  data() {
    return {
      // currentPcommentID和currentCcommentID是用来进行通知的跳转
      currentPcommentID: 0,
      currentCcommentID: 0,
      before: '',
      sortkind: 'Date',
      allComments: false,
      userTelephone: '',
      partition: '',
      showDeleteModal: false,
      showReportModal: false,
      showccmenu: false,
      reportReason: '',
      fileList: [],
      dialogVisible: false,
      post: {
        postID: 0,
        author: '',
        authorTelephone: '',
        title: '',
        content: '',
        like: '',
        comment: '',
        postTime: '',
        isSaved: '',
        isLiked: '',
        showMenu: false,
        showCommentForm: false,
        photos: '',
        browse: '',
      },
      comments: [],
      pcomment: {
        userTelephone: '',
        postID: '',
        content: '',
      },
      ccomment: {
        userTelephone: '',
        postID: 0,
        pcommentID: 0,
        content: '',
        userTargetName: '',
        ccommentID: 0,
      },
      // 用来记录是否悬停在子评论上
      isHovered: false,
      nowSubIndex: 0,
      nowIndex: 0,
      replyshow: false,
      // 用来修复发表评论后页面的跳转问题（尚未修复）
      showcommentsindex: 0, // 当先评论的回复所对应的帖子评论
      nowReplyComment: -1, // 当前想要回复的评论的评论
      showRepliesModal: false, // 显示窗口
      showEmoji: false,
      buttonLoading: false,
      showloading: true,
    };
  },
  computed: {
    ...mapState({
      userInfo: (state) => state.userModule.userInfo,
      tagTypeMap() {
        return {
          大厂: 'primary',
          高工资: 'success',
          实习: 'danger',
        };
      },
    }),
    commentsNum() {
      let num = len(this.comments);
      for (let i = 0; i < this.comments.length; i += 1) {
        const comment = this.comments[i];
        num += len(comment.subComments);
      }
      return num;
    },
    fileListGet() {
      console.log(this.post.photos.split('|'));
      if (this.post.photos === '') return [];
      return this.post.photos.split('|');
    },
    isNightStyle() {
      if (JSON.parse(localStorage.getItem('Style')) === 'night') {
        return true;
      }
      return false;
    },
    // 计算属性，根据当前展示的评论数和每次展示的评论数量，返回可见的评论
    visibleComments() {
      if (this.allComments === true) {
        return this.comments;
      }
      return this.comments.slice(0, 10);
    },
  },

  created() {
    if (typeof localStorage !== 'undefined') {
      if (!localStorage.getItem('Style')) {
        localStorage.setItem('Style', JSON.stringify('day'));
      }
    }
    console.log(this.$route.query);
    if (this.$route.query.before) {
      this.before = this.$route.query.before;
      // 将postID保存在本地缓存中
      localStorage.setItem('Before', JSON.stringify(this.$route.query.before));
    } else {
      // 在本地缓存在直接读取postID
      this.before = JSON.parse(localStorage.getItem('Before'));
    }
    if (this.$route.query.id) {
      this.post.postID = parseInt(this.$route.query.id, 10);
      // 将postID保存在本地缓存中
      localStorage.setItem('PostID', JSON.stringify(this.$route.query.id));
    } else {
      // 在本地缓存在直接读取postID
      this.post.postID = JSON.parse(localStorage.getItem('PostID'));
    }
    if (this.$route.query.partition) {
      this.partition = this.$route.query.partition;
      // 将partition保存在本地缓存中
      localStorage.setItem('Partition', JSON.stringify(this.$route.query.partition));
    } else {
      this.partition = '';
    }
    this.userTelephone = this.userInfo.phone;
    // 获取当前评论ID
    this.currentPcommentID = this.$route.query.pcommentID;
    this.currentCcommentID = this.$route.query.ccommentID;
    // 根据该id向后端发送请求，获取该帖子的详细信息，并展示在页面上
    this.postShowDetails({
      userTelephone: this.userTelephone,
      postID: this.post.postID,
    })
      .then((post) => {
        this.post.postID = post.data.PostID;
        this.post.author = post.data.UserName;
        this.post.authorTitle = this.getUserTitle(post.data.UserScore);
        this.post.authorTelephone = post.data.UserTelephone;
        this.post.authorAvatar = post.data.UserAvatar;
        this.post.authorIdentity = post.data.UserIdentity;
        this.post.title = post.data.Title;
        this.post.content = post.data.Content;
        this.post.like = post.data.Like;
        this.post.comment = post.data.Comment;
        this.post.postTime = post.data.PostTime;
        this.post.tag = post.data.Tag ? post.data.Tag.split(',')
          .map((tagText) => ({
            type: this.tagTypeMap[tagText.trim()],
            label: tagText.trim(),
          })) : [];
        this.post.isSaved = post.data.IsSaved;
        this.post.isLiked = post.data.IsLiked;
        this.post.showMenu = false;
        this.post.photos = post.data.Photos;
        this.post.browse = post.data.Browse;
        this.showloading = false;
      })
      .catch((err) => {
        this.showloading = false;
        this.$toast.fail(`加载失败\n${err.response.data.msg}`);
        console.error(err.msg);
      });
    this.pcommentsShow();
    // 这里或许有比setTimeout更好的写法，但是暂时写不出来，
    // 所以先用setTimeout的方法来确保pcommentsShow执行完成之后再执行this.scrollToComment()
    // 已解决，在pcommentsShow中实现了
    // setTimeout(() => {
    //   // this.scrollToComment();
    // }, 500);
  },
  methods: {
    addEmojiToCcomment(emoji) {
      const textarea = document.getElementById('ccommentInput');// Get the textarea element
      const startPos = textarea.selectionStart; // Get the cursor's start position
      const endPos = textarea.selectionEnd; // Get the cursor's end position
      // Insert the emoji at the cursor position
      this.ccomment.content = this.ccomment.content.slice(0, startPos)
          + emoji.native + this.ccomment.content.slice(endPos);

      // Update the cursor position to be after the inserted emoji
      const newCursorPos = startPos + emoji.native.length;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
      // this.ccomment.content += emoji.native;
    },
    addEmojiToPcomment(emoji) {
      const textarea = document.getElementById('pcommentInput');
      // const textarea = this.$refs.pcommentTextarea.$el; // Get the textarea element
      const startPos = textarea.selectionStart; // Get the cursor's start position
      const endPos = textarea.selectionEnd; // Get the cursor's end position

      // Insert the emoji at the cursor position
      this.pcomment.content = this.pcomment.content.slice(0, startPos)
          + emoji.native + this.pcomment.content.slice(endPos);

      // Update the cursor position to be after the inserted emoji
      const newCursorPos = startPos + emoji.native.length;
      textarea.setSelectionRange(newCursorPos, newCursorPos);
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
    getTagColor(authorTitle) {
      switch (authorTitle) {
        case '菜鸟':
          return '#007BFF';
        case '大虾':
          return '#28A745';
        case '码农':
          return '#DC3545';
        case '程序猿':
          return '#FF6600';
        case '工程师':
          return '#FFC107';
        case '大牛':
          return '#17A2B8';
        case '专家':
          return '#C428eb';
        case '大神':
          return '#6C757D';
        case '祖师爷':
          return '#343A40';
        default:
          return 'default-color';
      }
    },
    cclike(index, subIndex) {
      this.ccommentlike({
        userTelephone: this.userTelephone,
        ccommentID: this.comments[index].subComments[subIndex].ccommentID,
        isLiked: this.comments[index].subComments[subIndex].isLiked,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
      // 更新点赞状态及点赞数
      // eslint-disable-next-line max-len
      this.comments[index].subComments[subIndex].isLiked = !this.comments[index].subComments[subIndex].isLiked;
      // eslint-disable-next-line max-len
      if (this.comments[index].subComments[subIndex].isLiked) {
        this.comments[index].subComments[subIndex].likeNum += 1;
      } else {
        this.comments[index].subComments[subIndex].likeNum -= 1;
      }
    },
    // 发表评论的评论或者回复评论的评论
    ccommentPost(index, author, ccommentID) {
      this.buttonLoading = true; // 按钮进入加载状态
      const comment = this.comments[index];
      this.ccomment.postID = this.post.postID;
      this.ccomment.pcommentID = comment.pcommentID;
      this.ccomment.userTelephone = this.userTelephone;
      this.ccomment.userTargetName = author;
      this.ccomment.ccommentID = ccommentID;
      this.postCcomment(this.ccomment)
        .then(() => {
          // 弹窗提示
          this.$toast.success('回复成功');
          setTimeout(() => {
            this.nowReplyComment = '';
            this.replyshow = false;
            this.pcommentsShow();
            // 清空输入的内容
            this.ccomment.content = '';
            this.ccomment.ccommentID = 0;
            this.buttonLoading = false; // 按钮恢复
          }, 1000);
        })
        .catch((err) => {
          this.$toast.fail(`回复失败${err.response.data.msg}`);
          this.buttonLoading = false; // 按钮恢复
        });
    },
    ccommentdelete(SubComment) {
      this.deleteCcomment({
        ccommentID: SubComment.ccommentID,
      })
        .then(() => {
          this.$toast.success('删除成功');
          this.$router.go(0);
        })
        .catch((err) => {
          console.error(err);
        });
    },
    clearReportReason() {
      this.reportReason = '';
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
    goback() {
      // 返回采用直接关闭的形式
      // window.close();
      // eslint-disable-next-line no-restricted-globals
      history.back();
      localStorage.setItem('NowPostID', JSON.stringify(this.post.postID));
      console.log(this.before);
      if (this.before === 'home') {
        this.$router.replace({
          name: 'home',
          query: { partition: this.partition },
        });
      } else if (this.before === 'save') {
        this.$router.replace({ name: 'save' });
      } else if (this.before === 'history') {
        this.$router.replace({ name: 'history' });
      } else if (this.before === 'notice') {
        this.$router.replace({ name: 'notice' });
      } else if (this.before === 'discover') {
        this.$router.replace({ name: 'discover' });
      }
    },
    handleKeyboardEvent() {
      // 处理键盘事件，即使是一个空的处理程序
    },
    handlePictureCardPreview(url) {
      console.log(url);
      ImagePreview({
        images: [this.getImage(url)],
        closeable: true,
      });
    },
    getImage(url) {
      const replacedUrl = url.replace('/resized/', '/uploads/');
      return replacedUrl;
    },
    len,
    like() {
      const userTelephone = this.userInfo.phone;
      // 请求
      this.postLike({
        userTelephone,
        postID: this.post.postID,
        isLiked: this.post.isLiked,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
      // 更新点赞状态及点赞数
      this.post.isLiked = !this.post.isLiked;
      if (this.post.isLiked) {
        this.post.like += 1;
      } else {
        this.post.like -= 1;
      }
    },
    pclike(index) {
      this.pcommentlike({
        userTelephone: this.userTelephone,
        pcommentID: this.comments[index].pcommentID,
        isLiked: this.comments[index].isLiked,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
      // 更新点赞状态及点赞数
      this.comments[index].isLiked = !this.comments[index].isLiked;
      if (this.comments[index].isLiked) {
        this.comments[index].likeNum += 1;
      } else {
        this.comments[index].likeNum -= 1;
      }
    },
    // 发表帖子评论
    pcommentPost() {
      this.buttonLoading = true; // 按钮进入加载状态
      this.pcomment.postID = this.post.postID;
      this.pcomment.userTelephone = this.userTelephone;
      this.postPcomment(this.pcomment)
        .then(() => {
          this.$toast.success('评论成功');
          setTimeout(() => {
            this.pcommentsShow();
            this.pcomment.content = '';
            this.buttonLoading = false; // 按钮恢复
          }, 1000);
        })
        .catch((err) => {
          this.$toast.fail(`评论失败${err.response.data.msg}`);
          this.buttonLoading = false; // 按钮恢复
        });
    },
    pcommentdelete(comment) {
      this.deletePcomment({
        pcommentID: comment.pcommentID,
      })
        .then(() => {
          this.$toast.success('删除成功');
          this.$router.go(0);
        })
        .catch((err) => {
          console.error(err);
        });
    },
    pcommentsShow() {
      const postid = this.post.postID;
      // 请求
      try {
        this.showPcomments({
          userTelephone: this.userTelephone,
          postID: postid,
        })
          .then(({ data }) => {
            const comments = data.map((pcomment) => ({
              pcommentID: pcomment.PcommentID,
              author: pcomment.Author,
              authorTitle: this.getUserTitle(pcomment.AuthorScore),
              authorAvatar: pcomment.AuthorAvatar,
              authorIdentity: pcomment.AuthorIdentity,
              authorTelephone: pcomment.AuthorTelephone,
              commentTime: pcomment.CommentTime,
              content: pcomment.Content,
              likeNum: pcomment.LikeNum,
              subComments: pcomment.SubComments.map((subComment) => ({
                ...subComment,
                authorTitle: this.getUserTitle(subComment.authorScore),
              })),
              isLiked: pcomment.IsLiked,
              showMenu: false,
              showReplyForm: false,
              showAllReplies: false,
              heat: pcomment.LikeNum + len(pcomment.SubComments),
            }));
            this.comments = this.sortcomments(comments);
            // .sort((a, b) => new Date(b.commentTime) - new Date(a.commentTime))
            // this.comments = data;
          })
          .catch((err) => {
            console.error(err);
          })
          .then(() => {
            this.scrollToComment();
          })
          .catch((err) => {
            console.error(err);
          });
      } catch (err) {
        console.error(err);
      }
    },
    postdelete() {
      this.deletepost({
        postID: this.post.postID,
      })
        .then(() => {
          this.$toast.success('删除成功');
          this.$router.go(-1);
        })
        .catch((err) => {
          this.$toast.fail(`删除失败${err.response.date.msg}`);
          console.error(err);
        });
    },
    save() {
      const userTelephone = this.userInfo.phone;
      // 请求
      this.postSave({
        userTelephone,
        postID: this.post.postID,
        isSaved: this.post.isSaved,
      })
        .then(() => {
        })
        .catch((err) => {
          console.error(err);
        });
      // 更新点赞状态及点赞数
      this.post.isSaved = !this.post.isSaved;
    },
    scrollToComment() {
      console.log('scrollToComment');
      // 获取当前评论所在的元素
      let commentEl = document.getElementById(`comment-${this.currentPcommentID}`);
      console.log(this.currentPcommentID);
      console.log(commentEl);
      // const commentRef = this.$refs.commentRef[3];
      // if (commentRef) {
      //   commentRef.setAttribute('tabindex', '-1');
      //   commentRef.scrollIntoView({ behavior: 'smooth', duration: 500 });
      // }
      if (commentEl === null) {
        this.allComments = true;
        const get = new Promise((resolve, reject) => {
          setTimeout(() => {
            commentEl = document.getElementById(`comment-${this.currentPcommentID}`);
            if (commentEl) {
              resolve(commentEl);
            } else {
              reject(commentEl);
            }
          });
        });
        get.then(() => {
          // 使用vue-scrollto插件平滑滚动到元素所在位置
          this.$scrollTo(commentEl, {
            duration: 750, // 滚动动画持续时间，单位为毫秒
            offset: -150, // 滚动偏移量，用于调整滚动位置
          });
          commentEl.classList.add('blink');
          setTimeout(() => {
            commentEl.classList.remove('blink');
          }, 10000);
          if (this.currentCcommentID) {
            const parentEl = commentEl.parentNode;
            this.showcommentsindex = Array.prototype.indexOf.call(parentEl.children, commentEl);
            this.showRepliesModal = true;
            setTimeout(() => {
              let childEl = document.getElementById(`ccomment-${this.currentCcommentID}`);
              console.log(childEl, this.currentCcommentID);
              if (childEl === null) {
                this.showAllReplies(this.showcommentsindex);
                setTimeout(() => {
                  childEl = document.getElementById(`ccomment-${this.currentCcommentID}`);
                });
              }
              if (childEl) {
                // 使用vue-scrollto插件平滑滚动到元素所在位置
                childEl.focus();
                childEl.classList.add('blink');
                setTimeout(() => {
                  childEl.classList.remove('blink');
                }, 10000);
              }
            }, 500);
          }
        })
          .catch((err) => {
            console.log(err);
          });
      } else {
        // 使用vue-scrollto插件平滑滚动到元素所在位置
        this.$scrollTo(commentEl, {
          duration: 750, // 滚动动画持续时间，单位为毫秒
          offset: -150, // 滚动偏移量，用于调整滚动位置
        });
        commentEl.classList.add('blink');
        setTimeout(() => {
          commentEl.classList.remove('blink');
        }, 10000);
        if (this.currentCcommentID) {
          const parentEl = commentEl.parentNode;
          this.showcommentsindex = Array.prototype.indexOf.call(parentEl.children, commentEl);
          this.showRepliesModal = true;
          setTimeout(() => {
            let childEl = document.getElementById(`ccomment-${this.currentCcommentID}`);
            console.log(childEl, this.currentCcommentID);
            childEl.focus();
            childEl.classList.add('blink');
            setTimeout(() => {
              childEl.classList.remove('blink');
            }, 10000);
            if (childEl === null) {
              this.showAllReplies(this.showcommentsindex);
              const get = new Promise((resolve, reject) => {
                setTimeout(() => {
                  childEl = document.getElementById(`ccomment-${this.currentCcommentID}`);
                  if (childEl) {
                    resolve();
                  } else {
                    reject();
                  }
                });
              });
              get.then(() => {
                setTimeout(() => {
                  console.log('移动到对应位置');
                  childEl.focus();
                  childEl.classList.add('blink');
                  setTimeout(() => {
                    childEl.classList.remove('blink');
                  }, 10000);
                }, 500);
              })
                .catch((error) => {
                  console.log(error);
                });
            }
          });
        }
      }
    },
    showAllComments() {
      this.allComments = !this.allComments;// 将帖子所有评论都展示出来
    },
    ...mapActions('postModule', { postShowDetails: 'showDetails' }),
    ...mapActions('postModule', { postLike: 'like' }),
    ...mapActions('userModule', { postSave: 'save' }),
    ...mapActions('postModule', { deletepost: 'deletepost' }),
    ...mapActions('postModule', { submitreport: 'submitreport' }),
    ...mapActions('commentModule', { showPcomments: 'showPcomments' }),
    ...mapActions('commentModule', { deletePcomment: 'deletePcomment' }),
    ...mapActions('commentModule', { deleteCcomment: 'deleteCcomment' }),
    ...mapActions('commentModule', { postPcomment: 'postPcomment' }),
    ...mapActions('commentModule', { postCcomment: 'postCcomment' }),
    ...mapActions('commentModule', { pcommentlike: 'pcommentLike' }),
    ...mapActions('commentModule', { ccommentlike: 'ccommentLike' }),
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
    showAllReplies(index) {
      // 将评论的所有子评论都显示出来
      this.comments[index].showAllReplies = !this.comments[index].showAllReplies;
    },
    showEmojiStatus() {
      this.showEmoji = !this.showEmoji;
    },
    sortcomments(comments) {
      if (this.sortkind === 'Date') {
        return comments.sort((a, b) => new Date(b.commentTime) - new Date(a.commentTime));
      }
      return comments.sort((a, b) => b.heat - a.heat);
    },
    submitReport(type, idnum) {
      this.submitreport({
        TargetID: idnum,
        Targettype: type,
        userTelephone: this.userInfo.phone,
        Reason: this.reportReason,
      })
        .then(() => {
          this.$toast.success('举报发送成功');
        })
        .catch((err) => {
          this.$toast.fail(`举报发送失败${err.response.data.msg}`);
        });
    },
    toggleMenu() {
      this.post.showMenu = !this.post.showMenu;
    },
    visibleSubComments(index) {
      if (this.comments[index] === undefined) return [];
      if (this.comments[index].showAllReplies === true) {
        return this.comments[index].subComments;
      }
      return this.comments[index].subComments;
    },
  },
  beforeRouteLeave(to, from, next) {
    // 返回上一页面时清空本地缓存
    localStorage.removeItem('PostID');
    localStorage.removeItem('Before');
    next();
  },
};
</script>

<style lass="scss" scoped>
@import '@/style/css/PostDetailsView.css';

#div_boxs_lnteado /deep/ .van-nav-bar__title {
  color: white;
  font-size: 0.45rem;
}

.thumbnail-container {
  display: flex;
  flex-wrap: wrap;
}

.thumbnail-container div {
  width: calc(100% / 3);
  padding: 10px;
  box-sizing: border-box;
}

@keyframes blink {
  0% {
    opacity: 1;
    background-color: transparent;
  }
  50% {
    opacity: 0.5;
    background-color: rgb(128, 128, 128);
  }
  100% {
    opacity: 1;
    background-color: transparent;
  }
}

/* 使用更具体的选择器 */
.blink {
  animation-name: blink;
  animation-duration: 1.0s;
  animation-iteration-count: 3;
  background-color: transparent !important; /* 覆盖框架中的样式 */
}

.emoji-mart[data-v-7bc71df8] {
  font-family: -apple-system, BlinkMacSystemFont, "Helvetica Neue", sans-serif;
  display: -ms-flexbox;
  display: flex;
  -ms-flex-direction: column;
  flex-direction: column;
  height: 300px;
  color: #ffffff !important;
  border: 1px solid #d9d9d9;
  border-radius: 5px;
  background: #fff;
}

.post_title {
  vertical-align: middle;
  text-align: center;
  font-size: 35px;
  font-weight: bold;
  margin-top: 10px;
}

.post_content {
  vertical-align: middle;
  font-size: 30px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: anywhere;
}

.tag-group {
  margin-left: 500px;
  margin-top: -50px;
}

.photo {
  width: 200px;
  height: 200px;
}
</style>
