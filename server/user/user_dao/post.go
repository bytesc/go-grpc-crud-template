package user_dao

import (
	"go_crud/server/user/user_dao/service"
	"go_crud/utils/mysql_db"
	"gorm.io/gorm"
	"log"
	"time"
)

func RecordPasswordWrong(userData mysql_db.UserList, tries uint) bool {
	// 方式一 分布式锁 强一致性
	lock := service.GetRedLock(userData.Name)
	// 尝试获取锁
	if err := lock.Lock(); err != nil {
		log.Println("获取锁失败:", err)
		return false
	}
	defer lock.Unlock()
	go service.ContinueLock(lock)

	db := service.DataBase.Session(&gorm.Session{NewDB: true})
	userData.PasswordTry = tries
	if userData.PasswordTry >= 10 {
		userData.LockedUntil = time.Now().Add(time.Hour)
		userData.PasswordTry = 0
	}
	db.Save(&userData)
	go service.ClearNameRedisCache(userData.Name)
	return true
}

func SetUserStatus(userData mysql_db.UserList, status string) bool {
	// 方式一 分布式锁 强一致性
	lock := service.GetRedLock(userData.Name)
	// 尝试获取锁
	if err := lock.Lock(); err != nil {
		log.Println("获取锁失败:", err)
		return false
	}
	defer lock.Unlock()           // 函数结束时释放锁
	go service.ContinueLock(lock) // 启动一个协程来定期检查并延时锁

	db := service.DataBase.Session(&gorm.Session{NewDB: true})
	userData.Status = status
	db.Save(&userData)
	go service.ClearNameRedisCache(userData.Name)
	return true
}

func SetUserPwd(userData mysql_db.UserList, newPwd string) bool {
	db := service.DataBase.Session(&gorm.Session{NewDB: true})
	userData.Password = newPwd
	// 方式二 双删 弱一致性
	go service.ClearNameRedisCache(userData.Name)
	db.Save(&userData)
	go service.SendMsgToMq(userData.Name, time.Now().Add(1000*time.Millisecond))
	return true
}
