const { defineConfig } = require('@vue/cli-service');
// const CompressionPlugin = require('compression-webpack-plugin');//引入gzip压缩插件

module.exports = defineConfig({
  transpileDependencies: true,
});
const path = require('path');
const fs = require('fs');

module.exports = {
  devServer: {
    host: '',
    port: 8081,
    open: true,
    // https: {
    //   cert: fs.readFileSync(path.join(__dirname, 'src/ssl/cert.crt')),
    //   key: fs.readFileSync(path.join(__dirname, 'src/ssl/cert.key')),
    // },
  },
  publicPath: process.env.NODE_ENV === 'production'
    ? '/pc/'
    : '/',
  // configureWebpack: {
  //   plugins: [
  //     // gizp压缩
  //     new CompressionPlugin()({
  //       test: /\.(js|css|json|txt|html|ico|svg)(\?.*)?$/i,
  //       threshold: 10240, // 对10K以上的数据进行压缩
  //       deleteOriginalAssets: true
  //     })
  //   ]
  // },
};
