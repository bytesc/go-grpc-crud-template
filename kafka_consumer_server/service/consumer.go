package service

import (
	"container/heap"
	"context"
	"encoding/json"
	"log"
	"time"
)

var MsgSignal chan int

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

		//log.Printf("received message: %s\n",m.Name)

		//// 执行清除缓存操作
		//service.ClearNameRedisCache(m.Name)

		// 创建一个延时任务，并推入堆中
		task := &DelayedTask{
			Name:      m.Name,
			Timestamp: m.Timestamp,
		}
		heap.Push(&TaskHeap, task)
		MsgSignal <- 1

		if err := KfReader.CommitMessages(ctx, msg); err != nil {
			log.Printf("error while committing message offset: %s\n", err)
			continue
		}
		// 打印时间戳
		log.Printf("Message timestamp: %d\n", m.Timestamp)
	}
}

func TaskWorker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select{
		case <-MsgSignal:
			task()
		case <-ticker.C:
			task()
		}

	}
}


func task(){
	// 处理堆中的任务
	for TaskHeap.Len() > 0 {
		task := heap.Pop(&TaskHeap).(*DelayedTask)

		if time.Now().UnixNano() < task.Timestamp {
			heap.Push(&TaskHeap, task)
			continue
		}

		// 执行清除缓存操作
		ok := ClearNameRedisCache(task.Name)
		if !ok {
			// 删除失败，延迟一秒重新推入队列
			task.Timestamp = time.Now().Add(1 * time.Second).UnixNano()
			heap.Push(&TaskHeap, task)
			log.Printf("Failed to clear cache for %s, retrying in 1 second\n", task.Name)
			continue
		}
		log.Printf("Cleared cache for %s at timestamp: %d\n", task.Name, task.Timestamp)
	}
}

