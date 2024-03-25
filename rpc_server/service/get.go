package service

import (
	"context"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
	"gorm.io/gorm"
)

func (*GrpcServer) Query(c context.Context, req *crud_pb.QueryRequest) (res *crud_pb.QueryResponse, e error) {
	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db.CrudList
	var resultList []*crud_pb.CrudList
	request := req.GetName()
	result := db.Where("name = ?", request).Find(&dataList)

	for _, cl := range dataList {
		// 创建一个CrudList消息实例
		crudList := &crud_pb.CrudList{
			Id:       int32(cl.ID), // gorm.Model中的ID字段
			Name:     cl.Name,
			Level:    cl.Level,
			Email:    cl.Email,
			Phone:    cl.Phone,
			Birthday: cl.Birthday,
			Address:  cl.Address,
		}
		// 将CrudList消息添加到QueryResponse的列表中
		resultList = append(resultList, crudList)
	}

	if result.Error != nil {
		return &crud_pb.QueryResponse{
			Code:    1,
			Message: result.Error.Error(),
			List:    resultList,
		}, nil
	}
	return &crud_pb.QueryResponse{
		Code:    0,
		Message: "success",
		List:    resultList,
	}, nil
}
