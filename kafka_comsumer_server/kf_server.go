package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_crud/kafka_comsumer_server/service"
	"log"
)

func main() {

	service.Init()

	for {
		// 读取下一条消息
		ctx := context.Background()
		msg, err := service.KfReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while receiving message: %s\n", err)
			break
		}

		// 解析消息
		var m service.Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Printf("error while unmarshaling message: %s\n", err)
			continue
		}

		// 执行清除缓存操作
		service.ClearNameRedisCache(m.Name)

		// 打印时间戳
		fmt.Printf("Message timestamp: %d\n", m.Timestamp)
	}

	// 关闭reader
	if err := service.KfReader.Close(); err != nil {
		log.Println("failed to close reader:", err)
	}
}
