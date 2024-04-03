package service

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	mqService "go_crud/kafka_consumer_server/service"
	"log"
	"time"
)

func SendMsgToMq(name string, time time.Time) {
	msg := &mqService.Message{
		Name:      "user:" + name,
		Timestamp: time.Unix(),
	}
	// 将消息序列化为JSON
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
	}

	// 发送消息到Kafka
	err = KfWriter.WriteMessages(context.Background(), kafka.Message{
		Value: jsonBytes,
	})
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	log.Println("Message sent to Kafka successfully")
}
