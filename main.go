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
	"go_crud/server/crud_rpc"
	"go_crud/server/files"
	"go_crud/server/midware"
	"go_crud/server/user"
	"go_crud/server/user/user_dao"
	"go_crud/server/utils"
	"log"
)

func main() {
	//配置相关
	defer cmd.Clean()
	cmd.Start()

	//数据库相关
	userDb, err := mysql_db.ConnectToDatabase("user_db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = userDb.AutoMigrate(&mysql_db.UserList{})

	if err != nil {
		fmt.Println("Error init database:", err)
		return
	}

	user_dao.Init()

	// 服务相关
	r := server.CreateServer()
	r.Use(cors.Default()) //解决跨域

	serverLogger, _ := logger.InitLogger(zap.DebugLevel)
	defer serverLogger.Sync()
	r.Use(logger.GinLogger(serverLogger), logger.GinRecovery(serverLogger, true))

	utils.PingGET(r)

	Router := r.Group("api/refresh", midware.CheckLogin("refresh", userDb))
	utils.RefreshGET(Router)

	userRouter := r.Group("api/user")
	userRouter.Use(gin.Logger(), gin.Recovery())
	user.LoginPost(userRouter, userDb)
	user.SignUpPost(userRouter, userDb)
	user.LogoutGet(userRouter, userDb)
	user.ChangePwdPost(userRouter, userDb)
	user.GetPubKey(userRouter)

	//crudRpcRouter := r.Group("/api/crud")
	crudRpcRouter := r.Group("/api/crud", midware.CheckLogin("crud", userDb))
	crudRpcRouter.Use(gin.Logger(), gin.Recovery())
	crud_rpc.AddPOST(crudRpcRouter)
	crud_rpc.QueryGET(crudRpcRouter)
	crud_rpc.QueryPageGET(crudRpcRouter)
	crud_rpc.DeletePOST(crudRpcRouter)
	crud_rpc.UpdatePOST(crudRpcRouter)

	filesRouter := r.Group("/api/files")
	filesRouter.Use(gin.Logger(), gin.Recovery(), midware.CheckLogin("files", userDb))
	files.FileUploadPOST(filesRouter, nil)
	files.BigFileUploadPOST(filesRouter, nil)
	files.FileListGet(filesRouter, nil)
	files.FileDownload(filesRouter, nil)
	files.FileDelete(filesRouter, nil)

	//r.Run("0.0.0.0:8088") // 监听并在 0.0.0.0:8088 上启动服务
	// http://127.0.0.1:8088/ping
	//fmt.Println(r)

	err = r.Run(viper.GetString("server.addr") + ":" + viper.GetString("server.port"))
	if err != nil {
		log.Println(err.Error())
	}

}
