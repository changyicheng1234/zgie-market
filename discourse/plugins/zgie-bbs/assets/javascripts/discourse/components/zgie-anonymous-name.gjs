import Component from "@glimmer/component";

export default class ZgieAnonymousName extends Component {
  static shouldRender(args) {
    return args.post?.zgie_anonymous_author === true;
  }

  <template>
    <span class="zgie-anonymous-name">{{@outletArgs.name}}</span>
  </template>
}
