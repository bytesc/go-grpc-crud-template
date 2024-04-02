package redis_cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/spf13/viper"
)

func ConnectToRedis(dbName string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString(dbName + ".Addr"), // Redis地址
		Password: viper.GetString(dbName + ".Password"),
		DB:       0,
	})
	return rdb
}

func NewSync(rdb *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)
	return rs
}
