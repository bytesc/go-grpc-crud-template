package user_dao

import (
	"go_crud/mysql_db"
	"gorm.io/gorm"
	"time"
)

func RecordPasswordWrong(userData mysql_db.UserList, tries uint) {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.PasswordTry = tries
	if userData.PasswordTry >= 10 {
		userData.LockedUntil = time.Now().Add(time.Hour)
		userData.PasswordTry = 0
	}
	go clearRedisCache(userData.Name)
	db.Save(&userData)
	go clearRedisCache(userData.Name)
}

func SetUserStatus(userData mysql_db.UserList, status string) {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.Status = status
	go clearRedisCache(userData.Name)
	db.Save(&userData)
	go clearRedisCache(userData.Name)
}

func SetUserPwd(userData mysql_db.UserList, newPwd string) {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.Password = newPwd
	go clearRedisCache(userData.Name)
	db.Save(&userData)
	go clearRedisCache(userData.Name)
}
