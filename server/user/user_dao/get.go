package user_dao

import (
	"go_crud/mysql_db"
	"gorm.io/gorm"
)

func GetUserByName(name string, DB *gorm.DB) []mysql_db.UserList {
	db := DB.Session(&gorm.Session{NewDB: true})
	var adminDataList []mysql_db.UserList
	db.Where("name = ?", name).Find(&adminDataList)
	return adminDataList
}
