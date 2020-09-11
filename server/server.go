package main

import (
	"fmt"
	"google.golang.org/grpc/credentials"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "grpc-demo/proto/hello" // 引入编译生成的包
)

const (
	// Address gRPC服务地址
	Address = "caron:50052"
)

// 定义helloService并实现约定的接口
type helloService struct{}

// HelloService Hello服务
var HelloService = helloService{}

// SayHello 实现Hello服务接口
func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)

	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	tls := true
	var s *grpc.Server

	if tls {
		// TLS认证
		creds, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}

		// 实例化grpc Server, 并开启TLS认证
		s = grpc.NewServer(grpc.Creds(creds))
	} else {
		//无认证
		s = grpc.NewServer()
	}

	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)

	fmt.Println("Listen on " + Address)
	s.Serve(listen)
}
