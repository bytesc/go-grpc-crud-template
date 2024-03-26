package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/crud_rpc_server/service"
	"go_crud/mysql_db"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

func main() {
	var err error

	viper.AddConfigPath("./conf/")
	viper.SetConfigName("rpc_server_config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("配置文件错误 %s", err.Error()))
	}

	service.Database, err = mysql_db.ConnectToDatabase("crud_db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = service.Database.AutoMigrate(&mysql_db.CrudList{})
	if err != nil {
		fmt.Println("Error init database:", err)
		return
	}

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

	go etcd()

	// 开始处理客户端请求。
	err = s.Serve(listen)
}

func etcd() {
	// 初始化etcd客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			viper.GetString("etcd.Endpoint"),
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Error creating etcd client: %v \n", err)
	}
	defer cli.Close()

	// 获取服务配置
	serviceName := viper.GetString("server.Name")
	serviceAddr := viper.GetString("server.Addr")

	// 注册服务到etcd
	leaseResp, err := cli.Grant(context.Background(), 10)
	if err != nil {
		log.Printf("Error granting lease: %v \n", err)
	}

	// 设置键值对，其中键通常是服务名称，值是服务地址
	putResp, err := cli.Put(context.Background(), serviceName, serviceAddr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		log.Printf("Error putting service to etcd: %v \n", err)
	}
	log.Println(putResp, serviceName, serviceAddr)

	// 保持心跳，以续约租约
	keepAliveChan, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		log.Printf("Error keeping lease alive: %v \n", err)
	}

	for {
		select {
		case ka := <-keepAliveChan:
			if ka == nil {
				log.Println("Lease expired or KeepAlive channel closed")
				return
			}
			//log.Println(ka.String())
		}
	}
}