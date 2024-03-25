package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
	"go_crud/rpc_server/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

func main() {
	var err error

	viper.AddConfigPath("../conf/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("配置文件错误 %s", err.Error()))
	}

	service.Database, err = mysql_db.ConnectToDatabase()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = service.Database.AutoMigrate(&mysql_db.CrudList{})

	// 监听端口
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}

	// 创建一个gRPC服务器实例。
	s := grpc.NewServer()
	server := service.GrpcServer{}
	// 将server结构体注册为gRPC服务。
	crud_pb.RegisterCRUDServiceServer(s, &server)
	fmt.Println("grpc server running :8080")
	// 开始处理客户端请求。
	err = s.Serve(listen)
}
