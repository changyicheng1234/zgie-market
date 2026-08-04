![image-20240402094632607](https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/blog/image-20240402094632607.png)

# 软工集市，软工人定义的世界

### 项目介绍

见[README.md · Moochen/SSE_market_server](https://gitee.com/yang-peiyue/SSE_market_server/blob/master/README.md)

### 如何运行软工集市电脑前端

软工集市采用前后端分离的架构，所以需要配置node.js

1. 首先要配置<font color="red">node.js、npm（或yarn)</font>
2. 修改`.env.development`的VUE_APP_BASE_URL为你的后端地址
3. 输入`yarn`或`npm install`来安装vue项目的所有依赖
4. 使用<font color="red">yarn serve</font>或者<font color="red">npm run serve</font>命令运行vue项目

我们的项目yarn和npm在开发中混合使用，这并不是推荐的做法，建议只使用其中一个。

### 文件存放位置

 1.vue文件统一放在`src/views`里，然后创建一个文件夹，统一命名为小写开头比如说`register`，在文件夹里相应的放一个vue文件，统一为大写对应的驼峰命名法，以View结尾，比如上面`register`对应的就是`RegisterView.vue`  

2.css文件统一放在`src/style/css`里   

3.scss文件统一放在`src/style/scss`里  

4.图片资源统一放在`src/assets/image`里

### 各个文件夹的作用

1.`assets`文件夹用于放图片、字体等的静态资源文件  

2.`components`文件夹用于存放公共组件，可复用的组件应该放在这里。每个组件一般包含一个.vue文件、一个.js文件和一个.css文件  3.`helper`文件夹中用于存放一些辅助函数或工具类代码。这些函数和类通常不需要被直接导入到Vue组件中，而是可以在整个应用程序中共享  

4.`router`文件夹用于存放路由配置文件，通常是`index.js`文件。module里是封装了一些路由，如`user.js`封装了用户所用到的路由，如注册和登录  

5.`service`文件夹用于存放一些处理服务的文件。其中`storageService.js`文件用户处理本地缓存服务，`userService.js`里定义了与用户相关的给后端发送消息的服务函数  

6.`store`文件夹用于存放Vuex全局状态管理相关的文件，通常是`index.js`文件。module里的`user.js`定义了与用户相关的从后端接收消息的状态函数  

7.style文件夹用户存放css、scss等样式文件  

8.`utils`文件夹用于存放可复用的存放工具类函数或库，`request.js`中封装了axios(前后端数据交互)的函数  

9.`views`文件夹用户存放页面视图组件，每个页面对应一个vue文件。其中`layout/NavbarView.vue`用于写固定的外部框架，包含应用的头部、侧边栏、底部和主题内容等。并在`App.vue`中引用  

10.`App.vue`是Vue.js单文件组件中的根组件，主要用于定义整个应用程序的外部框架，包含应用的头部、侧边栏、底部和主题内容等。在`App.vue`中，可以定义全局的样式和共享的状态数据，以及处理路由导航和其他应用级别的事件  

11.`main.js`是Vue.js应用程序的入口文件，在main.js中，我们首先需要引入Vue框架和所需的插件和组件。然后创建Vue实例，并指定一个根组件，通常是`App.vue`组件。接着通过调用Vue实例的`$mount`方法将根组件挂载到一个HTML元素上，使其成为应用程序的显示部分  

12.下方的.`env.development`文件和`.env.development.local`文件是配置环境变量的。比如在后端中定义的路由是`http://localhost:8080/api/auth/register`、`http://localhost:8080/api/auth/login`等等我们把前面共同的前缀`http://localhost:8080/api`写在环境变量中，那么我们在`userService`中定义的时候只需要写后面，比如:

``` javascript
return request.post('auth/register', { name, telephone, password });
```

### 如何创建一个页面并与后端实现数据交互

 我们假设要新建一个personalHomepage的个人信息页面，需要放回给后端的信息是用户的telephone,而需要从后端获取返回的关于user结构体的全部信息

1.首先我们根据命名规则，在`src/views`中创建一个`personalHomepage`文件夹，再在其中创建一个`PersonalHomepageView.vue`的vue文件  2.然后vue文件大致可分为以下三个部分，具体代码的注解在`RegisterView.vue`、`ProfileView.vue`、`NavbarView.vue`文件中有相应注释

 ①模板（template）：模板包含了渲染视图的 HTML 代码，用于定义页面结构和显示数据  ②脚本（script）：脚本部分包含了 Vue 实例的逻辑代码和数据对象。可以在此处定义计算属性、监听器和方法等 ③样式（style）：样式部分定义了组件的 CSS 样式，它可以是普通的 CSS 代码或者使用预处理器如 Sass 或 Less 等

 2.然后在`router/module`里找到对应文件(`user.js`或者`post.js`或者是其他)，添加相应路由，如：

``` javascript
{
    path: '/personalHomepage',
    name: 'personalHomepage',
    component: () => import('../views/personalHomepage/'PersonalHomepageView.vue'),
},
```
 具体使用方法在`router/module/user.js`的`register`进行相应注释

3.然后在`service`文件夹中找到对应文件(`userService.js`或者`postService.js`或是其他)，添加相应给后端发送消息的服务函数代码，如：

``` javascript
const personalHomepage = ({ telephone }) => {
  return request.post('auth/personalHomepage', { telephone });
};
```
具体使用方法在`service/userService.js`文件的`register`函数进行相应注释

4.最后在`store/module`里找到对应文件(`user.js`或者`post.js`或者是其他),添加相应从后端接收消息的状态函数，如：

``` javascript
personalHomepage(context, { telephone }) { 
    return new Promise((resolve, reject) => {
        userService.personalHomepage({ telephone }).then((res) => {
          resolve(res);
        }).catch((err) => {
          reject(err);
        });
    });
},
```
 具体使用在'`store/module.user.js`的`actions`里面的`register`函数进行相应注释