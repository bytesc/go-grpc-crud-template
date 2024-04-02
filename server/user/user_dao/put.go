package user_dao

import (
	"go_crud/utils/mysql_db"
	"gorm.io/gorm"
)

func CreateUser(userData mysql_db.UserList) *gorm.DB {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	result := db.Create(&userData)
	return result
}
