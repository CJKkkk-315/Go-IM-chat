package utils

import (
	"github.com/go-redis/redis/v8"
	"net"
	"sync"
	"time"
)
var HeartLock sync.Mutex
var RB *redis.Client
var OnlineMap map[int] net.Conn
var HeartMap map[int] time.Time
func heartCheck() {
	for {
		time.Sleep(time.Second)
		HeartLock.Lock()
		for k,v := range HeartMap {
			if time.Now().Sub(v).Seconds() > 5 {
				delete(OnlineMap, k)
			}
		}
		HeartLock.Unlock()
	}

}
func InitAll() {
	RB = redis.NewClient(&redis.Options{
			Addr:               "127.0.0.1:6379",
		})
	OnlineMap = make(map[int] net.Conn)
	HeartMap = make(map[int] time.Time)

	go heartCheck()
}