package user_dao

import (
	"github.com/go-redis/redis/v8"
	"go_crud/redis_cache"
)

var RDB *redis.Client

func Init() {
	RDB = redis_cache.ConnectToRedis("user_redis")
}
