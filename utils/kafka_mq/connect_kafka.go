package kafka_mq

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func NewKafkaReader(mqName string) *kafka.Reader {
	// 创建Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   viper.GetStringSlice(mqName + ".broker"),
		Topic:     viper.GetString(mqName + ".topic"),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return r
}
