package kafka_mq

import (
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func NewKafkaConsumer(mqName string) *kafka.Reader {
	// 创建Kafka reader
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   viper.GetStringSlice(mqName + ".broker"),
		Topic:     viper.GetString(mqName + ".topic"),
		GroupID:   viper.GetString(mqName + ".group_id"),
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	return r
}

func NewKafkaProducer(mqName string) *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: viper.GetStringSlice(mqName + ".broker"),
		Topic:   viper.GetString(mqName + ".topic"),
	})

	return writer
}
