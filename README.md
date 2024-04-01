# go-grpc-crud-template
✨基于 golang, grpc, gin 和 redis, MySQL, etcd 和 vue3 的简单分布式信息管理系统✨📌含完整前后端，微服务：分布式信息管理系统模板，后台管理系统模板，数据库管理系统模板。实现 grpc 微服务远程过程调用，redis 缓存，etcd 服务发现，负载均衡。令牌签验，非对称加密。通过 Web 应用完成对数据库的增删改查(CRUD)，文件流的上传和下载。📌前后端分离

📌[在线演示链接](http://bytesc.top:8009)

📌[配套前端项目地址](https://github.com/bytesc/vue-crud-template)

[个人网站：www.bytesc.top](http://www.bytesc.top) 

[个人博客：blog.bytesc.top](http://blog.bytesc.top)

🔔 如有项目相关问题，欢迎在本项目提出`issue`，我一般会在 24 小时内回复。

## 效果展示

### 文件流

文件流上传

![](./docs/readme_img/imgf1.png)

文件流下载

![](./docs/readme_img/imgf2.png)

文件列表

![](./docs/readme_img/imgfl.png)

### CRUD

![](./docs/readme_img/img1.png)

完善的查询

![](./docs/readme_img/img2.png)

多选删除

![](./docs/readme_img/img3.png)

编辑行

![](./docs/readme_img/img4.png)


### 用户登录

![](./docs/readme_img/imgu.png)

面包屑导航

![](./docs/readme_img/img7.png)
![](./docs/readme_img/img8.png)

## 项目运行方法

### 后端运行环境

- `go1.20.5`
- `MySQL 8.0.31`
- `Redis 7.2.4`
- `etcd 3.4.31`

### 安装依赖
（非必要，后续运行时候也会自动安装）
```bash
# go mod download
# go get -u gorm.io/driver/sqlite
go get -u gorm.io/driver/mysql
go get -u gorm.io/gorm
go get -u github.com/gin-gonic/gin

go get -u github.com/golang-jwt/jwt/v5

go get -u go.uber.org/zap

go get github.com/go-playground/validator/v10

go get github.com/spf13/viper
go get github.com/gin-gonic/gin/binding@v1.9.1

go get -u github.com/gin-contrib/cors

go get -u github.com/go-redis/redis/v8

go get github.com/go-redsync/redsync/v4
go get github.com/go-redsync/redsync/v4/redis/goredis/v8

go get google.golang.org/grpc
go get -u google.golang.org/protobuf

go get -u go.etcd.io/etcd/client/v3

```

### 创建数据库

登录`mysql`终端，创建数据库：
```sql
create database crud_list default charset utf8mb4;
```

### 连接数据库

`./conf/config.yaml` 为 `api-gateway` 层配置文件

修改其中
```yaml
user_db:
  DriverName: mysql
  Database: crud_list
  Port: 3306
  UserName: root
  Password: 123456
  Host: 127.0.0.1 #host.docker.internal #
  Charset: utf8mb4
```

`./crud_rpc_server\conf\rpc_server_config.yaml` 为 `crud-service` 层配置文件

修改其中
```yaml
crud_db:
  DriverName: mysql
  Database: crud_list
  Port: 3306
  UserName: root
  Password: 123456
  Host: 127.0.0.1 #host.docker.internal #
  Charset: utf8mb4
```

上述两个数据库可以为两个完全独立的数据库，当然，也可以是同一个。如果需要使用其它数据库，例如 `PostgreSQL, SQLite, SQL Server`。`./mysql_db/connect_db.go` 为数据库配置文件。修改方法，参考 [grom 官方文档 数据库连接](https://gorm.io/zh_CN/docs/connecting_to_the_database.html)


### 连接 redis

`./conf/config.yaml` 

修改其中
```yaml
user_redis:
  Addr: "localhost:6379"
  Password: "123456"
```

### 连接 etcd

```bash
etcd
```

`./conf/config.yaml` 

修改其中  `Endpoint`
```yaml
etcd:
  Endpoint: "127.0.0.1:2379"
  keys:
    crud_rpc: crud_rpc
```

`./crud_rpc_server\conf\rpc_server_config.yaml`

修改其中
```yaml
etcd:
  Endpoint: "127.0.0.1:2379"
```

这里两个配置必须连接同一个`etcd`集群


### 配置端口

`./conf/config.yaml` 
```yaml
server:
  Addr: 0.0.0.0
  Port: 8008
```
这里 `0.0.0.0` 代表运行来自所有 ip 的访问
 
如果需要启动多个微服务，每个微服务以下配置的 `Addr` 要不同

`./crud_rpc_server\conf\rpc_server_config.yaml`
```yaml
server:
  Name: crud_rpc
  Listen: "0.0.0.0:8080"
  Addr: "127.0.0.1:8080"
```


### 运行

#### api-gateway

编译（会自动安装依赖）：
```bash
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go build main.go 
```

运行：
```bash
.\main
```

#### crud 微服务

编译（会自动安装依赖）：
```bash
cd ./crud_rpc_server
go build rpc_server.go 
```

运行：
```bash
.\rpc_server
```

#### 测试

浏览器输入 url:
```txt
http://localhost:8008/ping
```
看到以下内容代表服务运行成功
```js
{"message":"请求成功"}
```
如果希望看到界面，需要用到配套的前端项目📌[配套前端项目地址](https://github.com/bytesc/vue-crud-template)

### 修改 grpc proto

如果修改了`./crud_rpc_server/crud_rpc.proto` 需要重新生成代码

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc --go_out=. --go-grpc_out=. --proto_path=. *proto
```


### 框架官方文档
- https://gorm.io/zh_CN/docs
- https://gin-gonic.com/zh-cn/docs
- https://grpc.io/docs/languages/go/quickstart/
- https://doc.oschina.net/grpc?t=60133
- https://protobuf.dev/programming-guides/proto3/


# 开源许可证

此翻译版本仅供参考，以 LICENSE 文件中的英文版本为准

MIT 开源许可证：

版权所有 (c) 2023 bytesc

特此授权，免费向任何获得本软件及相关文档文件（以下简称“软件”）副本的人提供使用、复制、修改、合并、出版、发行、再许可和/或销售软件的权利，但须遵守以下条件：

上述版权声明和本许可声明应包含在所有副本或实质性部分中。

本软件按“原样”提供，不作任何明示或暗示的保证，包括但不限于适销性、特定用途适用性和非侵权性。在任何情况下，作者或版权持有人均不对因使用本软件而产生的任何索赔、损害或其他责任负责，无论是在合同、侵权或其他方面。