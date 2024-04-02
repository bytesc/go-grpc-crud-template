package service

import (
	"github.com/go-redsync/redsync/v4"
	"time"
)

func GetRedLock(name string) *redsync.Mutex {
	lock := RedSyncLock.NewMutex("lock:user:"+name,
		redsync.WithExpiry(10*time.Second),
		redsync.WithTries(15), // 设置尝试获取锁的最大次数为5次
		redsync.WithRetryDelay(1*time.Second))
	return lock
}

// ContinueLock 启动一个协程来定期检查并延时锁
func ContinueLock(lock *redsync.Mutex) {
	for {
		time.Sleep(5 * time.Second)
		// 尝试延时锁
		if ok, err := lock.Extend(); ok != true {
			if err != nil {
				//log.Println("延时锁失败:", err)
			}
			break
		}
	}
}
