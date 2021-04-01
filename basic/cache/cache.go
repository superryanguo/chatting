package cache

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/superryanguo/chatting/basic/rediser"
)

var (
	rc          *redis.Client
	m           sync.RWMutex
	ExpiredDate = 3600 * 24 * 30 * time.Second
)

func SaveToCache(key string, val []byte) (err error) {
	if rc == nil {
		log.Debug("redis.client un-init")
	}
	log.Debug("SaveCache key=", key, " val=", string(val))

	if err = rc.Set(key, val, ExpiredDate).Err(); err != nil {
		log.Debug("[saveToCache] err:", err.Error())
		return err
	}
	return
}

func SaveToCacheDate(key string, val []byte, date time.Duration) (err error) {
	if rc == nil {
		log.Debug("redis.client un-init")
	}
	log.Debug("SaveCache key=", key, " val=", string(val))

	if err = rc.Set(key, val, date).Err(); err != nil {
		log.Debug("[saveToCache] err:", err.Error())
		return err
	}
	return
}

func DelFromCache(key string) (err error) {
	if rc == nil {
		log.Debug("redis.client un-init")
	}
	if err = rc.Del(key).Err(); err != nil {
		log.Debug("[delFromCache] err:", err.Error())
		return err
	}
	return
}

func GetFromCache(key string) ([]byte, error) {
	if rc == nil {
		log.Debug("redis.client un-init")
	}
	//val, err := rc.Get(key).Result()
	//TODO: any bug here if the key is not exist?
	val, err := rc.Get(key).Bytes()
	if err != nil { //redis.Nil if key not exist
		log.Debug("[getFromCache]can't find key, err=", err.Error())
		return nil, err
	}

	return val, nil
}

func Init() {
	m.Lock()
	defer m.Unlock()

	rc = rediser.GetRedis()
}
