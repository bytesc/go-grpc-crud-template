package service

import (
	"context"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/utils/mysql_db"
	"gorm.io/gorm"
)

func (*GrpcServer) Delete(c context.Context, req *crud_pb.DeleteRequest) (res *crud_pb.DeleteResponse, e error) {
	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db.CrudList
	//var resultList []*crud_pb.CrudList
	id := req.GetId()

	var tmpDataList []mysql_db.CrudList

	db.Where("id = ?", id).Find(&tmpDataList) //数据库查找
	if len(tmpDataList) == 0 {                //没有查到
		return &crud_pb.DeleteResponse{
			Code:    1,
			Message: "更新失败，数据不存在",
		}, nil
	} else {
		result := db.Where("id = ?", id).Delete(&dataList) //删除对应数据
		if result.RowsAffected == 0 {
			return &crud_pb.DeleteResponse{
				Code:    1,
				Message: "删除失败",
			}, nil
		} else if result.Error != nil {
			return &crud_pb.DeleteResponse{
				Code:    1,
				Message: "删除失败" + result.Error.Error(),
			}, nil
		} else {
			return &crud_pb.DeleteResponse{
				Code:    0,
				Message: "删除成功",
			}, nil
		}
	}
}
