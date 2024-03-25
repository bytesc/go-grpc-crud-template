package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_crud/cmd"
	"go_crud/logger"
	"go_crud/mysql_db"
	"go_crud/server"
	"go_crud/server/crud"
	"go_crud/server/crud_rpc"
	"go_crud/server/files"
	"go_crud/server/midware"
	"go_crud/server/user"
	"go_crud/server/utils"
)

func main() {
	//配置相关
	defer cmd.Clean()
	cmd.Start()

	//数据库相关
	db, err := mysql_db.ConnectToDatabase()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = db.AutoMigrate(&mysql_db.CrudList{})
	err = db.AutoMigrate(&mysql_db.UserList{})

	// 服务相关
	r := server.CreateServer()
	r.Use(cors.Default()) //解决跨域

	log, _ := logger.InitLogger(zap.DebugLevel)
	defer log.Sync()
	r.Use(logger.GinLogger(log), logger.GinRecovery(log, true))

	utils.PingGET(r)

	Router := r.Group("api/refresh", midware.CheckLogin("refresh", db))
	utils.RefreshGET(Router)

	userRouter := r.Group("api/user")
	userRouter.Use(gin.Logger(), gin.Recovery())
	user.LoginPost(userRouter, db)
	user.SignUpPost(userRouter, db)
	user.LogoutGet(userRouter, db)
	user.ChangePwdPost(userRouter, db)
	user.GetPubKey(userRouter)

	crudRouter := r.Group("/api/crud")
	crudRouter.Use(gin.Logger(), gin.Recovery(), midware.CheckLogin("crud", db))
	crud.AddPOST(crudRouter, db)
	crud.DeletePOST(crudRouter, db)
	crud.UpdatePOST(crudRouter, db)
	crud.QueryGET(crudRouter, db)
	crud.QueryPageGET(crudRouter, db)

	crudRpcRouter := r.Group("/api/crud_rpc")
	//, midware.CheckLogin("crud_rpc", db)
	crudRpcRouter.Use(gin.Logger(), gin.Recovery())
	crud_rpc.AddPOST(crudRpcRouter)
	crud_rpc.QueryGET(crudRpcRouter)

	filesRouter := r.Group("/api/files")
	filesRouter.Use(gin.Logger(), gin.Recovery(), midware.CheckLogin("files", db))
	files.FileUploadPOST(filesRouter, db)
	files.BigFileUploadPOST(filesRouter, db)
	files.FileListGet(filesRouter, db)
	files.FileDownload(filesRouter, db)
	files.FileDelete(filesRouter, db)

	//r.Run("0.0.0.0:8088") // 监听并在 0.0.0.0:8088 上启动服务
	// http://127.0.0.1:8088/ping
	//fmt.Println(r)

	r.Run(viper.GetString("server.addr") + ":" + viper.GetString("server.port"))

}
