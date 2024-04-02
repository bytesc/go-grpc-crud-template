package crud_rpc

import (
	"fmt"
	"github.com/spf13/viper"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/server/crud_rpc/service"
	"go_crud/utils/etcd_center"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient() crud_pb.CRUDServiceClient {
	//addr := "127.0.0.1:8080"
	addr := etcd_center.GetAddrFromEtcd(viper.GetString("etcd.keys.crud_rpc"), service.ETCD)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	client := crud_pb.NewCRUDServiceClient(conn)
	return client
}
