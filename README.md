# 基于Beego的后台管理系统

## 后台模板：AdminLTE

### 使用步骤

1、先安装beego框架

2、获取项目
`$ go get https://github.com/huanzz/beego__admin`

3、修改 conf/app.conf 中mysql的设置

4、终端下进入到项目目录
```
$ go build      
$ beego__admin syncdb   //数据库初始化
$ bee run
```

5、访问http://localhost:8080


