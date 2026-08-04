module.exports = {
  plugins: {
    'postcss-pxtorem': {
      rootValue: 75, // 设计稿宽度 / 10，假设设计稿宽度为750px
      propList: ['*', '!min-width', '!max-width'],
      selectorBlackList: ['.van'], // 把vant2排除，因为vant本来就是对移动端设计的，一般不需要postcss-pxtorem进行转换。
      replace: true,
      exclude: '/node_modules',
    },
  },
};
