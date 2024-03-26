package crud_rpc

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_crud/crud_rpc_server/crud_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"time"
)

func NewClient() crud_pb.CRUDServiceClient {
	//addr := "127.0.0.1:8080"
	addr := etcd()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	client := crud_pb.NewCRUDServiceClient(conn)
	return client
}

func etcd() string {
	// 创建etcd客户端
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{
			viper.GetString("etcd.Endpoint"),
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Error creating etcd client: %v \n", err)
	}
	defer etcdClient.Close()

	// 从etcd获取服务地址
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcdClient.Get(ctx, viper.GetString("etcd.keys.crud_rpc"))
	//fmt.Println(resp)
	cancel()
	if err != nil {
		log.Printf("failed to get service address from etcd: %v \n", err)
	}

	if len(resp.Kvs) == 0 {
		log.Println("service address not found in etcd")
	}

	//// 解析服务地址
	//serviceAddr := string(resp.Kvs[0].Value)

	// 随机选择一个服务地址，负载均衡
	rand.Seed(time.Now().UnixNano())
	serviceAddr := string(resp.Kvs[rand.Intn(len(resp.Kvs))].Value)

	return serviceAddr
}
