package user_dao

import (
	"context"
	"fmt"
	"go_crud/mysql_db"
	"gorm.io/gorm"
	"log"
	"time"
)

func RecordPasswordWrong(userData mysql_db.UserList, DB *gorm.DB, tries uint) {
	db := DB.Session(&gorm.Session{NewDB: true})
	userData.PasswordTry = tries
	if userData.PasswordTry >= 10 {
		userData.LockedUntil = time.Now().Add(time.Hour)
		userData.PasswordTry = 0
	}
	db.Save(&userData)
	clearRedisCache(userData.Name)
}

func SetUserStatus(userData mysql_db.UserList, DB *gorm.DB, status string) {
	db := DB.Session(&gorm.Session{NewDB: true})
	userData.Status = status
	db.Save(&userData)
	clearRedisCache(userData.Name)
}

func SetUserPwd(userData mysql_db.UserList, DB *gorm.DB, newPwd string) {
	db := DB.Session(&gorm.Session{NewDB: true})
	userData.Password = newPwd
	db.Save(&userData)
	clearRedisCache(userData.Name)
}

func clearRedisCache(name string) {
	// 连接到 Redis
	rdb := RDB
	// 删除与用户 ID 相关的缓存
	redisKey := fmt.Sprintf("user:%s", name)
	ctx := context.Background()
	_, err := rdb.Del(ctx, redisKey).Result()
	if err != nil {
		log.Printf("Failed to clear Redis cache for user ID %s: %v", name, err)
	}
}
