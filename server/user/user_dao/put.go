package user_dao

import (
	"go_crud/mysql_db"
	"gorm.io/gorm"
)

func CreateUser(userData mysql_db.UserList, DB *gorm.DB) *gorm.DB {
	db := DB.Session(&gorm.Session{NewDB: true})
	result := db.Create(&userData)
	return result
}
