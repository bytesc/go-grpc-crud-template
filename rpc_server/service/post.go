package service

import (
	"context"
	"go_crud/rpc_server/crud_pb"
	"gorm.io/gorm"
)

func (*GrpcServer) Add(c context.Context, req *crud_pb.AddRequest) (res *crud_pb.AddResponse, e error) {

	db := Database.Session(&gorm.Session{NewDB: true})
	request := req.GetList()
	result := db.Create(&request)
	//fmt.Println(result)

	if result.Error != nil {
		return &crud_pb.AddResponse{
			Code:    1,
			Message: result.Error.Error(),
		}, nil
	}
	return &crud_pb.AddResponse{
		Code:    0,
		Message: "success",
	}, nil
}
