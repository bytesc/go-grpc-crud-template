package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/crud_rpc_server/service"
	"go_crud/utils/etcd_center"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	service.Init()

	// 监听端口
	listen, err := net.Listen("tcp", viper.GetString("server.Listen"))
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}

	// 创建一个gRPC服务器实例。
	s := grpc.NewServer()
	server := service.GrpcServer{}
	// 将server结构体注册为gRPC服务。
	crud_pb.RegisterCRUDServiceServer(s, &server)
	fmt.Println("grpc server running :" + viper.GetString("server.Listen"))

	go etcd_center.RegisterAddrToEtcd(
		viper.GetString("server.Name"),
		viper.GetString("server.Addr"),
		service.ETCD,
	)

	// 开始处理客户端请求。
	err = s.Serve(listen)
}
