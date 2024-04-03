package service

import (
	"fmt"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go_crud/crud_rpc_server/crud_pb"
	"go_crud/utils/etcd_center"
	"go_crud/utils/mysql_db"
	"gorm.io/gorm"
)

var Database *gorm.DB
var ETCD *clientv3.Client

type GrpcServer struct {
	crud_pb.UnimplementedCRUDServiceServer
}

func Init() {
	var err error

	viper.AddConfigPath("./conf/")
	viper.SetConfigName("rpc_server_config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("配置文件错误 %s", err.Error()))
	}

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

	ETCD = etcd_center.NewEtcdClient("etcd")
}
