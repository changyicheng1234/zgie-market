![image-20240402094632607](https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/blog/image-20240402094632607.png)

# 软工集市，软工人定义的世界

### 项目介绍

SSE_MARKET是一个跨校区的学院内部交流平台，主要以论坛的形式为软工师生提供自由交流和信息汇集的半匿名空间，可在上面看帖、搜索、发帖、评论、回复、点赞等，提供宽松、平等的交流氛围，致力于消除学院内部的信息差以及促进数字文化传承。

SSE_MARKET又称软工集市，主要由21级本科生组成的 SSE_MARKET小组负责设计、开发、部署和维护，它脱胎于软件中级实训课，并在学院的支持下逐步发展，现注册加入师生约 400人，学生包括大一到大四本科生以及各级研究生，教师包括行政老师、技术老师到专业老师等。

[软工集市成功交接，欢迎新任成员 – SSE_MARKET博客 (ssemarket.cn)](https://ssemarket.cn/2024/04/02/软工集市成功交接，欢迎新任成员/)

在2024.3.6，软工集市正式完成交接。现在的软工集市由22级和23级本科生组成的新小组负责开发、部署、优化、维护。

[官方博客](https://ssemarket.cn)

### 视频讲解

目前，软工集市现任技术顾问thinkerhui主导的B站软工集市教程已上线，将持续更新。

[thinkerhui的个人空间-thinkerhui个人主页-哔哩哔哩视频 (bilibili.com)](https://space.bilibili.com/448341656/channel/collectiondetail?sid=2596476)

![image-20240402095909217](https://sse-market-source-1320172928.cos.ap-guangzhou.myqcloud.com/blog/image-20240402095909217.png)


### 如何运行软工集市后端

SSEMARKET后端为go项目，需要配置go开发的基本环境。此外，还需要配置mysql和redis的数据库环境。

1. 首先需要在开发环境安装go语言，理论上安装最新版本就行。
   [Download and install - The Go Programming Language](https://golang.google.cn/doc/install)

2. 选择自己喜欢的IDE配置开发环境，这里推荐vscode,goland

3. 安装go项目依赖

4. 安装配置mysql8.0

5. 安装配置redis

6. 在`config/application.yml`改数据库配置

7. ```shell
   go run main.go
   ```

   输入该命令即可运行项目。

注意事项：1、2 、4、5主要是在配置基本环境，如果本来配置好了可以直接跳过。
这里并没有给出具体的安装配置步骤，而是直接贴出了参考教程，这个一方面是由于Windows、MacOS和linux等操作系统的配置过程会有所不同，因此给出教程也不一定适用（比如mysql，redis的教程是适用windows的）。

___

### 各个文件夹的作用

>## <font size=3> 1.`common`文件夹用于存放一些通用的功能模块或工具函数，如数据库 <br> 2.`config`文件夹用于存放配置文件 <br> 3.`controller`文件夹用于存放应用程序的控制器代码。控制器是应用程序的核心部分之一，它负责接收客户端请求并作出响应 <br> 4.`middleware`文件夹用于存放中间件代码，例如日志记录、认证、CORS等 <br> 5.`model`文件夹用于存放应用程序的数据模型定义，通常使用gorm库来实现对象关系映射 <br> 6.`util`文件夹用于存放工具函数文件 <br> 7.`main.go`是Gin应用程序的主入口文件 <br> 8.`routes.go`是Gin应用程序的路由定义文件

___

># <font size=4> 三、如何创建一个需求处理函数，以注册功能为例

>## <font size=3> 1.首先在`controller`里添加对于函数，函数名大写如Register(),函数具体写法参见`controller/userController.go`的代码注释<br> 2.然后在`routes`中创建一个路由，并对接到`controller`对应的函数，如：
``` go
r.POST("/api/auth/register", controller.Register)
```