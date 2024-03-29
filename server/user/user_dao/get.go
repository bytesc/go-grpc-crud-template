package user_dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go_crud/mysql_db"
	"gorm.io/gorm"
	"log"
	"time"
)

func GetUserByName(name string) []mysql_db.UserList {
	lockKey := "user_status_lock_" + name
	if !RedisLock(lockKey) {
		// 如果获取锁失败，则直接返回或者进行重试等策略
		log.Println("Failed to acquire lock for user status update")
		return nil
	}
	defer RedisUnLock(lockKey)

	rdb := RDB
	db := DataBase.Session(&gorm.Session{NewDB: true})
	var adminDataList []mysql_db.UserList

	// 使用Redis缓存
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", name)
	result, err := rdb.Get(ctx, key).Result()
	if err == nil {
		// 如果Redis中存在缓存，则直接返回
		//fmt.Println("从Redis缓存中获取数据")
		if err := json.Unmarshal([]byte(result), &adminDataList); err == nil {
			return adminDataList
		}
	}
	// 如果Redis中没有缓存，则查询MySQL数据库
	//fmt.Println("从MySQL数据库中获取数据")
	if err := db.Where("name = ?", name).Find(&adminDataList).Error; err != nil {
		return nil
	}
	// 将查询结果缓存到Redis
	data, err := json.Marshal(adminDataList)
	err = rdb.Set(ctx, key, data, 5*time.Minute).Err()
	if err != nil {
		log.Println("Redis缓存失败：", err.Error())
	}

	return adminDataList
}
