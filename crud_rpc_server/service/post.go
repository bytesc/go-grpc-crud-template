package service

import (
	"context"
	"go_crud/crud_rpc_server/crud_pb"
	mysql_db2 "go_crud/utils/mysql_db"
	"gorm.io/gorm"
)

func (*GrpcServer) Add(c context.Context, req *crud_pb.AddRequest) (res *crud_pb.AddResponse, e error) {

	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db2.CrudList
	var resultList []*crud_pb.CrudList
	resultList = append(resultList, req.GetList())
	mysql_db2.CrudListRpcToOrm(resultList, &dataList)
	result := db.Create(&dataList[0])
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
