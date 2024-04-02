package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/segmentio/kafka-go"
	"go_crud/utils/kafka_mq"
	"go_crud/utils/mysql_db"
	"go_crud/utils/redis_cache"
	"gorm.io/gorm"
)

var RDB *redis.Client
var DataBase *gorm.DB
var RedSyncLock *redsync.Redsync

var KfWriter *kafka.Writer

func Init() {
	RDB = redis_cache.ConnectToRedis("user_redis")
	RedSyncLock = redis_cache.NewSync(redis_cache.ConnectToRedis("lock_redis"))
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

	KfWriter = kafka_mq.NewKafkaProducer("user_cache_kafka")

	ClearEntireRedisCache()
}
