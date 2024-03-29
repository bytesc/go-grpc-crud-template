package user_dao

import (
	"go_crud/mysql_db"
	"gorm.io/gorm"
	"log"
	"time"
)

func RecordPasswordWrong(userData mysql_db.UserList, tries uint) bool {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.PasswordTry = tries
	if userData.PasswordTry >= 10 {
		userData.LockedUntil = time.Now().Add(time.Hour)
		userData.PasswordTry = 0
	}
	go clearRedisCache(userData.Name)
	db.Save(&userData)
	go clearRedisCache(userData.Name)
	return true
}

func SetUserStatus(userData mysql_db.UserList, status string) bool {
	// 定义锁的键值，通常使用被操作的数据的唯一标识
	lockKey := "user_status_lock_" + userData.Name
	// 尝试获取锁
	if !RedisLock(lockKey) {
		// 如果获取锁失败，则直接返回或者进行重试等策略
		log.Println("Failed to acquire lock for user status update")
		return false
	}
	// 处理完毕后一定要释放锁
	defer RedisUnLock(lockKey)

	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.Status = status
	db.Save(&userData)
	//time.Sleep(1 * time.Minute)
	clearRedisCache(userData.Name)
	return true
}

func SetUserPwd(userData mysql_db.UserList, newPwd string) bool {
	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.Password = newPwd
	go clearRedisCache(userData.Name)
	db.Save(&userData)
	go clearRedisCache(userData.Name)
	return true
}
