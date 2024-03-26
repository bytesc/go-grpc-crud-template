package service

import (
	"context"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/mysql_db"
	"gorm.io/gorm"
)

func (*GrpcServer) Query(c context.Context, req *crud_pb.QueryRequest) (res *crud_pb.QueryResponse, e error) {
	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db.CrudList
	var resultList []*crud_pb.CrudList
	name := req.GetName()
	result := db.Where("name = ?", name).Find(&dataList)

	mysql_db.CrudListOrmToRpc(dataList, &resultList)

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

func (*GrpcServer) QueryPage(c context.Context, req *crud_pb.QueryPageRequest) (res *crud_pb.QueryPageResponse, e error) {
	db := Database.Session(&gorm.Session{NewDB: true})
	var dataList []mysql_db.CrudList
	var resultList []*crud_pb.CrudList
	pageSize := int(req.GetPageSize())
	pageNum := int(req.GetPageNum())

	db = db.Model(dataList)
	allQueries := req.QueryParams
	for key, value := range allQueries {
		db.Where(key+" LIKE ?", "%"+value+"%")
	}

	var total int64 //保存数据条数
	db.Count(&total).Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&dataList)
	//limit(-1)表示查询全部数据，offset表示跳过多少条

	mysql_db.CrudListOrmToRpc(dataList, &resultList)
	//fmt.Println(resultList)
	if len(dataList) == 0 { //没有查到
		return &crud_pb.QueryPageResponse{
			Code:    1,
			Total:   0,
			Message: "没有查到",
			List:    resultList,
		}, nil
	} else {
		return &crud_pb.QueryPageResponse{
			Code:    0,
			Message: "查询成功",
			Total:   total,
			List:    resultList,
		}, nil
	}

}
