package main

import (
	pb "hellov1/proto/hello" // 引入proto包
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))

	// // 连接
	// conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// TLS连接
	creds, err := credentials.NewClientTLSFromFile("./keys/test.pem", "*.wuzhi555.cc")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))

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
	grpclog.Infoln(res.Message)
}
