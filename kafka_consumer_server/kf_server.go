package main

import (
	"go_crud/kafka_consumer_server/service"
	"log"
	"sync"
)

func main() {

	service.Init()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		service.MsgConsumer()
	}()

	go func() {
		defer wg.Done()
		service.TaskWorker()
	}()

	log.Println("kafka consumer running")
	wg.Wait()
}
