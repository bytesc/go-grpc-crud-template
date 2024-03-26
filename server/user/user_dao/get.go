package user_dao

import (
	"context"
	"encoding/json"
	"fmt"
	"go_crud/mysql_db"
	"go_crud/redis_cache"
	"gorm.io/gorm"
	"log"
	"time"
)

func GetUserByName(name string, DB *gorm.DB) []mysql_db.UserList {
	rdb := redis_cache.ConnectToRedis("user_redis")
	var adminDataList []mysql_db.UserList

	// 使用Redis缓存
	ctx := context.Background()
	key := fmt.Sprintf("user:%s", name)
	result, err := rdb.Get(ctx, key).Result()
	if err == nil {
		// 如果Redis中存在缓存，则直接返回
		fmt.Println("从Redis缓存中获取数据")
		if err := json.Unmarshal([]byte(result), &adminDataList); err == nil {
			return adminDataList
		}
	}

	// 如果Redis中没有缓存，则查询MySQL数据库
	fmt.Println("从MySQL数据库中获取数据")
	if err := DB.Where("name = ?", name).Find(&adminDataList).Error; err != nil {
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
