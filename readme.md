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

生成go文件，运行gen.sh文件，生成的文件会放./proto/hello文件下
```go
cd proto/hello
protoc --go_out=plugins=grpc:. *.proto
```

演示代码目录结构
```
|—— normal/
    |—— client/
        |—— client.go   // 客户端
    |—— server/
        |—— server.go   // 服务端
|—— keys/                 // 证书目录
    |—— server.key
    |—— server.pem
|—— proto/
    |—— hello/
        |—— hello.proto   // proto描述文件
        |—— hello.pb.go   // proto编译后文件
```

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

演示代码目录结构
```
|—— tls/
    |—— client/
        |—— client.go   // 客户端
    |—— server/
        |—— server.go   // 服务端
|—— keys/                 // 证书目录
    |—— server.key
    |—— server.pem
|—— proto/
    |—— hello/
        |—— hello.proto   // proto描述文件
        |—— hello.pb.go   // proto编译后文件
```

服务端代码
```go
// TLS认证
creds, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
if err != nil {
    grpclog.Fatalf("Failed to generate credentials %v", err)
}
// 实例化grpc Server, 并开启TLS认证
s = grpc.NewServer(grpc.Creds(creds))
```
客户端代码
```go
creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "xx")
if err != nil {
    grpclog.Fatalf("Failed to create TLS credentials, %v", err)
}
conn, err = grpc.Dial(Address, grpc.WithTransportCredentials(creds))
```

### 问题记录

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


## TLS+TOKEN认证

这里我们定义了一个customCredential结构，并实现了两个方法GetRequestMetadata和RequireTransportSecurity。
这是gRPC提供的自定义认证方式，每次RPC调用都会传输认证信息。customCredential其实是实现了grpc/credential
包内的PerRPCCredentials接口。每次调用，token信息会通过请求的metadata传输到服务端。下面具体看一下服务端如
何获取metadata中的信息。
```go

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
    // 解析metada中的信息并验证
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, grpc.Errorf(codes.Unauthenticated, "无Token认证信息")
    }

    var (
        appid  string
        appkey string
    )

    if val, ok := md["appid"]; ok {
        appid = val[0]
    }

    if val, ok := md["appkey"]; ok {
        appkey = val[0]
    }

    if appid != "101010" || appkey != "i am key" {
        return nil, grpc.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
    }

    resp := new(pb.HelloResponse)
    resp.Message = fmt.Sprintf("Hello %s.\nToken info: appid=%s,appkey=%s", in.Name, appid, appkey)

    return resp, nil
}
```


演示代码目录结构
```
|—— token/
    |—— client/
        |—— client.go   // 客户端
    |—— server/
        |—— server.go   // 服务端
|—— keys/                 // 证书目录
    |—— server.key
    |—— server.pem
|—— proto/
    |—— hello/
        |—— hello.proto   // proto描述文件
        |—— hello.pb.go   // proto编译后文件
```


# 拦截器校验
grpc服务端和客户端都提供了interceptor功能，功能类似middleware，很适合在这里处理验证、日志等流程。

在自定义Token认证的示例中，认证信息是由每个服务中的方法处理并认证的，如果有大量的接口方法，这种姿势就太不优雅了，
每个接口实现都要先处理认证信息。这个时候interceptor就可以用来解决了这个问题，在请求被转到具体接口之前处理认证
信息，一处认证，到处无忧。 在客户端，我们增加一个请求日志，记录请求相关的参数和耗时等等。

目录结构
```shell script
|—— interceptor/
    |—— client/
        |—— main.go   // 客户端
    |—— server/
        |—— main.go   // 服务端
|—— keys/             // 证书目录
    |—— server.key
    |—— server.pem
|—— proto/
    |—— hello/
        |—— hello.proto   // proto描述文件
        |—— hello.pb.go   // proto编译后文件
```


服务端代码：
```go
var opts []grpc.ServerOption
// TLS认证
creds, err := credentials.NewServerTLSFromFile("../../keys/server.pem", "../../keys/server.key")
if err != nil {
    grpclog.Fatalf("Failed to generate credentials %v", err)
}
opts = append(opts, grpc.Creds(creds))
// 注册interceptor
opts = append(opts, grpc.UnaryInterceptor(interceptor))
// 实例化grpc Server
s := grpc.NewServer(opts...)

...


// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    err := auth(ctx)
    if err != nil {
        return nil, err
    }
    // 继续处理请求
    return handler(ctx, req)
}

```


客户端代码
```go
var opts []grpc.DialOption
if OpenTLS {
    // TLS连接
    creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "xx")
    if err != nil {
        grpclog.Fatalf("Failed to create TLS credentials %v", err)
    }
    opts = append(opts, grpc.WithTransportCredentials(creds))
} else {
    opts = append(opts, grpc.WithInsecure())
}
// 指定自定义认证
opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
// 指定客户端interceptor
opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

conn, err := grpc.Dial(Address, opts...)
```


