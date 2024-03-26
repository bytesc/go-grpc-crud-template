package user_dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go_crud/redis_cache"
	"log"
)

var RDB *redis.Client

func Init() {
	RDB = redis_cache.ConnectToRedis("user_redis")
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
