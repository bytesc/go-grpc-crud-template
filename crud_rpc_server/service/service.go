package service

import (
	"fmt"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/mysql_db"
	"gorm.io/gorm"
)

var Database *gorm.DB

type GrpcServer struct {
	crud_pb.UnimplementedCRUDServiceServer
}

func Init() {
	var err error
	Database, err = mysql_db.ConnectToDatabase("crud_db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = Database.AutoMigrate(&mysql_db.CrudList{})
	if err != nil {
		fmt.Println("Error init database:", err)
		return
	}
}
