package user_dao

import (
	"go_crud/server/user/user_dao/service"
	"go_crud/utils/mysql_db"
	"gorm.io/gorm"
	"log"
)

func GetUserByName(name string) []mysql_db.UserList {

	lock := service.GetRedLock(name)
	if err := lock.Lock(); err != nil {
		log.Println("获取锁失败:", err)
		return nil
	}
	defer lock.Unlock()
	go service.ContinueLock(lock)

	db := service.DataBase.Session(&gorm.Session{NewDB: true})
	var adminDataList []mysql_db.UserList

	adminDataList = service.GetFromRedisByName(name)
	// 如果Redis中没有缓存，则查询MySQL数据库
	if adminDataList != nil {
		return adminDataList
	} else {
		//fmt.Println("从MySQL数据库中获取数据")
		err := db.Where("name = ?", name).Find(&adminDataList).Error
		if err != nil {
			return nil
		}
		service.SetNameToRedis(name, adminDataList)
	}

	return adminDataList
}
