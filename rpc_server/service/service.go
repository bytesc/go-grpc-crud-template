package service

import (
	"go_crud/rpc_server/crud_pb"
	"gorm.io/gorm"
)

var Database *gorm.DB

type GrpcServer struct {
	crud_pb.UnimplementedCRUDServiceServer
}
