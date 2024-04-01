package user_dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"go_crud/mysql_db"
	"go_crud/redis_cache"
	"gorm.io/gorm"
	"log"
)

var RDB *redis.Client
var DataBase *gorm.DB
var RedSyncLock *redsync.Redsync

func Init() {
	RDB = redis_cache.ConnectToRedis("user_redis")
	RedSyncLock = redis_cache.NewSync(RDB)
	var err error
	DataBase, err = mysql_db.ConnectToDatabase("user_db")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	err = DataBase.AutoMigrate(&mysql_db.UserList{})
	if err != nil {
		fmt.Println("Error init database:", err)
		return
	}
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

func clearEntireRedisCache() {
	// 使用context.Background()作为请求的上下文
	ctx := context.Background()

	// 使用FlushDB命令清空当前数据库
	err := RDB.FlushDB(ctx).Err()
	if err != nil {
		log.Printf("Failed to clear entire Redis cache: %v", err)
	} else {
		log.Println("Successfully cleared entire Redis cache")
	}
}
