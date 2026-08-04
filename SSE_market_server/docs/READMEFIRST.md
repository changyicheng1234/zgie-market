主要实现登录、注册、修改密码以及邮箱发送验证码操作

后端代码中验证码实现部分需要用到redis

包括common/redis.go, userController, main.go

如果想要实现邮箱发送验证码的功能，可以将someConst/mailConst中的EmailUsername和Password进行更改


Password使用的是授权码，不是密码