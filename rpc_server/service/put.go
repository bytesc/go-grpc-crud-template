package service

import (
	"context"
	"go_crud/mysql_db"
	"go_crud/rpc_server/crud_pb"
	"gorm.io/gorm"
)

func (*GrpcServer) Update(c context.Context, req *crud_pb.UpdateRequest) (res *crud_pb.UpdateResponse, e error) {

	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db.CrudList
	var resultList []*crud_pb.CrudList
	id := req.GetId()

	var tmpDataList []mysql_db.CrudList
	result := db.Where("id = ?", id).Find(&tmpDataList) //数据库查找
	//fmt.Println(tmpDataList)
	//db.Where("name = ?", "张三").Find(&dataList)
	//fmt.Println(dataList[0].ID)

	if result.Error != nil {
		return &crud_pb.UpdateResponse{
			Code:    1,
			Message: result.Error.Error(),
		}, nil
	}

	if len(tmpDataList) == 0 {
		return &crud_pb.UpdateResponse{
			Code:    1,
			Message: "更新失败，数据不存在",
		}, nil
	} else {
		resultList = append(resultList, req.List)
		mysql_db.CrudListRpcToOrm(resultList, &dataList)
		dataList[0].ID = uint(id)
		result := db.Where("id = ?", id).Updates(&dataList[0])
		if result.RowsAffected == 0 {
			return &crud_pb.UpdateResponse{
				Code:    1,
				Message: "更新失败",
			}, nil
		} else if result.Error != nil {
			return &crud_pb.UpdateResponse{
				Code:    1,
				Message: "更新失败" + result.Error.Error(),
			}, nil
		} else {
			return &crud_pb.UpdateResponse{
				Code:    0,
				Message: "更新成功",
			}, nil
		}

	}
}
