package service

import (
	"container/heap"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go_crud/utils/kafka_mq"
	"go_crud/utils/redis_cache"
	"log"
)

// Message 结构体用来解析Kafka消息中的JSON
type Message struct {
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

var RDB *redis.Client
var KfReader *kafka.Reader

var TaskHeap DelayedTaskHeap

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

	KfReader = kafka_mq.NewKafkaConsumer("user_cache_kafka")

	heap.Init(&TaskHeap)
}

func ClearNameRedisCache(name string) bool {
	// 连接到 Redis
	rdb := RDB
	ctx := context.Background()
	_, err := rdb.Del(ctx, name).Result()
	if err != nil {
		log.Printf("Failed to clear Redis cache for user ID %s: %v", name, err)
		return false
	}
	return true
}
