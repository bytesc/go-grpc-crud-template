package user_dao

import (
	"github.com/go-redsync/redsync/v4"
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

	lock := RedSyncLock.NewMutex("lock:user:"+userData.Name,
		redsync.WithExpiry(10*time.Second),
		redsync.WithTries(5), // 设置尝试获取锁的最大次数为5次
		redsync.WithRetryDelay(1*time.Second))
	// 尝试获取锁
	if err := lock.Lock(); err != nil {
		log.Println("获取锁失败:", err)
		return false
	}
	defer lock.Unlock() // 函数结束时释放锁
	// 启动一个协程来定期检查并延时锁
	//go func() {
	//	for {
	//		time.Sleep(3 * time.Second)
	//		// 尝试延时锁
	//		if ok, err := lock.Extend(); ok != true {
	//			log.Println("延时锁失败:", err)
	//			break
	//		}
	//	}
	//}()

	db := DataBase.Session(&gorm.Session{NewDB: true})
	userData.Status = status
	db.Save(&userData)
	//time.Sleep(time.Minute)
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
