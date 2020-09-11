# 简介
gRPC 是一个高性能、开源、通用的RPC框架，由Google推出，基于HTTP2协议标准设计开发，默认采用Protocol Buffers
数据序列化协议，支持多种开发语言。gRPC提供了一种简单的方法来精确的定义服务，并且为客户端和服务端自动生成可靠的功能库。

在gRPC客户端可以直接调用不同服务器上的远程程序，使用姿势看起来就像调用本地程序一样，很容易去构建分布式应用和服务。
和很多RPC系统一样，服务端负责实现定义好的接口并处理客户端的请求，客户端根据接口描述直接调用需要的服务。客户端和服
务端可以分别使用gRPC支持的不同语言实现。

# 特性
- 强大的IDL  
gRPC使用ProtoBuf来定义服务，ProtoBuf是由Google开发的一种数据序列化协议（类似于XML、JSON、hessian）。
ProtoBuf能够将数据进行序列化，并广泛应用在数据存储、通信协议等方面。

- 多语言支持  

gRPC支持多种语言，并能够基于语言自动生成客户端和服务端功能库。目前已提供了C版本grpc、Java版本grpc-java 和 Go版本grpc-go，其它语言的版本正
在积极开发中，其中，grpc支持C、C++、Node.js、Python、Ruby、Objective-C、PHP和C#等语言，grpc-java已经支持Android开发。

- HTTP2  
gRPC基于HTTP2标准设计，所以相对于其他RPC框架，gRPC带来了更多强大功能，如双向流、头部压缩、多复用请求等。
这些功能给移动设备带来重大益处，如节省带宽、降低TCP链接次数、节省CPU使用和延长电池寿命等。同时，gRPC还能够提高了云端服务和Web应用的性能。
gRPC既能够在客户端应用，也能够在服务器端应用，从而以透明的方式实现客户端和服务器端的通信和简化通信系统的构建。



# 安装
安装protoc工具，从https://github.com/protocolbuffers/protobuf/releases/tag/v3.9.0页面上选择
直接下载软件包，将protoc解压到$GOPATH/bin路径下
```shell script
go get github.com/golang/protobuf/proto
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
```
上面安装好后，会在GOPATH/bin下生成protoc-gen-go

# GRPC认证方式

## TLS认证示例
### 证书制作
制作私钥 (.key)
```shell script
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
$ openssl ecparam -genkey -name secp384r1 -out server.key
```

自签名公钥(x509) (PEM-encodings .pem|.crt)
```shell script
$ openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
```
自定义信息
```shell script
-----
Country Name (2 letter code) [AU]:CN  //这个是中国的缩写
State or Province Name (full name) [Some-State]:XxXx  //省份
Locality Name (eg, city) []:XxXx  //城市
Organization Name (eg, company) [Internet Widgits Pty Ltd]:XX Co. Ltd  //公司名称
Organizational Unit Name (eg, section) []:Dev   //部门名称
Common Name (e.g. server FQDN or YOUR name) []:server name   //服务名称 例如www.topgoer.com
Email Address []:xxx@xxx.com  //邮箱地址
```

目录结构
```
|—— hello-tls/
    |—— client/
        |—— main.go   // 客户端
    |—— server/
        |—— main.go   // 服务端
|—— keys/                 // 证书目录
    |—— server.key
    |—— server.pem
|—— proto/
    |—— hello/
        |—— hello.proto   // proto描述文件
        |—— hello.pb.go   // proto编译后文件
```


# 问题记录

客户端连接时报错
```shell script
x509: cannot validate certificate for 10.30.0.163 because it doesn't contain any IP SANs
```
解决方法：  
创建证书时使用IP别名
服务端创建证书时，使用IP别名（根据实际情况随便起一个，例如:caron）
修改host文件/etc/hosts
在文件中添加行：
```
10.30.0.163 caron
```


