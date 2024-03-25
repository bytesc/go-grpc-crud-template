package crud_rpc

import (
	"fmt"
	"go_crud/rpc_server/crud_pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewClient() crud_pb.CRUDServiceClient {
	addr := "127.0.0.1:8080"
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf(fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	client := crud_pb.NewCRUDServiceClient(conn)
	return client
}
