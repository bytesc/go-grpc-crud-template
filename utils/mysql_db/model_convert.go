package mysql_db

import (
	"go_crud/crud_rpc_server/crud_pb"
	"gorm.io/gorm"
)

func CrudListOrmToRpc(dataList []CrudList, resultList *[]*crud_pb.CrudList) {
	for _, cl := range dataList {
		// 创建一个CrudList消息实例
		crudList := &crud_pb.CrudList{
			Id:       int64(cl.ID), // gorm.Model中的ID字段
			Name:     cl.Name,
			Level:    cl.Level,
			Email:    cl.Email,
			Phone:    cl.Phone,
			Birthday: cl.Birthday,
			Address:  cl.Address,
		}
		// 将CrudList消息添加到QueryResponse的列表中
		*resultList = append(*resultList, crudList)
	}
}

func CrudListRpcToOrm(resultList []*crud_pb.CrudList, dataList *[]CrudList) {
	for _, ql := range resultList {
		// 创建一个CrudList结构体实例
		cl := CrudList{
			// gorm.Model中的ID字段，这里假设ID是int类型
			Model: gorm.Model{
				ID: uint(ql.Id),
			},
			Name:     ql.Name,
			Level:    ql.Level,
			Email:    ql.Email,
			Phone:    ql.Phone,
			Birthday: ql.Birthday,
			Address:  ql.Address,
		}
		// 将CrudList结构体添加到切片中
		*dataList = append(*dataList, cl)
	}
}
