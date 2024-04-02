package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go_crud/utils/kafka_mq"
	"go_crud/utils/redis_cache"
	"log"
)

var RDB *redis.Client
var KF_READER *kafka.Reader

func Init() {
	var err error

	viper.AddConfigPath("./conf/")
	viper.SetConfigName("kafka_server_config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("配置文件错误 %s", err.Error()))
	}

	RDB = redis_cache.ConnectToRedis("user_redis")

	KF_READER = kafka_mq.NewKafkaReader("user_cache_kafka")
}

func ClearNameRedisCache(name string) {
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
