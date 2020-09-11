package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	pb "grpc-demo/proto/hello" // 引入proto包
)

const (
	// Address gRPC服务地址
	Address = "caron:50052"
)

func main() {
	var conn *grpc.ClientConn
	var err error

	tls := true
	// TLS连接  记得把server name改成你写的服务器地址
	if tls {
		creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "xx")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials, %v", err)
		}

		conn, err = grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	} else {
		// 普通链接
		conn, err = grpc.Dial(Address, grpc.WithInsecure())
	}

	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用方法
	req := &pb.HelloRequest{Name: "gRPC"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(res.Message)
}
