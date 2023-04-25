# Order_Management_System

> 正在施工中，本身为简易Demo，着重帮助掌握Go语言语法以及Gin开发框架简单使用，喜欢就点个Star吧！

## 准备工作

### 数据库

本项目数据库为`mysql-8.0.29-winx64`，数据字段如下所示：



![image-20230423214159264](README/image-20230423214159264.png)

> 提供SQL语句一键建表

~~~mysql
DROP TABLE IF EXISTS userinfo;
CREATE TABLE userinfo (
  userid INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL,
  password VARCHAR(255) NOT NULL,
  registerAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  status INT DEFAULT 1,
  isdelete INT DEFAULT 0
);

DROP TABLE IF EXISTS shops;
CREATE TABLE shops (
  shopid INT AUTO_INCREMENT PRIMARY KEY,
  shopname VARCHAR(255) NOT NULL,
  rating FLOAT NOT NULL DEFAULT 5.0
);


DROP TABLE IF EXISTS orders;
CREATE TABLE orders (
  orderid INT AUTO_INCREMENT PRIMARY KEY,
  userid INT NOT NULL,
  shopid INT NOT NULL,
  status ENUM('待支付', '待发货', '待收货', '已完成', '已取消') DEFAULT '待支付',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (userid) REFERENCES userinfo(userid),
  FOREIGN KEY (shopid) REFERENCES shops(shopid)
);

DROP TABLE IF EXISTS reviews;
CREATE TABLE reviews (
  reviewid INT AUTO_INCREMENT PRIMARY KEY,
  orderid INT NOT NULL,
  userid INT NOT NULL,
  content TEXT NOT NULL,
  rating INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (orderid) REFERENCES orders(orderid),
  FOREIGN KEY (userid) REFERENCES userinfo(userid)
);
~~~

### Go语言环境

golang版本为`go1.20.2 windows/amd64`

提供一篇写的挺好的帖子，感谢大哥的开源精神！

[vs code配置go开发环境 - 知乎 (zhihu.com)](https://zhuanlan.zhihu.com/p/262984879)



若无法下载import依赖，请在命令行输入：

~~~shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
~~~

## 运行项目

> 开发框架采用Gin（一个用 Go (Golang) 编写的 HTTP Web 框架）

~~~shell
cd Order_Management_System
go mod init github.com/jjyaoao/Order_Management_System
go mod tidy
go run .\Order_Management_System.go
~~~

## 展望

目前只是实现一些初级的功能，还有一些未来的设想可以继续补充，例如：

- 用户名的注册去重

- 数据库中密文储存用户密码
- 更佳细腻的用户登录检测
- 实现与前端页面的交互连接
- 将购物车功能与订单功能分开
- 增加物流信息，实现`待发货`与`待收货`的功能。

这里也是一些未完成的设想，鉴于本人太菜，以后会慢慢填坑的（bushi



