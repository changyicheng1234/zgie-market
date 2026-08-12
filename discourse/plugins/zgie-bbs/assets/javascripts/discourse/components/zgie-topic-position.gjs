import Component from "@glimmer/component";
import { tracked } from "@glimmer/tracking";
import { service } from "@ember/service";
import { trustHTML } from "@ember/template";
import { modifier } from "ember-modifier";
import { i18n } from "discourse-i18n";

const POSITIONABLE_POSTS = [
  ".nested-view__op-article[data-post-number]",
  ".nested-post__article[data-post-number]",
  ".nested-post__collapsed-bar[data-post-number]",
  ".nested-post__placeholder[data-post-number]",
].join(",");

export default class ZgieTopicPosition extends Component {
  static shouldRender(args) {
    return args.model?.is_nested_view === true;
  }

  @service header;

  @tracked compact = false;
  @tracked current = 1;
  @tracked currentDate;
  @tracked dockLeft;

  setup = modifier((element) => {
    this._element = element;
    const nestedView = element.closest(".nested-view");
    if (!nestedView) {
      this._element = null;
      return;
    }

    window.addEventListener("scroll", this._onViewportChange, {
      passive: true,
    });
    window.addEventListener("resize", this._onViewportChange, {
      passive: true,
    });

    this._mutationObserver = new MutationObserver(this._onViewportChange);
    this._mutationObserver.observe(nestedView, {
      childList: true,
      subtree: true,
    });

    if (window.ResizeObserver) {
      this._resizeObserver = new ResizeObserver(this._onViewportChange);
      this._resizeObserver.observe(nestedView);
    }

    this._scheduleUpdate();

    return () => {
      window.removeEventListener("scroll", this._onViewportChange);
      window.removeEventListener("resize", this._onViewportChange);
      this._mutationObserver?.disconnect();
      this._resizeObserver?.disconnect();
      cancelAnimationFrame(this._frame);
      this._element = null;
    };
  });

  _onViewportChange = () => this._scheduleUpdate();

  get topic() {
    return this.args.outletArgs.model;
  }

  get total() {
    return Math.max(
      1,
      Number(this.topic?.posts_count) || 0,
      Number(this.topic?.highest_post_number) || 0
    );
  }

  get ariaLabel() {
    return i18n("zgie_bbs.topic_position", {
      current: this.current,
      total: this.total,
    });
  }

  get progressStyle() {
    const ratio =
      this.total <= 1
        ? 0
        : (Math.min(this.current, this.total) - 1) / (this.total - 1);
    const visualPercentage = 5 + Math.max(0, Math.min(1, ratio)) * 90;
    const declarations = [`--zgie-topic-position: ${visualPercentage}%`];
    if (!this.compact && this.dockLeft != null) {
      declarations.push(`left: ${this.dockLeft}px`);
    }

    return trustHTML(declarations.join("; "));
  }

  get classes() {
    return this.compact
      ? "zgie-topic-position --compact"
      : "zgie-topic-position";
  }

  get startedDate() {
    return this._formatDate(this.topic?.created_at);
  }

  get lastDate() {
    return this._formatDate(this.topic?.last_posted_at);
  }

  _scheduleUpdate() {
    if (this._frame) {
      return;
    }

    this._frame = requestAnimationFrame(() => {
      this._frame = null;
      this._updateDockPosition();
      this._updateCurrentPost();
    });
  }

  _updateDockPosition() {
    if (!this._element) {
      return;
    }

    const nestedView = this._element.closest(".nested-view");
    if (!nestedView) {
      return;
    }

    const topicContent = nestedView.querySelector(
      ".nested-view__op, .nested-view__controls, .nested-view__roots"
    );
    const viewRect = (topicContent || nestedView).getBoundingClientRect();
    const gap = 24;
    const requiredWidth = 128;
    const availableWidth = window.innerWidth - viewRect.right - gap;
    const compact = availableWidth < requiredWidth;

    this.compact = compact;
    this.dockLeft = compact ? null : Math.round(viewRect.right + gap);
  }

  _updateCurrentPost() {
    const posts = this._postsInReadingOrder();
    if (!posts.length) {
      return;
    }

    const headerOffset = Number(this.header?.headerOffset) || 60;
    const readingLine = Math.min(
      window.innerHeight * 0.36,
      headerOffset + 180
    );
    const post = posts.reduce((closest, candidate) => {
      const rect = candidate.getBoundingClientRect();
      const distance =
        rect.top <= readingLine && rect.bottom >= readingLine
          ? 0
          : Math.min(
              Math.abs(rect.top - readingLine),
              Math.abs(rect.bottom - readingLine)
            );

      return !closest || distance < closest.distance
        ? { element: candidate, distance }
        : closest;
    }, null)?.element;

    const readingPosition = posts.indexOf(post) + 1;
    if (readingPosition > 0) {
      this.current = readingPosition;
      this.currentDate = this._dateFromPost(post);
    }
  }

  _postsInReadingOrder() {
    const nestedView = this._element?.closest(".nested-view");
    if (!nestedView) {
      return [];
    }

    const seenPostNumbers = new Set();
    return [...nestedView.querySelectorAll(POSITIONABLE_POSTS)].filter(
      (post) => {
        const postNumber = Number(post.dataset.postNumber);
        if (postNumber <= 0 || seenPostNumbers.has(postNumber)) {
          return false;
        }

        seenPostNumbers.add(postNumber);
        return true;
      }
    );
  }

  _dateFromPost(post) {
    const dateElement = post.querySelector("[data-time]");
    const rawDate = dateElement?.dataset.time;
    if (!rawDate) {
      return;
    }

    const numericDate = Number(rawDate);
    const date = Number.isFinite(numericDate)
      ? new Date(
          numericDate < 1_000_000_000_000 ? numericDate * 1000 : numericDate
        )
      : new Date(rawDate);
    return this._formatDate(date);
  }

  _formatDate(value) {
    if (!value) {
      return "";
    }

    const date = value instanceof Date ? value : new Date(value);
    if (Number.isNaN(date.getTime())) {
      return "";
    }

    return new Intl.DateTimeFormat(document.documentElement.lang || "zh-CN", {
      month: "long",
      day: "numeric",
    }).format(date);
  }

  <template>
    <aside
      class={{this.classes}}
      aria-label={{this.ariaLabel}}
      style={{this.progressStyle}}
      {{this.setup}}
    >
      <time
        class="zgie-topic-position__date --start"
        datetime={{this.topic.created_at}}
      >{{this.startedDate}}</time>

      <div class="zgie-topic-position__rail">
        <div class="zgie-topic-position__handle">
          <span class="zgie-topic-position__thumb" aria-hidden="true"></span>
          <span class="zgie-topic-position__meta">
            <strong class="zgie-topic-position__count" aria-live="polite">
              {{this.current}} / {{this.total}}
            </strong>
            {{#if this.currentDate}}
              <span class="zgie-topic-position__current-date">
                {{this.currentDate}}
              </span>
            {{/if}}
          </span>
        </div>
      </div>

      <time
        class="zgie-topic-position__date --end"
        datetime={{this.topic.last_posted_at}}
      >{{this.lastDate}}</time>
    </aside>
  </template>
}
