package main

import (
	"fmt"
	pb "hellov1/proto/hello" // 引入编译生成的包
	"net"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
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
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("./keys/test.pem", "./keys/test.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	// 实例化grpc Server, 并开启TLS认证
	s := grpc.NewServer(grpc.Creds(creds))

	// 注册HelloService
	pb.RegisterHelloServer(s, HelloService)
	grpclog.Infoln("Listen on " + Address)
	s.Serve(listen)
}
