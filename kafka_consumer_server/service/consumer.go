package service

import (
	"container/heap"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func MsgConsumer() {
	for {
		// 读取下一条消息
		ctx := context.Background()
		msg, err := KfReader.ReadMessage(ctx)
		if err != nil {
			log.Printf("error while receiving message: %s\n", err)
			break
		}

		// 解析消息
		var m Message
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			log.Printf("error while unmarshaling message: %s\n", err)
			continue
		}

		//// 执行清除缓存操作
		//service.ClearNameRedisCache(m.Name)

		// 创建一个延时任务，并推入堆中
		task := &DelayedTask{
			Name:      m.Name,
			Timestamp: m.Timestamp,
		}
		heap.Push(&TaskHeap, task)

		// 打印时间戳
		fmt.Printf("Message timestamp: %d\n", m.Timestamp)
	}
}

func TaskWorker() {
	for {
		// 处理堆中的任务
		for TaskHeap.Len() > 0 {
			task := heap.Pop(&TaskHeap).(*DelayedTask)

			// 等待直到任务的时间戳到达
			now := time.Now().Unix()
			if now < task.Timestamp {
				time.Sleep(time.Unix(task.Timestamp, 0).Sub(time.Now()))
			}

			// 执行清除缓存操作
			ClearNameRedisCache(task.Name)
			fmt.Printf("Cleared cache for %s at timestamp: %d\n", task.Name, task.Timestamp)
		}
	}
}
